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
	"os"
	"testing"

	"github.com/regclient/regclient/types/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/mapper"
	"github.com/pulumi/pulumi/sdk/v3/go/property"
)

func TestIndexLifecycle(t *testing.T) {
	t.Parallel()
	realClient := func(t *testing.T) Client { return nil }

	tests := []struct {
		name string
		skip bool

		op     func(t *testing.T) integration.Operation
		client func(t *testing.T) Client
	}{
		{
			name:   "not pushed",
			client: realClient,
			op: func(t *testing.T) integration.Operation {
				return integration.Operation{
					Inputs: property.NewMap(map[string]property.Value{
						"tag": property.New(
							"docker.io/pulumibot/buildkit-e2e:manifest-unit",
						),
						"sources": property.New([]property.Value{
							property.New("docker.io/pulumibot/buildkit-e2e:arm64"),
							property.New("docker.io/pulumibot/buildkit-e2e:amd64"),
						}),
						"push": property.New(false),
					}),
				}
			},
		},
		{
			name:   "pushed",
			skip:   os.Getenv("DOCKER_HUB_PASSWORD") == "",
			client: realClient,
			op: func(t *testing.T) integration.Operation {
				return integration.Operation{
					Inputs: property.NewMap(map[string]property.Value{
						"tag": property.New(
							"docker.io/pulumibot/buildkit-e2e:manifest",
						),
						"sources": property.New([]property.Value{
							property.New("docker.io/pulumibot/buildkit-e2e:arm64"),
							property.New("docker.io/pulumibot/buildkit-e2e:amd64"),
						}),
						"push": property.New(true),
						"registry": property.New(map[string]property.Value{
							"address":  property.New("docker.io"),
							"username": property.New("pulumibot"),
							"password": property.New(os.Getenv("DOCKER_HUB_PASSWORD")).WithSecret(true),
						}),
					}),
				}
			},
		},
		{
			name: "expired credentials",
			client: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				c := NewMockClient(ctrl)
				c.EXPECT().ManifestCreate(gomock.Any(), true, gomock.Any(), gomock.Any())
				c.EXPECT().ManifestInspect(gomock.Any(), gomock.Any()).Return("", errs.ErrHTTPUnauthorized)
				c.EXPECT().ManifestDelete(gomock.Any(), gomock.Any()).Return(nil)
				return c
			},
			op: func(t *testing.T) integration.Operation {
				return integration.Operation{
					Inputs: property.NewMap(map[string]property.Value{
						"tag": property.New(
							"docker.io/pulumibot/buildkit-e2e:manifest",
						),
						"sources": property.New([]property.Value{
							property.New("docker.io/pulumibot/buildkit-e2e:arm64"),
							property.New("docker.io/pulumibot/buildkit-e2e:amd64"),
						}),
						"push": property.New(true),
						"registry": property.New(map[string]property.Value{
							"address":  property.New("docker.io"),
							"username": property.New("pulumibot"),
							"password": property.New(os.Getenv("DOCKER_HUB_PASSWORD")).WithSecret(true),
						}),
					}),
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.skip {
				t.Skip("missing environment variables")
			}
			lc := integration.LifeCycleTest{
				Resource: "docker-build:index:Index",
				Create:   tt.op(t),
			}
			s := newServer(tt.client(t))

			err := s.Configure(provider.ConfigureRequest{})
			require.NoError(t, err)

			lc.Run(t, s)
		})
	}
}

func TestIndexDiff(t *testing.T) {
	t.Parallel()
	urn := resource.NewURN("test", "provider", "a", "docker-build:index:Index", "test")
	baseArgs := IndexArgs{Sources: []string{"docker.io/nginx:latest"}}
	baseState := IndexState{IndexArgs: baseArgs}

	tests := []struct {
		name   string
		state  func(*testing.T, IndexState) IndexState
		inputs func(*testing.T, IndexArgs) IndexArgs

		wantChanges bool
	}{
		{
			name:        "no diff if no changes",
			state:       func(*testing.T, IndexState) IndexState { return baseState },
			inputs:      func(*testing.T, IndexArgs) IndexArgs { return baseArgs },
			wantChanges: false,
		},
		{
			name:  "diff if tag changes",
			state: func(*testing.T, IndexState) IndexState { return baseState },
			inputs: func(t *testing.T, a IndexArgs) IndexArgs {
				a.Tag = "new-tag"
				return a
			},
			wantChanges: true,
		},
		{
			name: "no diff if registry password changes",
			state: func(_ *testing.T, s IndexState) IndexState {
				s.Registry = &Registry{
					Address:  "foo",
					Username: "foo",
					Password: "foo",
				}
				return s
			},
			inputs: func(_ *testing.T, a IndexArgs) IndexArgs {
				a.Registry = &Registry{
					Address:  "foo",
					Username: "foo",
					Password: "DIFFERENT PASSWORD",
				}
				return a
			},
			wantChanges: false,
		},
		{
			name:  "diff if registry added",
			state: func(*testing.T, IndexState) IndexState { return baseState },
			inputs: func(_ *testing.T, a IndexArgs) IndexArgs {
				a.Registry = &Registry{Address: "foo.com", Username: "foo", Password: "foo"}
				return a
			},
			wantChanges: true,
		},
		{
			name: "diff if registry user changes",
			state: func(_ *testing.T, s IndexState) IndexState {
				s.Registry = &Registry{
					Address:  "foo",
					Username: "foo",
					Password: "foo",
				}
				return s
			},
			inputs: func(_ *testing.T, a IndexArgs) IndexArgs {
				a.Registry = &Registry{
					Address:  "DIFFERENT USER",
					Username: "foo",
					Password: "foo",
				}
				return a
			},
			wantChanges: true,
		},
	}

	s := newServer(nil)

	encode := func(t *testing.T, x any) property.Map {
		raw, err := mapper.New(&mapper.Opts{IgnoreMissing: true}).Encode(x)
		require.NoError(t, err)
		return resource.FromResourcePropertyMap(resource.NewPropertyMapFromMap(raw))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			resp, err := s.Diff(provider.DiffRequest{
				Urn:    urn,
				State:  encode(t, tt.state(t, baseState)),
				Inputs: encode(t, tt.inputs(t, baseArgs)),
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.wantChanges, resp.HasChanges, resp.DetailedDiff)
		})
	}
}
