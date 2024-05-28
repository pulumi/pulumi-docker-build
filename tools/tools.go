//go:build tools
// +build tools

// See https://play-with-go.dev/tools-as-dependencies_go119_en/ for an explanation of this file.
package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/pulumi/pulumi-dotnet/pulumi-language-dotnet"
	_ "github.com/pulumi/pulumi-java/pkg/cmd/pulumi-language-java"
	_ "github.com/pulumi/pulumi-yaml/cmd/pulumi-converter-yaml"
	_ "github.com/pulumi/pulumi-yaml/cmd/pulumi-language-yaml"
	_ "github.com/pulumi/pulumi/pkg/v3/cmd/pulumi"
	_ "github.com/pulumi/pulumi/sdk/go/pulumi-language-go/v3"
	_ "github.com/pulumi/pulumi/sdk/nodejs/cmd/pulumi-language-nodejs/v3"
	_ "github.com/pulumi/pulumi/sdk/python/cmd/pulumi-language-python/v3"
)
