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
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/regclient/regclient/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:gosec // G101: test fixtures, not real credentials.
const (
	resourceAWSUser           = "resource-aws-user"
	resourceAWSPassword       = "resource-aws-password"
	resourceDockerHubUser     = "resource-dockerhub-user"
	resourceDockerHubPassword = "resource-dockerhub-password"
)

func TestExec(t *testing.T) {
	t.Parallel()

	h, err := newHost(t.Context(), nil)
	require.NoError(t, err)
	cli, err := wrap(h)
	require.NoError(t, err)

	err = cli.exec(t.Context(), []string{buildxName, "version"}, nil)
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
			//nolint:gosec // G101: test fixture, not a real credential.
			ecr: {
				Username:      "host-aws-user",
				Password:      "host-aws-password",
				ServerAddress: ecr,
			},
			// Legacy config includes http/https scheme.
			//nolint:gosec // G101: test fixture, not a real credential.
			"https://misc": {
				Username:      "host-misc-user",
				Password:      "host-misc-password",
				ServerAddress: miscName,
			},
		},
	}

	registries := []Registry{
		//nolint:gosec // G101: test fixture, not a real credential.
		{
			Address:  awsECRAddress,
			Username: resourceAWSUser,
			Password: resourceAWSPassword,
		},
		//nolint:gosec // G101: test fixture, not a real credential.
		{
			Address:  dockerIO,
			Username: resourceDockerHubUser,
			Password: resourceDockerHubPassword,
		},
	}

	_, err = wrap(h, registries...)
	require.NoError(t, err)

	cli, err := wrap(h, registries...)
	require.NoError(t, err)

	expected := map[string]types.AuthConfig{
		//nolint:gosec // G101: test fixture, not a real credential.
		awsECRAddress: {
			Username:      resourceAWSUser,
			Password:      resourceAWSPassword,
			ServerAddress: awsECRAddress,
		},
		//nolint:gosec // G101: test fixture, not a real credential.
		config.DockerRegistryAuth: {
			Username:      resourceDockerHubUser,
			Password:      resourceDockerHubPassword,
			ServerAddress: config.DockerRegistryDNS,
		},
		//nolint:gosec // G101: test fixture, not a real credential.
		miscName: {
			Username:      "host-misc-user",
			Password:      "host-misc-password",
			ServerAddress: miscName,
		},
	}
	assert.Equal(t, expected, cli.auths)
	assert.Len(t, h.auths, 2) // In-memory host auth is unchanged.

	// Assert that our on-disk host's auth is untouched.
	realhostRefreshed, err := newHost(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, realhost.auths, realhostRefreshed.auths)
}

// TestSandboxedAuthProvider ensures the AuthConfigProvider we hand to buildkit
// resolves credentials solely from our in-memory (sandboxed) config, including
// DockerHub's configfile-key mapping. This guards the isolation wrap()
// documents: unlike buildx's dockerconfig.LoadAuthConfig, it must not consult
// ambient on-disk buildx config.
func TestSandboxedAuthProvider(t *testing.T) {
	t.Parallel()

	registries := []Registry{
		//nolint:gosec // G101: test fixture, not a real credential.
		{
			Address:  awsECRAddress,
			Username: resourceAWSUser,
			Password: resourceAWSPassword,
		},
		//nolint:gosec // G101: test fixture, not a real credential.
		{
			Address:  dockerIO,
			Username: resourceDockerHubUser,
			Password: resourceDockerHubPassword,
		},
	}

	c, err := wrap(&host{auths: map[string]types.AuthConfig{}}, registries...)
	require.NoError(t, err)

	provide := sandboxedAuthProvider(c)

	// DockerHub resolves via the configfile-key mapping.
	hub, err := provide(context.Background(), authprovider.DockerHubRegistryHost, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, resourceDockerHubUser, hub.Username)
	assert.Equal(t, resourceDockerHubPassword, hub.Password)

	// A scoped registry resolves to its in-memory credential.
	ecr, err := provide(context.Background(), awsECRAddress, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, resourceAWSUser, ecr.Username)
	assert.Equal(t, resourceAWSPassword, ecr.Password)

	// An unknown host yields no credential (nothing leaks in from elsewhere).
	unknown, err := provide(context.Background(), "unknown.example.com", nil, nil)
	require.NoError(t, err)
	assert.Empty(t, unknown.Username)
	assert.Empty(t, unknown.Password)
}
