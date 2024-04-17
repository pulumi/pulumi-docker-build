package examples

import (
	"crypto/rsa"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/crypto/ssh/agent"
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
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		if err := agent.ServeAgent(a, conn); err != nil {
			panic(err)
		}
	}()

	return sock
}
