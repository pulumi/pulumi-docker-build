//go:build dotnet || go || nodejs || python || yaml || all
// +build dotnet go nodejs python yaml all

package examples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
)

// assertSecretsInputIsSecret checks that the "secrets" Image resource in a deployed
// stack has its "secrets" input property stored as a Pulumi secret in state.
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
