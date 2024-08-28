package namespaces

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/unix"
)

// forks a new process from the system
func fork() (int, error) {
	syscall.ForkLock.Lock()
	pid, _, errno := unix.RawSyscall(unix.SYS_FORK, 0, 0, 0)
	syscall.ForkLock.Unlock()
	if errno != 0 {
		return 0, fmt.Errorf("syscall fork: %v", errno)
	}

	return int(pid), nil
}

// creates a new namespace
func unshare(flags int) error {
	_, _, errno := unix.RawSyscall(unix.SYS_UNSHARE, uintptr(flags), 0, 0)
	if errno != 0 {
		return fmt.Errorf("syscall unshare: %v", errno)
	}
	return nil
}

// joins an existing namespace
func setns(fd, nstype uintptr) error {
	_, _, errno := unix.RawSyscall(unix.SYS_SETNS, fd, nstype, 0)
	if errno != 0 {
		return fmt.Errorf("syscall setns: %v", errno)
	}
	return nil
}

// changes root
func pivotRoot(newRoot, putold string) error {
	return unix.PivotRoot(newRoot, putold)
}
