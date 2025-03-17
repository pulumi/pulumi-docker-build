package examples

import (
	"os"
	"path"
	"testing"

	"github.com/pulumi/providertest"
	"github.com/pulumi/providertest/providers"
	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/assertpreview"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/pulumi/pulumi-docker-build/provider"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
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
	pt := pulumitest.NewPulumiTest(t, "upgrade",
		opttest.AttachProviderServer("docker-build", providerServerFactory))
	previewResult := providertest.PreviewProviderUpgrade(t, pt, "docker-build", "0.0.1")

	assertpreview.HasNoChanges(t, previewResult)
}

func providerServerFactory(pt providers.PulumiTest) (pulumirpc.ResourceProviderServer, error) {
	return provider.New(nil)
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
