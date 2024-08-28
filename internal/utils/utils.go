package utils

import (
	"fmt"
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
