{{% examples %}}
## Example Usage
{{% example %}}
### Push to AWS ECR with caching

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import * as docker_build from "@pulumi/docker-build";

const ecrRepository = new aws.ecr.Repository("ecr-repository", {});
const authToken = aws.ecr.getAuthorizationTokenOutput({
    registryId: ecrRepository.registryId,
});
const myImage = new docker_build.Image("my-image", {
    cacheFrom: [{
        registry: {
            ref: pulumi.interpolate`${ecrRepository.repositoryUrl}:cache`,
        },
    }],
    cacheTo: [{
        registry: {
            imageManifest: true,
            ociMediaTypes: true,
            ref: pulumi.interpolate`${ecrRepository.repositoryUrl}:cache`,
        },
    }],
    context: {
        location: "./app",
    },
    push: true,
    registries: [{
        address: ecrRepository.repositoryUrl,
        password: authToken.apply(authToken => authToken.password),
        username: authToken.apply(authToken => authToken.userName),
    }],
    tags: [pulumi.interpolate`${ecrRepository.repositoryUrl}:latest`],
});
export const ref = myImage.ref;
```
```python
import pulumi
import pulumi_aws as aws
import pulumi_docker_build as docker_build

ecr_repository = aws.ecr.Repository("ecr-repository")
auth_token = aws.ecr.get_authorization_token_output(registry_id=ecr_repository.registry_id)
my_image = docker_build.Image("my-image",
    cache_from=[{
        "registry": {
            "ref": ecr_repository.repository_url.apply(lambda repository_url: f"{repository_url}:cache"),
        },
    }],
    cache_to=[{
        "registry": {
            "image_manifest": True,
            "oci_media_types": True,
            "ref": ecr_repository.repository_url.apply(lambda repository_url: f"{repository_url}:cache"),
        },
    }],
    context={
        "location": "./app",
    },
    push=True,
    registries=[{
        "address": ecr_repository.repository_url,
        "password": auth_token.password,
        "username": auth_token.user_name,
    }],
    tags=[ecr_repository.repository_url.apply(lambda repository_url: f"{repository_url}:latest")])
pulumi.export("ref", my_image.ref)
```
```csharp
using System.Collections.Generic;
using System.Linq;
using Pulumi;
using Aws = Pulumi.Aws;
using DockerBuild = Pulumi.DockerBuild;

return await Deployment.RunAsync(() => 
{
    var ecrRepository = new Aws.Ecr.Repository("ecr-repository");

    var authToken = Aws.Ecr.GetAuthorizationToken.Invoke(new()
    {
        RegistryId = ecrRepository.RegistryId,
    });

    var myImage = new DockerBuild.Image("my-image", new()
    {
        CacheFrom = new[]
        {
            new DockerBuild.Inputs.CacheFromArgs
            {
                Registry = new DockerBuild.Inputs.CacheFromRegistryArgs
                {
                    Ref = ecrRepository.RepositoryUrl.Apply(repositoryUrl => $"{repositoryUrl}:cache"),
                },
            },
        },
        CacheTo = new[]
        {
            new DockerBuild.Inputs.CacheToArgs
            {
                Registry = new DockerBuild.Inputs.CacheToRegistryArgs
                {
                    ImageManifest = true,
                    OciMediaTypes = true,
                    Ref = ecrRepository.RepositoryUrl.Apply(repositoryUrl => $"{repositoryUrl}:cache"),
                },
            },
        },
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        Push = true,
        Registries = new[]
        {
            new DockerBuild.Inputs.RegistryArgs
            {
                Address = ecrRepository.RepositoryUrl,
                Password = authToken.Apply(getAuthorizationTokenResult => getAuthorizationTokenResult.Password),
                Username = authToken.Apply(getAuthorizationTokenResult => getAuthorizationTokenResult.UserName),
            },
        },
        Tags = new[]
        {
            ecrRepository.RepositoryUrl.Apply(repositoryUrl => $"{repositoryUrl}:latest"),
        },
    });

    return new Dictionary<string, object?>
    {
        ["ref"] = myImage.Ref,
    };
});

