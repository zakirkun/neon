package deploy

type Config struct {
	Services []ServiceConfig `yaml:"services"`
}

type ServiceConfig struct {
	Name        string       `yaml:"name"`
	Image       string       `yaml:"image"`
	Replicas    uint64       `yaml:"replicas"`
	Ports       []PortConfig `yaml:"ports"`
	Environment []string     `yaml:"environment"`
	Networks    []string     `yaml:"networks"`
	Deploy      DeployConfig `yaml:"deploy"`
}

type PortConfig struct {
	Target    uint32 `yaml:"target"`
	Published uint32 `yaml:"published"`
}

type DeployConfig struct {
	UpdateConfig  UpdateConfig  `yaml:"update_config"`
	RestartPolicy RestartPolicy `yaml:"restart_policy"`
	Resources     Resources     `yaml:"resources"`
}

type UpdateConfig struct {
	Parallelism uint64 `yaml:"parallelism"`
	Delay       string `yaml:"delay"`
}

type RestartPolicy struct {
	Condition   string `yaml:"condition"`
	MaxAttempts uint64 `yaml:"max_attempts"`
}

type Resources struct {
	Limits ResourceLimit `yaml:"limits"`
}

type ResourceLimit struct {
	CPUs   string `yaml:"cpus"`
	Memory string `yaml:"memory"`
}
