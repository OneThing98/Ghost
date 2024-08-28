package container

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/OneThing98/Ghost/internal/utils"
	"github.com/OneThing98/Ghost/pkg/chroot"
	"github.com/OneThing98/Ghost/pkg/namespaces"
)

type ContainerRunner interface {
	Run(command []string) error
}

type childRunner struct{}

func (r *childRunner) Run(command []string) error {
	fmt.Printf("Running %v as child\n", command)

	newRoot := "/home/rojin/dev/jammy-base-amd64"
	fmt.Printf("Attempting to chroot to: %s\n", newRoot)

	if err := chroot.CreateChroot(newRoot); err != nil {
		return fmt.Errorf("failed to setup chroot: %v", err)
	}

	if err := utils.SetupHostname("container"); err != nil {
		return fmt.Errorf("failed to setup hostname: %v", err)
	}

	err := utils.RunContainerProcess(command)

	chroot.UnmountProc()

	return err
}

type mainRunner struct{}

func (r *mainRunner) Run(command []string) error {
	fmt.Printf("Running %v\n", command)

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, command...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = namespaces.SetupNamespaces()

	return cmd.Run()
}

func NewContainerRunner(mode string) (ContainerRunner, error) {
	switch mode {
	case "run":
		return &mainRunner{}, nil
	case "child":
		return &childRunner{}, nil
	default:
		return nil, fmt.Errorf("invalid mode: %s", mode)
	}
}
