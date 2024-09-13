package network

import (
	"fmt"
	"os"

	container "github.com/OneThing98/containerpkg"
	"github.com/OneThing98/namespaces"
	"golang.org/x/sys/unix"
)

func SetupVeth(config *container.Network) error {
	if err := InterfaceDown(config.TempVethName); err != nil {
		return fmt.Errorf("interface down %s %s", config.TempVethName, err)
	}
	if err := ChangeInterfaceName(config.TempVethName, "eth0"); err != nil {
		return fmt.Errorf("change %s to eth0 %s", config.TempVethName, err)
	}
	if err := SetInterfaceIp("eth0", config.IP); err != nil {
		return fmt.Errorf("set eth0 ip %s", err)
	}

	if err := SetMtu("eth0", config.Mtu); err != nil {
		return fmt.Errorf("set eth0 mtu to %d %s", config.Mtu, err)
	}
	if err := InterfaceUp("eth0"); err != nil {
		return fmt.Errorf("eth0 up %s", err)
	}

	if err := SetMtu("lo", config.Mtu); err != nil {
		return fmt.Errorf("set lo mtu to %d %s", config.Mtu, err)
	}
	if err := InterfaceUp("lo"); err != nil {
		return fmt.Errorf("lo up %s", err)
	}

	if config.Gateway != "" {
		if err := SetDefaultGateway(config.Gateway); err != nil {
			return fmt.Errorf("set gateway to %s %s", config.Gateway, err)
		}
	}
	return nil
}

func SetupNamespaceMountDir(root string) error {
	if err := os.MkdirAll(root, 0666); err != nil {
		return err
	}
	if err := unix.Mount("", root, "none", unix.MS_SHARED|unix.MS_REC, ""); err != nil && err != unix.EINVAL {
		return err
	}
	if err := unix.Mount(root, root, "none", unix.MS_BIND, ""); err != nil {
		return err
	}
	return nil
}

func CreateNetworkNamespace(bindingPath string) error {
	f, err := os.OpenFile(bindingPath, os.O_RDONLY|os.O_CREATE|os.O_EXCL, 0)
	if err != nil {
		return err
	}
	f.Close()

	if err := namespaces.CreateNewNamespace(container.CLONE_NEWNET, bindingPath); err != nil {
		return err
	}
	return nil
}

func DeleteNetworkNamespace(bindingPath string) error {
	if err := unix.Unmount(bindingPath, 0); err != nil {
		return err
	}
	return os.Remove(bindingPath)
}
