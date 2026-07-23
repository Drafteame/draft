package connect

import "time"

// ConnConfig is the top-level structure of the YAML config file (~/.draft/dbconnect.yml).
type ConnConfig struct {
	Defaults     map[string]DefaultConfig  `yaml:"defaults"`
	Environments map[string]EnvConfig      `yaml:"environments"`
	Connections  map[string]ConnTypeConfig `yaml:"connections"`
}

// DefaultConfig holds per-type defaults (e.g. remote port).
type DefaultConfig struct {
	RemotePort int `yaml:"remote_port"`
}

// EnvConfig holds the bastion and cluster suffixes for one environment.
type EnvConfig struct {
	Bastion  BastionConfig  `yaml:"bastion"`
	Clusters ClustersConfig `yaml:"clusters"`
}

// BastionConfig holds the SSM bastion access parameters.
type BastionConfig struct {
	Target  string `yaml:"target"`
	Profile string `yaml:"profile"`
	Region  string `yaml:"region"`
}

// ClustersConfig holds the host-suffix fragments for each engine in one environment.
type ClustersConfig struct {
	RDS   string `yaml:"rds"`
	Cache string `yaml:"cache"`
	DocDB string `yaml:"docdb"`
}

// ConnTypeConfig holds the instance list for one DB type.
type ConnTypeConfig struct {
	Instances []ServiceConfig `yaml:"instances"`
}

// ServiceConfig is a single instance entry inside a DB type.
// LocalPorts maps environment name → local port, e.g. {"dev": 56000, "prod": 56011}.
type ServiceConfig struct {
	Name       string         `yaml:"name"`
	LocalPorts map[string]int `yaml:"local_ports"`
}

// ResolvedConnection contains all values needed to open an SSM tunnel,
// computed at runtime from (dbType, name) where name = "{service}-{env}".
type ResolvedConnection struct {
	DBType     string
	Name       string // full CLI name, e.g. "turbo-dev"
	Service    string // e.g. "turbo"
	Env        string // "dev" or "prod"
	Bastion    BastionConfig
	Host       string
	RemotePort int
	LocalPort  int
}

// RuntimeState holds the persisted state of active tunnels.
type RuntimeState struct {
	Connections map[string]RuntimeEntry `json:"connections"`
}

// RuntimeEntry holds runtime metadata for an active tunnel.
type RuntimeEntry struct {
	DBType     string    `json:"db_type"`
	Name       string    `json:"name"`
	Env        string    `json:"env"`
	PID        int       `json:"pid"`
	LocalPort  int       `json:"local_port"`
	RemotePort int       `json:"remote_port"`
	Host       string    `json:"host"`
	StartedAt  time.Time `json:"started_at"`
}

// runtimeKey builds the state map key for a given type+name.
func runtimeKey(dbType, name string) string {
	return dbType + ":" + name
}
