{{% examples %}}
## Example Usage
{{% example %}}
### Multi-platform registry caching

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as dockerbuild from "@pulumi/dockerbuild";

const amd64 = new dockerbuild.Image("amd64", {
    cacheFrom: [{
        registry: {
            ref: "docker.io/pulumi/pulumi:cache-amd64",
        },
    }],
    cacheTo: [{
        registry: {
            mode: dockerbuild.CacheMode.Max,
            ref: "docker.io/pulumi/pulumi:cache-amd64",
        },
    }],
    context: {
        location: "app",
    },
    platforms: [dockerbuild.Platform.Linux_amd64],
    tags: ["docker.io/pulumi/pulumi:3.107.0-amd64"],
});
const arm64 = new dockerbuild.Image("arm64", {
    cacheFrom: [{
        registry: {
            ref: "docker.io/pulumi/pulumi:cache-arm64",
        },
    }],
    cacheTo: [{
        registry: {
            mode: dockerbuild.CacheMode.Max,
            ref: "docker.io/pulumi/pulumi:cache-arm64",
        },
    }],
    context: {
        location: "app",
    },
    platforms: [dockerbuild.Platform.Linux_arm64],
    tags: ["docker.io/pulumi/pulumi:3.107.0-arm64"],
});
const index = new dockerbuild.Index("index", {
    sources: [
        amd64.ref,
        arm64.ref,
    ],
    tag: "docker.io/pulumi/pulumi:3.107.0",
});
export const ref = index.ref;
```
```python
import pulumi
import pulumi_dockerbuild as dockerbuild

amd64 = dockerbuild.Image("amd64",
    cache_from=[dockerbuild.CacheFromArgs(
        registry=dockerbuild.CacheFromRegistryArgs(
            ref="docker.io/pulumi/pulumi:cache-amd64",
        ),
    )],
    cache_to=[dockerbuild.CacheToArgs(
        registry=dockerbuild.CacheToRegistryArgs(
            mode=dockerbuild.CacheMode.MAX,
            ref="docker.io/pulumi/pulumi:cache-amd64",
        ),
    )],
    context=dockerbuild.BuildContextArgs(
        location="app",
    ),
    platforms=[dockerbuild.Platform.LINUX_AMD64],
    tags=["docker.io/pulumi/pulumi:3.107.0-amd64"])
arm64 = dockerbuild.Image("arm64",
    cache_from=[dockerbuild.CacheFromArgs(
        registry=dockerbuild.CacheFromRegistryArgs(
            ref="docker.io/pulumi/pulumi:cache-arm64",
        ),
    )],
    cache_to=[dockerbuild.CacheToArgs(
        registry=dockerbuild.CacheToRegistryArgs(
            mode=dockerbuild.CacheMode.MAX,
            ref="docker.io/pulumi/pulumi:cache-arm64",
        ),
    )],
    context=dockerbuild.BuildContextArgs(
        location="app",
    ),
    platforms=[dockerbuild.Platform.LINUX_ARM64],
    tags=["docker.io/pulumi/pulumi:3.107.0-arm64"])
index = dockerbuild.Index("index",
    sources=[
        amd64.ref,
        arm64.ref,
    ],
    tag="docker.io/pulumi/pulumi:3.107.0")
pulumi.export("ref", index.ref)
```
```csharp
using System.Collections.Generic;
using System.Linq;
using Pulumi;
using Dockerbuild = Pulumi.Dockerbuild;