```
```go
package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ecr"
	"github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		ecrRepository, err := ecr.NewRepository(ctx, "ecr-repository", nil)
		if err != nil {
			return err
		}
		authToken := ecr.GetAuthorizationTokenOutput(ctx, ecr.GetAuthorizationTokenOutputArgs{
			RegistryId: ecrRepository.RegistryId,
		}, nil)
		myImage, err := dockerbuild.NewImage(ctx, "my-image", &dockerbuild.ImageArgs{
			CacheFrom: dockerbuild.CacheFromArray{
				&dockerbuild.CacheFromArgs{
					Registry: &dockerbuild.CacheFromRegistryArgs{
						Ref: ecrRepository.RepositoryUrl.ApplyT(func(repositoryUrl string) (string, error) {
							return fmt.Sprintf("%v:cache", repositoryUrl), nil
						}).(pulumi.StringOutput),
					},
				},
			},
			CacheTo: dockerbuild.CacheToArray{
				&dockerbuild.CacheToArgs{
					Registry: &dockerbuild.CacheToRegistryArgs{
						ImageManifest: pulumi.Bool(true),
						OciMediaTypes: pulumi.Bool(true),
						Ref: ecrRepository.RepositoryUrl.ApplyT(func(repositoryUrl string) (string, error) {
							return fmt.Sprintf("%v:cache", repositoryUrl), nil
						}).(pulumi.StringOutput),
					},
				},
			},
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("./app"),
			},
			Push: pulumi.Bool(true),
			Registries: dockerbuild.RegistryArray{
				&dockerbuild.RegistryArgs{
					Address: ecrRepository.RepositoryUrl,
					Password: authToken.ApplyT(func(authToken ecr.GetAuthorizationTokenResult) (*string, error) {
						return &authToken.Password, nil
					}).(pulumi.StringPtrOutput),
					Username: authToken.ApplyT(func(authToken ecr.GetAuthorizationTokenResult) (*string, error) {
						return &authToken.UserName, nil
					}).(pulumi.StringPtrOutput),
				},
			},
			Tags: pulumi.StringArray{
				ecrRepository.RepositoryUrl.ApplyT(func(repositoryUrl string) (string, error) {
					return fmt.Sprintf("%v:latest", repositoryUrl), nil
				}).(pulumi.StringOutput),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("ref", myImage.Ref)
		return nil
	})
}
```
```yaml
description: Push to AWS ECR with caching
name: ecr
outputs:
    ref: ${my-image.ref}
resources:
    ecr-repository:
        type: aws:ecr:Repository
    my-image:
        properties:
            cacheFrom:
                - registry:
                    ref: ${ecr-repository.repositoryUrl}:cache
            cacheTo:
                - registry:
                    imageManifest: true
                    ociMediaTypes: true
                    ref: ${ecr-repository.repositoryUrl}:cache
            context:
                location: ./app
            push: true
            registries:
                - address: ${ecr-repository.repositoryUrl}
                  password: ${auth-token.password}
                  username: ${auth-token.userName}
            tags:
                - ${ecr-repository.repositoryUrl}:latest
        type: docker-build:Image
runtime: yaml
variables:
    auth-token:
        fn::aws:ecr:getAuthorizationToken:
            registryId: ${ecr-repository.registryId}
