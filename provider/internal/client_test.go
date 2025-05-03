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
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	buildx "github.com/docker/buildx/build"
	"github.com/docker/buildx/builder"
	"github.com/docker/buildx/util/confutil"
	"github.com/docker/buildx/util/dockerutil"
	"github.com/docker/buildx/util/progress"
	"github.com/docker/docker/api/types/registry"
	"github.com/moby/buildkit/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuth(t *testing.T) {
	t.Parallel()
	user := "pulumibot"
	if u := os.Getenv("DOCKER_HUB_USER"); u != "" {
		user = u
	}
	password := os.Getenv("DOCKER_HUB_PASSWORD")
	address := "docker.io"

	cli := testcli(t, true, Registry{
		Address:  address,
		Username: user,
		Password: password,
	})

	_, err := cli.Client().
		RegistryLogin(context.Background(), registry.AuthConfig{ServerAddress: address})
	assert.NoError(t, err)
}

func TestCustomHost(t *testing.T) {
	socket := "unix:///foo/bar.sock"

	//nolint:paralleltest // not compatible with Setenv
	t.Run("env", func(t *testing.T) {
		t.Setenv("DOCKER_HOST", socket)

		h, err := newHost(context.Background(), nil)
		require.NoError(t, err)
		cli, err := wrap(h)
		require.NoError(t, err)

		assert.Equal(t, socket, cli.Client().DaemonHost())
		assert.Equal(t, socket, cli.DockerEndpoint().Host)
	})

	t.Run("config", func(t *testing.T) {
		t.Parallel()
		h, err := newHost(context.Background(), &Config{Host: socket})
		require.NoError(t, err)
		cli, err := wrap(h)
		require.NoError(t, err)

		assert.Equal(t, socket, cli.Client().DaemonHost())
		assert.Equal(t, socket, cli.DockerEndpoint().Host)
	})
}

