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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDockerfile(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		d       Dockerfile
		givenC  Context
		preview bool

		wantErr string
	}{
		{
			name: "relative",
			d: Dockerfile{
				Location: "../internal/../internal/testdata/noop/Dockerfile",
			},
		},
		{
			name: "missing file",
			d: Dockerfile{
				Location: "/does/not/exist/Dockerfile",
			},
			wantErr: "no such file",
		},
		{
			name: "invalid syntax",
			d: Dockerfile{
				Location: "testdata/Dockerfile.invalid",
			},
			wantErr: "unknown instruction: RUNN",
		},
		{
			name: "invalid syntax inline",
			d: Dockerfile{
				Inline: "RUNN it",
			},
			wantErr: "unknown instruction: RUNN",
		},
		{
			name: "valid syntax inline",
			d: Dockerfile{
				Inline: "FROM scratch",
			},
		},
		{
			name: "valid custom syntax inline",
			d: Dockerfile{
				Inline: `# syntax=docker.io/docker/dockerfile:1.7-labs
FROM public.ecr.aws/docker/library/node:22-alpine AS base

WORKDIR /app
COPY --parents ./package.json ./package-lock.json ./apps/*/package.json ./packages/*/package.json ./
`,
			},
		},
		{
			name:    "unset",
			d:       Dockerfile{},
			wantErr: "missing 'location' or 'inline'",
		},
		{
			name:   "unset with remote context",
			d:      Dockerfile{},
			givenC: Context{Location: "https://github.com/foobar"},
		},
		{
			name:    "preview",
			d:       Dockerfile{},
			preview: true,
		},
		{
			name:    "over-specified",
			d:       Dockerfile{Location: ".", Inline: "FROM scratch"},
			wantErr: `only specify "file" or "inline", not both`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.d.validate(tt.preview, &tt.givenC)

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.wantErr)
			}
		})
	}
}