```
```java
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.aws.ecr.Repository;
import com.pulumi.aws.ecr.EcrFunctions;
import com.pulumi.aws.ecr.inputs.GetAuthorizationTokenArgs;
import com.pulumi.dockerbuild.Image;
import com.pulumi.dockerbuild.ImageArgs;
import com.pulumi.dockerbuild.inputs.CacheFromArgs;
import com.pulumi.dockerbuild.inputs.CacheFromRegistryArgs;
import com.pulumi.dockerbuild.inputs.CacheToArgs;
import com.pulumi.dockerbuild.inputs.CacheToRegistryArgs;
import com.pulumi.dockerbuild.inputs.BuildContextArgs;
import com.pulumi.dockerbuild.inputs.RegistryArgs;
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
        var ecrRepository = new Repository("ecrRepository");

        final var authToken = EcrFunctions.getAuthorizationToken(GetAuthorizationTokenArgs.builder()
            .registryId(ecrRepository.registryId())
            .build());

        var myImage = new Image("myImage", ImageArgs.builder()
            .cacheFrom(CacheFromArgs.builder()
                .registry(CacheFromRegistryArgs.builder()
                    .ref(ecrRepository.repositoryUrl().applyValue(repositoryUrl -> String.format("%s:cache", repositoryUrl)))
                    .build())
                .build())
            .cacheTo(CacheToArgs.builder()
                .registry(CacheToRegistryArgs.builder()
                    .imageManifest(true)
                    .ociMediaTypes(true)
                    .ref(ecrRepository.repositoryUrl().applyValue(repositoryUrl -> String.format("%s:cache", repositoryUrl)))
                    .build())
                .build())
            .context(BuildContextArgs.builder()
                .location("./app")
                .build())
            .push(true)
            .registries(RegistryArgs.builder()
                .address(ecrRepository.repositoryUrl())
                .password(authToken.applyValue(getAuthorizationTokenResult -> getAuthorizationTokenResult).applyValue(authToken -> authToken.applyValue(getAuthorizationTokenResult -> getAuthorizationTokenResult.password())))
                .username(authToken.applyValue(getAuthorizationTokenResult -> getAuthorizationTokenResult).applyValue(authToken -> authToken.applyValue(getAuthorizationTokenResult -> getAuthorizationTokenResult.userName())))
                .build())
            .tags(ecrRepository.repositoryUrl().applyValue(repositoryUrl -> String.format("%s:latest", repositoryUrl)))
            .build());

        ctx.export("ref", myImage.ref());
    }
}
```
{{% /example %}}
{{% example %}}
### Multi-platform image

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as docker_build from "@pulumi/docker-build";

const image = new docker_build.Image("image", {
    context: {
        location: "app",
    },
    platforms: [
        docker_build.Platform.Plan9_amd64,
        docker_build.Platform.Plan9_386,
    ],
    push: false,
});
```
```python
import pulumi
import pulumi_docker_build as docker_build

image = docker_build.Image("image",
    context={
        "location": "app",
    },
    platforms=[
        docker_build.Platform.PLAN9_AMD64,
        docker_build.Platform.PLAN9_386,
    ],
    push=False)
```
```csharp
using System.Collections.Generic;
using System.Linq;
using Pulumi;
using DockerBuild = Pulumi.DockerBuild;

return await Deployment.RunAsync(() => 
{
    var image = new DockerBuild.Image("image", new()
    {
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "app",
        },
        Platforms = new[]
        {
            DockerBuild.Platform.Plan9_amd64,
            DockerBuild.Platform.Plan9_386,
        },
        Push = false,
    });

});

```
```go
package main

import (
	"github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := dockerbuild.NewImage(ctx, "image", &dockerbuild.ImageArgs{
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("app"),
			},
			Platforms: docker - build.PlatformArray{
				dockerbuild.Platform_Plan9_amd64,
				dockerbuild.Platform_Plan9_386,
			},
			Push: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```
```yaml
description: Multi-platform image
name: multi-platform
resources:
    image:
        properties:
            context:
                location: app
            platforms:
                - plan9/amd64
                - plan9/386
            push: false
        type: docker-build:Image
runtime: yaml
```
```java
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.dockerbuild.Image;
import com.pulumi.dockerbuild.ImageArgs;
import com.pulumi.dockerbuild.inputs.BuildContextArgs;
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
        var image = new Image("image", ImageArgs.builder()
            .context(BuildContextArgs.builder()
                .location("app")
                .build())
            .platforms(            
                "plan9/amd64",
                "plan9/386")
            .push(false)
            .build());

    }
}
```
{{% /example %}}
{{% example %}}
### Registry export

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as docker_build from "@pulumi/docker-build";

const image = new docker_build.Image("image", {
    context: {
        location: "app",
    },
    push: true,
    registries: [{
        address: "docker.io",
        password: dockerHubPassword,
        username: "pulumibot",
    }],
    tags: ["docker.io/pulumi/pulumi:3.107.0"],
});
export const ref = myImage.ref;
```
```python
import pulumi
import pulumi_docker_build as docker_build

image = docker_build.Image("image",
    context={
        "location": "app",
    },
    push=True,
    registries=[{
        "address": "docker.io",
        "password": docker_hub_password,
        "username": "pulumibot",
    }],
    tags=["docker.io/pulumi/pulumi:3.107.0"])
pulumi.export("ref", my_image["ref"])
```
```csharp
using System.Collections.Generic;
using System.Linq;
using Pulumi;
using DockerBuild = Pulumi.DockerBuild;

return await Deployment.RunAsync(() => 
{
    var image = new DockerBuild.Image("image", new()
    {
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "app",
        },
        Push = true,
        Registries = new[]
        {
            new DockerBuild.Inputs.RegistryArgs
            {
                Address = "docker.io",
                Password = dockerHubPassword,
                Username = "pulumibot",
            },
        },
        Tags = new[]
        {
            "docker.io/pulumi/pulumi:3.107.0",
        },
    });

    return new Dictionary<string, object?>
    {
        ["ref"] = myImage.Ref,
    };
});

