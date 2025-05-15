import pulumi
import pulumi_docker_build as docker_build

config = pulumi.Config()
docker_hub_password = config.require("dockerHubPassword")
multi_platform = docker_build.Image("multiPlatform",
    push=False,
    dockerfile={
        "location": "./app/Dockerfile.multiPlatform",
    },
    context={
        "location": "./app",
    },
    platforms=[
        docker_build.Platform.PLAN9_AMD64,
        docker_build.Platform.PLAN9_386,
    ])
registry_push = docker_build.Image("registryPush",
    push=False,
    context={
        "location": "./app",
    },
    tags=["docker.io/pulumibot/buildkit-e2e:example"],
    exports=[{
        "registry": {
            "oci_media_types": True,
            "push": False,
        },
    }],
    registries=[{
        "address": "docker.io",
        "username": "pulumibot",
        "password": docker_hub_password,
    }])
cached = docker_build.Image("cached",
    push=False,
    context={
        "location": "./app",
    },
    cache_to=[{
        "local": {
            "dest": "tmp/cache",
            "mode": docker_build.CacheMode.MAX,
        },
    }],
    cache_from=[{
        "local": {
            "src": "tmp/cache",
        },
    }])
build_args = docker_build.Image("buildArgs",
    push=False,
    dockerfile={
        "location": "./app/Dockerfile.buildArgs",
    },
    context={
        "location": "./app",
    },
    build_args={
        "SET_ME_TO_TRUE": "true",
    })
extra_hosts = docker_build.Image("extraHosts",
    push=False,
    dockerfile={
        "location": "./app/Dockerfile.extraHosts",
    },
    context={
        "location": "./app",
    },
    add_hosts=["metadata.google.internal:169.254.169.254"])
ssh_mount = docker_build.Image("sshMount",
    push=False,
    dockerfile={
        "location": "./app/Dockerfile.sshMount",
    },
    context={
        "location": "./app",
    },
    ssh=[{
        "id": "default",
    }])
secrets = docker_build.Image("secrets",
    push=False,
    dockerfile={
        "location": "./app/Dockerfile.secrets",
    },
    context={
        "location": "./app",
    },
    secrets={
        "password": "hunter2",
    })
labels = docker_build.Image("labels",
    push=False,
    context={
        "location": "./app",
    },
    labels={
        "description": "This image will get a descriptive label 👍",
    })
target = docker_build.Image("target",
    push=False,
    dockerfile={
        "location": "./app/Dockerfile.target",
    },
    context={
        "location": "./app",
    },
    target="build-me")
named_contexts = docker_build.Image("namedContexts",
    push=False,
    dockerfile={
        "location": "./app/Dockerfile.namedContexts",
    },
    context={
        "location": "./app",
        "named": {
            "golang:latest": {
                "location": "docker-image://golang@sha256:b8e62cf593cdaff36efd90aa3a37de268e6781a2e68c6610940c48f7cdf36984",
            },
        },
    })
remote_context = docker_build.Image("remoteContext",
    push=False,
    context={
        "location": "https://raw.githubusercontent.com/pulumi/pulumi-docker/api-types/provider/testdata/Dockerfile",
    })
remote_context_with_inline = docker_build.Image("remoteContextWithInline",
    push=False,
    dockerfile={
        "inline": """FROM busybox
COPY hello.c ./
""",
    },
    context={
        "location": "https://github.com/docker-library/hello-world.git",
    })
inline = docker_build.Image("inline",
    push=False,
    dockerfile={
        "inline": """FROM alpine
RUN echo "This uses an inline Dockerfile! 👍"
""",
    })
docker_load = docker_build.Image("dockerLoad",
    push=False,
    context={
        "location": "./app",
    },
    exports=[{
        "docker": {
            "tar": True,
        },
    }])
pulumi.export("platforms", multi_platform.platforms)
