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

package internal

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"

	// For examples/docs.
	_ "embed"
	// These imports are needed to register the drivers with buildkit.
	_ "github.com/docker/buildx/driver/docker-container"
	_ "github.com/docker/buildx/driver/kubernetes"
	_ "github.com/docker/buildx/driver/remote"

	"github.com/distribution/reference"
	controllerapi "github.com/docker/buildx/controller/pb"
	"github.com/docker/docker/errdefs"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/secrets/secretsprovider"
	"github.com/regclient/regclient/types/ref"

	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	_ infer.Annotated                             = (*Image)(nil)
	_ infer.Annotated                             = (*ImageArgs)(nil)
	_ infer.Annotated                             = (*ImageState)(nil)
	_ infer.CustomCheck[ImageArgs]                = (*Image)(nil)
	_ infer.CustomDelete[ImageState]              = (*Image)(nil)
	_ infer.CustomDiff[ImageArgs, ImageState]     = (*Image)(nil)
	_ infer.CustomRead[ImageArgs, ImageState]     = (*Image)(nil)
	_ infer.CustomResource[ImageArgs, ImageState] = (*Image)(nil)
	_ infer.CustomUpdate[ImageArgs, ImageState]   = (*Image)(nil)
)

//go:embed embed/image-examples.md
var _imageExamples string

//go:embed embed/image-migration.md
var _migration string

// Image is a Docker image build using buildkit.
type Image struct{}

// Annotate provides a description of the Image resource.
func (i *Image) Annotate(a infer.Annotator) {
	a.Describe(&i, dedent(`
		A Docker image built using buildx -- Docker's interface to the improved
		BuildKit backend.

		## Stability

		**This resource is pre-1.0 and in public preview.**

		We will strive to keep APIs and behavior as stable as possible, but we
		cannot guarantee stability until version 1.0.
	`)+
		"\n\n"+_migration+
		"\n\n"+_imageExamples,
	)
}

// ImageArgs instantiates a new Image.
type ImageArgs struct {
	AddHosts       []string          `pulumi:"addHosts,optional"`
	BuildArgs      map[string]string `pulumi:"buildArgs,optional"`
	BuildOnPreview *bool             `pulumi:"buildOnPreview,optional"`
	Builder        *BuilderConfig    `pulumi:"builder,optional"`
	CacheFrom      []CacheFrom       `pulumi:"cacheFrom,optional"`
	CacheTo        []CacheTo         `pulumi:"cacheTo,optional"`
	Context        *BuildContext     `pulumi:"context,optional"`
	Dockerfile     *Dockerfile       `pulumi:"dockerfile,optional"`
	Exports        []Export          `pulumi:"exports,optional"`
	Labels         map[string]string `pulumi:"labels,optional"`
	Load           bool              `pulumi:"load,optional"`
	Network        *NetworkMode      `pulumi:"network,optional"`
	NoCache        bool              `pulumi:"noCache,optional"`
	Platforms      []Platform        `pulumi:"platforms,optional"`
	Pull           bool              `pulumi:"pull,optional"`
	Push           bool              `pulumi:"push,optional"`
	Registries     []Registry        `pulumi:"registries,optional"`
	Secrets        map[string]string `pulumi:"secrets,optional"`
	SSH            []SSH             `pulumi:"ssh,optional"`
	Tags           []string          `pulumi:"tags,optional"`
	Target         string            `pulumi:"target,optional"`
	Exec           bool              `pulumi:"exec,optional"`
}

