import * as pulumi from "@pulumi/pulumi";
import * as dockerbuild from "@pulumi/dockerbuild";

const config = new pulumi.Config();
const dockerHubPassword = config.require("dockerHubPassword");
const multiPlatform = new dockerbuild.Image("multiPlatform", {
    dockerfile: {
        location: "./app/Dockerfile.multiPlatform",
    },
    context: {
        location: "./app",
    },
    platforms: [
        dockerbuild.Platform.Plan9_amd64,
        dockerbuild.Platform.Plan9_386,
    ],
});
const registryPush = new dockerbuild.Image("registryPush", {
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
const cached = new dockerbuild.Image("cached", {
    context: {
        location: "./app",
    },
    cacheTo: [{
        local: {
            dest: "tmp/cache",
            mode: dockerbuild.CacheMode.Max,
        },
    }],
    cacheFrom: [{
        local: {
            src: "tmp/cache",
        },
    }],
});
const buildArgs = new dockerbuild.Image("buildArgs", {
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
const extraHosts = new dockerbuild.Image("extraHosts", {
    dockerfile: {
        location: "./app/Dockerfile.extraHosts",
    },
    context: {
        location: "./app",
    },
    addHosts: ["metadata.google.internal:169.254.169.254"],
});
const sshMount = new dockerbuild.Image("sshMount", {
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
const secrets = new dockerbuild.Image("secrets", {
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
const labels = new dockerbuild.Image("labels", {
    context: {
        location: "./app",
    },
    labels: {
        description: "This image will get a descriptive label 👍",
    },
});
const target = new dockerbuild.Image("target", {
    dockerfile: {
        location: "./app/Dockerfile.target",
    },
    context: {
        location: "./app",
    },
    target: "build-me",
});
const namedContexts = new dockerbuild.Image("namedContexts", {
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
const remoteContext = new dockerbuild.Image("remoteContext", {context: {
    location: "https://raw.githubusercontent.com/pulumi/pulumi-docker/api-types/provider/testdata/Dockerfile",
}});
const remoteContextWithInline = new dockerbuild.Image("remoteContextWithInline", {
    dockerfile: {
        inline: `FROM busybox
COPY hello.c ./
`,
    },
    context: {
        location: "https://github.com/docker-library/hello-world.git",
    },
});
const inline = new dockerbuild.Image("inline", {
    dockerfile: {
        inline: `FROM alpine
RUN echo "This uses an inline Dockerfile! 👍"
`,
    },
    context: {
        location: "./app",
    },
});
const dockerLoad = new dockerbuild.Image("dockerLoad", {
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
