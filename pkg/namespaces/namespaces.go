package namespaces

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/unix"
)

func SetupNamespaces() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
	}
}

func CreateNewNamespace(namespace string, bindTo string) error {
	flag := namespaceMap[namespace]
	name := namespaceFileMap[namespace]
	nspath := filepath.Join("/proc/self/ns", name)

	pid, err := fork()
	if err != nil {
		return err
	}

	if pid == 0 {
		if err := unshare(flag); err != nil {
			return fmt.Errorf("unshare %s: %v", namespace, err)
		}
		if err := unix.Mount(nspath, bindTo, "none", unix.MS_BIND, ""); err != nil {
			return fmt.Errorf("bind mount %s: %v", nspath, err)
		}
		os.Exit(0)
	}

	_, err = unix.Wait4(pid, nil, 0, nil)
	return err
}

func joinExistingNamespace(fd uintptr, namespace string) error {
	flag := namespaceMap[namespace]
	return setns(fd, uintptr(flag))
}
