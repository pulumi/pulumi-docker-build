// Copyright 2024, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:generate go run go.uber.org/mock/mockgen -typed -package internal -source client.go -destination mockclient_test.go --self_package github.com/pulumi/pulumi-docker-build/provider/internal

package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/distribution/reference"
	buildx "github.com/docker/buildx/build"
	"github.com/docker/buildx/commands"
	controllerapi "github.com/docker/buildx/controller/pb"
	"github.com/docker/buildx/util/dockerutil"
	"github.com/docker/buildx/util/platformutil"
	"github.com/docker/buildx/util/progress"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/docker/api/types/image"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/util/progress/progressui"
	"github.com/regclient/regclient/types/descriptor"
	"github.com/regclient/regclient/types/errs"
	"github.com/regclient/regclient/types/manifest"
	"github.com/regclient/regclient/types/ref"

	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
)

// Client handles all our Docker API calls.
type Client interface {
	Build(ctx context.Context, b Build) (*client.SolveResponse, error)
	BuildKitEnabled() (bool, error)
	Inspect(ctx context.Context, id string) ([]descriptor.Descriptor, error)
	Delete(ctx context.Context, id string) error

	ManifestCreate(ctx context.Context, push bool, target string, refs ...string) error
	ManifestInspect(ctx context.Context, target string) (string, error)
	ManifestDelete(ctx context.Context, target string) error
}

// Build encapsulates all of the user-provider build parameters and options.
type Build interface {
	BuildOptions() controllerapi.BuildOptions
	Inline() string
	ShouldExec() bool
	Secrets() session.Attachable
}

var _ Client = (*cli)(nil)

func newDockerCLI(config *Config) (*command.DockerCli, error) {
	cli, err := command.NewDockerCli(
		command.WithDefaultContextStoreConfig(),
		command.WithContentTrustFromEnv(),
	)
	if err != nil {
		return nil, err
	}

	opts := flags.NewClientOptions()
	if config != nil && config.Host != "" {
		opts.Hosts = append(opts.Hosts, config.Host)
	}
	err = cli.Initialize(opts)
	if err != nil {
		return nil, err
	}

	// TODO: Log some version information for debugging.

	return cli, nil
}

