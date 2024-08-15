package namespaces

import (
    "syscall"
)

func SetupNamespaces() *syscall.SysProcAttr {
    return &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
    }
}
