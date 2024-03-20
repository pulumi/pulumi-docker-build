package main

import (
	"github.com/pulumi/pulumi-dockerbuild/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		myRandomResource, err := dockerbuild.NewRandom(ctx, "myRandomResource", &dockerbuild.RandomArgs{
			Length: pulumi.Int(24),
		})
		if err != nil {
			return err
		}
		ctx.Export("value", myRandomResource.Result)
		return nil
	})
}
