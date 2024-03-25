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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

// TestConfigure checks backwards-compatibility with SDKs that still send
// provider config as JSON-encoded strings. This test can be removed once we
// upgrade to a version of pulumi that no longer generates SDKs with that
// behavior.
func TestConfigure(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	p, err := New(nil)
	require.NoError(t, err)

	args, err := structpb.NewStruct(map[string]any{
		"registries": `[{"address": "docker.io"}]`,
	})
	require.NoError(t, err)

	_, err = p.Configure(ctx, &pulumirpc.ConfigureRequest{
		Args: args,
	})

	assert.NoError(t, err)
}