func TestBuild(t *testing.T) {
	t.Parallel()

	tmpdir := t.TempDir()
	Max := Max

	exampleContext := &BuildContext{Context: Context{Location: "../../examples/app"}}

	tests := []struct {
		name string
		skip bool
		args ImageArgs

		auths []Registry
	}{
		{
			name: "multiPlatform",
			args: ImageArgs{
				Context: exampleContext,
				Dockerfile: &Dockerfile{
					Location: "../../examples/app/Dockerfile.multiPlatform",
				},
				Platforms: []Platform{"plan9/amd64", "plan9/arm64"},
			},
		},
		{
			name: "registryPush",
			skip: os.Getenv("DOCKER_HUB_PASSWORD") == "",
			args: ImageArgs{
				Context: exampleContext,
				Tags:    []string{"docker.io/pulumibot/buildkit-e2e:unit"},
				Push:    true,
			},
			auths: []Registry{{
				Address:  "docker.io",
				Username: "pulumibot",
				Password: os.Getenv("DOCKER_HUB_PASSWORD"),
			}},
		},
		{
			name: "cached",
			args: ImageArgs{
				Context: exampleContext,
				Tags:    []string{"cached"},
				CacheTo: []CacheTo{{Local: &CacheToLocal{
					Dest:          filepath.Join(tmpdir, "cache"),
					CacheWithMode: CacheWithMode{Mode: &Max},
				}}},
				CacheFrom: []CacheFrom{{Local: &CacheFromLocal{
					Src: filepath.Join(tmpdir, "cache"),
				}}},
			},
		},
		{
			name: "buildArgs",
			args: ImageArgs{
				Context: exampleContext,
				Dockerfile: &Dockerfile{
					Location: "../../examples/app/Dockerfile.buildArgs",
				},
				BuildArgs: map[string]string{
					"SET_ME_TO_TRUE": "true",
				},
			},
		},
		{
			name: "extraHosts",
			args: ImageArgs{
				Context: exampleContext,
				Dockerfile: &Dockerfile{
					Location: "../../examples/app/Dockerfile.extraHosts",
				},
				AddHosts: []string{
					"metadata.google.internal:169.254.169.254",
				},
			},
		},
		{
			name: "sshMount",
			skip: os.Getenv("SSH_AUTH_SOCK") == "",
			args: ImageArgs{
				Context: exampleContext,
				Dockerfile: &Dockerfile{
					Location: "../../examples/app/Dockerfile.sshMount",
				},
				SSH: []SSH{{ID: "default"}},
			},
		},
		{
			name: "secrets",
			args: ImageArgs{
				Context: exampleContext,
				Dockerfile: &Dockerfile{
					Location: "../../examples/app/Dockerfile.secrets",
				},
				Secrets: map[string]string{
					"password": "hunter2",
				},
				NoCache: true,
			},
		},
		{
			name: "labels",
			args: ImageArgs{
				Context: exampleContext,
				Labels: map[string]string{
					"description": "foo",
				},
			},
		},
		{
			name: "target",
			args: ImageArgs{
				Context: exampleContext,
				Dockerfile: &Dockerfile{
					Location: "../../examples/app/Dockerfile.target",
				},
				Target: "build-me",
			},
		},
		{
			name: "namedContext",
			args: ImageArgs{
				Context: &BuildContext{
					Context: Context{
						Location: "../../examples/app",
					},
					Named: NamedContexts{
						"golang:latest": Context{
							Location: "docker-image://golang@sha256:b8e62cf593cdaff36efd90aa3a37de268e6781a2e68c6610940c48f7cdf36984",
						},
					},
				},
				Dockerfile: &Dockerfile{
					Location: "../../examples/app/Dockerfile.namedContexts",
				},
			},
		},
		{
			name: "remoteContext",
			args: ImageArgs{
				Context: &BuildContext{
					Context: Context{
						Location: "https://raw.githubusercontent.com/pulumi/pulumi-docker/api-types/provider/testdata/Dockerfile",
					},
				},
			},
		},
		{
			name: "remoteContextWithInline",
			args: ImageArgs{
				Context: &BuildContext{
					Context: Context{
						Location: "https://github.com/docker-library/hello-world.git",
					},
				},
				Dockerfile: &Dockerfile{
					Inline: dedent(`
					FROM busybox
					COPY hello.c ./
					`),
				},
			},
		},
		{
			name: "inline",
			args: ImageArgs{
				Context: exampleContext,
				Dockerfile: &Dockerfile{
					Inline: dedent(`
					FROM alpine
					RUN echo 👍
					`),
				},
			},
		},
		{
			name: "dockerLoad",
			args: ImageArgs{
				Context: exampleContext,
				Load:    true,
			},
		},
	}

	// Add an exec: true version for all of our test cases.
	for _, tt := range tests {
		tt.name = "exec-" + tt.name
		tt.args.Exec = true
		tmpdir := filepath.Join(t.TempDir(), "exec")
		for _, c := range tt.args.CacheTo {
			if c.Local != nil {
				c.Local.Dest = tmpdir
			}
		}
		for _, c := range tt.args.CacheFrom {
			if c.Local != nil {
				c.Local.Src = tmpdir
			}
		}
		tests = append(tests, tt)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.skip {
				t.Skip()
			}
			ctx := context.Background()
			cli := testcli(t, true, tt.auths...)

			build, err := tt.args.toBuild(ctx, true, false)
			require.NoError(t, err)

			_, err = cli.Build(ctx, build)
			assert.NoError(t, err, cli.err.String())
		})
	}
}

