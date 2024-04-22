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
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/api/types/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
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

		h, err := newHost(nil)
		require.NoError(t, err)
		cli, err := wrap(h)
		require.NoError(t, err)

		assert.Equal(t, socket, cli.Client().DaemonHost())
		assert.Equal(t, socket, cli.DockerEndpoint().Host)
	})

	t.Run("config", func(t *testing.T) {
		t.Parallel()
		h, err := newHost(&Config{Host: socket})
		require.NoError(t, err)
		cli, err := wrap(h)
		require.NoError(t, err)

		assert.Equal(t, socket, cli.Client().DaemonHost())
		assert.Equal(t, socket, cli.DockerEndpoint().Host)
	})
}

func TestBuild(t *testing.T) {
	t.Parallel()
	// Workaround for https://github.com/pulumi/pulumi-go-provider/issues/159
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	pctx := NewMockProviderContext(ctrl)
	pctx.EXPECT().Log(gomock.Any(), gomock.Any()).AnyTimes()
	pctx.EXPECT().LogStatus(gomock.Any(), gomock.Any()).AnyTimes()
	pctx.EXPECT().Done().Return(ctx.Done()).AnyTimes()
	pctx.EXPECT().
		Value(gomock.Any()).
		DoAndReturn(func(key any) any { return ctx.Value(key) }).
		AnyTimes()
	pctx.EXPECT().Err().Return(ctx.Err()).AnyTimes()
	pctx.EXPECT().Deadline().Return(ctx.Deadline()).AnyTimes()

	tmpdir := t.TempDir()
	max := Max

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
					CacheWithMode: CacheWithMode{Mode: &max},
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
		tt := tt
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.skip {
				t.Skip()
			}
			cli := testcli(t, true, tt.auths...)

			build, err := tt.args.toBuild(pctx, false)
			require.NoError(t, err)

			_, err = cli.Build(pctx, build)
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
		tt := tt
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

func TestBuildError(t *testing.T) {
	t.Parallel()

	if os.Getenv("CI") != "" {
		t.Skip("flaky on CI for some reason")
	}

	ctrl, ctx := gomock.WithContext(context.Background(), t)

	exampleContext := &BuildContext{Context: Context{Location: "../../examples/app"}}

	args := ImageArgs{
		Context: exampleContext,
		Dockerfile: &Dockerfile{
			Inline: "FROM alpine\nRUN echo hello\nRUN badcmd",
		},
	}
	logged := bytes.Buffer{}

	pctx := NewMockProviderContext(ctrl)
	pctx.EXPECT().Done().Return(ctx.Done()).AnyTimes()
	pctx.EXPECT().
		Value(gomock.Any()).
		DoAndReturn(func(key any) any { return ctx.Value(key) }).
		AnyTimes()
	pctx.EXPECT().Err().Return(ctx.Err()).AnyTimes()
	pctx.EXPECT().Deadline().Return(ctx.Deadline()).AnyTimes()

	pctx.EXPECT().LogStatus(gomock.Any(), gomock.Any()).AnyTimes()
	pctx.EXPECT().Log(gomock.Any(), gomock.Any()).DoAndReturn(func(_ diag.Severity, msg string) {
		logged.WriteString(msg)
	}).AnyTimes()

	cli := testcli(t, true)

	build, err := args.toBuild(pctx, false)
	require.NoError(t, err)

	_, err = cli.Build(pctx, build)
	assert.Error(t, err)

	want := []string{
		`RUN echo hello`,
		`/bin/sh: badcmd: not found`,
	}

	for _, want := range want {
		assert.Contains(t, logged.String(), want)
	}
	assert.ErrorContains(t, err,
		`process "/bin/sh -c badcmd" did not complete successfully: exit code: 127`,
	)
}

func TestBuildExecError(t *testing.T) {
	t.Parallel()
	ctrl, _ := gomock.WithContext(context.Background(), t)

	exampleContext := &BuildContext{Context: Context{Location: "../../examples/app"}}

	args := ImageArgs{
		Context: exampleContext,
		Dockerfile: &Dockerfile{
			Inline: "FROM alpine\nRUN echo hello\nRUN badcmd",
		},
		Exec: true,
	}

	pctx := NewMockProviderContext(ctrl)
	pctx.EXPECT().Log(
		diag.Warning,
		"No exports were specified so the build will only remain in the local build cache. "+
			"Use `push` to upload the image to a registry, or silence this warning with a `cacheonly` export.",
	)
	pctx.EXPECT().LogStatus(gomock.Any(), gomock.Any()).AnyTimes()

	cli := testcli(t, true)

	build, err := args.toBuild(pctx, false)
	require.NoError(t, err)

	_, err = cli.Build(pctx, build)
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

// testcli returns a new standalone CLI instance. Set ping to true if a live
// daemon is required -- the test will be skipped if the daemon is not available.
func testcli(t *testing.T, ping bool, auths ...Registry) *cli {
	h, err := newHost(nil)
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
