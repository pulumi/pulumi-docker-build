import pulumi
import pulumi_dockerbuild as dockerbuild

config = pulumi.Config()
docker_hub_password = config.require("dockerHubPassword")
multi_platform = dockerbuild.Image("multiPlatform",
    dockerfile=dockerbuild.DockerfileArgs(
        location="./app/Dockerfile.multiPlatform",
    ),
    context=dockerbuild.BuildContextArgs(
        location="./app",
    ),
    platforms=[
        dockerbuild.Platform.PLAN9_AMD64,
        dockerbuild.Platform.PLAN9_386,
    ])
registry_push = dockerbuild.Image("registryPush",
    context=dockerbuild.BuildContextArgs(
        location="./app",
    ),
    tags=["docker.io/pulumibot/buildkit-e2e:example"],
    exports=[dockerbuild.ExportArgs(
        registry=dockerbuild.ExportRegistryArgs(
            oci_media_types=True,
            push=False,
        ),
    )],
    registries=[dockerbuild.RegistryArgs(
        address="docker.io",
        username="pulumibot",
        password=docker_hub_password,
    )])
cached = dockerbuild.Image("cached",
    context=dockerbuild.BuildContextArgs(
        location="./app",
    ),
    cache_to=[dockerbuild.CacheToArgs(
        local=dockerbuild.CacheToLocalArgs(
            dest="tmp/cache",
            mode=dockerbuild.CacheMode.MAX,
        ),
    )],
    cache_from=[dockerbuild.CacheFromArgs(
        local=dockerbuild.CacheFromLocalArgs(
            src="tmp/cache",
        ),
    )])
build_args = dockerbuild.Image("buildArgs",
    dockerfile=dockerbuild.DockerfileArgs(
        location="./app/Dockerfile.buildArgs",
    ),
    context=dockerbuild.BuildContextArgs(
        location="./app",
    ),
    build_args={
        "SET_ME_TO_TRUE": "true",
    })
extra_hosts = dockerbuild.Image("extraHosts",
    dockerfile=dockerbuild.DockerfileArgs(
        location="./app/Dockerfile.extraHosts",
    ),
    context=dockerbuild.BuildContextArgs(
        location="./app",
    ),
    add_hosts=["metadata.google.internal:169.254.169.254"])
ssh_mount = dockerbuild.Image("sshMount",
    dockerfile=dockerbuild.DockerfileArgs(
        location="./app/Dockerfile.sshMount",
    ),
    context=dockerbuild.BuildContextArgs(
        location="./app",
    ),
    ssh=[dockerbuild.SSHArgs(
        id="default",
    )])
secrets = dockerbuild.Image("secrets",
    dockerfile=dockerbuild.DockerfileArgs(
        location="./app/Dockerfile.secrets",
    ),
    context=dockerbuild.BuildContextArgs(
        location="./app",
    ),
    secrets={
        "password": "hunter2",
    })
labels = dockerbuild.Image("labels",
    context=dockerbuild.BuildContextArgs(
        location="./app",
    ),
    labels={
        "description": "This image will get a descriptive label 👍",
    })
target = dockerbuild.Image("target",
    dockerfile=dockerbuild.DockerfileArgs(
        location="./app/Dockerfile.target",
    ),
    context=dockerbuild.BuildContextArgs(
        location="./app",
    ),
    target="build-me")
named_contexts = dockerbuild.Image("namedContexts",
    dockerfile=dockerbuild.DockerfileArgs(
        location="./app/Dockerfile.namedContexts",
    ),
    context=dockerbuild.BuildContextArgs(
        location="./app",
        named={
            "golang:latest": dockerbuild.ContextArgs(
                location="docker-image://golang@sha256:b8e62cf593cdaff36efd90aa3a37de268e6781a2e68c6610940c48f7cdf36984",
            ),
        },
    ))
remote_context = dockerbuild.Image("remoteContext", context=dockerbuild.BuildContextArgs(
    location="https://raw.githubusercontent.com/pulumi/pulumi-docker/api-types/provider/testdata/Dockerfile",
))
remote_context_with_inline = dockerbuild.Image("remoteContextWithInline",
    dockerfile=dockerbuild.DockerfileArgs(
        inline="""FROM busybox
COPY hello.c ./
""",
    ),
    context=dockerbuild.BuildContextArgs(
        location="https://github.com/docker-library/hello-world.git",
    ))
inline = dockerbuild.Image("inline",
    dockerfile=dockerbuild.DockerfileArgs(
        inline="""FROM alpine
RUN echo "This uses an inline Dockerfile! 👍"
""",
    ),
    context=dockerbuild.BuildContextArgs(
        location="./app",
    ))
docker_load = dockerbuild.Image("dockerLoad",
    context=dockerbuild.BuildContextArgs(
        location="./app",
    ),
    exports=[dockerbuild.ExportArgs(
        docker=dockerbuild.ExportDockerArgs(
            tar=True,
        ),
    )])
pulumi.export("platforms", multi_platform.platforms)
