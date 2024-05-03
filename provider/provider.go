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
	"context"
	"encoding/json"

	gp "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	rpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	"github.com/pulumi/pulumi-docker-build/provider/internal"
	"github.com/pulumi/pulumi-docker-build/provider/internal/deprecated"
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
	return gp.RawServer(Name, Version, configurableProvider(internal.NewBuildxProvider()))(host)
}

// configurableProvider is a workaround for
// https://github.com/pulumi/pulumi-go-provider/issues/171 and
// In short, our SDKs send provider Configure requests as simple strings
// instead of rich objects. We don't want to preserve this behavior in
// pulumi-go-provider, but we also haven't updated SDKs yet to send rich types.
//
// If you find yourself in a position where you need to copy this -- STOP!
// https://github.com/pulumi/pulumi/pull/15032 should be merged with this fix.
func configurableProvider(p gp.Provider) gp.Provider {
	configure := p.Configure

	p.Configure = func(ctx context.Context, req gp.ConfigureRequest) error {
		r, err := p.GetSchema(ctx, gp.GetSchemaRequest{Version: 0})
		if err != nil {
			return err
		}
		spec := schema.PackageSpec{}
		err = json.Unmarshal([]byte(r.Schema), &spec)
		if err != nil {
			return err
		}

		ce := deprecated.New(spec.Config)
		if props, err := ce.UnmarshalProperties(req.Args); err == nil {
			req.Args = props
		}
		return configure(ctx, req)
	}

	return p
}
