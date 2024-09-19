package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
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

func exec(container *libcontainer.Container) (int, error) {
	log.Println("Starting exec function")
	var (
		netFile *os.File
		err     error
	)
	container.NetNsFd = 0
	if userNet {
		log.Println("UserNet is enabled, opening network namespace file")
		netFile, err = os.Open("/root/nsroot/test")
		if err != nil {
			return 1, fmt.Errorf("failed to open network namespace file: %w", err)
		}
		container.NetNsFd = netFile.Fd()
	}
	log.Printf("Container Config: %+v\n", container)
	pid, err := namespaces.ContainerExec(container)
	if err != nil {
		return 1, fmt.Errorf("error executing container: %w", err)
	}
	if displayPid {
		fmt.Println(pid)
	}

	log.Println("Waiting on the process to complete")
	exitcode, err := utils.WaitOnPid(pid)
	if err != nil {
		return 1, fmt.Errorf("error waiting on child process: %w", err)
	}
	fmt.Println(exitcode)
	if userNet {
		log.Println("Closing network namespace file and deleting namespace")
		netFile.Close()
		if err := network.DeleteNetworkNamespace("/root/nsroot/test"); err != nil {
			return 1, fmt.Errorf("failed to delete network namespace: %w", err)
		}
	}
	log.Printf("Exiting with code: %d", exitcode)
	return exitcode, nil
}

func execIn(container *libcontainer.Container) (int, error) {
	log.Println("Starting execIn function")
	f, err := os.Open("/root/nsroot/test")
	if err != nil {
		return 1, fmt.Errorf("failed to open network namespace file: %w", err)
	}
	container.NetNsFd = f.Fd()
	pid, err := namespaces.ContainerExecIn(container, &libcontainer.Command{
		Env: container.Command.Env,
		Args: []string{
			newCommand,
		},
	})
	if err != nil {
		return 1, fmt.Errorf("error executing in container: %w", err)
	}
	log.Println("Waiting on the process to complete")
	exitcode, err := utils.WaitOnPid(pid)
	if err != nil {
		return 1, fmt.Errorf("error waiting on child process: %w", err)
	}
	log.Printf("Exiting with code: %d", exitcode)
	return exitcode, nil
}

func createNet(config *libcontainer.Network) (int, error) {
	log.Println("Starting createNet function")
	root := "/root/nsroot"
	if err := network.SetupNamespaceMountDir(root); err != nil {
		return 1, fmt.Errorf("failed to set up namespace mount directory: %w", err)
	}
	nspath := root + "/test"
	if err := network.CreateNetworkNamespace(nspath); err != nil {
		return 1, fmt.Errorf("failed to create network namespace: %w", err)
	}
	log.Println("Network namespace created, setting up veth pair")
	if err := network.CreateVethPair("veth0", config.TempVethName); err != nil {
		return 1, fmt.Errorf("failed to create veth pair: %w", err)
	}
	if err := network.SetInterfaceMaster("veth0", config.Bridge); err != nil {
		return 1, fmt.Errorf("failed to set interface master: %w", err)
	}
	if err := network.InterfaceUp("veth0"); err != nil {
		return 1, fmt.Errorf("failed to bring interface up: %w", err)
	}

	log.Println("Opening network namespace path")
	f, err := os.Open(nspath)
	if err != nil {
		return 1, fmt.Errorf("failed to open namespace path: %w", err)
	}
	defer f.Close()
	if err := network.SetInterfaceInNamespaceFd("veth1", int(f.Fd())); err != nil {
		return 1, fmt.Errorf("failed to set interface in namespace: %w", err)
	}
	log.Println("Network namespace setup complete")
	return 0, nil
}

func printErr(err error) {
	if err != nil {
		log.Printf("Error: %v", err)
	}
}

func main() {
	var (
		err      error
		cliCmd   = flag.Arg(0)
		config   = flag.Arg(1)
		exitCode int
	)

	log.Println("Opening container configuration file")
	f, err := os.Open(config)
	if err != nil {
		printErr(fmt.Errorf("failed to open container configuration file: %w", err))
		os.Exit(1)
	}

	log.Println("Decoding container configuration")
	dec := json.NewDecoder(f)
	var container *libcontainer.Container

	if err := dec.Decode(&container); err != nil {
		printErr(fmt.Errorf("failed to decode container configuration: %w", err))
		os.Exit(1)
	}
	f.Close()

	log.Printf("Command received: %s", cliCmd)
	switch cliCmd {
	case "exec":
		exitCode, err = exec(container)
	case "execin":
		exitCode, err = execIn(container)
	case "net":
		exitCode, err = createNet(&libcontainer.Network{
			TempVethName: "veth1",
			IP:           "172.17.0.100/16",
			Gateway:      "172.17.42.1",
			Mtu:          1500,
			Bridge:       "docker0",
		})
	default:
		err = fmt.Errorf("command not supported: %s", cliCmd)
		exitCode = 1
	}

	if err != nil {
		printErr(err)
	}

	log.Printf("Exiting with code: %d", exitCode)
	os.Exit(exitCode)
}