// Annotate describes inputs to the Image resource.
func (ia *ImageArgs) Annotate(a infer.Annotator) {
	a.Describe(&ia.AddHosts, dedent(`
		Custom "host:ip" mappings to use during the build.

		Equivalent to Docker's "--add-host" flag.
	`))
	a.Describe(&ia.BuildArgs, dedent(`
		"ARG" names and values to set during the build.

		These variables are accessed like environment variables inside "RUN"
		instructions.

		Build arguments are persisted in the image, so you should use "secrets"
		if these arguments are sensitive.

		Equivalent to Docker's "--build-arg" flag.
	`))
	a.Describe(&ia.BuildOnPreview, dedent(`
		Setting this to "false" will always skip image builds during previews,
		and setting it to "true" will always build images during previews.

		Images built during previews are never exported to registries, however
		cache manifests are still exported.

		On-disk Dockerfiles are always validated for syntactic correctness
		regardless of this setting.

		Defaults to "true" as a safeguard against broken images merging as part
		of CI pipelines.
	`))
	a.SetDefault(&ia.BuildOnPreview, pulumi.Bool(true))
	a.Describe(&ia.Builder, dedent(`
		Builder configuration.
	`))
	a.Describe(&ia.CacheFrom, dedent(`
		Cache export configuration.

		Equivalent to Docker's "--cache-from" flag.
	`))
	a.Describe(&ia.CacheTo, dedent(`
		Cache import configuration.

		Equivalent to Docker's "--cache-to" flag.
	`))
	a.Describe(&ia.Context, dedent(`
		Build context settings.

		Equivalent to Docker's "PATH | URL | -" positional argument.
	`))
	a.Describe(&ia.Dockerfile, dedent(`
		Dockerfile settings.

		Equivalent to Docker's "--file" flag.
	`))
	a.Describe(&ia.Exports, dedent(`
		Controls where images are persisted after building.

		Images are only stored in the local cache unless "exports" are
		explicitly configured.

		Exporting to multiple destinations requires a daemon running BuildKit
		0.13 or later.

		Equivalent to Docker's "--output" flag.
	`))
	a.Describe(&ia.Labels, dedent(`
		Attach arbitrary key/value metadata to the image.

		Equivalent to Docker's "--label" flag.
	`))
	a.Describe(&ia.Load, dedent(`
		When "true" the build will automatically include a "docker" export.

		Defaults to "false".

		Equivalent to Docker's "--load" flag.
	`))
	a.Describe(&ia.Network, dedent(`
		Set the network mode for "RUN" instructions. Defaults to "default".

		For custom networks, configure your builder with "--driver-opt network=...".

		Equivalent to Docker's "--network" flag.
	`))
	a.Describe(&ia.NoCache, dedent(`
		Do not import cache manifests when building the image.

		Equivalent to Docker's "--no-cache" flag.
	`))
	a.Describe(&ia.Platforms, dedent(`
		Set target platform(s) for the build. Defaults to the host's platform.

		Equivalent to Docker's "--platform" flag.
	`))
	a.Describe(&ia.Pull, dedent(`
		Always pull referenced images.

		Equivalent to Docker's "--pull" flag.
	`))
	a.Describe(&ia.Push, dedent(`
		When "true" the build will automatically include a "registry" export.

		Defaults to "false".

		Equivalent to Docker's "--push" flag.
	`))
	a.Describe(&ia.Secrets, dedent(`
		A mapping of secret names to their corresponding values.

		Unlike the Docker CLI, these can be passed by value and do not need to
		exist on-disk or in environment variables.

		Build arguments and environment variables are persistent in the final
		image, so you should use this for sensitive values.

		Similar to Docker's "--secret" flag.
	`))
	a.Describe(&ia.SSH, dedent(`
		SSH agent socket or keys to expose to the build.

		Equivalent to Docker's "--ssh" flag.
	`))
	a.Describe(&ia.Tags, dedent(`
		Name and optionally a tag (format: "name:tag").

		If exporting to a registry, the name should include the fully qualified
		registry address (e.g. "docker.io/pulumi/pulumi:latest").

		Equivalent to Docker's "--tag" flag.
	`))
	a.Describe(&ia.Target, dedent(`
		Set the target build stage(s) to build.

		If not specified all targets will be built by default.

		Equivalent to Docker's "--target" flag.
	`))
	a.Describe(&ia.Registries, dedent(`
		Registry credentials. Required if reading or exporting to private
		repositories.

		Credentials are kept in-memory and do not pollute pre-existing
		credentials on the host.

		Similar to "docker login".
	`))

	a.Describe(&ia.Exec, dedent(`
		Use "exec" mode to build this image.

		By default the provider embeds a v25 Docker client with v0.12 buildx
		support. This helps ensure consistent behavior across environments and
		is compatible with alternative build backends (e.g. "buildkitd"), but
		it may not be desirable if you require a specific version of buildx.
		For example you may want to run a custom "docker-buildx" binary with
		support for [Docker Build
		Cloud](https://docs.docker.com/build/cloud/setup/) (DBC).

		When this is set to "true" the provider will instead execute the
		"docker-buildx" binary directly to perform its operations. The user is
		responsible for ensuring this binary exists, with correct permissions
		and pre-configured builders, at a path Docker expects (e.g.
		"~/.docker/cli-plugins").

		Debugging "exec" mode may be more difficult as Pulumi will not be able
		to surface fine-grained errors and warnings. Additionally credentials
		are temporarily written to disk in order to provide them to the
		"docker-buildx" binary.
	`))

	d := Default
	a.SetDefault(&ia.Network, &d)
}

