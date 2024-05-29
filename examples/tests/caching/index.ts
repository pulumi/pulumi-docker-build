import * as buildx from "@pulumi/docker-build";
import * as pulumi from "@pulumi/pulumi";

const config = new pulumi.Config();

const start = new Date().getTime();

// docker buildx build \
//  -f Dockerfile \
//  --cache-to type=local,dest=tmp,mode=max,oci-mediatypes=true \
//  --cache-from type=local,src=tmp \
//  --build-arg SLEEP-MS=$SLEEP_MS \
//  -t not-pushed \
//  -f Dockerfile \
//  .
const img = new buildx.Image(`buildx-${config.require("name")}`, {
  tags: ["not-pushed"],
  dockerfile: { location: "Dockerfile" },
  push: false,
  context: { location: "." },
  buildArgs: {
    SLEEP_SECONDS: config.require("SLEEP_SECONDS"),
  },
  cacheTo: [{ raw: config.require("cacheTo") }],
  cacheFrom: [{ raw: config.require("cacheFrom") }],
  // Set registry auth if it was provided.
  registries: config.require("username")
    ? [
        {
          address: config.getSecret("address"),
          username: config.getSecret("username"),
          password: config.getSecret("password"),
        },
      ]
    : undefined,
});

export const durationSeconds = img.ref.apply(
  (_) => (new Date().getTime() - start) / 1000.0
);
