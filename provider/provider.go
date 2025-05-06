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

package provider

import (
	"github.com/pulumi/pulumi-docker-build/provider/internal"
	gp "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	rpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

// Version is initialized by the Go linker to contain the semver of this build.
var Version = "0.0.1"

// Name needs to match $PACK in Makefile.
const Name string = "docker-build"

// Serve launches the gRPC server for the resource provider.
func Serve() error {
	return provider.Main(Name, New)
}

// New creates a new provider.
func New(host *provider.HostClient) (rpc.ResourceProviderServer, error) {
	return gp.RawServer(Name, Version, internal.NewBuildxProvider())(host)
}
