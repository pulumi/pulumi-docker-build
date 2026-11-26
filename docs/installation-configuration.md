---
title: Docker Build Installation & Configuration
meta_desc: Provides an overview on how to configure the Pulumi Docker Build Provider.
layout: package
---

The Pulumi Docker Build provider builds modern Docker images with [buildx](https://docs.docker.com/reference/cli/docker/buildx/) and [BuildKit](https://docs.docker.com/build/buildkit/).

## Installation

The Docker Build provider is available as a package in all Pulumi languages:

* JavaScript/TypeScript: [`@pulumi/docker-build`](https://www.npmjs.com/package/@pulumi/docker-build)
* Python: [`pulumi-docker-build`](https://pypi.org/project/pulumi-docker-build/)
* Go: [`github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild`](https://pkg.go.dev/github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild)
* .NET: [`Pulumi.DockerBuild`](https://www.nuget.org/packages/Pulumi.DockerBuild)
* Java: [`com.pulumi/docker-build`](https://central.sonatype.com/artifact/com.pulumi/docker-build)

## Configuring the provider

### Host

The `DOCKER_HOST` environment variable can be used to specify a custom build daemon's location.

```bash
$ export DOCKER_HOST=tcp://127.0.0.1:2376/
```

This can also be specified in your stack's configuration:

```bash
$ pulumi config set docker-build:host tcp://127.0.0.1:2376/
```

### Registries

The `registries` key configures credentials for one or more registries the provider can authenticate with. It is an array of objects with the following properties:

* `address` (required): The registry's address (e.g. `docker.io`).
* `username`: Username for the registry.
* `password`: Password or token for the registry.

Stack configuration suits long-lived credentials, such as a Docker Hub or GitHub Container Registry access token. Set it with `--path`, and use `--secret` for the password so it is encrypted:

```bash
$ pulumi config set --path 'docker-build:registries[0].address' docker.io
$ pulumi config set --path 'docker-build:registries[0].username' my-user
$ pulumi config set --path --secret 'docker-build:registries[0].password' my-token
```

Additional registries can be added at subsequent indexes (`registries[1]`, `registries[2]`, and so on). A stored token like this has to be rotated by hand.

### Pushing to a cloud registry

For a cloud registry, don't store a password in configuration. Have the program fetch a short-lived push token from the cloud provider your stack is already logged in to, and pass it to the `Image` resource's `registries` input. That leaves nothing to store or rotate.

| Cloud | Registry | Push credentials |
| --- | --- | --- |
| Google Cloud | Artifact Registry `Repository` | Username `oauth2accesstoken`, with `getClientConfig().accessToken` as the password |
| AWS | `awsx.ecr.Repository` | `aws.ecr.getAuthorizationToken` |
| Azure | Container Registry `Registry` with the admin user enabled | `listRegistryCredentials` |

#### Google Cloud Artifact Registry (TypeScript)

```typescript
import * as dockerBuild from "@pulumi/docker-build";
import * as gcp from "@pulumi/gcp";
import * as pulumi from "@pulumi/pulumi";

const repository = new gcp.artifactregistry.Repository("images", {
    location: "us-central1",
    repositoryId: "my-images",
    format: "DOCKER",
});

// A short-lived OAuth token from the stack's Google Cloud login.
const clientConfig = gcp.organizations.getClientConfigOutput();

const image = new dockerBuild.Image("app", {
    context: { location: "./app" },
    tags: [pulumi.interpolate`${repository.registryUri}/app:latest`],
    push: true,
    registries: [{
        // registryUri is "<host>/<project>/<repository>"; log in to the host.
        address: repository.registryUri.apply(uri => uri.split("/")[0]),
        username: "oauth2accesstoken",
        password: pulumi.secret(clientConfig.accessToken),
    }],
});

export const imageRef = image.ref;
```

#### Amazon ECR (Python)

```python
import pulumi
import pulumi_aws as aws
import pulumi_awsx as awsx
import pulumi_docker_build as docker_build

repository = awsx.ecr.Repository("app", force_delete=True)

# A short-lived ECR login from the stack's AWS credentials.
token = aws.ecr.get_authorization_token_output(registry_id=repository.repository.registry_id)

image = docker_build.Image(
    "app",
    context=docker_build.BuildContextArgs(location="./app"),
    tags=[repository.url.apply(lambda url: f"{url}:latest")],
    push=True,
    registries=[
        docker_build.RegistryArgs(
            address=repository.url,
            username=token.user_name,
            password=pulumi.Output.secret(token.password),
        )
    ],
)

pulumi.export("imageRef", image.ref)
```

#### Azure Container Registry (C#)

```csharp
using System.Collections.Generic;
using Pulumi;
using Pulumi.AzureNative.ContainerRegistry;
using Pulumi.AzureNative.ContainerRegistry.Inputs;
using ResourceGroup = Pulumi.AzureNative.Resources.ResourceGroup;
using DockerBuild = Pulumi.DockerBuild;

return await Deployment.RunAsync(() =>
{
    var resourceGroup = new ResourceGroup("rg");

    // The admin user lets the push authenticate with a username and password.
    var registry = new Registry("registry", new RegistryArgs
    {
        ResourceGroupName = resourceGroup.Name,
        Sku = new SkuArgs { Name = SkuName.Basic },
        AdminUserEnabled = true,
    });

    var credentials = ListRegistryCredentials.Invoke(new ListRegistryCredentialsInvokeArgs
    {
        ResourceGroupName = resourceGroup.Name,
        RegistryName = registry.Name,
    });

    var image = new DockerBuild.Image("app", new DockerBuild.ImageArgs
    {
        Context = new DockerBuild.Inputs.BuildContextArgs { Location = "./app" },
        Tags = { Output.Format($"{registry.LoginServer}/app:latest") },
        Push = true,
        Registries =
        {
            new DockerBuild.Inputs.RegistryArgs
            {
                Address = registry.LoginServer,
                Username = credentials.Apply(c => c.Username!),
                Password = Output.CreateSecret(credentials.Apply(c => c.Passwords[0].Value!)),
            },
        },
    });

    return new Dictionary<string, object?>
    {
        ["imageRef"] = image.Ref,
    };
});
```
