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
	"strings"

	// For examples/docs.
	_ "embed"

	"github.com/regclient/regclient/types/errs"

	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

var (
	_ infer.Annotated                             = (*Index)(nil)
	_ infer.Annotated                             = (*IndexArgs)(nil)
	_ infer.Annotated                             = (*IndexState)(nil)
	_ infer.CustomCheck[IndexArgs]                = (*Index)(nil)
	_ infer.CustomResource[IndexArgs, IndexState] = (*Index)(nil)
	_ infer.CustomDelete[IndexState]              = (*Index)(nil)
	_ infer.CustomDiff[IndexArgs, IndexState]     = (*Index)(nil)
	_ infer.CustomRead[IndexArgs, IndexState]     = (*Index)(nil)
	_ infer.CustomUpdate[IndexArgs, IndexState]   = (*Index)(nil)
)

//go:embed embed/index-examples.md
var _indexExamples string

// Index is an OCI index or manifest list on a remote registry.
type Index struct{}

// IndexArgs instantiate an Index.
type IndexArgs struct {
	Tag      string    `pulumi:"tag"`
	Sources  []string  `pulumi:"sources"`
	Push     *bool     `pulumi:"push,optional"`
	Registry *Registry `pulumi:"registry,optional"`
}

func (i IndexArgs) isPushed() bool {
	if i.Push == nil {
		return true // default
	}
	return *i.Push
}

// IndexState captures the state of an Index.
type IndexState struct {
	IndexArgs

	Ref string `pulumi:"ref" provider:"output"`
}

// Annotate sets docstrings and defaults on Index.
func (i *Index) Annotate(a infer.Annotator) {
	a.Describe(&i, dedent(`
		A wrapper around "docker buildx imagetools create" to create an index
		(or manifest list) referencing one or more existing images.

		In most cases you do not need an "Index" to build a multi-platform
		image -- specifying multiple platforms on the "Image" will handle this
		for you automatically.

		However, as of April 2024, building multi-platform images _with
		caching_ will only export a cache for one platform at a time (see [this
		discussion](https://github.com/docker/buildx/discussions/1382) for more
		details).

		Therefore this resource can be helpful if you are building
		multi-platform images with caching: each platform can be built and
		cached separately, and an "Index" can join them all together. An
		example of this is shown below.

		This resource creates an OCI image index or a Docker manifest list
		depending on the media types of the source images.
	`)+
		"\n\n"+_indexExamples,
	)
}

// Annotate sets docstrings and defaults on IndexArgs.
func (i *IndexArgs) Annotate(a infer.Annotator) {
	a.Describe(&i.Registry, dedent(`
		Authentication for the registry where the tagged index will be pushed.

		Credentials can also be included with the provider's configuration.
	`))
	a.Describe(&i.Sources, dedent(`
		Existing images to include in the index.
	`))
	a.Describe(&i.Tag, dedent(`
		The tag to apply to the index.
	`))
	a.Describe(&i.Push, dedent(`
		If true, push the index to the target registry.

		Defaults to "true".
	`))

	a.SetDefault(&i.Push, true)
}

// Annotate sets docstrings on IndexState.
func (i *IndexState) Annotate(a infer.Annotator) {
	a.Describe(&i.Ref, dedent(`
		The pushed tag with digest.

		Identical to the tag if the index was not pushed.
	`))
}

// Create is a passthrough to Update.
func (i *Index) Create(
	ctx context.Context,
	req infer.CreateRequest[IndexArgs],
) (infer.CreateResponse[IndexState], error) {
	resp, err := i.Update(ctx,
		infer.UpdateRequest[IndexArgs, IndexState]{
			ID:      req.Name,
			Olds:    IndexState{},
			News:    req.Inputs,
			Preview: req.Preview,
		},
	)
	return infer.CreateResponse[IndexState]{ID: req.Name, Output: resp.Output}, err
}

// Update performs `buildx imagetools create` to create a new OCI index /
// manifest list.
func (i *Index) Update(
	ctx context.Context,
	req infer.UpdateRequest[IndexArgs, IndexState],
) (infer.UpdateResponse[IndexState], error) {
	state, input := req.Olds, req.News

	state.IndexArgs = input
	state.Ref = input.Tag

	cli, err := i.client(ctx, state, input)
	if err != nil {
		return infer.UpdateResponse[IndexState]{Output: state}, err
	}

	if req.Preview {
		return infer.UpdateResponse[IndexState]{Output: state}, nil
	}

	provider.GetLogger(ctx).
		Debugf("creating index with tag %s and sources %s", input.Tag, input.Sources)

	err = cli.ManifestCreate(ctx, input.isPushed(), input.Tag, input.Sources...)
	if err != nil {
		return infer.UpdateResponse[IndexState]{Output: state}, fmt.Errorf("creating: %w", err)
	}

	// Read remote manifest information, if it exists.
	live, err := i.Read(ctx,
		infer.ReadRequest[IndexArgs, IndexState]{ID: req.ID, Inputs: input, State: state},
	)
	if err != nil {
		return infer.UpdateResponse[IndexState]{Output: state}, fmt.Errorf("reading: %w", err)
	}
	return infer.UpdateResponse[IndexState]{Output: live.State}, nil
}

