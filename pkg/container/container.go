package container

type Container struct {
	ID           string       `json:"id,omitempty"`                   // unique identifier for the container
	NSPid        int          `json:"namespace_pid,omitempty"`        // process id in the container's namespace
	Command      *Command     `json:"command,omitempty"`              //command to run inside the container
	RootFs       string       `json:"rootfs,omitempty"`               //path to the root file system
	ReadonlyFs   bool         `json:"readonly_fs,omitempty"`          //whether the the file system is read only
	NetNsFd      uintptr      `json:"network_namespace_fd,omitempty"` //file descriptor for network namespace
	User         string       `json:"user,omitempty"`                 //user to run as inside the container
	WorkingDir   string       `json:"working_dir,omitempty"`          //working directory inside the container
	Namespaces   Namespaces   `json:"namespaces,omitempty"`           //namespace to apply to the container
	Capabilities Capabilities `json:"capabilities,omitempty"`         //capabilities to apply to the container
	CgroupsPath  string       `json:"cgroups_path,omitempty"`         //Path to cgroups settings
	Network      *Network     `json:"network,omitempty"`              //Network settings for container
}

type Command struct {
	Args []string `json:"args,omitempty"`
	Env  []string `json:"environment,omitempty"`
}

type Network struct {
	TempVethName string `json:"temp_veth,omitempty"` //temporary veth pair name
	IP           string `json:"ip,omitempty"`        //IP address to assign the container
	Gateway      string `json:"gateway,omitempty"`   //Network gateway for the container
	Bridge       string `json:"bridge,omitempty"`    //Network bridge name
	Mtu          int    `json:"mtu,omitempty"`       //MTU size for the network interface

}

type Namespace string

type Namespaces []Namespace

type Capability string

type Capabilities []Capability

//methods to check if a namespace or capability is contained inside their respective slices.

func (n Namespaces) Contains(ns Namespace) bool {
	for _, nns := range n {
		if nns == ns {
			return true
		}
	}
	return false
}

func (c Capabilities) Contains(cc Capability) bool {
	for _, ccs := range c {
		if ccs == cc {
			return true
		}
	}
	return false
}
