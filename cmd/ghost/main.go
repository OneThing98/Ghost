package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/OneThing98/Ghost/internal/utils"
	"github.com/OneThing98/Ghost/pkg/network"
	libcontainer "github.com/OneThing98/containerpkg"
	"github.com/OneThing98/namespaces"
)

var (
	displayPid bool
	newCommand string
	userNet    bool
)

func init() {
	flag.BoolVar(&displayPid, "pid", false, "display the pid before waiting")
	flag.StringVar(&newCommand, "cmd", "/bin/bash", "command to run in the existing namespace")
	flag.BoolVar(&userNet, "net", false, "use a net namespace")
	flag.Parse()
}

func exec(container *libcontainer.Container) error {
	var (
		netFile *os.File
		err     error
	)
	container.NetNsFd = 0
	if userNet {
		netFile, err = os.Open("/root/nsroot/test")
		if err != nil {
			return err
		}
		container.NetNsFd = netFile.Fd()
	}
	pid, err := namespaces.ContainerExec(container)
	if err != nil {
		return fmt.Errorf("error exec container %s", err)
	}
	if displayPid {
		fmt.Println(pid)
	}

	exitcode, err := utils.WaitOnPid(pid)
	if err != nil {
		return fmt.Errorf("error waiting on child %s", err)
	}
	fmt.Println(exitcode)
	if userNet {
		netFile.Close()
		if err := network.DeleteNetworkNamespace("/root/nsroot/test"); err != nil {
			return err
		}
	}
	os.Exit(exitcode)
	return nil
}

func execIn(container *libcontainer.Container) error {
	f, err := os.Open("/root/nsroot/test")
	if err != nil {
		return nil
	}
	container.NetNsFd = f.Fd()
	pid, err := namespaces.ContainerExecIn(container, &libcontainer.Command{
		Env: container.Command.Env,
		Args: []string{
			newCommand,
		},
	})
	if err != nil {
		return fmt.Errorf("error execin container %s", err)
	}
	exitcode, err := utils.WaitOnPid(pid)
	if err != nil {
		return fmt.Errorf("error waiting on child %s", err)
	}
	os.Exit(exitcode)
	return nil
}

func createNet(config *libcontainer.Network) error {
	root := "/root/nsroot"
	if err := network.SetupNamespaceMountDir(root); err != nil {
		return err
	}
	nspath := root + "/test"
	if err := network.CreateNetworkNamespace(nspath); err != nil {
		return nil
	}
	if err := network.CreateVethPair("veth0", config.TempVethName); err != nil {
		return err
	}
	if err := network.SetInterfaceMaster("veth0", config.Bridge); err != nil {
		return err
	}
	if err := network.InterfaceUp("veth0"); err != nil {
		return err
	}

	f, err := os.Open(nspath)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := network.SetInterfaceInNamespaceFd("veth1", int(f.Fd())); err != nil {
		return err
	}
	return nil
}

func printErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
	os.Exit(1)
}

func main() {
	var (
		err    error
		cliCmd = flag.Arg(0)
		config = flag.Arg(1)
	)

	f, err := os.Open(config)
	if err != nil {
		printErr(err)
	}

	dec := json.NewDecoder(f)
	var container *libcontainer.Container

	if err := dec.Decode(&container); err != nil {
		printErr(nil)
	}
	f.Close()
	switch cliCmd {
	case "exec":
		err = exec(container)
	case "execin":
		err = execIn(container)
	case "net":
		err = createNet(&libcontainer.Network{
			TempVethName: "veth1",
			IP:           "172.17.0.100/16",
			Gateway:      "172.17.42.1",
			Mtu:          1500,
			Bridge:       "docker0",
		})
	default:
		err = fmt.Errorf("command not supported: %s", cliCmd)
	}
	if err != nil {
		printErr(err)
	}
}

// func main() {
//     if len(os.Args) < 3 {
//         log.Fatal("Usage: run <command> or child <command>")
//     }

//     runner, err := container.NewContainerRunner(os.Args[1])
//     if err != nil {
//         log.Fatalf("Error creating container runner: %v", err)
//     }

//     if err := runner.Run(os.Args[2:]); err != nil {
//         log.Fatalf("Error running command: %v", err)
//     }
// }
