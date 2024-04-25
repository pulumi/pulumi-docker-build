package main

import (
	"github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		dockerHubPassword := cfg.Require("dockerHubPassword")
		multiPlatform, err := dockerbuild.NewImage(ctx, "multiPlatform", &dockerbuild.ImageArgs{
			Push: pulumi.Bool(false),
			Dockerfile: &dockerbuild.DockerfileArgs{
				Location: pulumi.String("./app/Dockerfile.multiPlatform"),
			},
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("./app"),
			},
			Platforms: dockerbuild.PlatformArray{
				dockerbuild.Platform_Plan9_amd64,
				dockerbuild.Platform_Plan9_386,
			},
		})
		if err != nil {
			return err
		}
		_, err = dockerbuild.NewImage(ctx, "registryPush", &dockerbuild.ImageArgs{
			Push: pulumi.Bool(false),
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("./app"),
			},
			Tags: pulumi.StringArray{
				pulumi.String("docker.io/pulumibot/buildkit-e2e:example"),
			},
			Exports: dockerbuild.ExportArray{
				&dockerbuild.ExportArgs{
					Registry: &dockerbuild.ExportRegistryArgs{
						OciMediaTypes: pulumi.Bool(true),
						Push:          pulumi.Bool(false),
					},
				},
			},
			Registries: dockerbuild.RegistryArray{
				&dockerbuild.RegistryArgs{
					Address:  pulumi.String("docker.io"),
					Username: pulumi.String("pulumibot"),
					Password: pulumi.String(dockerHubPassword),
				},
			},
		})
		if err != nil {
			return err
		}
		_, err = dockerbuild.NewImage(ctx, "cached", &dockerbuild.ImageArgs{
			Push: pulumi.Bool(false),
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("./app"),
			},
			CacheTo: dockerbuild.CacheToArray{
				&dockerbuild.CacheToArgs{
					Local: &dockerbuild.CacheToLocalArgs{
						Dest: pulumi.String("tmp/cache"),
						Mode: dockerbuild.CacheModeMax,
					},
				},
			},
			CacheFrom: dockerbuild.CacheFromArray{
				&dockerbuild.CacheFromArgs{
					Local: &dockerbuild.CacheFromLocalArgs{
						Src: pulumi.String("tmp/cache"),
					},
				},
			},
		})
		if err != nil {
			return err
		}
		_, err = dockerbuild.NewImage(ctx, "buildArgs", &dockerbuild.ImageArgs{
			Push: pulumi.Bool(false),
			Dockerfile: &dockerbuild.DockerfileArgs{
				Location: pulumi.String("./app/Dockerfile.buildArgs"),
			},
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("./app"),
			},
			BuildArgs: pulumi.StringMap{
				"SET_ME_TO_TRUE": pulumi.String("true"),
			},
		})
		if err != nil {
			return err
		}
		_, err = dockerbuild.NewImage(ctx, "extraHosts", &dockerbuild.ImageArgs{
			Push: pulumi.Bool(false),
			Dockerfile: &dockerbuild.DockerfileArgs{
				Location: pulumi.String("./app/Dockerfile.extraHosts"),
			},
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("./app"),
			},
			AddHosts: pulumi.StringArray{
				pulumi.String("metadata.google.internal:169.254.169.254"),
			},
		})
		if err != nil {
			return err
		}
		_, err = dockerbuild.NewImage(ctx, "sshMount", &dockerbuild.ImageArgs{
			Push: pulumi.Bool(false),
			Dockerfile: &dockerbuild.DockerfileArgs{
				Location: pulumi.String("./app/Dockerfile.sshMount"),
			},
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("./app"),
			},
			Ssh: dockerbuild.SSHArray{
				&dockerbuild.SSHArgs{
					Id: pulumi.String("default"),
				},
			},
		})
		if err != nil {
			return err
		}
		_, err = dockerbuild.NewImage(ctx, "secrets", &dockerbuild.ImageArgs{
			Push: pulumi.Bool(false),
			Dockerfile: &dockerbuild.DockerfileArgs{
				Location: pulumi.String("./app/Dockerfile.secrets"),
			},
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("./app"),
			},
			Secrets: pulumi.StringMap{
				"password": pulumi.String("hunter2"),
			},
		})
		if err != nil {
			return err
		}
		_, err = dockerbuild.NewImage(ctx, "labels", &dockerbuild.ImageArgs{
			Push: pulumi.Bool(false),
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("./app"),
			},
			Labels: pulumi.StringMap{
				"description": pulumi.String("This image will get a descriptive label 👍"),
			},
		})
		if err != nil {
			return err
		}
		_, err = dockerbuild.NewImage(ctx, "target", &dockerbuild.ImageArgs{
			Push: pulumi.Bool(false),
			Dockerfile: &dockerbuild.DockerfileArgs{
				Location: pulumi.String("./app/Dockerfile.target"),
			},
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("./app"),
			},
			Target: pulumi.String("build-me"),
		})
		if err != nil {
			return err
		}
		_, err = dockerbuild.NewImage(ctx, "namedContexts", &dockerbuild.ImageArgs{
			Push: pulumi.Bool(false),
			Dockerfile: &dockerbuild.DockerfileArgs{
				Location: pulumi.String("./app/Dockerfile.namedContexts"),
			},
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("./app"),
				Named: dockerbuild.ContextMap{
					"golang:latest": &dockerbuild.ContextArgs{
						Location: pulumi.String("docker-image://golang@sha256:b8e62cf593cdaff36efd90aa3a37de268e6781a2e68c6610940c48f7cdf36984"),
					},
				},
			},
		})
		if err != nil {
			return err
		}
		_, err = dockerbuild.NewImage(ctx, "remoteContext", &dockerbuild.ImageArgs{
			Push: pulumi.Bool(false),
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("https://raw.githubusercontent.com/pulumi/pulumi-docker/api-types/provider/testdata/Dockerfile"),
			},
		})
		if err != nil {
			return err
		}
		_, err = dockerbuild.NewImage(ctx, "remoteContextWithInline", &dockerbuild.ImageArgs{
			Push: pulumi.Bool(false),
			Dockerfile: &dockerbuild.DockerfileArgs{
				Inline: pulumi.String("FROM busybox\nCOPY hello.c ./\n"),
			},
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("https://github.com/docker-library/hello-world.git"),
			},
		})
		if err != nil {
			return err
		}
		_, err = dockerbuild.NewImage(ctx, "inline", &dockerbuild.ImageArgs{
			Push: pulumi.Bool(false),
			Dockerfile: &dockerbuild.DockerfileArgs{
				Inline: pulumi.String("FROM alpine\nRUN echo \"This uses an inline Dockerfile! 👍\"\n"),
			},
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("./app"),
			},
		})
		if err != nil {
			return err
		}
		_, err = dockerbuild.NewImage(ctx, "dockerLoad", &dockerbuild.ImageArgs{
			Push: pulumi.Bool(false),
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("./app"),
			},
			Exports: dockerbuild.ExportArray{
				&dockerbuild.ExportArgs{
					Docker: &dockerbuild.ExportDockerArgs{
						Tar: pulumi.Bool(true),
					},
				},
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("platforms", multiPlatform.Platforms)
		return nil
	})
}
