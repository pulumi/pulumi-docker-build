package examples

import (
	"os"
	"path"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/require"
)

func TestPythonExample(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	test := integration.ProgramTestOptions{
		Dir:             path.Join(cwd, "python"),
		RelativeWorkDir: ".",
		Dependencies: []string{
			path.Join("..", "sdk", "python", "bin"),
		},
		Secrets: map[string]string{
			"dockerHubPassword": os.Getenv("DOCKER_HUB_PASSWORD"),
		},
	}

	integration.ProgramTest(t, &test)
}
