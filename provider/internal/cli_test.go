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
	"io"
	"testing"

	"github.com/docker/cli/cli/config/types"
	"github.com/regclient/regclient/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	t.Parallel()

	h, err := newHost(t.Context(), nil)
	require.NoError(t, err)
	cli, err := wrap(h)
	require.NoError(t, err)

	err = cli.exec(t.Context(), []string{"buildx", "version"}, nil)
	assert.NoError(t, err)

	out, err := io.ReadAll(cli.r)
	require.NoError(t, err)
	assert.Contains(t, string(out), "github.com/docker/buildx")
}

func TestWrappedAuth(t *testing.T) {
	t.Parallel()
	ecr := "https://1234.dkr.ecr.us-west-2.amazonaws.com"

	realhost, err := newHost(context.Background(), nil)
	require.NoError(t, err)

	h := &host{
		auths: map[string]types.AuthConfig{
			ecr: {
				Username:      "host-aws-user",
				Password:      "host-aws-password",
				ServerAddress: ecr,
			},
			"https://misc": { // Legacy config includes http/https scheme.
				Username:      "host-misc-user",
				Password:      "host-misc-password",
				ServerAddress: "misc",
			},
		},
	}

	registries := []Registry{
		{
			Address:  "1234.dkr.ecr.us-west-2.amazonaws.com",
			Username: "resource-aws-user",
			Password: "resource-aws-password",
		},
		{
			Address:  "docker.io",
			Username: "resource-dockerhub-user",
			Password: "resource-dockerhub-password",
		},
	}

	_, err = wrap(h, registries...)
	require.NoError(t, err)

	cli, err := wrap(h, registries...)
	require.NoError(t, err)

	expected := map[string]types.AuthConfig{
		"1234.dkr.ecr.us-west-2.amazonaws.com": {
			Username:      "resource-aws-user",
			Password:      "resource-aws-password",
			ServerAddress: "1234.dkr.ecr.us-west-2.amazonaws.com",
		},
		config.DockerRegistryAuth: {
			Username:      "resource-dockerhub-user",
			Password:      "resource-dockerhub-password",
			ServerAddress: config.DockerRegistryDNS,
		},
		"misc": {
			Username:      "host-misc-user",
			Password:      "host-misc-password",
			ServerAddress: "misc",
		},
	}
	assert.Equal(t, expected, cli.auths)
	assert.Len(t, h.auths, 2) // In-memory host auth is unchanged.

	// Assert that our on-disk host's auth is untouched.
	realhostRefreshed, err := newHost(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, realhost.auths, realhostRefreshed.auths)
}
