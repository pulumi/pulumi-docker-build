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

package internal

import (
	"context"
	"fmt"

	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	pschema "github.com/pulumi/pulumi-go-provider/middleware/schema"
	"github.com/pulumi/pulumi-java/pkg/codegen/java"
	csgen "github.com/pulumi/pulumi/pkg/v3/codegen/dotnet"
	gogen "github.com/pulumi/pulumi/pkg/v3/codegen/go"
	tsgen "github.com/pulumi/pulumi/pkg/v3/codegen/nodejs"
	pygen "github.com/pulumi/pulumi/pkg/v3/codegen/python"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"github.com/pulumi/pulumi/sdk/v3/go/property"
)

var (
	_ infer.CustomConfigure = (*Config)(nil)
	_ infer.Annotated       = (*Config)(nil)
	_ infer.Annotated       = (*Registry)(nil)
)

// Config configures the buildx provider.
type Config struct {
	Host       string     `pulumi:"host,optional"`
	Registries []Registry `pulumi:"registries,optional"`

	host *host
}

// _mockClientKey is used by tests to inject a mock Docker client.
var _mockClientKey any = "mock-client"

// Annotate provides user-facing descriptions and defaults for Config's fields.
func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.Host, "The build daemon's address.")
	a.SetDefault(&c.Host, "", "DOCKER_HOST")
}

// Configure validates and processes user-provided configuration values.
func (c *Config) Configure(ctx context.Context) error {
	h, err := newHost(ctx, c)
	if err != nil {
		return fmt.Errorf("getting host: %w", err)
	}
	c.host = h
	return nil
}

// NewBuildxProvider returns a new buildx provider.
func NewBuildxProvider() provider.Provider {
	prov := infer.Provider(
		infer.Options{
			Metadata: pschema.Metadata{
				DisplayName: "docker-build",
				Keywords:    []string{"docker", "buildkit", "buildx", "kind/native"},
				Description: "A Pulumi provider for building modern Docker images with buildx and BuildKit.",
				Homepage:    "https://pulumi.com",
				Publisher:   "Pulumi",
				License:     "Apache-2.0",
				Repository:  "https://github.com/pulumi/pulumi-docker-build",
				LanguageMap: map[string]any{
					"go": gogen.GoPackageInfo{
						// GenerateResourceContainerTypes: true,
						RespectSchemaVersion: true,
						Generics:             gogen.GenericsSettingSideBySide,
						PackageImportAliases: map[string]string{
							"github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild": "dockerbuild",
						},
						ImportBasePath: "github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild",
					},
					"csharp": csgen.CSharpPackageInfo{
						RespectSchemaVersion: true,
						PackageReferences: map[string]string{
							"Pulumi": "3.*",
						},
					},
					"java": java.PackageInfo{
						BuildFiles:                      "gradle",
						GradleNexusPublishPluginVersion: "1.1.0",
						Dependencies: map[string]string{
							"com.pulumi:pulumi":               "0.20.0",
							"com.google.code.gson:gson":       "2.8.9",
							"com.google.code.findbugs:jsr305": "3.0.2",
						},
					},
					"nodejs": tsgen.NodePackageInfo{
						RespectSchemaVersion: true,
					},
					"python": pygen.PackageInfo{
						RespectSchemaVersion: true,
						PyProject: struct {
							Enabled bool `json:"enabled,omitempty"`
						}{Enabled: true},
					},
				},
			},
			Resources: []infer.InferredResource{
				infer.Resource[*Image](),
				infer.Resource[*Index](),
			},
			ModuleMap: map[tokens.ModuleName]tokens.ModuleName{
				"internal": "index",
			},
			Config: infer.Config[*Config](),
		},
	)

	prov.DiffConfig = diffConfigIgnoreInternal(prov.DiffConfig)

	return prov
}

// TODO(pulumi/pulumi-docker-build#404): Remove this function once the bug is fixed in either
// upstream pu/pu or pulumi-go-provider.

// diffConfigInternalIgnore is a custom DiffConfig implementation for the buildx provider. This is required to
// circumvent the bug identified in https://github.com/pulumi/pulumi-docker-build/issues/404.
// Since `__internal` is currently populated in new inputs, but stripped in old state, we need to
// ignore this field in the diff. There is no easy way to override DiffConfig to compare inputs only.
func diffConfigIgnoreInternal(
	diffConfig func(ctx context.Context, req provider.DiffRequest) (provider.DiffResponse, error),
) func(ctx context.Context, req provider.DiffRequest) (provider.DiffResponse, error) {
	return func(ctx context.Context, req provider.DiffRequest) (provider.DiffResponse, error) {
		m := req.News.AsMap()
		delete(m, "__internal")
		delete(m, "__pulumi-go-provider-infer")
		delete(m, "__pulumi-go-provider-version")
		req.News = property.NewMap(m)

		return diffConfig(ctx, req)
	}
}

// Schema returns our package specification.
func Schema(ctx context.Context, version string) schema.PackageSpec {
	p := NewBuildxProvider()
	spec, err := provider.GetSchema(ctx, "docker-build", version, p)
	contract.AssertNoErrorf(err, "missing schema")
	return spec
}
