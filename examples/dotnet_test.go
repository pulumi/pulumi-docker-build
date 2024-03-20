//go:build dotnet || all
// +build dotnet all

package examples

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/require"
)

func TestDotNetExample(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	nuget := filepath.Join(cwd, "../nuget")
	t.Setenv("PULUMI_LOCAL_NUGET", nuget)

	cmd := exec.Command("dotnet", "nuget", "add", "source", nuget)
	_ = cmd.Run()

	test := integration.ProgramTestOptions{
		Dir: path.Join(cwd, "dotnet"),
		Dependencies: []string{
			"Pulumi.Dockerbuild",
		},
		Secrets: map[string]string{
			"dockerHubPassword": os.Getenv("DOCKER_HUB_PASSWORD"),
		},
		NoParallel: true,
	}

	integration.ProgramTest(t, &test)
}
