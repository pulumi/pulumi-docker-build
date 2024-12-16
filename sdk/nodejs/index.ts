// This line should be fixed up by pulumi-bot.

import * as pulumi from "@pulumi/pulumi";
import * as utilities from "./utilities";

// Export members:
export { ImageArgs } from "./image";
export type Image = import("./image").Image;
export const Image: typeof import("./image").Image = null as any;
utilities.lazyLoad(exports, ["Image"], () => require("./image"));

export { IndexArgs } from "./index_";
export type Index = import("./index_").Index;
export const Index: typeof import("./index_").Index = null as any;
utilities.lazyLoad(exports, ["Index"], () => require("./index_"));

export { ProviderArgs } from "./provider";
export type Provider = import("./provider").Provider;
export const Provider: typeof import("./provider").Provider = null as any;
utilities.lazyLoad(exports, ["Provider"], () => require("./provider"));

// Export enums:
export * from "./types/enums";

// Export sub-modules:
import * as config from "./config";
import * as types from "./types";

export { config, types };

const _module = {
  version: utilities.getVersion(),
  construct: (name: string, type: string, urn: string): pulumi.Resource => {
    switch (type) {
      case "docker-build:index:Image":
        return new Image(name, <any>undefined, { urn });
      case "docker-build:index:Index":
        return new Index(name, <any>undefined, { urn });
      default:
        throw new Error(`unknown resource type ${type}`);
    }
  },
};
pulumi.runtime.registerResourceModule("docker-build", "index", _module);
pulumi.runtime.registerResourcePackage("docker-build", {
  version: utilities.getVersion(),
  constructProvider: (
    name: string,
    type: string,
    urn: string
  ): pulumi.ProviderResource => {
    if (type !== "pulumi:providers:docker-build") {
      throw new Error(`unknown provider type ${type}`);
    }
    return new Provider(name, <any>undefined, { urn });
  },
});
