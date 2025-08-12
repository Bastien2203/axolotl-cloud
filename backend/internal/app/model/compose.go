package model

type ComposeService struct {
	Image       string            `yaml:"image"`
	Ports       []string          `yaml:"ports"`
	Env         map[string]string `yaml:"environment"`
	Volumes     []string          `yaml:"volumes"`
	Networks    []string          `yaml:"networks"`
	NetworkMode string            `yaml:"network_mode" default:"bridge"`
	Build       *ComposeBuild     `yaml:"build,omitempty"`
}

type ComposeFile struct {
	Services map[string]ComposeService `yaml:"services"`
}

type ComposeBuild struct {
	Context    string            `yaml:"context"`
	Dockerfile string            `yaml:"dockerfile,omitempty"`
	Args       map[string]string `yaml:"args,omitempty"`
}
