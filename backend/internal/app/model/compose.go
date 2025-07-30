package model

type ComposeService struct {
	Image       string            `yaml:"image"`
	Ports       []string          `yaml:"ports"`
	Env         map[string]string `yaml:"environment"`
	Volumes     []string          `yaml:"volumes"`
	Networks    []string          `yaml:"networks"`
	NetworkMode string            `yaml:"network_mode" default:"bridge"`
}

type ComposeFile struct {
	Services map[string]ComposeService `yaml:"services"`
}
