import * as pulumi from "@pulumi/pulumi";
import * as docker_build from "@pulumi/docker-build";

const inline = new docker_build.Image("inline", {
  push: false,
  dockerfile: {
    inline: `FROM alpine
RUN echo "This uses an inline Dockerfile! 👍"
`,
  },
});
