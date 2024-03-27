//go:build nodejs || python || dotnet || java || all
// +build nodejs python dotnet java all

package examples

import (
	"net"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/crypto/ssh/agent"
)

func TestMain(m *testing.M) {
	if os.Getenv("SSH_AUTH_SOCK") == "" {
		sock := sshagent()
		os.Setenv("SSH_AUTH_SOCK", sock)
	}

	os.Exit(m.Run())
}

func sshagent() string {
	dir := os.TempDir()
	sock := filepath.Join(dir, "test.sock")

	_ = os.Remove(sock)

	l, err := net.Listen("unix", sock)
	if err != nil {
		panic(err)
	}

	go func() {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		a := agent.NewKeyring()
		agent.ServeAgent(a, conn)
	}()

	return sock
}