func TestBuildkitEnabled(t *testing.T) {
	t.Parallel()
	cli := testcli(t, false)
	ok, err := cli.BuildKitEnabled()
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestInspect(t *testing.T) {
	t.Parallel()
	cli := testcli(t, false)
	descriptors, err := cli.Inspect(context.Background(), "pulumibot/myapp:buildx")
	require.NoError(t, err)
	assert.Equal(
		t,
		"application/vnd.docker.distribution.manifest.v2+json",
		descriptors[0].MediaType,
	)
}

func TestNormalizeReference(t *testing.T) {
	t.Parallel()
	tests := []struct {
		ref     string
		want    string
		wantErr string
	}{
		{
			ref:  "foo",
			want: "docker.io/library/foo:latest",
		},
		{
			ref:  "pulumi/pulumi:v3.100.0",
			want: "docker.io/pulumi/pulumi:v3.100.0",
		},
		{
			ref:     "invalid:ref:format",
			wantErr: "invalid reference format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.ref, func(t *testing.T) {
			t.Parallel()
			ref, err := normalizeReference(tt.ref)
			if err != nil {
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.Equal(t, ref.String(), tt.want)
			}
		})
	}
}

//nolint:paralleltest // Overrides default logger.
func TestBuildError(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("flaky on CI for some reason")
	}

	l := slog.Default()
	defer slog.SetDefault(l)

	// Override go-provider's default logger to capture and tee to stdout.
	logger := &bytes.Buffer{}
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(io.MultiWriter(logger, os.Stdout), nil),
		),
	)

	exampleContext := &BuildContext{Context: Context{Location: "../../examples/app"}}

	args := ImageArgs{
		Context: exampleContext,
		Dockerfile: &Dockerfile{
			Inline: "FROM alpine\nRUN echo hello\nRUN badcmd",
		},
	}

	ctx := context.Background()
	cli := testcli(t, true)

	build, err := args.toBuild(ctx, true, false)
	require.NoError(t, err)

	_, err = cli.Build(ctx, build)
	assert.Error(t, err)

	want := []string{
		`RUN echo hello`,
		`/bin/sh: badcmd: not found`,
	}

	for _, want := range want {
		assert.Contains(t, logger.String(), want)
	}
	assert.ErrorContains(t, err,
		`process "/bin/sh -c badcmd" did not complete successfully: exit code: 127`,
	)
}

func TestBuildExecError(t *testing.T) {
	t.Parallel()

	exampleContext := &BuildContext{Context: Context{Location: "../../examples/app"}}

	args := ImageArgs{
		Context: exampleContext,
		Dockerfile: &Dockerfile{
			Inline: "FROM alpine\nRUN echo hello\nRUN badcmd",
		},
		Exec: true,
	}

	ctx := context.Background()
	cli := testcli(t, true)

	build, err := args.toBuild(ctx, true, false)
	require.NoError(t, err)

	_, err = cli.Build(ctx, build)
	assert.Error(t, err)

	want := []string{
		`RUN echo hello`,
		`/bin/sh: badcmd: not found`,
		`process "/bin/sh -c badcmd" did not complete successfully: exit code: 127`,
	}

	for _, want := range want {
		assert.Contains(t, cli.err.String(), want)
	}
}

func TestBuildCancelation(t *testing.T) {
	t.Parallel()
	cli := testcli(t, true)

	ctrl := gomock.NewController(t)

	ctx, cancel := context.WithCancel(context.Background())

	b := NewMockBuilder(ctrl)
	b.EXPECT().Build(
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
	).DoAndReturn(func(_ context.Context, _ []builder.Node, _ map[string]buildx.Options, _ *dockerutil.Client, _ *confutil.Config, _ progress.Writer) (map[string]*client.SolveResponse, error) {
		cancel()
		return nil, fmt.Errorf("cancel wasn't respected")
	})
	cli.builder = b

	resp, err := cli.Build(ctx, &build{})
	assert.ErrorIs(t, err, context.Canceled)
	assert.Nil(t, resp)
}

// testcli returns a new standalone CLI instance. Set ping to true if a live
// daemon is required -- the test will be skipped if the daemon is not available.
func testcli(t *testing.T, ping bool, auths ...Registry) *cli {
	h, err := newHost(context.Background(), nil)
	require.NoError(t, err)

	cli, err := wrap(h, auths...)
	require.NoError(t, err)

	if ping {
		_, err := cli.Client().Ping(context.Background())
		if err != nil {
			t.Skip(err)
		}
	}

	return cli
}
