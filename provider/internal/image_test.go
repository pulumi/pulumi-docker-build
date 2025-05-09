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
	"os"
	"path/filepath"
	"testing"

	_ "github.com/docker/buildx/driver/docker-container"

	"github.com/distribution/reference"
	pb "github.com/docker/buildx/controller/pb"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/regclient/regclient/types/descriptor"
	"github.com/regclient/regclient/types/platform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/mapper"
	"github.com/pulumi/pulumi/sdk/v3/go/property"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var _fakeURN = resource.NewURN("test", "provider", "a", "docker-build:index:Image", "test")

func TestImageLifecycle(t *testing.T) {
	t.Parallel()
	noClient := func(t *testing.T) Client {
		ctrl := gomock.NewController(t)
		return NewMockClient(ctrl)
	}

	_, err := reference.ParseNamed("docker.io/pulumibot/buildkit-e2e")
	require.NoError(t, err)

	tests := []struct {
		name string

		op     func(t *testing.T) integration.Operation
		client func(t *testing.T) Client
	}{
		{
			name: "happy path builds",
			client: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				c := NewMockClient(ctrl)
				c.EXPECT().BuildKitEnabled().Return(true, nil).AnyTimes()
				c.EXPECT().SupportsMultipleExports().Return(true).AnyTimes()
				c.EXPECT().Build(gomock.Any(), gomock.AssignableToTypeOf(&build{})).DoAndReturn(
					func(_ context.Context, b Build) (*client.SolveResponse, error) {
						assert.Equal(t, "testdata/noop/Dockerfile", b.BuildOptions().DockerfileName)
						return &client.SolveResponse{
							ExporterResponse: map[string]string{
								exptypes.ExporterImageDigestKey: "sha256:98ea6e4f216f2fb4b69fff9b3a44842c38686ca685f3f55dc48c5d3fb1107be4",
							},
						}, nil
					},
				).AnyTimes()
				c.EXPECT().Delete(gomock.Any(),
					"docker.io/pulumibot/buildkit-e2e@sha256:98ea6e4f216f2fb4b69fff9b3a44842c38686ca685f3f55dc48c5d3fb1107be4",
				).
					Return(nil)
				return c
			},
			op: func(t *testing.T) integration.Operation {
				return integration.Operation{
					Inputs: property.NewMap(map[string]property.Value{
						"push": property.New(false),
						"tags": property.New(
							[]property.Value{
								property.New("docker.io/pulumibot/buildkit-e2e"),
								property.New("docker.io/pulumibot/buildkit-e2e:main"),
							},
						),
						"platforms": property.New(
							[]property.Value{
								property.New("linux/arm64"),
								property.New("linux/amd64"),
							},
						),
						"context": property.New(map[string]property.Value{
							"location": property.New("testdata/noop"),
						}),
						"dockerfile": property.New(map[string]property.Value{
							"location": property.New("testdata/noop/Dockerfile"),
						}),
						"exports": property.New(
							[]property.Value{
								property.New(map[string]property.Value{
									"raw": property.New("type=registry"),
								},
								),
							},
						),
						"registries": property.New(
							[]property.Value{
								property.New(map[string]property.Value{
									"address":  property.New("fakeaddress"),
									"username": property.New("fakeuser"),
									"password": property.New("password").WithSecret(true),
								}),
							},
						),
					}),
				}
			},
		},
		{
			name:   "tags are required when pushing",
			client: noClient,
			op: func(t *testing.T) integration.Operation {
				return integration.Operation{
					Inputs: property.NewMap(map[string]property.Value{
						"push": property.New(false),
						"tags": property.New([]property.Value{}),
						"context": property.New(map[string]property.Value{
							"location": property.New("testdata/noop"),
						}),
						"exports": property.New(
							[]property.Value{
								property.New(map[string]property.Value{
									"raw": property.New("type=registry"),
								}),
							},
						),
					}),
					ExpectFailure: true,
					CheckFailures: []provider.CheckFailure{
						{
							Property: "exports[0]",
							Reason:   "at least one tag or export name is needed when pushing to a registry",
						},
					},
				}
			},
		},
		{
			name:   "invalid exports",
			client: noClient,
			op: func(t *testing.T) integration.Operation {
				return integration.Operation{
					Inputs: property.NewMap(map[string]property.Value{
						"push": property.New(false),
						"tags": property.New(
							[]property.Value{property.New("invalid-exports")},
						),
						"exports": property.New(
							[]property.Value{
								property.New(map[string]property.Value{
									"raw": property.New("type="),
								}),
							},
						),
					}),
					ExpectFailure: true,
					CheckFailures: []provider.CheckFailure{{
						Property: "exports[0]",
						Reason:   "type is required for output",
					}},
				}
			},
		},
		{
			name: "requires buildkit",
			client: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				c := NewMockClient(ctrl)
				gomock.InOrder(
					c.EXPECT().BuildKitEnabled().Return(false, nil), // Preview.
				)
				return c
			},
			op: func(t *testing.T) integration.Operation {
				return integration.Operation{
					Inputs: property.NewMap(map[string]property.Value{
						"push": property.New(false),
						"tags": property.New(
							[]property.Value{property.New("foo")},
						),
						"context": property.New(map[string]property.Value{
							"location": property.New("testdata/noop"),
						}),
					}),
					ExpectFailure: true,
				}
			},
		},
		{
			name: "error reading DOCKER_BUILDKIT",
			client: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				c := NewMockClient(ctrl)
				gomock.InOrder(
					c.EXPECT().
						BuildKitEnabled().
						Return(false, errors.New("invalid DOCKER_BUILDKIT")), // Preview.
				)
				return c
			},
			op: func(t *testing.T) integration.Operation {
				return integration.Operation{
					Inputs: property.NewMap(map[string]property.Value{
						"push": property.New(false),
						"tags": property.New(
							[]property.Value{property.New("foo")},
						),
						"context": property.New(map[string]property.Value{
							"location": property.New("testdata/noop"),
						}),
					}),
					ExpectFailure: true,
				}
			},
		},
		{
			name: "file defaults to Dockerfile",
			client: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				c := NewMockClient(ctrl)
				c.EXPECT().BuildKitEnabled().Return(true, nil).AnyTimes()
				c.EXPECT().SupportsMultipleExports().Return(true).AnyTimes()
				c.EXPECT().Build(gomock.Any(), gomock.AssignableToTypeOf(&build{})).DoAndReturn(
					func(_ context.Context, b Build) (*client.SolveResponse, error) {
						assert.Equal(t, "testdata/noop/Dockerfile", b.BuildOptions().DockerfileName)
						return &client.SolveResponse{
							ExporterResponse: map[string]string{"image.name": "test:latest"},
						}, nil
					},
				).AnyTimes()
				c.EXPECT().Delete(gomock.Any(), "default-dockerfile").Return(nil)
				return c
			},
			op: func(t *testing.T) integration.Operation {
				return integration.Operation{
					Inputs: property.NewMap(map[string]property.Value{
						"push": property.New(false),
						"tags": property.New(
							[]property.Value{
								property.New("default-dockerfile"),
							},
						),
						"context": property.New(map[string]property.Value{
							"location": property.New("testdata/noop"),
						}),
					}),
					Hook: func(_, output property.Map) {
						dockerfile := output.Get("dockerfile")
						require.NotNil(t, dockerfile)
						require.True(t, dockerfile.IsMap())
						location := dockerfile.AsMap().Get("location")
						require.True(t, location.IsString())
						assert.Equal(t, "testdata/noop/Dockerfile", location.AsString())
					},
				}
			},
		},
		{
			name: "context defaults to current directory (pulumi-docker-build#78)",
			client: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				c := NewMockClient(ctrl)
				c.EXPECT().BuildKitEnabled().Return(true, nil).AnyTimes()
				c.EXPECT().SupportsMultipleExports().Return(true).AnyTimes()
				c.EXPECT().Build(gomock.Any(), gomock.AssignableToTypeOf(&build{})).DoAndReturn(
					func(_ context.Context, b Build) (*client.SolveResponse, error) {
						assert.Equal(t, "FROM alpine:latest", b.Inline())
						return &client.SolveResponse{
							ExporterResponse: map[string]string{"image.name": "alpine:latest"},
						}, nil
					},
				).AnyTimes()
				c.EXPECT().Delete(gomock.Any(), "inline-dockerfile").Return(nil)
				return c
			},
			op: func(t *testing.T) integration.Operation {
				return integration.Operation{
					Inputs: property.NewMap(map[string]property.Value{
						"push": property.New(false),
						"tags": property.New(
							[]property.Value{
								property.New("inline-dockerfile"),
							},
						),
						"buildOnPreview": property.New(true),
						"dockerfile": property.New(map[string]property.Value{
							"inline": property.New("FROM alpine:latest"),
						}),
					}),
					Hook: func(_, output property.Map) {
						context := output.Get("context")
						require.NotNil(t, context)
						require.True(t, context.IsMap())
						location := context.AsMap().Get("location")
						require.True(t, location.IsString())
						assert.Equal(t, ".", location.AsString())
					},
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			lc := integration.LifeCycleTest{
				Resource: "docker-build:index:Image",
				Create:   tt.op(t),
			}
			s := newServer(t.Context(), t, tt.client(t))

			err := s.Configure(provider.ConfigureRequest{})
			require.NoError(t, err)

			lc.Run(t, s)
		})
	}
}

