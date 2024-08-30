package specs

import (
	"encoding/json"
	"fmt"
	"os"
)

// OCI runtime specification configuration
type Config struct {
	ID           string       `json:"id,omitempty"`
	NamespacePid int          `json:"namespace_pid,omitempty"`
	Command      Command      `json:"command,omitempty"`
	RootFs       string       `json:"rootfs,omitempty"`
	ReadonlyFs   bool         `json:"readonly_fs,omitempty"`
	NetNsFd      uintptr      `json:"network_namespace_fd,omitempty"`
	User         string       `json:"user,omitempty"`
	WorkingDir   string       `json:"working_dir,omitempty"`
	Namespaces   []Namespace  `json:"namespaces,omitempty"`
	Capabilities []Capability `json:"capabilities,omitempty"`
	CgroupsPath  string       `json:"cgroups_path,omitempty"`
	Network      *Network     `json:"network,omitempty"`
}

// commands to run insde the container
type Command struct {
	Args []string `json:"args,omitempty"`
	Env  []string `json:"environment,omitempty"`
}

// linux namespaces that should be created or joined
type Namespace struct {
	Type string `json:"type,omitempty"` // Example types: "pid", "network", "ipc", etc.
	Path string `json:"path,omitempty"` // Path to bind mount for namespace
}

// POSIX capability that can be applied to a process
type Capability string

// network configurations
type Network struct {
	TempVethName string `json:"temp_veth,omitempty"`
	IP           string `json:"ip,omitempty"`
	Gateway      string `json:"gateway,omitempty"`
	Bridge       string `json:"bridge,omitempty"`
	Mtu          int    `json:"mtu,omitempty"`
}

//reads json config file and unmarshalls it into config struct

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failure to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("Failed to decode config json:%w", err)
	}
	return &config, nil
}

func (c *Config) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("ID is required")
	}

	if c.RootFs == "" {
		return fmt.Errorf("Rootfs path is required")
	}

	if len(c.Command.Args) == 0 {
		return fmt.Errorf("Command arguments are required")
	}

	if len(c.Namespaces) == 0 {
		return fmt.Errorf("At least one namespace is required")
	}

	if len(c.Capabilities) == 0 {
		return fmt.Errorf("At least one capability is required")
	}

	if c.Network != nil {
		if c.Network.IP == "" {
			return fmt.Errorf("Network IP is required while configuring network")
		}

		if c.Network.Bridge == "" {
			return fmt.Errorf("Network Bridge is required while configuring network")
		}
	}

	return nil

}
