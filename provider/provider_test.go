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

	"github.com/pulumi/pulumi-docker-build/provider/internal"
	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
)

// TestConfigure checks backwards-compatibility with SDKs that still send
// provider config as JSON-encoded strings. This test can be removed once we
// upgrade to a version of pulumi that no longer generates SDKs with that
// behavior.
func TestConfigure(t *testing.T) {
	t.Parallel()

	p := internal.NewBuildxProvider()

	args, err := structpb.NewStruct(map[string]any{
		"registries": `[{"address": "docker.io"}]`,
	})
	require.NoError(t, err)
	argsMap, err := plugin.UnmarshalProperties(args, plugin.MarshalOptions{})
	require.NoError(t, err)

	s := integration.NewServer("docker-build", semver.Version{Major: 0}, p)
	resp, err := s.CheckConfig(provider.CheckRequest{
		News: argsMap,
	})
	require.NoError(t, err)
	require.Empty(t, resp.Failures)
	assert.Equal(t, resource.PropertyMap{
		"registries": resource.NewProperty([]resource.PropertyValue{
			resource.NewProperty(resource.PropertyMap{
				"address":  resource.NewProperty("docker.io"),
				"username": resource.NewProperty(""),
				"password": resource.MakeSecret(resource.NewProperty("")),
			}),
		}),
		"host": resource.NewProperty(""),
	}, resp.Inputs)
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
