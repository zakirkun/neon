package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Github struct {
		Token string `yaml:"token"`
	} `yaml:"github"`

	Docker struct {
		Registry string `yaml:"registry"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"docker"`

	Swarm struct {
		ManagerNode string `yaml:"manager_node"`
		NetworkName string `yaml:"network_name"`
	} `yaml:"swarm"`

	Deploy struct {
		Replicas      int    `yaml:"replicas"`
		UpdateDelay   string `yaml:"update_delay"`
		RollbackDelay string `yaml:"rollback_delay"`
		FailureAction string `yaml:"failure_action"`
	} `yaml:"deploy"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