```
```go
package main

import (
	"github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := dockerbuild.NewImage(ctx, "image", &dockerbuild.ImageArgs{
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("app"),
			},
			Push: pulumi.Bool(true),
			Registries: dockerbuild.RegistryArray{
				&dockerbuild.RegistryArgs{
					Address:  pulumi.String("docker.io"),
					Password: pulumi.Any(dockerHubPassword),
					Username: pulumi.String("pulumibot"),
				},
			},
			Tags: pulumi.StringArray{
				pulumi.String("docker.io/pulumi/pulumi:3.107.0"),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("ref", myImage.Ref)
		return nil
	})
}
```
```yaml
description: Registry export
name: registry
outputs:
    ref: ${my-image.ref}
resources:
    image:
        properties:
            context:
                location: app
            push: true
            registries:
                - address: docker.io
                  password: ${dockerHubPassword}
                  username: pulumibot
            tags:
                - docker.io/pulumi/pulumi:3.107.0
        type: docker-build:Image
runtime: yaml
```
```java
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.dockerbuild.Image;
import com.pulumi.dockerbuild.ImageArgs;
import com.pulumi.dockerbuild.inputs.BuildContextArgs;
import com.pulumi.dockerbuild.inputs.RegistryArgs;
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
        var image = new Image("image", ImageArgs.builder()
            .context(BuildContextArgs.builder()
                .location("app")
                .build())
            .push(true)
            .registries(RegistryArgs.builder()
                .address("docker.io")
                .password(dockerHubPassword)
                .username("pulumibot")
                .build())
            .tags("docker.io/pulumi/pulumi:3.107.0")
            .build());

        ctx.export("ref", myImage.ref());
    }
}
```
{{% /example %}}
{{% example %}}
### Caching

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as docker_build from "@pulumi/docker-build";

const image = new docker_build.Image("image", {
    cacheFrom: [{
        local: {
            src: "tmp/cache",
        },
    }],
    cacheTo: [{
        local: {
            dest: "tmp/cache",
            mode: docker_build.CacheMode.Max,
        },
    }],
    context: {
        location: "app",
    },
    push: false,
});
```
```python
import pulumi
import pulumi_docker_build as docker_build

image = docker_build.Image("image",
    cache_from=[{
        "local": {
            "src": "tmp/cache",
        },
    }],
    cache_to=[{
        "local": {
            "dest": "tmp/cache",
            "mode": docker_build.CacheMode.MAX,
        },
    }],
    context={
        "location": "app",
    },
    push=False)
```
```csharp
using System.Collections.Generic;
using System.Linq;
using Pulumi;
using DockerBuild = Pulumi.DockerBuild;

return await Deployment.RunAsync(() => 
{
    var image = new DockerBuild.Image("image", new()
    {
        CacheFrom = new[]
        {
            new DockerBuild.Inputs.CacheFromArgs
            {
                Local = new DockerBuild.Inputs.CacheFromLocalArgs
                {
                    Src = "tmp/cache",
                },
            },
        },
        CacheTo = new[]
        {
            new DockerBuild.Inputs.CacheToArgs
            {
                Local = new DockerBuild.Inputs.CacheToLocalArgs
                {
                    Dest = "tmp/cache",
                    Mode = DockerBuild.CacheMode.Max,
                },
            },
        },
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "app",
        },
        Push = false,
    });

});

```
```go
package main

import (
	"github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := dockerbuild.NewImage(ctx, "image", &dockerbuild.ImageArgs{
			CacheFrom: dockerbuild.CacheFromArray{
				&dockerbuild.CacheFromArgs{
					Local: &dockerbuild.CacheFromLocalArgs{
						Src: pulumi.String("tmp/cache"),
					},
				},
			},
			CacheTo: dockerbuild.CacheToArray{
				&dockerbuild.CacheToArgs{
					Local: &dockerbuild.CacheToLocalArgs{
						Dest: pulumi.String("tmp/cache"),
						Mode: dockerbuild.CacheModeMax,
					},
				},
			},
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("app"),
			},
			Push: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```
