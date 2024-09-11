package capabilities

import (
	"fmt"
	"os"

	"github.com/OneThing98/capability"
	container "github.com/OneThing98/containerpkg"
)

var CapabilityMap = map[container.Capability]capability.Cap{
	"CAP_SETPCAP":        capability.CAP_SETPCAP,
	"CAP_SYS_MODULE":     capability.CAP_SYS_MODULE,
	"CAP_SYS_RAWIO":      capability.CAP_SYS_RAWIO,
	"CAP_SYS_PACCT":      capability.CAP_SYS_PACCT,
	"CAP_SYS_ADMIN":      capability.CAP_SYS_ADMIN,
	"CAP_SYS_NICE":       capability.CAP_SYS_NICE,
	"CAP_SYS_RESOURCE":   capability.CAP_SYS_RESOURCE,
	"CAP_SYS_TIME":       capability.CAP_SYS_TIME,
	"CAP_SYS_TTY_CONFIG": capability.CAP_SYS_TTY_CONFIG,
	"CAP_MKNOD":          capability.CAP_MKNOD,
	"CAP_AUDIT_WRITE":    capability.CAP_AUDIT_WRITE,
	"CAP_AUDIT_CONTROL":  capability.CAP_AUDIT_CONTROL,
	"CAP_MAC_OVERRIDE":   capability.CAP_MAC_OVERRIDE,
	"CAP_MAC_ADMIN":      capability.CAP_MAC_ADMIN,
}

func DropCapabilities(container *container.Container) error {
	caps, err := getCapabilities(container)

	if err != nil {
		return err
	}

	c, err := capability.NewPid(os.Getpid())
	if err != nil {
		return err
	}

	c.Unset(capability.CAPS|capability.BOUNDS, caps...)

	if err := c.Apply(capability.CAPS | capability.BOUNDS); err != nil {
		return err
	}

	return nil
}

func getCapabilities(container *container.Container) ([]capability.Cap, error) {
	var caps []capability.Cap

	for _, name := range container.Capabilities {
		if cap, ok := CapabilityMap[name]; ok {
			caps = append(caps, cap)
		} else {
			return nil, fmt.Errorf("unknown capability: %s", name)
		}
	}
	return caps, nil
}