type errNotFound struct{}

func (errNotFound) NotFound()     {}
func (errNotFound) Error() string { return "not found " }

func TestDelete(t *testing.T) {
	t.Parallel()
	t.Run("image was already deleted", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		client := NewMockClient(ctrl)
		client.EXPECT().
			Delete(gomock.Any(), "docker.io/pulumi/test@sha256:foo").
			Return(errNotFound{})

		i := &Image{docker: client}

		_, err := i.Delete(t.Context(), infer.DeleteRequest[ImageState]{
			ID: "foo,bar",
			State: ImageState{
				ImageArgs: ImageArgs{
					Tags: []string{"docker.io/pulumi/test:foo"},
					Push: true,
				},
				Digest: "sha256:foo",
			},
		})
		assert.NoError(t, err)
	})
}

func TestRead(t *testing.T) {
	t.Parallel()
	tag := "docker.io/pulumi/pulumitest"
	digest := "sha256:3be99cafdcd80a8e620da56bdc215acab6213bb608d3d492c0ba1807128786a1"

	ctrl := gomock.NewController(t)
	client := NewMockClient(ctrl)
	client.EXPECT().Inspect(gomock.Any(), fmt.Sprintf("%s:latest@%s", tag, digest)).Return(
		[]descriptor.Descriptor{
			{
				Platform: &platform.Platform{Architecture: "arm64"},
			},
			{
				Platform: &platform.Platform{Architecture: "unknown"},
			},
		}, nil)

	i := &Image{docker: client}

	resp, err := i.Read(t.Context(), infer.ReadRequest[ImageArgs, ImageState]{
		ID: "my-image",
		State: ImageState{
			ImageArgs: ImageArgs{
				Exports: []Export{{Raw: "type=registry"}},
				Tags:    []string{tag},
			},
			Digest: digest,
		},
	})

	require.NoError(t, err)
	assert.Equal(t, []string{tag}, resp.State.Tags)
}