```yaml
description: Caching
name: caching
resources:
    image:
        properties:
            cacheFrom:
                - local:
                    src: tmp/cache
            cacheTo:
                - local:
                    dest: tmp/cache
                    mode: max
            context:
                location: app
            push: false
        type: docker-build:Image
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
import com.pulumi.dockerbuild.inputs.CacheFromLocalArgs;
import com.pulumi.dockerbuild.inputs.CacheToArgs;
import com.pulumi.dockerbuild.inputs.CacheToLocalArgs;
import com.pulumi.dockerbuild.inputs.BuildContextArgs;
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
        var image = new Image("image", ImageArgs.builder()
            .cacheFrom(CacheFromArgs.builder()
                .local(CacheFromLocalArgs.builder()
                    .src("tmp/cache")
                    .build())
                .build())
            .cacheTo(CacheToArgs.builder()
                .local(CacheToLocalArgs.builder()
                    .dest("tmp/cache")
                    .mode("max")
                    .build())
                .build())
            .context(BuildContextArgs.builder()
                .location("app")
                .build())
            .push(false)
            .build());

    }
}
```
{{% /example %}}
{{% example %}}
### Docker Build Cloud

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as docker_build from "@pulumi/docker-build";

const image = new docker_build.Image("image", {
    builder: {
        name: "cloud-builder-name",
    },
    context: {
        location: "app",
    },
    exec: true,
    push: false,
});
```
```python
import pulumi
import pulumi_docker_build as docker_build

image = docker_build.Image("image",
    builder={
        "name": "cloud-builder-name",
    },
    context={
        "location": "app",
    },
    exec_=True,
    push=False)
```
```csharp
using System.Collections.Generic;
using System.Linq;
using Pulumi;
using DockerBuild = Pulumi.DockerBuild;

return await Deployment.RunAsync(() => 
{
    var image = new DockerBuild.Image("image", new()
    {
        Builder = new DockerBuild.Inputs.BuilderConfigArgs
        {
            Name = "cloud-builder-name",
        },
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "app",
        },
        Exec = true,
        Push = false,
    });

});

```
```go
package main

import (
	"github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := dockerbuild.NewImage(ctx, "image", &dockerbuild.ImageArgs{
			Builder: &dockerbuild.BuilderConfigArgs{
				Name: pulumi.String("cloud-builder-name"),
			},
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("app"),
			},
			Exec: pulumi.Bool(true),
			Push: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```
```yaml
description: Docker Build Cloud
name: dbc
resources:
    image:
        properties:
            builder:
                name: cloud-builder-name
            context:
                location: app
            exec: true
            push: false
        type: docker-build:Image
runtime: yaml
```
```java
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.dockerbuild.Image;
import com.pulumi.dockerbuild.ImageArgs;
import com.pulumi.dockerbuild.inputs.BuilderConfigArgs;
import com.pulumi.dockerbuild.inputs.BuildContextArgs;
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
        var image = new Image("image", ImageArgs.builder()
            .builder(BuilderConfigArgs.builder()
                .name("cloud-builder-name")
                .build())
            .context(BuildContextArgs.builder()
                .location("app")
                .build())
            .exec(true)
            .push(false)
            .build());

    }
}
```
{{% /example %}}
{{% example %}}
### Build arguments

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as docker_build from "@pulumi/docker-build";

const image = new docker_build.Image("image", {
    buildArgs: {
        SET_ME_TO_TRUE: "true",
    },
    context: {
        location: "app",
    },
    push: false,
});
```
```python
import pulumi
import pulumi_docker_build as docker_build

image = docker_build.Image("image",
    build_args={
        "SET_ME_TO_TRUE": "true",
    },
    context={
        "location": "app",
    },
    push=False)
```
```csharp
using System.Collections.Generic;
using System.Linq;
using Pulumi;
using DockerBuild = Pulumi.DockerBuild;

return await Deployment.RunAsync(() => 
{
    var image = new DockerBuild.Image("image", new()
    {
        BuildArgs = 
        {
            { "SET_ME_TO_TRUE", "true" },
        },
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "app",
        },
        Push = false,
    });

});

```
```go
package main

import (
	"github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := dockerbuild.NewImage(ctx, "image", &dockerbuild.ImageArgs{
			BuildArgs: pulumi.StringMap{
				"SET_ME_TO_TRUE": pulumi.String("true"),
			},
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("app"),
			},
			Push: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```
