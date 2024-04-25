// Copyright 2024, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package internal contains our clients, validation, and provider
// implementation for interacting with Docker's buildx APIs.
//
// The provider has two primary modes of operation when building an image. The
// default behavior is to use an embedded Docker CLI, which does not require to
// actually be installed on a host in order to perform builds (a build daemon
// must still be accessible locally or remotely). The second mode execs a
// "docker-buildx" binary on the host to perform builds. This second mode was
// added primarily for compatibility with Docker Build Cloud, which requires a
// custom docker-buildx binary.
//
// # CLIs
//
// In both execution modes we have several CLI clients. The first client is
// scoped to the host and initialized as part of the provider's Configure call.
// We use this CLI for host-level operations, in particular when we potentially
// initialize a new buildx builder and when we fetch existing credentials on
// the host.
//
// Each operation then has a CLI instance scoped to the life of the operation.
// This allows us to layer resource-scoped credentials on top of the host's
// existing credentials, and in practice Docker seems to handle these
// connections more reliably than a single CLI for all operations.
//
// # Credentials
//
// When using the embedded Docker client, secrets are communicated to the build
// daemon natively via gRPC callbacks. When running in exec mode, credentials
// must be communicated to the buildx binary via a configuration file. In order
// to not pollute the host's existing credentials with e.g. short-lived ECR
// tokens, we copy a small subset of the host's Docker config to a temporary
// directory and use that for the lifetime of the exec operation.
//
// # Preview mode
//
// The pulumi-go-provider primarily operates on simple Go structs and doesn't
// currently have a way to distinguish whether a value is unknown or empty. We
// ignore anything that is a zero value during previews before we apply
// validation or perform builds.
//
// # Diffs
//
// Another limitation of pulumi-go-provider is that it doesn't currently allow
// us to override the default Diff behavior. We intentionally apply
// "ignoreChanges" semantics to registry passwords, in order to reduce noise
// and unnecessary updates, but as a result we have to re-implement Diff from
// the ground up. This implementation is not nearly as rich as the default
// experience and should be replaced when an alternative is available.
package internal
