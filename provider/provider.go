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
	"fmt"

	"github.com/pulumi/pulumi-docker-build/provider/internal"
	"github.com/pulumi/pulumi-docker-build/provider/internal/deprecated"
	gp "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
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
	server, err := gp.RawServer(Name, Version, internal.NewBuildxProvider())(host)
	if err != nil {
		return nil, fmt.Errorf("building raw server: %w", err)
	}
	return &configurableProvider{ResourceProviderServer: server}, nil
}

// configurableProvider is a workaround for
// https://github.com/pulumi/pulumi-go-provider/issues/171 and
// In short, our SDKs send provider Configure requests as simple strings
// instead of rich objects. We don't want to preserve this behavior in
// pulumi-go-provider, but we also haven't updated SDKs yet to send rich types.
//
// If you find yourself in a position where you need to copy this -- STOP!
// https://github.com/pulumi/pulumi/pull/15032 should be merged with this fix.
type configurableProvider struct {
	rpc.ResourceProviderServer
}

func (p configurableProvider) Configure(
	ctx context.Context,
	request *rpc.ConfigureRequest,
) (*rpc.ConfigureResponse, error) {
	schema := internal.Schema(ctx, Version)
	ce := deprecated.New(schema.Config)
	buildxReq := request
	if props, err := ce.UnmarshalProperties(request.Args); err == nil {
		args, _ := plugin.MarshalProperties(props, plugin.MarshalOptions{
			Label:        "config",
			KeepUnknowns: true,
			SkipNulls:    true,
			KeepSecrets:  true,
			RejectAssets: true,
		})
		buildxReq.Args = args
	}
	return p.ResourceProviderServer.Configure(ctx, buildxReq)
}
