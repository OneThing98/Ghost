package main

import (
	"flag"
	"os"

	container "github.com/OneThing98/containerpkg"
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

func exec(container *container.Container) error {
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
