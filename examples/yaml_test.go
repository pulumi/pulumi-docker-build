package examples

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/pulumi/providertest"
	"github.com/pulumi/pulumi-docker-build/provider"
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

func TestYAMLExampleUpgrade(t *testing.T) {
	// t.Setenv("PULUMI_PROVIDER_TEST_MODE", "snapshot")

	cwd, err := os.Getwd()
	require.NoError(t, err)

	bin, err := filepath.Abs("../bin")
	require.NoError(t, err)

	t.Setenv("PATH", bin+":"+os.Getenv("PATH"))
	p, err := provider.New(nil)
	require.NoError(t, err)

	test := providertest.NewProviderTest(path.Join(cwd, "upgrade"),
		providertest.WithProviderName("docker-build"),
		providertest.WithBaselineVersion("0.0.1"),
		providertest.WithResourceProviderServer(p),
		// providertest.WithConfig("dockerHubPassword", os.Getenv("DOCKER_HUB_PASSWORD")), // Doesn't support secrets yet.
	)
	test.Run(t)
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

func TestDockerHub(t *testing.T) {
	if os.Getenv("DOCKER_HUB_PASSWORD") == "" {
		t.Skip("Missing DockerHub credentials")
	}

	cwd, err := os.Getwd()
	require.NoError(t, err)

	test := integration.ProgramTestOptions{
		Dir: path.Join(cwd, "tests/dockerhub"),
		Secrets: map[string]string{
			"dockerHubPassword": os.Getenv("DOCKER_HUB_PASSWORD"),
		},
	}

	integration.ProgramTest(t, &test)
}

func TestDockerHubUnauthenticated(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	test := integration.ProgramTestOptions{
		Dir: path.Join(cwd, "tests/unauthenticated"),
	}

	integration.ProgramTest(t, &test)
}
