//go:build java || all
// +build java all

package examples

import (
	"os"
	"path"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/require"
)

func TestYAMLExample(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	test := integration.ProgramTestOptions{
		Dir: path.Join(cwd, "yaml"),
		Secrets: map[string]string{
			"dockerHubPassword": os.Getenv("DOCKER_HUB_PASSWORD"),
		},
	}

	integration.ProgramTest(t, &test)
}

func TestECR(t *testing.T) {
	if os.Getenv("AWS_SESSION_TOKEN") == "" {
		t.Skip("Missing AWS credentials")
	}

	cwd, err := os.Getwd()
	require.NoError(t, err)

	test := integration.ProgramTestOptions{
		Dir: path.Join(cwd, "tests/ecr"),
	}

	integration.ProgramTest(t, &test)
}