func TestImageDiff(t *testing.T) {
	t.Parallel()
	emptyDir := t.TempDir()
	host := Host

	hash, err := hashBuildContext(emptyDir, "", nil)
	require.NoError(t, err)
	baseArgs := ImageArgs{
		Context:    &BuildContext{Context: Context{Location: emptyDir}},
		Dockerfile: &Dockerfile{Location: "testdata/noop"},
		Tags:       []string{},
	}
	baseState := ImageState{
		ContextHash: hash,
		ImageArgs:   baseArgs,
	}

	tests := []struct {
		name   string
		state  func(*testing.T, ImageState) ImageState
		inputs func(*testing.T, ImageArgs) ImageArgs

		wantChanges bool
	}{
		{
			name:        "no diff if build context is unchanged",
			state:       func(*testing.T, ImageState) ImageState { return baseState },
			inputs:      func(*testing.T, ImageArgs) ImageArgs { return baseArgs },
			wantChanges: false,
		},
		{
			name: "no diff if registry password changes",
			state: func(_ *testing.T, s ImageState) ImageState {
				s.Registries = []Registry{{
					Address:  "foo",
					Username: "foo",
					Password: "foo",
				}}
				return s
			},
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Registries = []Registry{{
					Address:  "foo",
					Username: "foo",
					Password: "DIFFERENT PASSWORD",
				}}
				return a
			},
			wantChanges: false,
		},
		{
			name: "no diff if pull=true but no exports",
			state: func(_ *testing.T, is ImageState) ImageState {
				is.Pull = true
				return is
			},
			inputs: func(t *testing.T, ia ImageArgs) ImageArgs {
				ia.Pull = true
				return ia
			},
			wantChanges: false,
		},
		{
			name: "diff if pull=true with exports",
			state: func(_ *testing.T, is ImageState) ImageState {
				is.Pull = true
				is.Load = true
				return is
			},
			inputs: func(t *testing.T, ia ImageArgs) ImageArgs {
				ia.Pull = true
				ia.Load = true
				return ia
			},
			wantChanges: true,
		},
		{
			name:  "diff if build context changes",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(t *testing.T, a ImageArgs) ImageArgs {
				tmp := filepath.Join(a.Context.Location, "tmp")
				err := os.WriteFile(tmp, []byte{}, 0o600)
				require.NoError(t, err)
				t.Cleanup(func() { _ = os.Remove(tmp) })
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if registry added",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Registries = []Registry{{}}
				return a
			},
			wantChanges: true,
		},
		{
			name: "diff if registry user changes",
			state: func(_ *testing.T, s ImageState) ImageState {
				s.Registries = []Registry{{
					Address:  "foo",
					Username: "foo",
					Password: "foo",
				}}
				return s
			},
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Registries = []Registry{{
					Address:  "DIFFERENT USER",
					Username: "foo",
					Password: "foo",
				}}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if buildArgs changes",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.BuildArgs = map[string]string{
					"foo": "bar",
				}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if pull changes",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(t *testing.T, ia ImageArgs) ImageArgs {
				ia.Pull = true
				return ia
			},
			wantChanges: true,
		},
		{
			name:  "diff if load changes",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(t *testing.T, ia ImageArgs) ImageArgs {
				ia.Load = true
				return ia
			},
			wantChanges: true,
		},
		{
			name:  "diff if push changes",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(t *testing.T, ia ImageArgs) ImageArgs {
				ia.Push = true
				return ia
			},
			wantChanges: true,
		},
		{
			name:  "diff if buildOnPreview doesn't change",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(t *testing.T, ia ImageArgs) ImageArgs {
				val := true
				ia.BuildOnPreview = &val
				return ia
			},
			wantChanges: true,
		},
		{
			name:  "diff if buildOnPreview changes",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(t *testing.T, ia ImageArgs) ImageArgs {
				val := false
				ia.BuildOnPreview = &val
				return ia
			},
			wantChanges: true,
		},
		{
			name:  "diff if ssh changes",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(t *testing.T, ia ImageArgs) ImageArgs {
				ia.SSH = []SSH{{ID: "default"}}
				return ia
			},
			wantChanges: true,
		},
		{
			name:  "diff if hosts change",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(t *testing.T, ia ImageArgs) ImageArgs {
				ia.AddHosts = []string{"localhost"}
				return ia
			},
			wantChanges: true,
		},
		{
			name:  "diff if cacheFrom changes",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.CacheFrom = []CacheFrom{{Raw: "a"}}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if cacheTo changes",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.CacheTo = []CacheTo{{Raw: "a"}}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if context changes",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Context = &BuildContext{Context: Context{Location: "testdata/ignores"}}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if named context changes",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Context = &BuildContext{Named: NamedContexts{"foo": Context{Location: "bar"}}}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if network changes",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Network = &host
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if dockerfile location changes",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Dockerfile = &Dockerfile{Location: "testdata/ignores/basedir/Dockerfile"}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if dockerfile inline changes",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Dockerfile = &Dockerfile{Inline: "FROM scratch"}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if platforms change",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Platforms = []Platform{"linux/amd64"}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if pull changes",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Pull = true
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if builder changes",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Builder = &BuilderConfig{Name: "foo"}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if tags change",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Tags = []string{"foo"}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if exports change",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Exports = []Export{{Raw: "foo"}}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if target changes",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Target = "foo"
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if pulling",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Pull = true
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if noCache changes",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.NoCache = true
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if labels change",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Labels = map[string]string{"foo": "bar"}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if secrets change",
			state: func(*testing.T, ImageState) ImageState { return baseState },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Secrets = map[string]string{"foo": "bar"}
				return a
			},
			wantChanges: true,
		},
		{
			name: "diff if local export doesn't exist",
			state: func(t *testing.T, state ImageState) ImageState {
				state.Exports = []Export{
					{Local: &ExportLocal{Dest: "not-real"}},
				}
				return state
			},
			inputs: func(_ *testing.T, args ImageArgs) ImageArgs {
				args.Exports = []Export{
					{Local: &ExportLocal{Dest: "not-real"}},
				}
				return args
			},
			wantChanges: true,
		},
		{
			name: "diff if tar export doesn't exist",
			state: func(t *testing.T, state ImageState) ImageState {
				state.Exports = []Export{
					{Tar: &ExportTar{ExportLocal: ExportLocal{Dest: "not-real"}}},
				}
				return state
			},
			inputs: func(_ *testing.T, args ImageArgs) ImageArgs {
				args.Exports = []Export{
					{Tar: &ExportTar{ExportLocal: ExportLocal{Dest: "not-real"}}},
				}
				return args
			},
			wantChanges: true,
		},
	}

	s := newServer(t.Context(), t, nil)

	encode := func(t *testing.T, x any) property.Map {
		raw, err := mapper.New(&mapper.Opts{IgnoreMissing: true}).Encode(x)
		require.NoError(t, err)
		return resource.FromResourcePropertyMap(resource.NewPropertyMapFromMap(raw))
	}

	for _, tt := range tests {
		baseState := baseState
		baseArgs := baseArgs
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			resp, err := s.Diff(provider.DiffRequest{
				Urn:    _fakeURN,
				State:  encode(t, tt.state(t, baseState)),
				Inputs: encode(t, tt.inputs(t, baseArgs)),
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.wantChanges, resp.HasChanges, resp.DetailedDiff)
		})
	}
}

