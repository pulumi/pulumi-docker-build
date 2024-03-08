import * as pulumi from "@pulumi/pulumi";
import * as docker_native from "@pulumi/docker-native";

const myRandomResource = new docker_native.Random("myRandomResource", {length: 24});
export const output = {
    value: myRandomResource.result,
};
