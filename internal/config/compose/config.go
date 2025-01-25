package compose

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
	Networks map[string]Network `yaml:"networks"`
	Volumes  map[string]Volume  `yaml:"volumes"`
}

type Service struct {
	Image       string            `yaml:"image"`
	Build       *BuildConfig      `yaml:"build"`
	Command     string            `yaml:"command"`
	Environment map[string]string `yaml:"environment"`
	Ports       []string          `yaml:"ports"`
	Networks    []string          `yaml:"networks"`
	Volumes     []string          `yaml:"volumes"`
	Deploy      DeployConfig      `yaml:"deploy"`
}

type BuildConfig struct {
	Context    string            `yaml:"context"`
	Dockerfile string            `yaml:"dockerfile"`
	Args       map[string]string `yaml:"args"`
}

type DeployConfig struct {
	Replicas     int            `yaml:"replicas"`
	Resources    ResourceConfig `yaml:"resources"`
	UpdateConfig UpdateConfig   `yaml:"update_config"`
	Restart      RestartConfig  `yaml:"restart_policy"`
}

type ResourceConfig struct {
	Limits struct {
		CPUs   string `yaml:"cpus"`
		Memory string `yaml:"memory"`
	} `yaml:"limits"`
}

type UpdateConfig struct {
	Parallelism int    `yaml:"parallelism"`
	Delay       string `yaml:"delay"`
}

type RestartConfig struct {
	Condition   string `yaml:"condition"`
	MaxAttempts int    `yaml:"max_attempts"`
}

type Network struct {
	External bool   `yaml:"external"`
	Name     string `yaml:"name"`
}

type Volume struct {
	External bool   `yaml:"external"`
	Name     string `yaml:"name"`
}

func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file compose: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("gagal parse compose file: %v", err)
	}

	return &config, nil
}