// ImageState is serialized to the program's state file.
type ImageState struct {
	ImageArgs

	Digest      string `pulumi:"digest"      provider:"output"`
	ContextHash string `pulumi:"contextHash" provider:"output"`
	Ref         string `pulumi:"ref"         provider:"output"`
}

// Annotate describes outputs of the Image resource.
func (is *ImageState) Annotate(a infer.Annotator) {
	is.ImageArgs.Annotate(a)

	a.Describe(&is.Digest, dedent(`
		A SHA256 digest of the image if it was exported to a registry or
		elsewhere.

		Empty if the image was not exported.

		Registry images can be referenced precisely as "<tag>@<digest>". The
		"ref" output provides one such reference as a convenience.
		`,
	))
	a.Describe(&is.ContextHash, dedent(`
		A preliminary hash of the image's build context.

		Pulumi uses this to determine if an image _may_ need to be re-built.
	`))
	a.Describe(&is.Ref, dedent(`
		If the image was pushed to any registries then this will contain a
		single fully-qualified tag including the build's digest.

		If the image had tags but was not exported, this will take on a value
		of one of those tags.

		This will be empty if the image had no exports and no tags.

		This is only for convenience and may not be appropriate for situations
		where multiple tags or registries are involved. In those cases this
		output is not guaranteed to be stable.

		For more control over tags consumed by downstream resources you should
		use the "digest" output.
	`))
}

// client produces a CLI client with scoped to this resource and layered on top
// of any host-level credentials.
func (i *Image) client(pctx provider.Context, state ImageState, args ImageArgs) (Client, error) {
	ctx := context.Context(pctx)

	cfg := infer.GetConfig[Config](pctx)

	if cli, ok := ctx.Value(_mockClientKey).(Client); ok {
		return cli, nil
	}

	// We prefer auth from args, the provider, and state in that order. We
	// build a slice in reverse order because wrap() will overwrite earlier
	// entries with later ones.
	auths := []Registry{}
	auths = append(auths, cfg.Registries...)
	auths = append(auths, args.Registries...)

	return wrap(cfg.host, auths...)
}

// Check validates ImageArgs, sets defaults, and ensures our client is
// authenticated.
func (i *Image) Check(
	_ provider.Context,
	_ string,
	_ resource.PropertyMap,
	news resource.PropertyMap,
) (ImageArgs, []provider.CheckFailure, error) {
	args, failures, err := infer.DefaultCheck[ImageArgs](news)
	if err != nil || len(failures) != 0 {
		return args, failures, err
	}

	// :(
	preview := news.ContainsUnknowns()

	if _, berr := args.validate(preview); berr != nil {
		errs := berr.(interface{ Unwrap() []error }).Unwrap()
		for _, e := range errs {
			if cf, ok := e.(checkFailure); ok {
				failures = append(failures, cf.CheckFailure)
			}
		}
	}

	return args, failures, err
}

type checkFailure struct {
	provider.CheckFailure
}

func (cf checkFailure) Error() string {
	return cf.Reason
}

func newCheckFailure(err error, format string, args ...any) error {
	return checkFailure{
		provider.CheckFailure{Property: fmt.Sprintf(format, args...), Reason: err.Error()},
	}
}