```yaml
description: Build arguments
name: build-args
resources:
    image:
        properties:
            buildArgs:
                SET_ME_TO_TRUE: "true"
            context:
                location: app
            push: false
        type: docker-build:Image
runtime: yaml
```
```java
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.dockerbuild.Image;
import com.pulumi.dockerbuild.ImageArgs;
import com.pulumi.dockerbuild.inputs.BuildContextArgs;
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
        var image = new Image("image", ImageArgs.builder()
            .buildArgs(Map.of("SET_ME_TO_TRUE", "true"))
            .context(BuildContextArgs.builder()
                .location("app")
                .build())
            .push(false)
            .build());

    }
}
```
{{% /example %}}
{{% example %}}
### Build target

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as docker_build from "@pulumi/docker-build";

const image = new docker_build.Image("image", {
    context: {
        location: "app",
    },
    push: false,
    target: "build-me",
});
```
```python
import pulumi
import pulumi_docker_build as docker_build

image = docker_build.Image("image",
    context={
        "location": "app",
    },
    push=False,
    target="build-me")
```
```csharp
using System.Collections.Generic;
using System.Linq;
using Pulumi;
using DockerBuild = Pulumi.DockerBuild;

return await Deployment.RunAsync(() => 
{
    var image = new DockerBuild.Image("image", new()
    {
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "app",
        },
        Push = false,
        Target = "build-me",
    });

});

```
```go
package main

import (
	"github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := dockerbuild.NewImage(ctx, "image", &dockerbuild.ImageArgs{
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("app"),
			},
			Push:   pulumi.Bool(false),
			Target: pulumi.String("build-me"),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```
```yaml
description: Build target
name: build-target
resources:
    image:
        properties:
            context:
                location: app
            push: false
            target: build-me
        type: docker-build:Image
runtime: yaml
```
```java
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.dockerbuild.Image;
import com.pulumi.dockerbuild.ImageArgs;
import com.pulumi.dockerbuild.inputs.BuildContextArgs;
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
        var image = new Image("image", ImageArgs.builder()
            .context(BuildContextArgs.builder()
                .location("app")
                .build())
            .push(false)
            .target("build-me")
            .build());

    }
}
```
{{% /example %}}
{{% example %}}
### Named contexts

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as docker_build from "@pulumi/docker-build";

const image = new docker_build.Image("image", {
    context: {
        location: "app",
        named: {
            "golang:latest": {
                location: "docker-image://golang@sha256:b8e62cf593cdaff36efd90aa3a37de268e6781a2e68c6610940c48f7cdf36984",
            },
        },
    },
    push: false,
});
```
```python
import pulumi
import pulumi_docker_build as docker_build

image = docker_build.Image("image",
    context={
        "location": "app",
        "named": {
            "golang:latest": {
                "location": "docker-image://golang@sha256:b8e62cf593cdaff36efd90aa3a37de268e6781a2e68c6610940c48f7cdf36984",
            },
        },
    },
    push=False)
```
```csharp
using System.Collections.Generic;
using System.Linq;
using Pulumi;
using DockerBuild = Pulumi.DockerBuild;

return await Deployment.RunAsync(() => 
{
    var image = new DockerBuild.Image("image", new()
    {
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "app",
            Named = 
            {
                { "golang:latest", new DockerBuild.Inputs.ContextArgs
                {
                    Location = "docker-image://golang@sha256:b8e62cf593cdaff36efd90aa3a37de268e6781a2e68c6610940c48f7cdf36984",
                } },
            },
        },
        Push = false,
    });

});

```
```go
package main

import (
	"github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := dockerbuild.NewImage(ctx, "image", &dockerbuild.ImageArgs{
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("app"),
				Named: dockerbuild.ContextMap{
					"golang:latest": &dockerbuild.ContextArgs{
						Location: pulumi.String("docker-image://golang@sha256:b8e62cf593cdaff36efd90aa3a37de268e6781a2e68c6610940c48f7cdf36984"),
					},
				},
			},
			Push: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```
```yaml
description: Named contexts
name: named-contexts
resources:
    image:
        properties:
            context:
                location: app
                named:
                    golang:latest:
                        location: docker-image://golang@sha256:b8e62cf593cdaff36efd90aa3a37de268e6781a2e68c6610940c48f7cdf36984
            push: false
        type: docker-build:Image
runtime: yaml
```
```java
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.dockerbuild.Image;
import com.pulumi.dockerbuild.ImageArgs;
import com.pulumi.dockerbuild.inputs.BuildContextArgs;
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
        var image = new Image("image", ImageArgs.builder()
            .context(BuildContextArgs.builder()
                .location("app")
                .named(Map.of("golang:latest", Map.of("location", "docker-image://golang@sha256:b8e62cf593cdaff36efd90aa3a37de268e6781a2e68c6610940c48f7cdf36984")))
                .build())
            .push(false)
            .build());

    }
}
```
{{% /example %}}
{{% example %}}
### Remote context

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as docker_build from "@pulumi/docker-build";

