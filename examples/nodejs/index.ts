import * as pulumi from "@pulumi/pulumi";
import * as dockerbuild from "@pulumi/dockerbuild";

const myRandomResource = new dockerbuild.Random("myRandomResource", {length: 24});
export const value = myRandomResource.result;
