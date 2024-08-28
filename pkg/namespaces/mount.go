package namespaces

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func SetUpNewMountNameSpace(rootfs string, readonly bool) error {

	//ensures mount events do not propagate outwards to other namespaces but propagate inwards
	if err := unix.Mount("", "/", "", unix.MS_SLAVE|unix.MS_REC, ""); err != nil {
		return fmt.Errorf("mounting / as slave: %v", err)
	}

	//bind mounts the new root filesystem onto itself to prepare for the new environment

	if err := unix.Mount(rootfs, rootfs, "bind", unix.MS_BIND|unix.MS_REC, ""); err != nil {
		return fmt.Errorf("binding mount: %s, %v", rootfs, err)
	}

	//if the read only flag is true then new root filesystem remounts as read-only
	if readonly {
		if err := unix.Mount(rootfs, rootfs, "bind", unix.MS_BIND|unix.MS_REMOUNT|unix.MS_RDONLY|unix.MS_REC, ""); err != nil {
			return fmt.Errorf("remounting %s as readonly: %v", rootfs, err)
		}
	}

	//custom wrapper around pivotRoot system call
	if err := pivotRoot(rootfs, filepath.Join(rootfs, ".old_root")); err != nil {
		return fmt.Errorf("pivot_root %s: %v", rootfs, err)
	}

	//unmount old root

	if err := unix.Unmount("/.old_root", unix.MNT_DETACH); err != nil {
		return fmt.Errorf("unmounting old root: %v", err)
	}

	return os.Chdir("/")

}
