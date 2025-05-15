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
)

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
