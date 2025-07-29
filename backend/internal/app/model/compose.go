package model

type ComposeService struct {
	Image    string            `yaml:"image"`
	Ports    []string          `yaml:"ports"`
	Env      map[string]string `yaml:"environment"`
	Volumes  []string          `yaml:"volumes"`
	Networks []string          `yaml:"networks"`
}

type ComposeFile struct {
	Services map[string]ComposeService `yaml:"services"`
}