func TestValidateImageArgs(t *testing.T) {
	t.Run("invalid inputs", func(t *testing.T) {
		t.Parallel()
		args := ImageArgs{
			Tags:      []string{"a/bad:tag:format"},
			Exports:   []Export{{Raw: "badexport,-"}},
			Context:   &BuildContext{Context: Context{Location: "./testdata"}},
			Platforms: []Platform{","},
			CacheFrom: []CacheFrom{{Raw: "=badcachefrom"}},
			CacheTo:   []CacheTo{{Raw: "=badcacheto"}},
		}

		_, err := args.validate(true, false)
		assert.ErrorContains(t, err, "invalid value badexport")
		assert.ErrorContains(t, err, "OSAndVersion specifier component must matc")
		assert.ErrorContains(t, err, "badcachefrom")
		assert.ErrorContains(t, err, "badcacheto")
		assert.ErrorContains(t, err, "invalid reference format")
		assert.ErrorContains(t, err, "testdata/Dockerfile")
	})

	t.Run("buildOnPreview", func(t *testing.T) {
		t.Parallel()
		args := ImageArgs{
			Context: &BuildContext{Context: Context{Location: "testdata/noop"}},
			Tags:    []string{"my-tag"},
			Exports: []Export{{Registry: &ExportRegistry{ExportImage{Push: pulumi.BoolRef(true)}}}},
		}
		actual, err := args.validate(true, true)
		assert.NoError(t, err)
		assert.Equal(t, "image", actual.Exports[0].Type)
		assert.Equal(t, "false", actual.Exports[0].Attrs["push"])

		actual, err = args.validate(true, false)
		assert.NoError(t, err)
		assert.Equal(t, "image", actual.Exports[0].Type)
		assert.Equal(t, "true", actual.Exports[0].Attrs["push"])
	})

	t.Run("unknowns", func(t *testing.T) {
		t.Parallel()
		// pulumi-go-provider gives us zero-values when a property is unknown.
		// We can't distinguish this from user-provided zero-values, but we
		// should:
		// - not fail previews due to these zero values,
		// - not attempt builds with invalid zero values,
		// - not allow invalid zero values in non-preview operations.
		unknowns := ImageArgs{
			BuildArgs: map[string]string{
				"known": "value",
				"":      "",
			},
			Builder:    nil,
			CacheFrom:  []CacheFrom{{GHA: &CacheFromGitHubActions{}}, {Raw: ""}},
			CacheTo:    []CacheTo{{GHA: &CacheToGitHubActions{}}, {Raw: ""}},
			Context:    nil,
			Exports:    []Export{{Raw: ""}},
			Dockerfile: nil,
			Platforms:  []Platform{"linux/amd64", ""},
			Registries: []Registry{
				{
					Address:  "",
					Password: "",
					Username: "",
				},
			},
			Tags: []string{"known", ""},
		}

		_, err := unknowns.validate(true, true)
		assert.NoError(t, err)
		assert.False(t, unknowns.buildable())

		_, err = unknowns.validate(true, false)
		assert.Error(t, err)
	})

	t.Run("disabled caches", func(t *testing.T) {
		t.Parallel()
		args := ImageArgs{
			Context:   &BuildContext{Context: Context{Location: "testdata/noop"}},
			CacheFrom: []CacheFrom{{Raw: "type=registry", Disabled: true}},
			CacheTo:   []CacheTo{{Raw: "type=registry", Disabled: true}},
			Exports:   []Export{{Raw: "type=registry", Disabled: true}},
		}

		opts, err := args.validate(true, true)
		assert.NoError(t, err)
		assert.Len(t, opts.CacheTo, 0)
		assert.Len(t, opts.CacheFrom, 0)
		assert.Len(t, opts.Exports, 0)

		opts, err = args.validate(true, false)
		assert.NoError(t, err)
		assert.Len(t, opts.CacheTo, 0)
		assert.Len(t, opts.CacheFrom, 0)
		assert.Len(t, opts.Exports, 0)
	})

	t.Run("environment variables", func(t *testing.T) {
		tests := []struct {
			name          string
			envs          map[string]string
			args          ImageArgs
			wantCacheFrom *pb.CacheOptionsEntry
			wantCacheTo   *pb.CacheOptionsEntry
		}{
			{
				name: "gha environment",
				envs: map[string]string{
					"ACTIONS_CACHE_URL":     "test-cache-url",
					"ACTIONS_RUNTIME_TOKEN": "test-runtime-token",
				},
				args: ImageArgs{
					Context:   &BuildContext{Context: Context{Location: "testdata/noop"}},
					CacheFrom: []CacheFrom{{GHA: &CacheFromGitHubActions{}}},
					CacheTo: []CacheTo{{GHA: &CacheToGitHubActions{
						CacheFromGitHubActions: CacheFromGitHubActions{},
					}}},
				},
				wantCacheFrom: &pb.CacheOptionsEntry{
					Type: "gha",
					Attrs: map[string]string{
						"token": "test-runtime-token",
						"url":   "test-cache-url",
					},
				},
				wantCacheTo: &pb.CacheOptionsEntry{
					Type: "gha",
					Attrs: map[string]string{
						"token": "test-runtime-token",
						"url":   "test-cache-url",
					},
				},
			},
			{
				name: "non-gha environment",
				envs: map[string]string{
					"ACTIONS_CACHE_URL":     "",
					"ACTIONS_RUNTIME_TOKEN": "",
				},
				args: ImageArgs{
					Context:   &BuildContext{Context: Context{Location: "testdata/noop"}},
					CacheFrom: []CacheFrom{{GHA: &CacheFromGitHubActions{}}},
					CacheTo: []CacheTo{{GHA: &CacheToGitHubActions{
						CacheFromGitHubActions: CacheFromGitHubActions{},
					}}},
				},
				wantCacheFrom: nil,
				wantCacheTo:   nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				for k, v := range tt.envs {
					t.Setenv(k, v)
				}
				validate := func(preview bool) {
					opts, err := tt.args.validate(true, preview)
					require.NoError(t, err)
					if tt.wantCacheFrom != nil {
						assert.Equal(t, tt.wantCacheFrom, opts.CacheFrom[0])
					} else {
						assert.Len(t, opts.CacheFrom, 0)
					}
					if tt.wantCacheTo != nil {
						assert.Equal(t, tt.wantCacheTo, opts.CacheTo[0])
					} else {
						assert.Len(t, opts.CacheTo, 0)
					}
				}
				validate(true)
				validate(false)
			})
		}
	})

	t.Run("multiple exports pre-0.13", func(t *testing.T) {
		t.Parallel()
		args := ImageArgs{
			Exports: []Export{{Raw: "type=local"}, {Raw: "type=tar"}},
		}
		_, err := args.validate(false, false)
		assert.ErrorContains(t, err, "multiple exports require a v0.13 buildkit daemon or newer")
	})

	t.Run("cache and export entries are union-ish", func(t *testing.T) {
		t.Parallel()
		args := ImageArgs{
			Exports:   []Export{{Tar: &ExportTar{}, Local: &ExportLocal{}}},
			CacheTo:   []CacheTo{{Raw: "type=tar", Local: &CacheToLocal{Dest: "/foo"}}},
			CacheFrom: []CacheFrom{{Raw: "type=tar", Registry: &CacheFromRegistry{}}},
		}
		_, err := args.validate(true, false)
		assert.ErrorContains(t, err, "exports should only specify one export type")
		assert.ErrorContains(t, err, "cacheFrom should only specify one cache type")
		assert.ErrorContains(t, err, "cacheTo should only specify one cache type")
	})

	t.Run("dockerfile parsing", func(t *testing.T) {
		t.Parallel()
		path := "./testdata/Dockerfile.invalid"
		data, err := os.ReadFile(path)
		require.NoError(t, err)

		for _, d := range []Dockerfile{
			{Location: path}, {Inline: string(data)},
		} {
			args := ImageArgs{Dockerfile: &d}
			_, err := args.validate(true, false)
			assert.ErrorContains(t, err, "unknown instruction: RUNN (did you mean RUN?)")
		}
	})
}

