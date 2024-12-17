[![Slack](http://www.pulumi.com/images/docs/badges/slack.svg)](https://slack.pulumi.com)
[![NPM version](https://badge.fury.io/js/%40pulumi%2fdocker-build.svg)](https://www.npmjs.com/package/@pulumi/docker-build)
[![Python version](https://badge.fury.io/py/pulumi-docker-build.svg)](https://pypi.org/project/pulumi-docker-build)
[![NuGet version](https://badge.fury.io/nu/pulumi.dockerbuild.svg)](https://badge.fury.io/nu/pulumi.dockerbuild)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/pulumi/pulumi-docker-build/sdk/go)](https://pkg.go.dev/github.com/pulumi/pulumi-docker-build/sdk/go)
[![License](https://img.shields.io/npm/l/%40pulumi%2Fpulumi.svg)](https://github.com/pulumi/pulumi-docker-build/blob/main/LICENSE)

# Docker-Build Resource Provider

A [Pulumi](http://pulumi.com) provider for building modern Docker images with [buildx](https://docs.docker.com/build/architecture/) and [BuildKit](https://docs.docker.com/build/buildkit/).

Not to be confused with the earlier
[Docker](http://github.com/pulumi/pulumi-docker) provider, which is still
appropriate for managing resources unrelated to building images.

| Provider               | Use cases                                                                                                                                                |
| ----------------       | -------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `@pulumi/docker-build` | Anything related to building images with `docker build`.                                                                                                 |
| `@pulumi/docker`       | Everything else -- including running containers and creating networks.                                                                                   |
