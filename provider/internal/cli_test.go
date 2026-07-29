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

	dockerconfig "github.com/docker/cli/cli/config"
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

// TestManifestCreateDoesNotLeakHostAuth is a regression test for a
// credential-isolation break in buildx v0.31+: commands.NewRootCmd's
// non-plugin PersistentPreRunE now calls dockerCli.Initialize(), which
// unconditionally reloads $DOCKER_CONFIG/config.json into cli.configFile,
// discarding wrap()'s in-memory, scoped AuthConfigs. Once that happens,
// buildx's auth provider (storeutil.GetImageConfig ->
// dockerconfig.LoadAuthConfig) reads whatever's left in ConfigFile(), so an
// Index resource's ManifestCreate call can silently authenticate with
// whatever's in the ambient host Docker config instead of the credentials
// the Index resource was scoped to.
//
// This mirrors TestWrappedAuth, but populates an actual on-disk
// $DOCKER_CONFIG (like a real `docker login`'d host or credential helper
// would) with credentials that differ from the ones we scope the operation
// to, then exercises ManifestCreate -- rather than just asserting on
// wrap()'s output -- since ManifestCreate is what triggers the second,
// leaking Initialize() call.
//
//nolint:paralleltest // Mutates the process-global docker/cli config directory.
func TestManifestCreateDoesNotLeakHostAuth(t *testing.T) {
	// Simulate an ambient host Docker config, e.g. from `docker login` or an
	// ECR/GCR credential helper, with different credentials than what we
	// scope this operation to.
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

	// Point $DOCKER_CONFIG at our fake ambient host config for the duration
	// of this test, and restore it afterwards.
	orig := dockerconfig.Dir()
	dockerconfig.SetDir(tmp)
	t.Cleanup(func() { dockerconfig.SetDir(orig) })

	h, err := newHost(context.Background(), nil)
	require.NoError(t, err)

	registries := []Registry{
		//nolint:gosec // G101: test fixture, not a real credential.
		{
			Address:  awsECRAddress,
			Username: "scoped-user",
			Password: "scoped-password",
		},
	}

	c, err := wrap(h, registries...)
	require.NoError(t, err)

	if _, err := c.Client().Ping(context.Background(), mobyclient.PingOptions{}); err != nil {
		t.Skip(err)
	}

	// Sanity check: wrap() scoped the credential as expected, same as
	// TestWrappedAuth.
	before, err := c.ConfigFile().GetAuthConfig(awsECRAddress)
	require.NoError(t, err)
	assert.Equal(t, "scoped-user", before.Username)

	// ManifestCreate will ultimately fail because these refs don't exist,
	// but commands.NewRootCmd's PersistentPreRunE (and thus
	// dockerCli.Initialize()) always runs first regardless of that later
	// failure -- so we don't need the refs to resolve for the regression to
	// manifest.
	_ = c.ManifestCreate(context.Background(), false, /* push */
		"127.0.0.1:1/target:latest", "127.0.0.1:1/source:latest")

	after, err := c.ConfigFile().GetAuthConfig(awsECRAddress)
	require.NoError(t, err)

	// This is the crux of the regression: ManifestCreate must not allow the
	// ambient host's config to replace the credentials wrap() scoped this
	// operation to.
	assert.Equal(t, "scoped-user", after.Username,
		"ManifestCreate must not let ambient host credentials silently replace wrap()'s sandboxed credentials")
}