// normalize returns a copy of ImageArgs after accounting for unknown (preview)
// values and CLI push/load shorthand.
//
// During preview the go-provider sends unknown inputs as zero values. In order
// to enable build-on-preview behavior, we omit zero values during previews.
func (ia *ImageArgs) normalize(preview bool) ImageArgs {
	normalized := ImageArgs{
		AddHosts:       filter(stringKeeper{preview}, ia.AddHosts...),
		BuildArgs:      mapKeeper{preview}.keep(ia.BuildArgs),
		BuildOnPreview: ia.BuildOnPreview,
		Builder:        ia.Builder,
		CacheFrom:      filter(stringerKeeper[CacheFrom]{preview}, ia.CacheFrom...),
		CacheTo:        filter(stringerKeeper[CacheTo]{preview}, ia.CacheTo...),
		Context:        contextKeeper{preview}.keep(ia.Context),
		Dockerfile:     ia.Dockerfile,
		Exports:        filter(stringerKeeper[Export]{preview}, ia.Exports...),
		Labels:         mapKeeper{preview}.keep(ia.Labels),
		Load:           ia.Load,
		Network:        ia.Network,
		NoCache:        ia.NoCache,
		Platforms:      filter(stringerKeeper[Platform]{preview}, ia.Platforms...),
		Pull:           ia.Pull,
		Push:           ia.Push,
		Registries:     filter(registryKeeper{preview}, ia.Registries...),
		SSH:            filter(stringerKeeper[SSH]{preview}, ia.SSH...),
		Secrets:        mapKeeper{preview}.keep(ia.Secrets),
		Tags:           filter(stringKeeper{preview}, ia.Tags...),
		Target:         ia.Target,
	}

	// Handle --push/--load shorthand.
	if normalized.Push {
		normalized.Exports = append(normalized.Exports, Export{Raw: "type=registry"})
	}
	if normalized.Load {
		normalized.Exports = append(normalized.Exports, Export{Raw: "type=docker"})
	}

	return normalized
}

// buildable returns true if the ImageArgs has no unknown values and can
// therefore be built during previews.
func (ia *ImageArgs) buildable() bool {
	// We can build the given inputs if filtering unknowns is a no-op.
	filtered := ia.normalize(true)
	return reflect.DeepEqual(ia, &filtered)
}

// isExported returns true if the args include a registry export.
func (ia *ImageArgs) isExported() bool {
	if ia.Push {
		return true
	}
	for _, e := range ia.Exports {
		if e.pushed() {
			return true
		}
	}
	return false
}

// shouldBuildOnPreview returns true if we should build this image during
// previews.
func (ia *ImageArgs) shouldBuildOnPreview() bool {
	if ia.BuildOnPreview != nil {
		return *ia.BuildOnPreview
	}
	return true
}

type build struct {
	opts    controllerapi.BuildOptions
	secrets map[string]string
	inline  string
	exec    bool
}

func (b build) BuildOptions() controllerapi.BuildOptions {
	return b.opts
}

func (b build) Inline() string {
	return b.inline
}

func (b build) Secrets() session.Attachable {
	m := map[string][]byte{}
	for k, v := range b.secrets {
		m[k] = []byte(v)
	}
	return secretsprovider.FromMap(m)
}

func (b build) ShouldExec() bool {
	return b.exec
}

func (ia ImageArgs) toBuild(
	ctx provider.Context,
	preview bool,
) (Build, error) {
	opts, err := ia.validate(preview)
	if err != nil {
		return nil, err
	}

	if len(ia.Exports) == 0 && !ia.Push && !ia.Load {
		ctx.Log(diag.Warning,
			"No exports were specified so the build will only remain in the local build cache. "+
				"Use `push` to upload the image to a registry, or silence this warning with a `cacheonly` export.",
		)
	}

	if len(opts.Platforms) > 1 && len(opts.CacheTo) > 0 {
		ctx.Log(
			diag.Warning,
			"Caching doesn't work reliably with multi-platform builds (https://github.com/docker/buildx/discussions/1382). "+
				"Instead, perform one cached build per platform and create an Index to join them all together.",
		)
	}

	return build{
		opts:    opts,
		inline:  ia.Dockerfile.Inline,
		secrets: ia.Secrets,
		exec:    ia.Exec,
	}, nil
}

