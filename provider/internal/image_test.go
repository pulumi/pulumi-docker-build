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
	"github.com/docker/buildx/util/buildflags"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	"github.com/regclient/regclient/types/descriptor"
	"github.com/regclient/regclient/types/errs"
	"github.com/regclient/regclient/types/platform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
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
			op: func(_ *testing.T) integration.Operation {
				return integration.Operation{
					Inputs: property.NewMap(map[string]property.Value{
						pushKey: property.New(false),
						tagsKey: property.New(
							[]property.Value{
								property.New("docker.io/pulumibot/buildkit-e2e"),
								property.New("docker.io/pulumibot/buildkit-e2e:main"),
							},
						),
						"platforms": property.New(
							[]property.Value{
								property.New("linux/arm64"),
								property.New(platformLinuxAMD64),
							},
						),
						contextKey: property.New(map[string]property.Value{
							locationKey: property.New(testdataNoop),
						}),
						"dockerfile": property.New(map[string]property.Value{
							locationKey: property.New("testdata/noop/Dockerfile"),
						}),
						exportsKey: property.New(
							[]property.Value{
								property.New(map[string]property.Value{
									rawKey: property.New(typeRegistry),
								},
								),
							},
						),
						"registries": property.New(
							[]property.Value{
								property.New(map[string]property.Value{
									addressKey:  property.New("fakeaddress"),
									usernameKey: property.New("fakeuser"),
									passwordKey: property.New(passwordKey).WithSecret(true),
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
			op: func(_ *testing.T) integration.Operation {
				return integration.Operation{
					Inputs: property.NewMap(map[string]property.Value{
						pushKey: property.New(false),
						tagsKey: property.New([]property.Value{}),
						contextKey: property.New(map[string]property.Value{
							locationKey: property.New(testdataNoop),
						}),
						exportsKey: property.New(
							[]property.Value{
								property.New(map[string]property.Value{
									rawKey: property.New(typeRegistry),
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
			op: func(_ *testing.T) integration.Operation {
				return integration.Operation{
					Inputs: property.NewMap(map[string]property.Value{
						pushKey: property.New(false),
						tagsKey: property.New(
							[]property.Value{property.New("invalid-exports")},
						),
						exportsKey: property.New(
							[]property.Value{
								property.New(map[string]property.Value{
									rawKey: property.New("type="),
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
			op: func(_ *testing.T) integration.Operation {
				return integration.Operation{
					Inputs: property.NewMap(map[string]property.Value{
						pushKey: property.New(false),
						tagsKey: property.New(
							[]property.Value{property.New(fooName)},
						),
						contextKey: property.New(map[string]property.Value{
							locationKey: property.New(testdataNoop),
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
			op: func(_ *testing.T) integration.Operation {
				return integration.Operation{
					Inputs: property.NewMap(map[string]property.Value{
						pushKey: property.New(false),
						tagsKey: property.New(
							[]property.Value{property.New(fooName)},
						),
						contextKey: property.New(map[string]property.Value{
							locationKey: property.New(testdataNoop),
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
			op: func(_ *testing.T) integration.Operation {
				return integration.Operation{
					Inputs: property.NewMap(map[string]property.Value{
						pushKey: property.New(false),
						tagsKey: property.New(
							[]property.Value{
								property.New("default-dockerfile"),
							},
						),
						contextKey: property.New(map[string]property.Value{
							locationKey: property.New(testdataNoop),
						}),
					}),
					Hook: func(_, output property.Map) {
						dockerfile := output.Get("dockerfile")
						require.NotNil(t, dockerfile)
						require.True(t, dockerfile.IsMap())
						location := dockerfile.AsMap().Get(locationKey)
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
			op: func(_ *testing.T) integration.Operation {
				return integration.Operation{
					Inputs: property.NewMap(map[string]property.Value{
						pushKey: property.New(false),
						tagsKey: property.New(
							[]property.Value{
								property.New("inline-dockerfile"),
							},
						),
						"buildOnPreview": property.New(true),
						"dockerfile": property.New(map[string]property.Value{
							inlineName: property.New("FROM alpine:latest"),
						}),
					}),
					Hook: func(_, output property.Map) {
						context := output.Get(contextKey)
						require.NotNil(t, context)
						require.True(t, context.IsMap())
						location := context.AsMap().Get(locationKey)
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
			s := newServer(t.Context(), t, mockClientF(tt.client(t)))

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

		i := &Image{clientF: mockClientF(client)}

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

	i := &Image{clientF: mockClientF(client)}

	resp, err := i.Read(t.Context(), infer.ReadRequest[ImageArgs, ImageState]{
		ID: "my-image",
		State: ImageState{
			ImageArgs: ImageArgs{
				Exports: []Export{{Raw: typeRegistry}},
				Tags:    []string{tag},
			},
			Digest: digest,
		},
	})

	require.NoError(t, err)
	assert.Equal(t, []string{tag}, resp.State.Tags)
}

// Read deletes state only on a definitive not-found, not on auth/network errors.
func TestReadInspectError(t *testing.T) {
	t.Parallel()
	tag := "docker.io/pulumi/pulumitest"
	digest := "sha256:3be99cafdcd80a8e620da56bdc215acab6213bb608d3d492c0ba1807128786a1"
	ref := fmt.Sprintf("%s:latest@%s", tag, digest)

	newImage := func(inspectErr error) (*Image, infer.ReadRequest[ImageArgs, ImageState]) {
		ctrl := gomock.NewController(t)
		client := NewMockClient(ctrl)
		client.EXPECT().Inspect(gomock.Any(), ref).Return(nil, inspectErr)
		args := ImageArgs{
			Exports: []Export{{Raw: "type=registry"}},
			Tags:    []string{tag},
		}
		return &Image{clientF: mockClientF(client)}, infer.ReadRequest[ImageArgs, ImageState]{
			ID:     "my-image",
			Inputs: args, // deletion guard reads Inputs.Tags
			State:  ImageState{ImageArgs: args, Digest: digest},
		}
	}

	t.Run("indeterminate error preserves state", func(t *testing.T) {
		t.Parallel()
		// Stands in for ECR's HTTP 403 expired-token and any transient failure.
		i, req := newImage(errors.New("authorization token has expired"))
		resp, err := i.Read(t.Context(), req)
		require.NoError(t, err)
		assert.Equal(t, "my-image", resp.ID, "resource must not be deleted on an indeterminate error")
		assert.Equal(t, []string{tag}, resp.State.Tags)
	})

	t.Run("not found deletes resource", func(t *testing.T) {
		t.Parallel()
		i, req := newImage(errs.ErrNotFound)
		resp, err := i.Read(t.Context(), req)
		require.NoError(t, err)
		assert.Empty(t, resp.ID, "a genuine not-found must still delete the resource")
	})
}

func TestImageDiff(t *testing.T) {
	t.Parallel()
	host := Host

	baseArgs := ImageArgs{
		Dockerfile: &Dockerfile{Location: testdataNoop},
		Tags:       []string{},
	}
	baseState := ImageState{
		ImageArgs: baseArgs,
	}

	tests := []struct {
		name   string
		state  func(*testing.T, ImageState) ImageState
		inputs func(*testing.T, ImageArgs) ImageArgs

		wantChanges bool
	}{
		{
			name:        "no diff if build context is unchanged",
			state:       func(_ *testing.T, s ImageState) ImageState { return s },
			inputs:      func(_ *testing.T, a ImageArgs) ImageArgs { return a },
			wantChanges: false,
		},
		{
			name: "no diff if registry password changes",
			state: func(_ *testing.T, s ImageState) ImageState {
				s.Registries = []Registry{{
					Address:  fooName,
					Username: fooName,
					Password: fooName,
				}}
				return s
			},
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Registries = []Registry{{
					Address:  fooName,
					Username: fooName,
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
			inputs: func(_ *testing.T, ia ImageArgs) ImageArgs {
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
			inputs: func(_ *testing.T, ia ImageArgs) ImageArgs {
				ia.Pull = true
				ia.Load = true
				return ia
			},
			wantChanges: true,
		},
		{
			name:  "diff if build context changes (same location, changed contents)",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(t *testing.T, a ImageArgs) ImageArgs {
				err := os.WriteFile(filepath.Join(a.Context.Location, "tmp"), []byte{}, 0o600)
				require.NoError(t, err)
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if registry added",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
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
					Address:  fooName,
					Username: fooName,
					Password: fooName,
				}}
				return s
			},
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Registries = []Registry{{
					Address:  "DIFFERENT USER",
					Username: fooName,
					Password: fooName,
				}}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if buildArgs changes",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.BuildArgs = map[string]string{
					fooName: barName,
				}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if pull changes",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, ia ImageArgs) ImageArgs {
				ia.Pull = true
				return ia
			},
			wantChanges: true,
		},
		{
			name:  "diff if load changes",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, ia ImageArgs) ImageArgs {
				ia.Load = true
				return ia
			},
			wantChanges: true,
		},
		{
			name:  "diff if push changes",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, ia ImageArgs) ImageArgs {
				ia.Push = true
				return ia
			},
			wantChanges: true,
		},
		{
			name:  "diff if buildOnPreview doesn't change",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, ia ImageArgs) ImageArgs {
				val := true
				ia.BuildOnPreview = &val
				return ia
			},
			wantChanges: true,
		},
		{
			name:  "diff if buildOnPreview changes",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, ia ImageArgs) ImageArgs {
				val := false
				ia.BuildOnPreview = &val
				return ia
			},
			wantChanges: true,
		},
		{
			name:  "diff if ssh changes",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, ia ImageArgs) ImageArgs {
				ia.SSH = []SSH{{ID: defaultSSHID}}
				return ia
			},
			wantChanges: true,
		},
		{
			name:  "diff if hosts change",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, ia ImageArgs) ImageArgs {
				ia.AddHosts = []string{"localhost"}
				return ia
			},
			wantChanges: true,
		},
		{
			name:  "diff if cacheFrom changes",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.CacheFrom = []CacheFrom{{Raw: "a"}}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if cacheTo changes",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.CacheTo = []CacheTo{{Raw: "a"}}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if context changes",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Context = &BuildContext{Context: Context{Location: "testdata/ignores"}}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if named context changes",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Context = &BuildContext{Named: NamedContexts{fooName: Context{Location: barName}}}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if network changes",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Network = &host
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if dockerfile location changes",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Dockerfile = &Dockerfile{Location: "testdata/ignores/basedir/Dockerfile"}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if dockerfile inline changes",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Dockerfile = &Dockerfile{Inline: fromScratch}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if platforms change",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Platforms = []Platform{platformLinuxAMD64}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if pull changes",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Pull = true
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if builder changes",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Builder = &BuilderConfig{Name: fooName}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if tags change",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Tags = []string{fooName}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if exports change",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Exports = []Export{{Raw: fooName}}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if target changes",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Target = fooName
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if pulling",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Pull = true
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if noCache changes",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.NoCache = true
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if labels change",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Labels = map[string]string{fooName: barName}
				return a
			},
			wantChanges: true,
		},
		{
			name:  "diff if secrets change",
			state: func(_ *testing.T, s ImageState) ImageState { return s },
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Secrets = map[string]string{fooName: barName}
				return a
			},
			wantChanges: true,
		},
		{
			name: "diff if secrets change but ignoreSecretsInDiffCalculation is set",
			state: func(_ *testing.T, s ImageState) ImageState {
				s.Secrets = map[string]string{fooName: oldBar}
				s.IgnoreSecretsInDiffCalculation = []string{fooName}
				return s
			},
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Secrets = map[string]string{fooName: newBar}
				a.IgnoreSecretsInDiffCalculation = []string{fooName}
				return a
			},
			wantChanges: false,
		},
		{
			name: "diff if secrets change but ignoreSecretsInDiffCalculation is set for another secret",
			state: func(_ *testing.T, s ImageState) ImageState {
				s.Secrets = map[string]string{fooName: oldBar}
				s.IgnoreSecretsInDiffCalculation = []string{"not_foo"}
				return s
			},
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Secrets = map[string]string{fooName: newBar}
				a.IgnoreSecretsInDiffCalculation = []string{"not_foo"}
				return a
			},
			wantChanges: true,
		},
		{
			name: "diff if secrets change but ignoreSecretsInDiffCalculation is set and secret is added",
			state: func(_ *testing.T, s ImageState) ImageState {
				s.IgnoreSecretsInDiffCalculation = []string{fooName}
				return s
			},
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Secrets = map[string]string{fooName: barName}
				a.IgnoreSecretsInDiffCalculation = []string{fooName}
				return a
			},
			wantChanges: true,
		},
		{
			name: "diff if secrets change but ignoreSecretsInDiffCalculation is set and secret is removed",
			state: func(_ *testing.T, s ImageState) ImageState {
				s.Secrets = map[string]string{fooName: barName}
				s.IgnoreSecretsInDiffCalculation = []string{fooName}
				return s
			},
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.IgnoreSecretsInDiffCalculation = []string{fooName}
				return a
			},
			wantChanges: true,
		},
		{
			name: "diff if an ignored secret and a non-ignored secret both change",
			state: func(_ *testing.T, s ImageState) ImageState {
				s.Secrets = map[string]string{fooName: "old_foo", barName: oldBar}
				s.IgnoreSecretsInDiffCalculation = []string{fooName}
				return s
			},
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Secrets = map[string]string{fooName: "new_foo", barName: newBar}
				a.IgnoreSecretsInDiffCalculation = []string{fooName}
				return a
			},
			wantChanges: true,
		},
		{
			name: "no diff if only ignoreSecretsInDiffCalculation changes",
			state: func(_ *testing.T, s ImageState) ImageState {
				s.Secrets = map[string]string{fooName: barName}
				return s
			},
			inputs: func(_ *testing.T, a ImageArgs) ImageArgs {
				a.Secrets = map[string]string{fooName: barName}
				a.IgnoreSecretsInDiffCalculation = []string{fooName}
				return a
			},
			wantChanges: false,
		},
		{
			name: "diff if local export doesn't exist",
			state: func(_ *testing.T, state ImageState) ImageState {
				state.Exports = []Export{
					{Local: &ExportLocal{Dest: notReal}},
				}
				return state
			},
			inputs: func(_ *testing.T, args ImageArgs) ImageArgs {
				args.Exports = []Export{
					{Local: &ExportLocal{Dest: notReal}},
				}
				return args
			},
			wantChanges: true,
		},
		{
			name: "diff if tar export doesn't exist",
			state: func(_ *testing.T, state ImageState) ImageState {
				state.Exports = []Export{
					{Tar: &ExportTar{ExportLocal: ExportLocal{Dest: notReal}}},
				}
				return state
			},
			inputs: func(_ *testing.T, args ImageArgs) ImageArgs {
				args.Exports = []Export{
					{Tar: &ExportTar{ExportLocal: ExportLocal{Dest: notReal}}},
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Per-subtest context dir so parallel subtests never share one.
			dir := t.TempDir()
			hash, err := hashBuildContext(dir, "", nil)
			require.NoError(t, err)

			baseState := baseState
			baseState.Context = &BuildContext{Context: Context{Location: dir}}
			baseState.ContextHash = hash
			baseArgs := baseArgs
			baseArgs.Context = &BuildContext{Context: Context{Location: dir}}

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
			Context: &BuildContext{Context: Context{Location: testdataNoop}},
			Tags:    []string{"my-tag"},
			Exports: []Export{{Registry: &ExportRegistry{ExportImage{Push: pulumi.BoolRef(true)}}}},
		}
		actual, err := args.validate(true, true)
		assert.NoError(t, err)
		assert.Equal(t, exportTypeImage, actual.Exports[0].Type)
		assert.Equal(t, falseLiteral, actual.Exports[0].Attrs[pushKey])

		actual, err = args.validate(true, false)
		assert.NoError(t, err)
		assert.Equal(t, exportTypeImage, actual.Exports[0].Type)
		assert.Equal(t, "true", actual.Exports[0].Attrs[pushKey])
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
				knownKey: "value",
				"":       "",
			},
			Builder:    nil,
			CacheFrom:  []CacheFrom{{GHA: &CacheFromGitHubActions{}}, {Raw: ""}},
			CacheTo:    []CacheTo{{GHA: &CacheToGitHubActions{}}, {Raw: ""}},
			Context:    nil,
			Exports:    []Export{{Raw: ""}},
			Dockerfile: nil,
			Platforms:  []Platform{platformLinuxAMD64, ""},
			Registries: []Registry{
				{
					Address:  "",
					Password: "",
					Username: "",
				},
			},
			Tags: []string{knownKey, ""},
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
			Context:   &BuildContext{Context: Context{Location: testdataNoop}},
			CacheFrom: []CacheFrom{{Raw: typeRegistry, Disabled: true}},
			CacheTo:   []CacheTo{{Raw: typeRegistry, Disabled: true}},
			Exports:   []Export{{Raw: typeRegistry, Disabled: true}},
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
			wantCacheFrom *buildflags.CacheOptionsEntry
			wantCacheTo   *buildflags.CacheOptionsEntry
		}{
			{
				name: "gha environment",
				envs: map[string]string{
					"ACTIONS_CACHE_URL":        testCacheURL,
					"ACTIONS_RUNTIME_TOKEN":    testRuntimeToken,
					"ACTIONS_RESULTS_URL":      testResultsURL,
					"ACTIONS_CACHE_SERVICE_V2": "true",
				},
				args: ImageArgs{
					Context:   &BuildContext{Context: Context{Location: testdataNoop}},
					CacheFrom: []CacheFrom{{GHA: &CacheFromGitHubActions{}}},
					CacheTo: []CacheTo{{GHA: &CacheToGitHubActions{
						CacheFromGitHubActions: CacheFromGitHubActions{},
					}}},
				},
				wantCacheFrom: &buildflags.CacheOptionsEntry{
					Type: cacheTypeGHA,
					Attrs: map[string]string{
						"token":  testRuntimeToken,
						"url":    testCacheURL,
						"url_v2": testResultsURL,
					},
				},
				wantCacheTo: &buildflags.CacheOptionsEntry{
					Type: cacheTypeGHA,
					Attrs: map[string]string{
						"token":  testRuntimeToken,
						"url":    testCacheURL,
						"url_v2": testResultsURL,
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
					Context:   &BuildContext{Context: Context{Location: testdataNoop}},
					CacheFrom: []CacheFrom{{GHA: &CacheFromGitHubActions{}}},
					CacheTo: []CacheTo{{GHA: &CacheToGitHubActions{
						CacheFromGitHubActions: CacheFromGitHubActions{},
					}}},
				},
				wantCacheFrom: nil,
				wantCacheTo:   nil,
			},
			{
				name: "s3 environment",
				envs: map[string]string{
					// Env creds resolve first in the AWS chain, so
					// addAwsCredentials injects these without reaching IMDS.
					"AWS_ACCESS_KEY_ID":         "test-access-key",
					"AWS_SECRET_ACCESS_KEY":     "test-secret-key", //nolint:gosec // test fixture, not a real credential.
					"AWS_SESSION_TOKEN":         "test-session-token",
					"AWS_EC2_METADATA_DISABLED": trueLiteral, // never touch IMDS
				},
				args: ImageArgs{
					Context:   &BuildContext{Context: Context{Location: testdataNoop}},
					CacheFrom: []CacheFrom{{S3: &CacheFromS3{Bucket: "my-bucket", Name: barName}}},
				},
				wantCacheFrom: &buildflags.CacheOptionsEntry{
					Type: "s3",
					Attrs: map[string]string{
						"bucket":            "my-bucket",
						"name":              barName,
						"access_key_id":     "test-access-key",
						"secret_access_key": "test-secret-key",
						"session_token":     "test-session-token",
					},
				},
				wantCacheTo: nil,
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
			Exports: []Export{{Raw: "type=local"}, {Raw: typeTar}},
		}
		_, err := args.validate(false, false)
		assert.ErrorContains(t, err, "multiple exports require a v0.13 buildkit daemon or newer")
	})

	t.Run("cache and export entries are union-ish", func(t *testing.T) {
		t.Parallel()
		args := ImageArgs{
			Exports:   []Export{{Tar: &ExportTar{}, Local: &ExportLocal{}}},
			CacheTo:   []CacheTo{{Raw: typeTar, Local: &CacheToLocal{Dest: fooDest}}},
			CacheFrom: []CacheFrom{{Raw: typeTar, Registry: &CacheFromRegistry{}}},
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
				Tags:    []string{knownKey},
				Exports: []Export{{Raw: ""}},
			},
			want: false,
		},
		{
			name: "unknown registry",
			args: ImageArgs{
				Tags:    []string{knownKey},
				Exports: []Export{{Docker: &ExportDocker{}}},
				Registries: []Registry{
					{
						Address:  dockerIO,
						Username: fooName,
						Password: "",
					},
				},
			},
			want: false,
		},
		{
			name: "known tags",
			args: ImageArgs{
				Tags: []string{knownKey},
			},
			want: true,
		},
		{
			name: "known exports",
			args: ImageArgs{
				Tags:    []string{knownKey},
				Exports: []Export{{Registry: &ExportRegistry{}}},
			},
			want: true,
		},
		{
			name: "known registry",
			args: ImageArgs{
				Tags:    []string{knownKey},
				Exports: []Export{{Registry: &ExportRegistry{}}},
				Registries: []Registry{
					{
						Address:  dockerIO,
						Username: fooName,
						Password: barName,
					},
				},
			},
			want: true,
		},
		{
			name: "known with ignoreSecretsInDiffCalculation set",
			args: ImageArgs{
				Tags:                           []string{knownKey},
				Secrets:                        map[string]string{fooName: barName},
				IgnoreSecretsInDiffCalculation: []string{fooName},
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
		Tags:      []string{fooName, barName},
		Platforms: []Platform{platformLinuxAMD64},
		Context:   &BuildContext{Context: Context{Location: testdataNoop}},
		CacheTo: []CacheTo{
			{GHA: &CacheToGitHubActions{CacheWithMode: CacheWithMode{&Max}}},
			{
				Registry: &CacheToRegistry{
					CacheFromRegistry: CacheFromRegistry{Ref: dockerIORef},
				},
			},
			{
				Registry: &CacheToRegistry{
					CacheFromRegistry: CacheFromRegistry{Ref: dockerIOTaggedRef},
				},
			},
		},
		CacheFrom: []CacheFrom{
			{S3: &CacheFromS3{Name: barName}},
			{Registry: &CacheFromRegistry{Ref: dockerIORef}},
			{Registry: &CacheFromRegistry{Ref: dockerIOTaggedRef}},
		},
	}

	_, err := ia.toBuild(context.Background(), true, false)
	assert.NoError(t, err)
}
