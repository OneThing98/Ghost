package namespaces

import (
	"fmt"
	"os"
	"strings"

	container "github.com/OneThing98/containerpkg"
)

func writeError(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
	os.Exit(1)
}

func setupEnvironment(container *container.Container) {
	addEnvIfNotSet(container, "container", "docker")
	addEnvIfNotSet(container, "TERM", "xterm")
	addEnvIfNotSet(container, "USER", "root")
	addEnvIfNotSet(container, "LOGNAME", "root")
}

func addEnvIfNotSet(container *container.Container, key, value string) {
	jv := fmt.Sprintf("%s=%s", key, value)
	if len(container.Command.Env) == 0 {
		container.Command.Env = []string{jv}
		return
	}

	for _, v := range container.Command.Env {
		parts := strings.Split(v, "=")
		if parts[0] == key {
			return
		}
	}

	container.Command.Env = append(container.Command.Env, jv)

}
