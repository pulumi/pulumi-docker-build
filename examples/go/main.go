package main

import (
	"github.com/pulumi/pulumi-dockerbuild/sdk/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		myRandomResource, err := docker.NewRandom(ctx, "myRandomResource", &docker.RandomArgs{
			Length: pulumi.Int(24),
		})
		if err != nil {
			return err
		}
		ctx.Export("value", myRandomResource.Result)
		return nil
	})
}