return await Deployment.RunAsync(() => 
{
    var amd64 = new Dockerbuild.Image("amd64", new()
    {
        CacheFrom = new[]
        {
            new Dockerbuild.Inputs.CacheFromArgs
            {
                Registry = new Dockerbuild.Inputs.CacheFromRegistryArgs
                {
                    Ref = "docker.io/pulumi/pulumi:cache-amd64",
                },
            },
        },
        CacheTo = new[]
        {
            new Dockerbuild.Inputs.CacheToArgs
            {
                Registry = new Dockerbuild.Inputs.CacheToRegistryArgs
                {
                    Mode = Dockerbuild.CacheMode.Max,
                    Ref = "docker.io/pulumi/pulumi:cache-amd64",
                },
            },
        },
        Context = new Dockerbuild.Inputs.BuildContextArgs
        {
            Location = "app",
        },
        Platforms = new[]
        {
            Dockerbuild.Platform.Linux_amd64,
        },
        Tags = new[]
        {
            "docker.io/pulumi/pulumi:3.107.0-amd64",
        },
    });

    var arm64 = new Dockerbuild.Image("arm64", new()
    {
        CacheFrom = new[]
        {
            new Dockerbuild.Inputs.CacheFromArgs
            {
                Registry = new Dockerbuild.Inputs.CacheFromRegistryArgs
                {
                    Ref = "docker.io/pulumi/pulumi:cache-arm64",
                },
            },
        },
        CacheTo = new[]
        {
            new Dockerbuild.Inputs.CacheToArgs
            {
                Registry = new Dockerbuild.Inputs.CacheToRegistryArgs
                {
                    Mode = Dockerbuild.CacheMode.Max,
                    Ref = "docker.io/pulumi/pulumi:cache-arm64",
                },
            },
        },
        Context = new Dockerbuild.Inputs.BuildContextArgs
        {
            Location = "app",
        },
        Platforms = new[]
        {
            Dockerbuild.Platform.Linux_arm64,
        },
        Tags = new[]
        {
            "docker.io/pulumi/pulumi:3.107.0-arm64",
        },
    });

    var index = new Dockerbuild.Index("index", new()
    {
        Sources = new[]
        {
            amd64.Ref,
            arm64.Ref,
        },
        Tag = "docker.io/pulumi/pulumi:3.107.0",
    });

    return new Dictionary<string, object?>
    {
        ["ref"] = index.Ref,
    };
});

```
```go
package main

import (
	"github.com/pulumi/pulumi-dockerbuild/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		amd64, err := dockerbuild.NewImage(ctx, "amd64", &dockerbuild.ImageArgs{
			CacheFrom: dockerbuild.CacheFromArray{
				&dockerbuild.CacheFromArgs{
					Registry: &dockerbuild.CacheFromRegistryArgs{
						Ref: pulumi.String("docker.io/pulumi/pulumi:cache-amd64"),
					},
				},
			},
			CacheTo: dockerbuild.CacheToArray{
				&dockerbuild.CacheToArgs{
					Registry: &dockerbuild.CacheToRegistryArgs{
						Mode: dockerbuild.CacheModeMax,
						Ref:  pulumi.String("docker.io/pulumi/pulumi:cache-amd64"),
					},
				},
			},
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("app"),
			},
			Platforms: dockerbuild.PlatformArray{
				dockerbuild.Platform_Linux_amd64,
			},
			Tags: pulumi.StringArray{
				pulumi.String("docker.io/pulumi/pulumi:3.107.0-amd64"),
			},
		})
		if err != nil {
			return err
		}
		arm64, err := dockerbuild.NewImage(ctx, "arm64", &dockerbuild.ImageArgs{
			CacheFrom: dockerbuild.CacheFromArray{
				&dockerbuild.CacheFromArgs{
					Registry: &dockerbuild.CacheFromRegistryArgs{
						Ref: pulumi.String("docker.io/pulumi/pulumi:cache-arm64"),
					},
				},
			},
			CacheTo: dockerbuild.CacheToArray{
				&dockerbuild.CacheToArgs{
					Registry: &dockerbuild.CacheToRegistryArgs{
						Mode: dockerbuild.CacheModeMax,
						Ref:  pulumi.String("docker.io/pulumi/pulumi:cache-arm64"),
					},
				},
			},
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("app"),
			},
			Platforms: dockerbuild.PlatformArray{
				dockerbuild.Platform_Linux_arm64,
			},
			Tags: pulumi.StringArray{
				pulumi.String("docker.io/pulumi/pulumi:3.107.0-arm64"),
			},
		})
		if err != nil {
			return err
		}
		index, err := dockerbuild.NewIndex(ctx, "index", &dockerbuild.IndexArgs{
			Sources: pulumi.StringArray{
				amd64.Ref,
				arm64.Ref,
			},
			Tag: pulumi.String("docker.io/pulumi/pulumi:3.107.0"),
		})
		if err != nil {
			return err
		}
		ctx.Export("ref", index.Ref)
		return nil
	})
}
```
```yaml
description: Multi-platform registry caching
name: registry-caching
outputs:
    ref: ${index.ref}
