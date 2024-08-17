package examples

import (
	"os"
	"path"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/require"
)

func TestJavaExample(t *testing.T) {
	t.Skip("not working yet")

	cwd, err := os.Getwd()
	require.NoError(t, err)

	test := integration.ProgramTestOptions{
		Dir: path.Join(cwd, "java"),
		Secrets: map[string]string{
			"dockerHubPassword": os.Getenv("DOCKER_HUB_PASSWORD"),
		},
	}

	integration.ProgramTest(t, &test)
}