func (i *Index) Read(
	ctx context.Context,
	req infer.ReadRequest[IndexArgs, IndexState],
) (infer.ReadResponse[IndexArgs, IndexState], error) {
	state, input := req.State, req.Inputs

	state.IndexArgs = input
	state.Ref = input.Tag

	if !input.isPushed() {
		provider.GetLogger(ctx).Debug("skipping read because index was not pushed")
		return infer.ReadResponse[IndexArgs, IndexState]{
			ID:     req.ID,
			Inputs: input,
			State:  state,
		}, nil // Nothing to read.
	}

	cli, err := i.client(ctx, state, input)
	if err != nil {
		return infer.ReadResponse[IndexArgs, IndexState]{
			ID:     req.ID,
			Inputs: input,
			State:  state,
		}, err
	}

	provider.GetLogger(ctx).Debug("reading index with tag " + input.Tag)

	digest, err := cli.ManifestInspect(ctx, input.Tag)
	if errors.Is(err, errs.ErrNotFound) {
		// A remote tag was expected but isn't there -- delete the resource.
		return infer.ReadResponse[IndexArgs, IndexState]{ID: "", Inputs: input, State: state}, nil
	}
	if errors.Is(err, errs.ErrHTTPUnauthorized) {
		provider.GetLogger(ctx).Warning("invalid credentials, skipping")
		return infer.ReadResponse[IndexArgs, IndexState]{
			ID:     req.ID,
			Inputs: input,
			State:  state,
		}, nil
	}
	if err != nil {
		return infer.ReadResponse[IndexArgs, IndexState]{
			ID:     req.ID,
			Inputs: input,
			State:  state,
		}, err
	}

	if ref, ok := addDigest(input.Tag, digest); ok {
		state.Ref = ref
	}

	return infer.ReadResponse[IndexArgs, IndexState]{ID: req.ID, Inputs: input, State: state}, nil
}

// Check confirms the Index's tag and source refs are all valid. This doesn't
// fully capture input requirements -- for example buildx requires refs to all
// exist on the same registry. This is sufficient to handle the most common
// cases for now.
func (i *Index) Check(
	ctx context.Context,
	req infer.CheckRequest,
) (infer.CheckResponse[IndexArgs], error) {
	args, failures, err := infer.DefaultCheck[IndexArgs](ctx, req.NewInputs)
	if err != nil {
		return infer.CheckResponse[IndexArgs]{Failures: failures, Inputs: args}, err
	}

	if _, err := normalizeReference(args.Tag); args.Tag != "" && err != nil {
		failures = append(
			failures,
			provider.CheckFailure{
				Property: "target",
				Reason:   err.Error(),
			},
		)
	}

	for idx, s := range args.Sources {
		if _, err := normalizeReference(s); s != "" && err != nil {
			failures = append(
				failures,
				provider.CheckFailure{
					Property: fmt.Sprintf("refs[%d]", idx),
					Reason:   err.Error(),
				},
			)
		}
	}

	return infer.CheckResponse[IndexArgs]{Failures: failures, Inputs: args}, nil
}

// Delete attempts to delete the remote manifest.
func (i *Index) Delete(
	ctx context.Context,
	req infer.DeleteRequest[IndexState],
) (infer.DeleteResponse, error) {
	state := req.State
	if !state.isPushed() {
		return infer.DeleteResponse{}, nil // Nothing to delete.
	}

	cli, err := i.client(ctx, state, state.IndexArgs)
	if err != nil {
		return infer.DeleteResponse{}, err
	}

	err = cli.ManifestDelete(ctx, state.Ref)
	// TODO: Upstream buildx swallows the error types we'd like to test for
	// here.
	if err != nil && strings.Contains(err.Error(), "No such manifest:") {
		return infer.DeleteResponse{}, nil
	}
	return infer.DeleteResponse{}, err
}

// Diff returns a diff of proposed changes against current state. Ideally we
// wouldn't need to implement all of this, but we currently have to in order to
// force `ignoreChanges`-style behavior on our registry password (which can
// change all the time due to short-lived AWS credentials).
func (i *Index) Diff(
	_ context.Context,
	req infer.DiffRequest[IndexArgs, IndexState],
) (provider.DiffResponse, error) {
	olds, news := req.Olds, req.News

	diff := map[string]provider.PropertyDiff{}
	update := provider.PropertyDiff{Kind: provider.Update}
	replace := provider.PropertyDiff{Kind: provider.UpdateReplace}

	if olds.Tag != news.Tag {
		diff["tag"] = replace
	}
	if !reflect.DeepEqual(olds.Sources, news.Sources) {
		diff["sources"] = update
	}
	if olds.Registry != nil && news.Registry != nil {
		if olds.Registry.Address != news.Registry.Address {
			diff["registry.address"] = update
			if olds.Registry.Address != "" {
				diff["registry.address"] = replace
			}
		}
		if olds.Registry.Username != news.Registry.Username {
			diff["registry.username"] = update
		}
	}
	if (olds.Registry == nil && news.Registry != nil) ||
		(olds.Registry != nil && news.Registry == nil) {
		diff["registry"] = update
	}
	// Intentionally ignore changes to registry.password

	return provider.DiffResponse{
		HasChanges:   len(diff) > 0,
		DetailedDiff: diff,
	}, nil
}

// client produces a CLI client scoped to this resource and layered on top of
// any host-level credentials.
func (i *Index) client(
	ctx context.Context,
	_ IndexState,
	args IndexArgs,
) (Client, error) {
	cfg := infer.GetConfig[Config](ctx)

	if cli, ok := ctx.Value(_mockClientKey).(Client); ok {
		return cli, nil
	}

	// We prefer auth from args, the provider, and state in that order. We
	// build a slice in reverse order because wrap() will overwrite earlier
	// entries with later ones.
	auths := []Registry{}
	auths = append(auths, cfg.Registries...)
	if args.Registry != nil {
		auths = append(auths, *args.Registry)
	}

	return wrap(cfg.host, auths...)
}