// validate confirms the ImageArgs are valid and returns BuildOptions
// appropriate for passing to builders.
func (ia *ImageArgs) validate(preview bool) (controllerapi.BuildOptions, error) {
	var multierr error

	if len(ia.Exports) > 1 {
		multierr = errors.Join(multierr,
			newCheckFailure(errors.New("multiple exports are currently unsupported"), "exports"),
		)
	}
	if ia.Push && ia.Load {
		multierr = errors.Join(
			multierr,
			newCheckFailure(
				errors.New("push and load may not be set together at the moment"),
				"push",
			),
		)
	}
	if len(ia.Exports) > 0 && (ia.Push || ia.Load) {
		multierr = errors.Join(multierr,
			newCheckFailure(errors.New("exports can't be provided with push or load"), "exports"),
		)
	}

	dockerfile, context, err := ia.Context.validate(preview, ia.Dockerfile)
	if err != nil {
		multierr = errors.Join(multierr, err)
	}
	ia.Dockerfile = dockerfile

	if err := ia.Dockerfile.validate(preview, context); err != nil {
		multierr = errors.Join(multierr, err)
	}

	// Discard any unknown inputs if this is a preview -- we don't want them to
	// cause validation errors.
	normalized := ia.normalize(preview)

	exports := []*controllerapi.ExportEntry{}
	for idx, e := range normalized.Exports {
		if e.Disabled {
			continue
		}
		exp, err := e.validate(preview, ia.Tags)
		if err != nil {
			multierr = errors.Join(multierr, newCheckFailure(err, "exports[%d]", idx))
			continue
		}
		exports = append(exports, exp)
	}

	platforms := []string{}
	for idx, p := range normalized.Platforms {
		platform, err := p.validate(preview)
		if err != nil {
			multierr = errors.Join(multierr, newCheckFailure(err, "platforms[%d]", idx))
			continue
		}
		platforms = append(platforms, platform)
	}

	cacheFrom := []*controllerapi.CacheOptionsEntry{}
	for idx, c := range normalized.CacheFrom {
		if c.String() == "" {
			continue // Disabled or unknown/preview.
		}
		cache, err := c.validate(preview)
		if err != nil {
			multierr = errors.Join(multierr, newCheckFailure(err, "cacheFrom[%d]", idx))
			continue
		}
		cacheFrom = append(cacheFrom, cache)
	}

	cacheTo := []*controllerapi.CacheOptionsEntry{}
	for idx, c := range normalized.CacheTo {
		if c.String() == "" {
			continue // Disabled or unknown/preview.
		}
		cache, err := c.validate(preview)
		if err != nil {
			multierr = errors.Join(multierr, newCheckFailure(err, "cacheTo[%d]", idx))
			continue
		}
		cacheTo = append(cacheTo, cache)
	}

	ssh := []*controllerapi.SSH{}
	for idx, s := range normalized.SSH {
		ss, err := s.validate()
		if err != nil {
			multierr = errors.Join(multierr, newCheckFailure(err, "ssh[%d]", idx))
			continue
		}
		ssh = append(ssh, ss)
	}

	for idx, t := range normalized.Tags {
		if _, err := reference.Parse(t); err != nil {
			multierr = errors.Join(multierr, newCheckFailure(err, "tags[%d]", idx))
		}
	}

	secrets := []*controllerapi.Secret{}
	for k, v := range normalized.Secrets {
		// We abuse the pb.Secret proto by stuffing the secret's value in
		// XXX_unrecognized. We never serialize this proto so this is tolerable.
		secrets = append(secrets, &controllerapi.Secret{
			ID:               k,
			XXX_unrecognized: []byte(v),
		})
	}

	builder := BuilderConfig{}
	if normalized.Builder != nil {
		builder = *normalized.Builder
	}

	opts := controllerapi.BuildOptions{
		BuildArgs:      normalized.BuildArgs,
		Builder:        builder.Name,
		CacheFrom:      cacheFrom,
		CacheTo:        cacheTo,
		ContextPath:    context.Location,
		DockerfileName: dockerfile.Location,
		Exports:        exports,
		ExtraHosts:     normalized.AddHosts,
		Labels:         normalized.Labels,
		NetworkMode:    normalized.Network.String(),
		NoCache:        normalized.NoCache,
		NamedContexts:  normalized.Context.namedMap(),
		Platforms:      platforms,
		Pull:           normalized.Pull,
		Secrets:        secrets,
		SSH:            ssh,
		Tags:           normalized.Tags,
		Target:         normalized.Target,
	}

	return opts, multierr
}

