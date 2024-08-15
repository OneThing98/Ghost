package chroot

import (
	"fmt"
	"os"
	"syscall"
)

func CreateChroot(newRoot string) error {
	if err := syscall.Chroot(newRoot); err != nil {
		return fmt.Errorf("Failed to chroot to %s: %v", newRoot, err)
	}

	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("Failed to change directory to root: %v", err)
	}

	procDir := "/proc"
	if err := syscall.Mount("proc", procDir, "proc", 0, ""); err != nil {
		return fmt.Errorf("Failed to mount /proc: %v", err)
	}

	return nil
}

func UnmountProc() {
    procDir := "/proc"
    if err := syscall.Unmount(procDir, 0); err != nil {
        fmt.Printf("Warning: Failed to unmount /proc: %v\n", err)
    } else {
        fmt.Println("/proc unmounted successfully")
    }
}
