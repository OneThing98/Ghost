package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/OneThing98/Ghost/pkg/namespaces"
)

func RunContainerProcess(command []string) error {
	cmdParts := strings.Fields(strings.Join(command, " "))
	if len(cmdParts) == 0 {
		return logError("no command provided")
	}

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.SysProcAttr = namespaces.SetupNamespaces()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return logError("failed to start process: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return logError("process finished with error: %v", err)
	}

	return nil
}

func SetupHostname(hostname string) error {
	return syscall.Sethostname([]byte(hostname))
}

func logError(format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	log.Print(err)
	return err
}

func WaitOnPid(pid int) (exitcode int, err error) {
	child, err := os.FindProcess(pid)
	if err != nil {
		return -1, err
	}
	state, err := child.Wait()
	if err != nil {
		return -1, err
	}
	return getExitCode(state), nil
}

func getExitCode(state *os.ProcessState) int {
	return state.Sys().(syscall.WaitStatus).ExitStatus()
}

func GenerateRandomName(size int) (string, error) {
	id := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		return "", err
	}
	return hex.EncodeToString(id), nil
}
