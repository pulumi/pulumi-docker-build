import * as pulumi from "@pulumi/pulumi";
import * as docker_build from "@pulumi/docker-build";

const config = new pulumi.Config();
const dockerHubPassword = config.require("dockerHubPassword");
const multiPlatform = new docker_build.Image("multiPlatform", {
    push: false,
    dockerfile: {
        location: "./app/Dockerfile.multiPlatform",
    },
    context: {
        location: "./app",
    },
    platforms: [
        docker_build.Platform.Plan9_amd64,
        docker_build.Platform.Plan9_386,
    ],
});
const registryPush = new docker_build.Image("registryPush", {
    push: false,
    context: {
        location: "./app",
    },
    tags: ["docker.io/pulumibot/buildkit-e2e:example"],
    exports: [{
        registry: {
            ociMediaTypes: true,
            push: false,
        },
    }],
    registries: [{
        address: "docker.io",
        username: "pulumibot",
        password: dockerHubPassword,
    }],
});
const cached = new docker_build.Image("cached", {
    push: false,
    context: {
        location: "./app",
    },
    cacheTo: [{
        local: {
            dest: "tmp/cache",
            mode: docker_build.CacheMode.Max,
        },
    }],
    cacheFrom: [{
        local: {
            src: "tmp/cache",
        },
    }],
});
const buildArgs = new docker_build.Image("buildArgs", {
    push: false,
    dockerfile: {
        location: "./app/Dockerfile.buildArgs",
    },
    context: {
        location: "./app",
    },
    buildArgs: {
        SET_ME_TO_TRUE: "true",
    },
});
const extraHosts = new docker_build.Image("extraHosts", {
    push: false,
    dockerfile: {
        location: "./app/Dockerfile.extraHosts",
    },
    context: {
        location: "./app",
    },
    addHosts: ["metadata.google.internal:169.254.169.254"],
});
const sshMount = new docker_build.Image("sshMount", {
    push: false,
    dockerfile: {
        location: "./app/Dockerfile.sshMount",
    },
    context: {
        location: "./app",
    },
    ssh: [{
        id: "default",
    }],
});
const secrets = new docker_build.Image("secrets", {
    push: false,
    dockerfile: {
        location: "./app/Dockerfile.secrets",
    },
    context: {
        location: "./app",
    },
    secrets: {
        password: "hunter2",
    },
});
const labels = new docker_build.Image("labels", {
    push: false,
    context: {
        location: "./app",
    },
    labels: {
        description: "This image will get a descriptive label 👍",
    },
});
const target = new docker_build.Image("target", {
    push: false,
    dockerfile: {
        location: "./app/Dockerfile.target",
    },
    context: {
        location: "./app",
    },
    target: "build-me",
});
const namedContexts = new docker_build.Image("namedContexts", {
    push: false,
    dockerfile: {
        location: "./app/Dockerfile.namedContexts",
    },
    context: {
        location: "./app",
        named: {
            "golang:latest": {
                location: "docker-image://golang@sha256:b8e62cf593cdaff36efd90aa3a37de268e6781a2e68c6610940c48f7cdf36984",
            },
        },
    },
});
const remoteContext = new docker_build.Image("remoteContext", {
    push: false,
    context: {
        location: "https://raw.githubusercontent.com/pulumi/pulumi-docker/api-types/provider/testdata/Dockerfile",
    },
});
const remoteContextWithInline = new docker_build.Image("remoteContextWithInline", {
    push: false,
    dockerfile: {
        inline: `FROM busybox
COPY hello.c ./
`,
    },
    context: {
        location: "https://github.com/docker-library/hello-world.git",
    },
});
const inline = new docker_build.Image("inline", {
    push: false,
    dockerfile: {
        inline: `FROM alpine
RUN echo "This uses an inline Dockerfile! 👍"
`,
    },
    context: {
        location: "./app",
    },
});
const dockerLoad = new docker_build.Image("dockerLoad", {
    push: false,
    context: {
        location: "./app",
    },
    exports: [{
        docker: {
            tar: true,
        },
    }],
});
export const platforms = multiPlatform.platforms;