// Create builds an image using buildkit.
func (i *Image) Create(
	ctx provider.Context,
	name string,
	input ImageArgs,
	preview bool,
) (string, ImageState, error) {
	state := ImageState{ImageArgs: input}
	id := name

	// Default our ref to one of our tags.
	for _, tag := range state.Tags {
		if _, err := normalizeReference(tag); err != nil {
			continue
		}
		state.Ref = tag
		break
	}

	cli, err := i.client(ctx, state, input)
	if err != nil {
		return id, state, err
	}

	ok, err := cli.BuildKitEnabled()
	if err != nil {
		return id, state, fmt.Errorf("checking buildkit compatibility: %w", err)
	}
	if !ok {
		return id, state, errors.New("buildkit is not supported on this host")
	}

	build, err := input.toBuild(ctx, preview)
	if err != nil {
		return id, state, fmt.Errorf("preparing: %w", err)
	}

	hash, err := hashBuildContext(
		input.Context.Location,
		input.Dockerfile.Location,
		input.Context.Named.Map(),
	)
	if err != nil {
		return id, state, fmt.Errorf("hashing build context: %w", err)
	}
	state.ContextHash = hash

	if preview && !input.shouldBuildOnPreview() {
		return id, state, nil
	}
	if preview && !input.buildable() {
		ctx.Log(diag.Warning, "Skipping preview build because some inputs are unknown.")
		return id, state, nil
	}

	result, err := cli.Build(ctx, build)
	if err != nil {
		return id, state, err
	}

	if d, ok := result.ExporterResponse[exptypes.ExporterImageDigestKey]; ok {
		state.Digest = d
		id = d
	}

	if state.Digest == "" {
		// Can't construct a ref, nothing else to do.
		return id, state, nil
	}

	// Take the first registry tag we find and add a digest to it. That becomes
	// our simplified "ref" output.
	for _, tag := range state.Tags {
		ref, ok := addDigest(tag, state.Digest)
		if !ok {
			continue
		}

		state.Ref = ref
		break
	}

	return id, state, nil
}

// Update builds a new image. Normally we create-replace resources, but for
// images built locally there is nothing to delete. We treat those cases as
// updates and simply re-build the image without deleting anything.
func (i *Image) Update(
	ctx provider.Context,
	name string,
	_ ImageState,
	input ImageArgs,
	preview bool,
) (ImageState, error) {
	_, state, err := i.Create(ctx, name, input, preview)
	return state, err
}

// Read attempts to read manifests from an image's exports. An image without
// exports will have no manifests.
func (i *Image) Read(
	ctx provider.Context,
	name string,
	input ImageArgs,
	state ImageState,
) (
	string, // id
	ImageArgs, // normalized inputs
	ImageState, // normalized state
	error,
) {
	cli, err := i.client(ctx, state, input)
	if err != nil {
		return name, input, state, err
	}

	if !state.isExported() {
		// Nothing was pushed -- all done.
		return name, input, state, nil
	}

	tagsToKeep := []string{}

	// Do a lookup on all of the tags at the digests we expect to see.
	for _, tag := range state.Tags {
		ref, ok := addDigest(tag, state.Digest)
		if !ok {
			// Not a pushed tag.
			tagsToKeep = append(tagsToKeep, tag)
			break
		}

		// Does a tag with this digest exist?
		descriptors, err := cli.Inspect(ctx, ref)
		if err != nil {
			ctx.Log(diag.Warning, err.Error())
			continue
		}

		for _, d := range descriptors {
			if d.Platform != nil && d.Platform.Architecture == "unknown" {
				// Ignore cache manifests.
				continue
			}

			tagsToKeep = append(tagsToKeep, tag)
			break
		}
	}

	// If we couldn't find the tags we expected then return an empty ID to
	// delete the resource.
	if len(input.Tags) > 0 && len(tagsToKeep) == 0 {
		return "", input, state, nil
	}

	state.Tags = tagsToKeep

	return name, input, state, nil
}

// Delete deletes an Image. If the Image was already deleted out-of-band it is
// treated as a success.
func (i *Image) Delete(
	ctx provider.Context,
	_ string,
	state ImageState,
) error {
	cli, err := i.client(ctx, state, state.ImageArgs)
	if err != nil {
		return err
	}

	if state.Digest == "" {
		// Nothing was exported. Just try to delete the local image.
		return cli.Delete(ctx, state.Ref)
	}

	digests := []string{}

	// Construct a ref with digest for each repository we pushed to.
	for _, tag := range state.Tags {
		ref, err := ref.New(tag)
		if err != nil {
			continue
		}
		digested := ref.SetDigest(state.Digest)
		digests = append(digests, digested.CommonName())
	}

	slices.Sort(digests)
	digests = slices.Compact(digests)

	var multierr error
	for _, digested := range digests {
		err = cli.Delete(context.Context(ctx), digested)
		if errdefs.IsNotFound(err) {
			ctx.Log(diag.Warning, digested+" not found")
			continue // Nothing to do.
		}
		multierr = errors.Join(multierr, err)
	}

	return multierr
}

