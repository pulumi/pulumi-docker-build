//go:build yaml || all
// +build yaml all

package examples

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pulumi/providertest"
	"github.com/pulumi/providertest/providers"
	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/assertpreview"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	"github.com/pulumi/pulumi-docker-build/provider"
)

func TestYAMLExample(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	test := integration.ProgramTestOptions{
		Dir: path.Join(cwd, "yaml"),
		Secrets: map[string]string{
			"dockerHubPassword": os.Getenv("DOCKER_HUB_PASSWORD"),
		},
		ExtraRuntimeValidation: assertSecretsInputIsSecret,
	}

	integration.ProgramTest(t, &test)
}

func assertSecretsInputIsSecret(t *testing.T, stack integration.RuntimeValidationStackInfo) {
	t.Helper()
	for _, res := range stack.Deployment.Resources {
		if res.Type != "docker-build:index:Image" || res.URN.Name() != "secrets" {
			continue
		}
		secretsInput, ok := res.Inputs["secrets"]
		if !assert.True(t, ok, "secrets input not found on 'secrets' Image resource") {
			return
		}
		m, ok := secretsInput.(map[string]any)
		if !assert.True(t, ok, "secrets input is not a map") {
			return
		}
		sig := m["4dabf18193072939515e22adb298388d"]
		assert.Equal(t, "1b47061264138c4ac30d75fd1eb44270", sig,
			"secrets input should be marked as a Pulumi secret in state")
		return
	}
	t.Error("'secrets' Image resource not found in deployment")
}

func TestHCLExample(t *testing.T) {
	t.Skip("Skipping until HCL is stable")
	cwd, err := os.Getwd()
	require.NoError(t, err)

	test := integration.ProgramTestOptions{
		Dir: path.Join(cwd, "hcl"),
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