resources:
    amd64:
        properties:
            cacheFrom:
                - registry:
                    ref: docker.io/pulumi/pulumi:cache-amd64
            cacheTo:
                - registry:
                    mode: max
                    ref: docker.io/pulumi/pulumi:cache-amd64
            context:
                location: app
            platforms:
                - linux/amd64
            tags:
                - docker.io/pulumi/pulumi:3.107.0-amd64
        type: dockerbuild:Image
    arm64:
        properties:
            cacheFrom:
                - registry:
                    ref: docker.io/pulumi/pulumi:cache-arm64
            cacheTo:
                - registry:
                    mode: max
                    ref: docker.io/pulumi/pulumi:cache-arm64
            context:
                location: app
            platforms:
                - linux/arm64
            tags:
                - docker.io/pulumi/pulumi:3.107.0-arm64
        type: dockerbuild:Image
    index:
        properties:
            sources:
                - ${amd64.ref}
                - ${arm64.ref}
            tag: docker.io/pulumi/pulumi:3.107.0
        type: dockerbuild:Index
runtime: yaml
```
```java
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.dockerbuild.Image;
import com.pulumi.dockerbuild.ImageArgs;
import com.pulumi.dockerbuild.inputs.CacheFromArgs;
import com.pulumi.dockerbuild.inputs.CacheFromRegistryArgs;
import com.pulumi.dockerbuild.inputs.CacheToArgs;
import com.pulumi.dockerbuild.inputs.CacheToRegistryArgs;
import com.pulumi.dockerbuild.inputs.BuildContextArgs;
import com.pulumi.dockerbuild.Index;
import com.pulumi.dockerbuild.IndexArgs;
import java.util.List;
import java.util.ArrayList;
import java.util.Map;
import java.io.File;
import java.nio.file.Files;
import java.nio.file.Paths;

public class App {
    public static void main(String[] args) {
        Pulumi.run(App::stack);
    }

    public static void stack(Context ctx) {
        var amd64 = new Image("amd64", ImageArgs.builder()        
            .cacheFrom(CacheFromArgs.builder()
                .registry(CacheFromRegistryArgs.builder()
                    .ref("docker.io/pulumi/pulumi:cache-amd64")
                    .build())
                .build())
            .cacheTo(CacheToArgs.builder()
                .registry(CacheToRegistryArgs.builder()
                    .mode("max")
                    .ref("docker.io/pulumi/pulumi:cache-amd64")
                    .build())
                .build())
            .context(BuildContextArgs.builder()
                .location("app")
                .build())
            .platforms("linux/amd64")
            .tags("docker.io/pulumi/pulumi:3.107.0-amd64")
            .build());

        var arm64 = new Image("arm64", ImageArgs.builder()        
            .cacheFrom(CacheFromArgs.builder()
                .registry(CacheFromRegistryArgs.builder()
                    .ref("docker.io/pulumi/pulumi:cache-arm64")
                    .build())
                .build())
            .cacheTo(CacheToArgs.builder()
                .registry(CacheToRegistryArgs.builder()
                    .mode("max")
                    .ref("docker.io/pulumi/pulumi:cache-arm64")
                    .build())
                .build())
            .context(BuildContextArgs.builder()
                .location("app")
                .build())
            .platforms("linux/arm64")
            .tags("docker.io/pulumi/pulumi:3.107.0-arm64")
            .build());

        var index = new Index("index", IndexArgs.builder()        
            .sources(            
                amd64.ref(),
                arm64.ref())
            .tag("docker.io/pulumi/pulumi:3.107.0")
            .build());

        ctx.export("ref", index.ref());
    }
}
```
{{% /example %}}
{{% /examples %}}