// Diff re-implements most of the default diff behavior, with the exception of
// ignoring "password" changes on registry inputs.
func (*Image) Diff(
	_ provider.Context,
	_ string,
	olds ImageState,
	news ImageArgs,
) (provider.DiffResponse, error) {
	diff := map[string]provider.PropertyDiff{}
	update := provider.PropertyDiff{Kind: provider.Update}

	if !reflect.DeepEqual(olds.AddHosts, news.AddHosts) {
		diff["addHosts"] = update
	}
	if !reflect.DeepEqual(olds.BuildArgs, news.BuildArgs) {
		diff["buildArgs"] = update
	}
	if !reflect.DeepEqual(olds.BuildOnPreview, news.BuildOnPreview) {
		diff["buildOnPreview"] = update
	}
	if !reflect.DeepEqual(olds.Builder, news.Builder) {
		diff["builder"] = update
	}
	if !reflect.DeepEqual(olds.CacheFrom, news.CacheFrom) {
		diff["cacheFrom"] = update
	}
	if !reflect.DeepEqual(olds.CacheTo, news.CacheTo) {
		diff["cacheTo"] = update
	}
	if olds.Context.Location != news.Context.Location {
		diff["context.location"] = update
	}
	if !reflect.DeepEqual(olds.Context.Named, news.Context.Named) {
		diff["context.named"] = update
	}
	if olds.Dockerfile.Location != news.Dockerfile.Location {
		diff["dockerfile.location"] = update
	}
	if olds.Dockerfile.Inline != news.Dockerfile.Inline {
		diff["dockerfile.inline"] = update
	}
	// Use string comparison to ignore any manifests attached to the export.
	if fmt.Sprint(olds.Exports) != fmt.Sprint(news.Exports) {
		diff["exports"] = update
	}
	if !reflect.DeepEqual(olds.Labels, news.Labels) {
		diff["labels"] = update
	}
	if olds.Load != news.Load {
		diff["load"] = update
	}
	if !reflect.DeepEqual(olds.Network, news.Network) {
		diff["network"] = update
	}
	if !reflect.DeepEqual(olds.NoCache, news.NoCache) {
		diff["noCache"] = update
	}
	if !reflect.DeepEqual(olds.Platforms, news.Platforms) {
		diff["platforms"] = update
	}
	if olds.Pull != news.Pull {
		diff["pull"] = update
	}
	if olds.Push != news.Push {
		diff["push"] = update
	}
	if !reflect.DeepEqual(olds.Secrets, news.Secrets) {
		diff["secrets"] = update
	}
	if !reflect.DeepEqual(olds.SSH, news.SSH) {
		diff["ssh"] = update
	}
	if !reflect.DeepEqual(olds.Tags, news.Tags) {
		diff["tags"] = update
	}
	if !reflect.DeepEqual(olds.Target, news.Target) {
		diff["target"] = update
	}

	// pull=true indicates that we want to keep base layers up-to-date. In this
	// case we'll always perform the build.
	if news.Pull && (len(news.Exports) > 0 || news.Push || news.Load) {
		diff["contextHash"] = update
	}

	// Check if anything has changed in our build context.
	hash, err := hashBuildContext(
		news.Context.Location,
		news.Dockerfile.Location,
		news.Context.Named.Map(),
	)
	if err != nil {
		return provider.DiffResponse{}, err
	}
	if hash != olds.ContextHash {
		diff["contextHash"] = update
	}

	// Registries need special handling because we ignore "password" changes to not introduce unnecessary changes.
	if len(olds.Registries) != len(news.Registries) {
		diff["registries"] = update
	} else {
		for idx, oldr := range olds.Registries {
			newr := news.Registries[idx]
			if (oldr.Username == newr.Username) && (oldr.Address == newr.Address) {
				continue
			}
			diff[fmt.Sprintf("registries[%d]", idx)] = update
			break
		}
	}

	return provider.DiffResponse{
		HasChanges:   len(diff) > 0,
		DetailedDiff: diff,
	}, nil
}

// addDigest constructs a tagged ref with an "@<digest>" suffix.
//
// Returns false if the given ref was not fully qualified.
func addDigest(ref, digest string) (string, bool) {
	named, err := reference.ParseNamed(ref)
	if err != nil {
		return "", false
	}
	tag := "latest"
	if tagged, ok := named.(reference.Tagged); ok {
		tag = tagged.Tag()
	}

	full, err := reference.Parse(
		fmt.Sprintf("%s:%s@%s", named.Name(), tag, digest),
	)
	if err != nil {
		return "", false
	}

	return full.String(), true
}