const image = new docker_build.Image("image", {
    context: {
        location: "https://raw.githubusercontent.com/pulumi/pulumi-docker/api-types/provider/testdata/Dockerfile",
    },
    push: false,
});
```
```python
import pulumi
import pulumi_docker_build as docker_build

image = docker_build.Image("image",
    context={
        "location": "https://raw.githubusercontent.com/pulumi/pulumi-docker/api-types/provider/testdata/Dockerfile",
    },
    push=False)
```
```csharp
using System.Collections.Generic;
using System.Linq;
using Pulumi;
using DockerBuild = Pulumi.DockerBuild;

return await Deployment.RunAsync(() => 
{
    var image = new DockerBuild.Image("image", new()
    {
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "https://raw.githubusercontent.com/pulumi/pulumi-docker/api-types/provider/testdata/Dockerfile",
        },
        Push = false,
    });

});

```
```go
package main

import (
	"github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := dockerbuild.NewImage(ctx, "image", &dockerbuild.ImageArgs{
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("https://raw.githubusercontent.com/pulumi/pulumi-docker/api-types/provider/testdata/Dockerfile"),
			},
			Push: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```
```yaml
description: Remote context
name: remote-context
resources:
    image:
        properties:
            context:
                location: https://raw.githubusercontent.com/pulumi/pulumi-docker/api-types/provider/testdata/Dockerfile
            push: false
        type: docker-build:Image
runtime: yaml
```
```java
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.dockerbuild.Image;
import com.pulumi.dockerbuild.ImageArgs;
import com.pulumi.dockerbuild.inputs.BuildContextArgs;
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
        var image = new Image("image", ImageArgs.builder()
            .context(BuildContextArgs.builder()
                .location("https://raw.githubusercontent.com/pulumi/pulumi-docker/api-types/provider/testdata/Dockerfile")
                .build())
            .push(false)
            .build());

    }
}
```
{{% /example %}}
{{% example %}}
### Inline Dockerfile

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as docker_build from "@pulumi/docker-build";

const image = new docker_build.Image("image", {
    context: {
        location: "app",
    },
    dockerfile: {
        inline: `FROM busybox
COPY hello.c ./
`,
    },
    push: false,
});
```
```python
import pulumi
import pulumi_docker_build as docker_build

image = docker_build.Image("image",
    context={
        "location": "app",
    },
    dockerfile={
        "inline": """FROM busybox
COPY hello.c ./
""",
    },
    push=False)
```
```csharp
using System.Collections.Generic;
using System.Linq;
using Pulumi;
using DockerBuild = Pulumi.DockerBuild;

return await Deployment.RunAsync(() => 
{
    var image = new DockerBuild.Image("image", new()
    {
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "app",
        },
        Dockerfile = new DockerBuild.Inputs.DockerfileArgs
        {
            Inline = @"FROM busybox
COPY hello.c ./
",
        },
        Push = false,
    });

});

```
```go
package main

import (
	"github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := dockerbuild.NewImage(ctx, "image", &dockerbuild.ImageArgs{
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("app"),
			},
			Dockerfile: &dockerbuild.DockerfileArgs{
				Inline: pulumi.String("FROM busybox\nCOPY hello.c ./\n"),
			},
			Push: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```
```yaml
description: Inline Dockerfile
name: inline
resources:
    image:
        properties:
            context:
                location: app
            dockerfile:
                inline: |
                    FROM busybox
                    COPY hello.c ./
            push: false
        type: docker-build:Image
runtime: yaml
```
```java
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.dockerbuild.Image;
import com.pulumi.dockerbuild.ImageArgs;
import com.pulumi.dockerbuild.inputs.BuildContextArgs;
import com.pulumi.dockerbuild.inputs.DockerfileArgs;
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
        var image = new Image("image", ImageArgs.builder()
            .context(BuildContextArgs.builder()
                .location("app")
                .build())
            .dockerfile(DockerfileArgs.builder()
                .inline("""
FROM busybox
COPY hello.c ./
                """)
                .build())
            .push(false)
            .build());

    }
}
```
{{% /example %}}
{{% example %}}
### Remote context

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as docker_build from "@pulumi/docker-build";