// Build performs a BuildKit build. Returns a map of target names (or one name,
// "default", if no targets were specified) to SolveResponses, which capture
// the build's digest and tags (if any).
func (c *cli) Build(
	ctx context.Context,
	build Build,
) (*client.SolveResponse, error) {
	opts := build.BuildOptions()

	go c.tail(ctx)
	defer contract.IgnoreClose(c)

	if build.ShouldExec() {
		return c.execBuild(build)
	}

	b, err := c.host.builderFor(build)
	if err != nil {
		return nil, err
	}
	printer, err := progress.NewPrinter(ctx, c.w,
		progressui.PlainMode,
		progress.WithDesc(
			fmt.Sprintf("building with %q instance using %s driver", b.name, b.driver),
			fmt.Sprintf("%s:%s", b.driver, b.name),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating printer: %w", err)
	}
	defer func() {
		// Log any warnings when we're done.
		_ = printer.Wait()
		for _, w := range printer.Warnings() {
			b := &bytes.Buffer{}
			fmt.Fprint(b, w.Short)
			for _, d := range w.Detail {
				fmt.Fprintf(b, "\n%s", d)
			}
			provider.GetLogger(ctx).Warning(b.String())
		}
	}()

	cacheFrom := []client.CacheOptionsEntry{}
	for _, c := range opts.CacheFrom {
		if c == nil {
			continue
		}
		cacheFrom = append(cacheFrom, client.CacheOptionsEntry{
			Type:  c.Type,
			Attrs: c.Attrs,
		})
	}
	cacheTo := []client.CacheOptionsEntry{}
	for _, c := range opts.CacheTo {
		if c == nil {
			continue
		}
		cacheTo = append(cacheTo, client.CacheOptionsEntry{
			Type:  c.Type,
			Attrs: c.Attrs,
		})
	}
	exports := []client.ExportEntry{}
	for _, e := range opts.Exports {
		if e == nil {
			continue
		}
		exports = append(exports, client.ExportEntry{
			Type:      e.Type,
			Attrs:     e.Attrs,
			OutputDir: e.Destination,
		})
	}
	platforms, _ := platformutil.Parse(opts.Platforms)
	platforms = platformutil.Dedupe(platforms)

	namedContexts := map[string]buildx.NamedContext{}
	for k, v := range opts.NamedContexts {
		ref, err := reference.ParseNormalizedNamed(k)
		if err != nil {
			return nil, err
		}
		name := strings.TrimSuffix(reference.FamiliarString(ref), ":latest")
		namedContexts[name] = buildx.NamedContext{Path: v}
	}

	ssh, err := controllerapi.CreateSSH(opts.SSH)
	if err != nil {
		return nil, err
	}

	target := opts.Target
	if target == "" {
		target = "default"
	}
	payload := map[string]buildx.Options{
		target: {
			Inputs: buildx.Inputs{
				ContextPath:      opts.ContextPath,
				DockerfilePath:   opts.DockerfileName,
				DockerfileInline: build.Inline(),
				NamedContexts:    namedContexts,
				InStream:         strings.NewReader(""),
			},
			// Disable default provenance for now. Docker's `manifest create`
			// doesn't handle manifests with provenance included; more reason
			// to use imagetools instead.
			Attests:     map[string]*string{"provenance": nil},
			BuildArgs:   opts.BuildArgs,
			CacheFrom:   cacheFrom,
			CacheTo:     cacheTo,
			Exports:     exports,
			ExtraHosts:  opts.ExtraHosts,
			NetworkMode: opts.NetworkMode,
			NoCache:     opts.NoCache,
			Labels:      opts.Labels,
			Platforms:   platforms,
			Pull:        opts.Pull,
			Tags:        opts.Tags,
			Target:      opts.Target,

			Session: []session.Attachable{
				ssh,
				authprovider.NewDockerAuthProvider(c.ConfigFile(), nil),
				build.Secrets(),
			},
		},
	}

	// Perform the build.
	results, err := buildx.Build(
		ctx,
		b.nodes,
		payload,
		dockerutil.NewClient(c),
		filepath.Dir(c.ConfigFile().Filename),
		printer,
	)
	if err != nil {
		c.dumplogs = true
		return nil, err
	}

	return results[target], err
}

// BuildKitEnabled returns true if the client supports buildkit.
func (c *cli) BuildKitEnabled() (bool, error) {
	return c.Cli.BuildKitEnabled()
}

func (c *cli) ManifestCreate(ctx context.Context, push bool, target string, refs ...string) error {
	go c.tail(ctx)
	defer contract.IgnoreClose(c)

	args := []string{
		// "buildx",
		"imagetools",
		"create",
		"--progress=plain",
		"--tag", target,
	}

	if !push {
		args = append(args, "--dry-run")
	}

	args = append(args, refs...)

	cmd := commands.NewRootCmd(os.Args[0], false, c)

	cmd.SetArgs(args)
	cmd.SetErr(c.Err())
	cmd.SetOut(c.Out())

	provider.GetLogger(ctx).Debug(fmt.Sprint("creating manifest with args", args))
	return cmd.ExecuteContext(ctx)
}

func (c *cli) ManifestInspect(ctx context.Context, target string) (string, error) {
	rc := c.rc()

	ref, err := ref.New(target)
	if err != nil {
		return "", err
	}

	m, err := rc.ManifestHead(ctx, ref)
	if err != nil {
		return "", fmt.Errorf("fetching %q: %w", ref, err)
	}

	return string(m.GetDescriptor().Digest), nil
}

func (c *cli) ManifestDelete(ctx context.Context, target string) error {
	rc := c.rc()

	ref, err := ref.New(target)
	if err != nil {
		return err
	}

	err = rc.ManifestDelete(ctx, ref)
	if errors.Is(err, errs.ErrHTTPStatus) {
		provider.GetLogger(ctx).Warning("this registry does not support deletions")
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}

// Inspect inspects an image.
func (c *cli) Inspect(ctx context.Context, r string) ([]descriptor.Descriptor, error) {
	ref, err := ref.New(r)
	if err != nil {
		return nil, err
	}
	rc := c.rc()

	m, err := rc.ManifestGet(ctx, ref)
	if err != nil {
		return nil, err
	}

	if mi, ok := m.(manifest.Indexer); ok {
		return mi.GetManifestList()
	}

	return []descriptor.Descriptor{m.GetDescriptor()}, nil
}

// Delete attempts to delete an image with the given ref. Many registries don't
// support the DELETE API yet, so this operation is not guaranteed to work.
func (c *cli) Delete(ctx context.Context, r string) error {
	// Attempt to delete the ref locally if it exists.
	_, _ = c.Client().ImageRemove(ctx, r, image.RemoveOptions{
		Force: true, // Needed in case the image has multiple tags.
	})

	// Attempt to delete the ref remotely if it was pushed -- requires a
	// digest.
	ref, err := ref.New(r)
	if err != nil || ref.Digest == "" {
		return nil
	}

	rc := c.rc()

	// TODO: Multi-platform manifests are left dangling on ECR.

	_ = rc.ManifestDelete(ctx, ref)

	return nil
}

func normalizeReference(ref string) (reference.Named, error) {
	namedRef, err := reference.ParseNormalizedNamed(ref)
	if err != nil {
		return nil, err
	}
	if _, isDigested := namedRef.(reference.Canonical); !isDigested {
		return reference.TagNameOnly(namedRef), nil
	}
	return namedRef, nil
}
