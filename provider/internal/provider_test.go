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
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-go-provider/integration"
	mwcontext "github.com/pulumi/pulumi-go-provider/middleware/context"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

func TestConfigure(t *testing.T) {
	t.Parallel()

	s := newServer(t.Context(), t, nil)

	err := s.Configure(
		provider.ConfigureRequest{},
	)
	assert.NoError(t, err)
}

// TestAnnotate sanity checks that our annotations don't panic.
func TestAnnotate(t *testing.T) {
	t.Parallel()

	for _, tt := range []infer.Annotated{
		&Config{},
		&Image{},
		&ImageArgs{},
		&ImageState{},
		&Index{},
		&IndexArgs{},
		&IndexState{},
	} {
		tt.Annotate(annotator{})
	}
}

// TestSchema sanity checks that our schema doesn't panic.
func TestSchema(t *testing.T) {
	t.Parallel()

	s := newServer(t.Context(), t, nil)

	_, err := s.GetSchema(provider.GetSchemaRequest{Version: 0})
	assert.NoError(t, err)
}

type annotator struct{}

func (annotator) Deprecate(_ any, _ string)                   {}
func (annotator) Describe(_ any, _ string)                    {}
func (annotator) SetDefault(_, _ any, _ ...string)            {}
func (annotator) SetToken(tokens.ModuleName, tokens.TypeName) {}
func (annotator) AddAlias(tokens.ModuleName, tokens.TypeName) {}
func (annotator) SetResourceDeprecationMessage(_ string)      {}

func newServer(ctx context.Context, t *testing.T, client Client) integration.Server {
	t.Helper()

	p := NewBuildxProvider()

	// Inject a mock client if provided.
	if client != nil {
		p = mwcontext.Wrap(p, func(ctx context.Context) context.Context {
			return context.WithValue(ctx, _mockClientKey, client)
		})
	}

	s, err := integration.NewServer(
		ctx,
		"docker-build", semver.Version{Major: 0},
		integration.WithProvider(p),
	)
	require.NoError(t, err)
	return s
}
