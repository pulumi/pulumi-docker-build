import pulumi
import pulumi_docker_build as docker_build

config = pulumi.Config()
docker_hub_password = config.require("dockerHubPassword")
multi_platform = docker_build.Image("multiPlatform",
    dockerfile=docker_build.DockerfileArgs(
        location="./app/Dockerfile.multiPlatform",
    ),
    context=docker_build.BuildContextArgs(
        location="./app",
    ),
    platforms=[
        docker_build.Platform.PLAN9_AMD64,
        docker_build.Platform.PLAN9_386,
    ])
registry_push = docker_build.Image("registryPush",
    context=docker_build.BuildContextArgs(
        location="./app",
    ),
    tags=["docker.io/pulumibot/buildkit-e2e:example"],
    exports=[docker_build.ExportArgs(
        registry=docker_build.ExportRegistryArgs(
            oci_media_types=True,
            push=False,
        ),
    )],
    registries=[docker_build.RegistryArgs(
        address="docker.io",
        username="pulumibot",
        password=docker_hub_password,
    )])
cached = docker_build.Image("cached",
    context=docker_build.BuildContextArgs(
        location="./app",
    ),
    cache_to=[docker_build.CacheToArgs(
        local=docker_build.CacheToLocalArgs(
            dest="tmp/cache",
            mode=docker_build.CacheMode.MAX,
        ),
    )],
    cache_from=[docker_build.CacheFromArgs(
        local=docker_build.CacheFromLocalArgs(
            src="tmp/cache",
        ),
    )])
build_args = docker_build.Image("buildArgs",
    dockerfile=docker_build.DockerfileArgs(
        location="./app/Dockerfile.buildArgs",
    ),
    context=docker_build.BuildContextArgs(
        location="./app",
    ),
    build_args={
        "SET_ME_TO_TRUE": "true",
    })
extra_hosts = docker_build.Image("extraHosts",
    dockerfile=docker_build.DockerfileArgs(
        location="./app/Dockerfile.extraHosts",
    ),
    context=docker_build.BuildContextArgs(
        location="./app",
    ),
    add_hosts=["metadata.google.internal:169.254.169.254"])
ssh_mount = docker_build.Image("sshMount",
    dockerfile=docker_build.DockerfileArgs(
        location="./app/Dockerfile.sshMount",
    ),
    context=docker_build.BuildContextArgs(
        location="./app",
    ),
    ssh=[docker_build.SSHArgs(
        id="default",
    )])
secrets = docker_build.Image("secrets",
    dockerfile=docker_build.DockerfileArgs(
        location="./app/Dockerfile.secrets",
    ),
    context=docker_build.BuildContextArgs(
        location="./app",
    ),
    secrets={
        "password": "hunter2",
    })
labels = docker_build.Image("labels",
    context=docker_build.BuildContextArgs(
        location="./app",
    ),
    labels={
        "description": "This image will get a descriptive label 👍",
    })
target = docker_build.Image("target",
    dockerfile=docker_build.DockerfileArgs(
        location="./app/Dockerfile.target",
    ),
    context=docker_build.BuildContextArgs(
        location="./app",
    ),
    target="build-me")
named_contexts = docker_build.Image("namedContexts",
    dockerfile=docker_build.DockerfileArgs(
        location="./app/Dockerfile.namedContexts",
    ),
    context=docker_build.BuildContextArgs(
        location="./app",
        named={
            "golang:latest": docker_build.ContextArgs(
                location="docker-image://golang@sha256:b8e62cf593cdaff36efd90aa3a37de268e6781a2e68c6610940c48f7cdf36984",
            ),
        },
    ))
remote_context = docker_build.Image("remoteContext", context=docker_build.BuildContextArgs(
    location="https://raw.githubusercontent.com/pulumi/pulumi-docker/api-types/provider/testdata/Dockerfile",
))
remote_context_with_inline = docker_build.Image("remoteContextWithInline",
    dockerfile=docker_build.DockerfileArgs(
        inline="""FROM busybox
COPY hello.c ./
""",
    ),
    context=docker_build.BuildContextArgs(
        location="https://github.com/docker-library/hello-world.git",
    ))
inline = docker_build.Image("inline",
    dockerfile=docker_build.DockerfileArgs(
        inline="""FROM alpine
RUN echo "This uses an inline Dockerfile! 👍"
""",
    ),
    context=docker_build.BuildContextArgs(
        location="./app",
    ))
docker_load = docker_build.Image("dockerLoad",
    context=docker_build.BuildContextArgs(
        location="./app",
    ),
    exports=[docker_build.ExportArgs(
        docker=docker_build.ExportDockerArgs(
            tar=True,
        ),
    )])
pulumi.export("platforms", multi_platform.platforms)
