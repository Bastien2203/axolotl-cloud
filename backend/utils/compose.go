package utils

import (
	"axolotl-cloud/internal/app/model"
	"strings"

	"gopkg.in/yaml.v3"
)

func ParseComposeFileFromBytes(bytes []byte, project *model.Project) (model.ComposeFile, []model.Container, error) {
	var compose model.ComposeFile
	if err := yaml.Unmarshal(bytes, &compose); err != nil {
		return model.ComposeFile{}, nil, err
	}
	containers, err := ParseComposeFile(compose, project)
	if err != nil {
		return model.ComposeFile{}, nil, err
	}
	return compose, containers, nil
}

func ParseComposeFile(content model.ComposeFile, project *model.Project) ([]model.Container, error) {
	containers := make([]model.Container, 0, len(content.Services))
	for name, service := range content.Services {
		container := model.Container{
			ProjectID:   project.ID,
			Name:        FormatContainerName(project.Name, name),
			DockerImage: service.Image,
			Ports:       parsePorts(service.Ports),
			Env:         service.Env,
			Volumes:     parseVolumes(service.Volumes),
			Networks:    parseNetworks(service.Networks),
			NetworkMode: service.NetworkMode,
		}
		containers = append(containers, container)
	}
	return containers, nil
}

func parsePorts(portDefs []string) map[string]string {
	ports := make(map[string]string)
	for _, def := range portDefs {
		parts := strings.Split(def, ":")
		if len(parts) == 2 {
			ports[parts[0]] = parts[1]
		}
	}
	return ports
}

func parseVolumes(volumeDefs []string) map[string]string {
	volumes := make(map[string]string)
	for _, def := range volumeDefs {
		parts := strings.Split(def, ":")
		if len(parts) == 2 {
			volumes[parts[0]] = parts[1]
		}
	}
	return volumes
}

func parseNetworks(networkDefs []string) []string {
	networks := make([]string, 0, len(networkDefs))
	for _, def := range networkDefs {
		if def != "" {
			networks = append(networks, def)
		}
	}
	return networks
}
