package examples

import (
	"crypto/rsa"
	"errors"
	"io"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh/agent"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
)

func TestMain(m *testing.M) {
	sock := sshagent()
	os.Setenv("SSH_AUTH_SOCK", sock)

	os.Exit(m.Run())
}

// sshagent crates an in-memory SSH agent with one identity.
func sshagent() string {
	dir, err := os.MkdirTemp(os.TempDir(), "docker-test-*")
	if err != nil {
		panic(err)
	}

	sock := filepath.Join(dir, "test.sock")

	l, err := net.Listen("unix", sock)
	if err != nil {
		panic(err)
	}

	a := agent.NewKeyring()
	//nolint:gosec
	key, err := rsa.GenerateKey(rand.New(rand.NewSource(42)), 2048)
	if err != nil {
		panic(err)
	}
	err = a.Add(agent.AddedKey{PrivateKey: key})
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				panic(err)
			}
			if err := agent.ServeAgent(a, conn); err != nil && !errors.Is(err, io.EOF) {
				panic(err)
			}
		}
	}()

	return sock
}

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
		sig, _ := m["4dabf18193072939515e22adb298388d"]
		assert.Equal(t, "1b47061264138c4ac30d75fd1eb44270", sig,
			"secrets input should be marked as a Pulumi secret in state")
		return
	}
	t.Error("'secrets' Image resource not found in deployment")
}
