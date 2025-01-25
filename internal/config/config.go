package config

import (
	"flag"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	configPath string
	cfg        Config
)

func init() {
	homeDir, _ := os.UserHomeDir()
	defaultConfig := filepath.Join(homeDir, ".neon", "config.yaml")
	flag.StringVar(&configPath, "config", defaultConfig, "Path to config file")
}

type Config struct {
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

func Load() error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &cfg)
}

func Get() *Config {
	return &cfg
}
