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

package provider

import (
	"context"
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"

	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	rpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

// TestConfigure checks backwards-compatibility with SDKs that still send
// provider config as JSON-encoded strings. This test can be removed once we
// upgrade to a version of pulumi that no longer generates SDKs with that
// behavior.
func TestConfigure(t *testing.T) {
	t.Parallel()

	p, err := New(nil)
	require.NoError(t, err)

	args, err := structpb.NewStruct(map[string]any{
		"registries": `[{"address": "docker.io"}]`,
	})
	require.NoError(t, err)

	// Ideally we would just call p.Configure directly, but we need the
	// integration server to inject runtime info for us.
	argsMap, err := plugin.UnmarshalProperties(args, plugin.MarshalOptions{})
	require.NoError(t, err)
	s := integration.NewServer("docker-build", semver.Version{Major: 0}, provider.Provider{
		// Roundabout way to get the integration server to invoke our outermost
		// Configure RPC endpoint.
		Configure: func(ctx context.Context, req provider.ConfigureRequest) error {
			args, err := plugin.MarshalProperties(req.Args, plugin.MarshalOptions{})
			require.NoError(t, err)
			_, err = p.Configure(ctx, &rpc.ConfigureRequest{
				Variables: req.Variables,
				Args:      args,
			})
			return err
		},
	})
	err = s.Configure(provider.ConfigureRequest{
		Args: argsMap,
	})

	assert.NoError(t, err)
}

func TestVersion(t *testing.T) {
	t.Parallel()

	_, err := semver.Parse(Version)
	assert.NoError(t, err)

	p, err := New(nil)
	require.NoError(t, err)

	info, err := p.GetPluginInfo(context.Background(), &emptypb.Empty{})
	assert.NoError(t, err)

	require.NotEqual(t, "", Version)
	assert.Equal(t, Version, info.Version)
}