func TestBuildable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args ImageArgs

		want bool
	}{
		{
			name: "unknown tags",
			args: ImageArgs{Tags: []string{""}},
			want: false,
		},
		{
			name: "unknown exports",
			args: ImageArgs{
				Tags:    []string{"known"},
				Exports: []Export{{Raw: ""}},
			},
			want: false,
		},
		{
			name: "unknown registry",
			args: ImageArgs{
				Tags:    []string{"known"},
				Exports: []Export{{Docker: &ExportDocker{}}},
				Registries: []Registry{
					{
						Address:  "docker.io",
						Username: "foo",
						Password: "",
					},
				},
			},
			want: false,
		},
		{
			name: "known tags",
			args: ImageArgs{
				Tags: []string{"known"},
			},
			want: true,
		},
		{
			name: "known exports",
			args: ImageArgs{
				Tags:    []string{"known"},
				Exports: []Export{{Registry: &ExportRegistry{}}},
			},
			want: true,
		},
		{
			name: "known registry",
			args: ImageArgs{
				Tags:    []string{"known"},
				Exports: []Export{{Registry: &ExportRegistry{}}},
				Registries: []Registry{
					{
						Address:  "docker.io",
						Username: "foo",
						Password: "bar",
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.args.buildable()
			assert.Equal(t, tt.want, actual)
		})
	}
}

func TestToBuild(t *testing.T) {
	t.Parallel()
	Max := Max

	ia := ImageArgs{
		Tags:      []string{"foo", "bar"},
		Platforms: []Platform{"linux/amd64"},
		Context:   &BuildContext{Context: Context{Location: "testdata/noop"}},
		CacheTo: []CacheTo{
			{GHA: &CacheToGitHubActions{CacheWithMode: CacheWithMode{&Max}}},
			{
				Registry: &CacheToRegistry{
					CacheFromRegistry: CacheFromRegistry{Ref: "docker.io/foo/bar"},
				},
			},
			{
				Registry: &CacheToRegistry{
					CacheFromRegistry: CacheFromRegistry{Ref: "docker.io/foo/bar:baz"},
				},
			},
		},
		CacheFrom: []CacheFrom{
			{S3: &CacheFromS3{Name: "bar"}},
			{Registry: &CacheFromRegistry{Ref: "docker.io/foo/bar"}},
			{Registry: &CacheFromRegistry{Ref: "docker.io/foo/bar:baz"}},
		},
	}

	_, err := ia.toBuild(context.Background(), true, false)
	assert.NoError(t, err)
}
