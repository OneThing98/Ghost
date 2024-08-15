package main

import (
    "log"
    "os"

    "github.com/OneThing98/Ghost/pkg/container"
)

func main() {
    if len(os.Args) < 3 {
        log.Fatal("Usage: run <command> or child <command>")
    }

    runner, err := container.NewContainerRunner(os.Args[1])
    if err != nil {
        log.Fatalf("Error creating container runner: %v", err)
    }

    if err := runner.Run(os.Args[2:]); err != nil {
        log.Fatalf("Error running command: %v", err)
    }
}
