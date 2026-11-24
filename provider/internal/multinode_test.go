// Copyright 2026, Pulumi Corporation.
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
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/docker/buildx/builder"
	"github.com/docker/buildx/util/imagetools"
	"github.com/docker/buildx/util/resolver"
	clicfg "github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/types"
	"github.com/moby/buildkit/session/auth/authprovider"
	mobyclient "github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

const scopedUser = "scoped-user"

// dockerCmd runs the docker CLI with DOCKER_CONFIG pointed at dir and returns
// its trimmed combined output.
func dockerCmd(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.CommandContext(t.Context(), "docker", args...)
	cmd.Env = append(os.Environ(), "DOCKER_CONFIG="+dir)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "docker %s: %s", strings.Join(args, " "), out)
	return strings.TrimSpace(string(out))
}

// TestBuildMultiNodeUsesRegistriesAuth builds a multi-platform image on a
// builder with one node per platform and pushes it to a registry that
// requires authentication. The host's Docker config has no credentials: they
// are only provided through the resource's registries. Each node pushes its
// platform's image and buildx then merges them into a manifest list in the
// provider's process, so the merge must authenticate with the resource's
// registries rather than with the host's Docker CLI.
//
// Needs a Docker daemon and the registry:2 and moby/buildkit images.
//
//nolint:paralleltest // Mutates the process-global docker/cli config directory.
func TestBuildMultiNodeUsesRegistriesAuth(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow test requiring a Docker daemon")
	}

	// A host config dir with no credentials. It also holds the buildx store,
	// so the builder we create below is isolated from the real host.
	hostDir := t.TempDir()
	orig := clicfg.Dir()
	clicfg.SetDir(hostDir)
	t.Cleanup(func() { clicfg.SetDir(orig) })
	require.NoError(t, configfile.New(filepath.Join(hostDir, "config.json")).Save())

	h, err := newHost(t.Context(), nil)
	require.NoError(t, err)
	pinger, err := wrap(h)
	require.NoError(t, err)
	if _, err := pinger.Client().Ping(t.Context(), mobyclient.PingOptions{}); err != nil {
		t.Skip(err)
	}

	// An authenticated registry on a random localhost port.
	const user, password = "multinode-user", "multinode-password" //nolint:gosec // test fixture.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	require.NoError(t, err)
	htpasswd := filepath.Join(t.TempDir(), "htpasswd")
	require.NoError(t, os.WriteFile(htpasswd, []byte(user+":"+string(hash)+"\n"), 0o600))

	registryID := dockerCmd(t, hostDir, "run", "-d", "--rm",
		"-p", "127.0.0.1::5000",
		"-v", htpasswd+":/auth/htpasswd:ro",
		"-e", "REGISTRY_AUTH=htpasswd",
		"-e", "REGISTRY_AUTH_HTPASSWD_REALM=test",
		"-e", "REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd",
		"registry:2",
	)
	t.Cleanup(func() {
		_ = exec.Command("docker", "rm", "-f", registryID).Run()
	})
	hostPort := dockerCmd(t, hostDir, "port", registryID, "5000/tcp")
	address := "localhost:" + hostPort[strings.LastIndex(hostPort, ":")+1:]
	require.Eventually(t, func() bool {
		resp, err := http.Get("http://" + address + "/v2/") //nolint:noctx,gosec // test helper.
		if err != nil {
			return false
		}
		_ = resp.Body.Close()
		return resp.StatusCode == http.StatusUnauthorized
	}, 30*time.Second, 200*time.Millisecond, "registry did not start")

	// A builder with one node per platform. buildx rejects two nodes with the
	// same endpoint, so each node gets its own Docker context pointing at the
	// same daemon. The nodes share the host's network so they can push to the
	// registry on localhost.
	daemon := dockerCmd(t, hostDir, "context", "inspect", "default", "--format", "{{.Endpoints.docker.Host}}")
	const builderName = "multinode"
	for i, platform := range []string{"linux/amd64", "linux/arm64"} {
		node := strings.TrimPrefix(platform, "linux/")
		dockerCmd(t, hostDir, "context", "create", node, "--docker", "host="+daemon)
		args := []string{
			"buildx", "create", "--name", builderName,
			"--driver", "docker-container", "--driver-opt", "network=host",
			"--platform", platform, "--node", node,
		}
		if i > 0 {
			args = append(args, "--append")
		}
		dockerCmd(t, hostDir, append(args, node)...)
	}
	t.Cleanup(func() {
		cmd := exec.Command("docker", "buildx", "rm", "--force", builderName)
		cmd.Env = append(os.Environ(), "DOCKER_CONFIG="+hostDir)
		_ = cmd.Run()
	})

	contextDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(contextDir, "Dockerfile"),
		[]byte("FROM scratch\nCOPY hello /hello\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(contextDir, "hello"), []byte("hello"), 0o600))

	args := ImageArgs{
		Context:   &BuildContext{Context: Context{Location: contextDir}},
		Builder:   &BuilderConfig{Name: builderName},
		Platforms: []Platform{"linux/amd64", "linux/arm64"},
		Tags:      []string{address + "/multinode:test"},
		Push:      true,
	}

	c, err := wrap(h, Registry{Address: address, Username: user, Password: password})
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Close() })

	build, err := args.toBuild(t.Context(), true, false)
	require.NoError(t, err)

	_, err = c.Build(t.Context(), build)
	require.NoError(t, err, c.err.String())
}

// TestNodesWithAuth confirms that the nodes we hand to buildx authenticate
// with the wrapped CLI's config (the host's credentials plus the resource's
// registries) and keep the rest of their image options.
func TestNodesWithAuth(t *testing.T) {
	t.Parallel()

	const address = "registry.example.com"
	h := &host{auths: map[string]types.AuthConfig{}}
	c, err := wrap(h, Registry{
		Address:  address,
		Username: scopedUser,
		Password: "scoped-password", //nolint:gosec // test fixture.
	})
	require.NoError(t, err)

	insecure := true
	registryConfig := map[string]resolver.RegistryConfig{address: {Insecure: &insecure}}
	hostAuth := authprovider.LoadAuthConfig(configfile.New(""))
	cached := []builder.Node{
		{Builder: "a", ImageOpt: imagetools.Opt{Auth: hostAuth, RegistryConfig: registryConfig}},
		{Builder: "b", ImageOpt: imagetools.Opt{Auth: hostAuth}},
	}

	nodes := c.nodesWithAuth(cached)
	require.Len(t, nodes, len(cached))
	for i, n := range nodes {
		assert.Equal(t, cached[i].Builder, n.Builder)
		require.NotNil(t, n.ImageOpt.Auth)
		got, err := n.ImageOpt.Auth(t.Context(), address, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, scopedUser, got.Username)
	}
	assert.Equal(t, registryConfig, nodes[0].ImageOpt.RegistryConfig)

	// The cached nodes are shared between resources, so they must not change.
	for _, n := range cached {
		got, err := n.ImageOpt.Auth(t.Context(), address, nil, nil)
		require.NoError(t, err)
		assert.Empty(t, got.Username)
	}
}
