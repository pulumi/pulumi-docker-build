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
	"path/filepath"
	"testing"

	clicfg "github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/types"
	mobyclient "github.com/moby/moby/client"
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
			Username: "resource-aws-user",
			Password: "resource-aws-password",
		},
		//nolint:gosec // G101: test fixture, not a real credential.
		{
			Address:  dockerIO,
			Username: "resource-dockerhub-user",
			Password: "resource-dockerhub-password",
		},
	}

	_, err = wrap(h, registries...)
	require.NoError(t, err)

	cli, err := wrap(h, registries...)
	require.NoError(t, err)

	expected := map[string]types.AuthConfig{
		//nolint:gosec // G101: test fixture, not a real credential.
		awsECRAddress: {
			Username:      "resource-aws-user",
			Password:      "resource-aws-password",
			ServerAddress: awsECRAddress,
		},
		//nolint:gosec // G101: test fixture, not a real credential.
		config.DockerRegistryAuth: {
			Username:      "resource-dockerhub-user",
			Password:      "resource-dockerhub-password",
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

// TestManifestCreateAuthSandbox guards the isolation regression from buildx
// v0.31+: NewRootCmd's PersistentPreRunE calls dockerCli.Initialize(), which
// reloads the host's Docker config and would discard wrap()'s sandboxed
// credentials before an `imagetools create` push authenticates. manifestCmd
// must re-apply the sandbox so the host's ambient credentials never leak.
func TestManifestCreateAuthSandbox(t *testing.T) {
	// Not parallel: mutates the process-global docker config dir to simulate a
	// host with ambient credentials.
	const hostRegistry = "https://host.example.com"

	// Seed a host docker config dir with an ambient credential.
	hostDir := t.TempDir()
	hostCfg := configfile.New(filepath.Join(hostDir, "config.json"))
	hostCfg.AuthConfigs = map[string]types.AuthConfig{
		//nolint:gosec // G101: test fixture, not a real credential.
		hostRegistry: {Username: "host-user", Password: "host-pass", ServerAddress: hostRegistry},
	}
	require.NoError(t, hostCfg.Save())

	// Point docker/cli's config dir at the host dir so dockerCli.Initialize()
	// reloads these credentials when the command runs.
	orig := clicfg.Dir()
	clicfg.SetDir(hostDir)
	t.Cleanup(func() { clicfg.SetDir(orig) })

	// wrap() a cli scoped to a different registry credential.
	h := &host{auths: map[string]types.AuthConfig{}}
	c, err := wrap(h, Registry{
		Address:  "scoped.example.com",
		Username: "scoped-user",
		//nolint:gosec // G101: test fixture, not a real credential.
		Password: "scoped-pass",
	})
	require.NoError(t, err)

	// Run the manifest command's PersistentPreRunE: Initialize() reloads the
	// host config, then the sandbox must be re-applied.
	cmd := c.manifestCmd([]string{"imagetools", "create", "--dry-run", "--tag", "example.com/x"})
	require.NoError(t, cmd.PersistentPreRunE(cmd, nil))

	got := c.ConfigFile().AuthConfigs
	assert.NotContains(t, got, hostRegistry, "host credential leaked through ManifestCreate")
	assert.Equal(t, c.auths, got, "config should hold only the scoped sandbox creds")
	assert.Nil(t, c.ConfigFile().CredentialHelpers)
	assert.Empty(t, c.ConfigFile().CredentialsStore)
}

// TestManifestCreateDoesNotLeakHostAuth is the end-to-end counterpart to
// TestManifestCreateAuthSandbox: it drives the real ManifestCreate against an
// on-disk $DOCKER_CONFIG holding a different credential for the same registry,
// proving the scoped credential still wins after Initialize() reloads. Needs a
// live Docker daemon (Ping gate).
//
//nolint:paralleltest // Mutates the process-global docker/cli config directory.
func TestManifestCreateDoesNotLeakHostAuth(t *testing.T) {
	// Ambient host config with a credential we did NOT scope this op to.
	tmp := t.TempDir()
	hostCfg := configfile.New(filepath.Join(tmp, "config.json"))
	hostCfg.AuthConfigs = map[string]types.AuthConfig{
		//nolint:gosec // G101: test fixture, not a real credential.
		awsECRAddress: {
			Username:      "ambient-host-user",
			Password:      "ambient-host-password",
			ServerAddress: awsECRAddress,
		},
	}
	require.NoError(t, hostCfg.Save())

	orig := clicfg.Dir()
	clicfg.SetDir(tmp)
	t.Cleanup(func() { clicfg.SetDir(orig) })

	h, err := newHost(context.Background(), nil)
	require.NoError(t, err)

	registries := []Registry{
		//nolint:gosec // G101: test fixture, not a real credential.
		{Address: awsECRAddress, Username: "scoped-user", Password: "scoped-password"},
	}
	c, err := wrap(h, registries...)
	require.NoError(t, err)

	if _, err := c.Client().Ping(context.Background(), mobyclient.PingOptions{}); err != nil {
		t.Skip(err)
	}

	before, err := c.ConfigFile().GetAuthConfig(awsECRAddress)
	require.NoError(t, err)
	assert.Equal(t, "scoped-user", before.Username)

	// The refs are unreachable, so this fails — but PersistentPreRunE (and
	// thus Initialize()) runs first, which is all the regression needs.
	_ = c.ManifestCreate(context.Background(), false, /* push */
		"127.0.0.1:1/target:latest", "127.0.0.1:1/source:latest")

	after, err := c.ConfigFile().GetAuthConfig(awsECRAddress)
	require.NoError(t, err)
	assert.Equal(t, "scoped-user", after.Username,
		"ManifestCreate must not let ambient host credentials replace wrap()'s sandboxed credentials")
}
