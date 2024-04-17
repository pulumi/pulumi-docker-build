package examples

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"io/fs"
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
	dir := os.TempDir()
	sock := filepath.Join(dir, "test.sock")

	// In case it already exists.
	err := os.Remove(sock)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		panic(err)
	}

	l, err := net.Listen("unix", sock)
	if err != nil {
		panic(err)
	}

	a := agent.NewKeyring()
	key, err := rsa.GenerateKey(rand.Reader, 4096)
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
		agent.ServeAgent(a, conn)
	}()

	return sock
}
