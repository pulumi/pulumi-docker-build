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
	"io"
	"testing"

	"github.com/docker/cli/cli/config/types"
	"github.com/regclient/regclient/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	t.Parallel()

	h, err := newHost(nil)
	require.NoError(t, err)
	cli, err := wrap(h)
	require.NoError(t, err)

	err = cli.exec([]string{"buildx", "version"}, nil)
	assert.NoError(t, err)

	out, err := io.ReadAll(cli.r)
	require.NoError(t, err)
	assert.Contains(t, string(out), "github.com/docker/buildx")
}

func TestWrappedAuth(t *testing.T) {
	ecr := "https://1234.dkr.ecr.us-west-2.amazonaws.com"
	h := &host{
		auths: map[string]types.AuthConfig{
			ecr: {
				Username:      "host-aws-user",
				Password:      "host-aws-password",
				ServerAddress: ecr,
			},
			"misc": {
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

	cli, err := wrap(h, registries...)
	require.NoError(t, err)

	require.Contains(t, cli.auths, ecr)
	aws := cli.auths[ecr]
	assert.Equal(t, "resource-aws-user", aws.Username)
	assert.Equal(t, "resource-aws-password", aws.Password)
	assert.Equal(t, "1234.dkr.ecr.us-west-2.amazonaws.com", aws.ServerAddress)

	require.Contains(t, cli.auths, config.DockerRegistryAuth)
	dockerhub := cli.auths[config.DockerRegistryAuth]
	assert.Equal(t, "resource-dockerhub-user", dockerhub.Username)
	assert.Equal(t, "resource-dockerhub-password", dockerhub.Password)
	assert.Equal(t, config.DockerRegistryDNS, dockerhub.ServerAddress)

	// Auths derived from the host should be untouched, e.g. no scheme added, etc.
	require.Contains(t, cli.auths, "misc")
	misc := cli.auths["misc"]
	assert.Equal(t, "host-misc-user", misc.Username)
	assert.Equal(t, "host-misc-password", misc.Password)
	assert.Equal(t, "misc", misc.ServerAddress)
}