const image = new docker_build.Image("image", {
    context: {
        location: "https://github.com/docker-library/hello-world.git",
    },
    dockerfile: {
        location: "app/Dockerfile",
    },
    push: false,
});
```
```python
import pulumi
import pulumi_docker_build as docker_build

image = docker_build.Image("image",
    context={
        "location": "https://github.com/docker-library/hello-world.git",
    },
    dockerfile={
        "location": "app/Dockerfile",
    },
    push=False)
```
```csharp
using System.Collections.Generic;
using System.Linq;
using Pulumi;
using DockerBuild = Pulumi.DockerBuild;

return await Deployment.RunAsync(() => 
{
    var image = new DockerBuild.Image("image", new()
    {
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "https://github.com/docker-library/hello-world.git",
        },
        Dockerfile = new DockerBuild.Inputs.DockerfileArgs
        {
            Location = "app/Dockerfile",
        },
        Push = false,
    });

});

```
```go
package main

import (
	"github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := dockerbuild.NewImage(ctx, "image", &dockerbuild.ImageArgs{
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("https://github.com/docker-library/hello-world.git"),
			},
			Dockerfile: &dockerbuild.DockerfileArgs{
				Location: pulumi.String("app/Dockerfile"),
			},
			Push: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```
```yaml
description: Remote context
name: remote-context
resources:
    image:
        properties:
            context:
                location: https://github.com/docker-library/hello-world.git
            dockerfile:
                location: app/Dockerfile
            push: false
        type: docker-build:Image
runtime: yaml
```
```java
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.dockerbuild.Image;
import com.pulumi.dockerbuild.ImageArgs;
import com.pulumi.dockerbuild.inputs.BuildContextArgs;
import com.pulumi.dockerbuild.inputs.DockerfileArgs;
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
        var image = new Image("image", ImageArgs.builder()
            .context(BuildContextArgs.builder()
                .location("https://github.com/docker-library/hello-world.git")
                .build())
            .dockerfile(DockerfileArgs.builder()
                .location("app/Dockerfile")
                .build())
            .push(false)
            .build());

    }
}
```
{{% /example %}}
{{% example %}}
### Local export

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as docker_build from "@pulumi/docker-build";

const image = new docker_build.Image("image", {
    context: {
        location: "app",
    },
    exports: [{
        docker: {
            tar: true,
        },
    }],
    push: false,
});
```
```python
import pulumi
import pulumi_docker_build as docker_build

image = docker_build.Image("image",
    context={
        "location": "app",
    },
    exports=[{
        "docker": {
            "tar": True,
        },
    }],
    push=False)
```
```csharp
using System.Collections.Generic;
using System.Linq;
using Pulumi;
using DockerBuild = Pulumi.DockerBuild;

return await Deployment.RunAsync(() => 
{
    var image = new DockerBuild.Image("image", new()
    {
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "app",
        },
        Exports = new[]
        {
            new DockerBuild.Inputs.ExportArgs
            {
                Docker = new DockerBuild.Inputs.ExportDockerArgs
                {
                    Tar = true,
                },
            },
        },
        Push = false,
    });

});

```
```go
package main

import (
	"github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := dockerbuild.NewImage(ctx, "image", &dockerbuild.ImageArgs{
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("app"),
			},
			Exports: dockerbuild.ExportArray{
				&dockerbuild.ExportArgs{
					Docker: &dockerbuild.ExportDockerArgs{
						Tar: pulumi.Bool(true),
					},
				},
			},
			Push: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```
```yaml
description: Local export
name: docker-load
resources:
    image:
        properties:
            context:
                location: app
            exports:
                - docker:
                    tar: true
            push: false
        type: docker-build:Image
runtime: yaml
```
```java
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.dockerbuild.Image;
import com.pulumi.dockerbuild.ImageArgs;
import com.pulumi.dockerbuild.inputs.BuildContextArgs;
import com.pulumi.dockerbuild.inputs.ExportArgs;
import com.pulumi.dockerbuild.inputs.ExportDockerArgs;
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
        var image = new Image("image", ImageArgs.builder()
            .context(BuildContextArgs.builder()
                .location("app")
                .build())
            .exports(ExportArgs.builder()
                .docker(ExportDockerArgs.builder()
                    .tar(true)
                    .build())
                .build())
            .push(false)
            .build());

    }
}
```
{{% /example %}}
{{% /examples %}}