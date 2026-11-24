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

Set it in your stack's configuration with `--path`, and use `--secret` for the password so it is encrypted:

```bash
$ pulumi config set --path 'docker-build:registries[0].address' docker.io
$ pulumi config set --path 'docker-build:registries[0].username' my-user
$ pulumi config set --path --secret 'docker-build:registries[0].password' my-token
```

Additional registries can be added at subsequent indexes (`registries[1]`, `registries[2]`, and so on).
