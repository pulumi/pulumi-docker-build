import * as buildx from "@pulumi/docker-build";

new buildx.Provider("with-structured-config", {
  registries: [{ username: "foo", password: "bar", address: "docker.io" }],
});
