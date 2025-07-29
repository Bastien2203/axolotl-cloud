package docker

import (
	"axolotl-cloud/infra/shared"
	"axolotl-cloud/utils"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	dImage "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
)

func (dc *DockerClient) ContainerExists(name string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cli := dc.cli
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return false, fmt.Errorf("failed to list containers: %w", err)
	}
	for _, c := range containers {
		if c.Names[0] == "/"+name {
			return true, nil
		}
	}
	return false, nil
}

func (dc *DockerClient) StartContainer(name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cli := dc.cli

	err := cli.ContainerStart(ctx, name, container.StartOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to start container %s: %w", name, err)
	}

	return name, nil
}

func (dc *DockerClient) StopContainer(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cli := dc.cli

	err := cli.ContainerStop(ctx, name, container.StopOptions{})
	if err != nil {
		return fmt.Errorf("failed to stop container %s: %w", name, err)
	}

	return nil
}

func (dc *DockerClient) RemoveContainer(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cli := dc.cli

	err := cli.ContainerRemove(ctx, name, container.RemoveOptions{Force: true})
	if err != nil {
		return fmt.Errorf("failed to remove container %s: %w", name, err)
	}

	return nil
}

func (dc *DockerClient) CreateContainer(name string, image string, ports map[string]string, env map[string]string, volumes map[string]string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cli := dc.cli

	// pull image
	reader, err := cli.ImagePull(ctx, image, dImage.PullOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to pull image %s: %w", image, err)
	}
	io.Copy(io.Discard, reader)
	reader.Close()

	// ports
	exposed := nat.PortSet{}
	bindings := nat.PortMap{}
	for host, cont := range ports {
		port := nat.Port(cont + "/tcp")
		exposed[port] = struct{}{}
		bindings[port] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: host}}
	}

	// env
	var envVars []string
	for k, v := range env {
		envVars = append(envVars, k+"="+v)
	}

	// volumes
	var mounts []mount.Mount
	volumesPathHost := shared.GetEnv("VOLUMES_PATH_HOST")
	volumesPathContainer := shared.GetEnv("VOLUMES_PATH_CONTAINER")

	for hostPath, containerPath := range volumes {
		sourceHost := hostPath
		if !utils.IsAbsolutePath(hostPath) {
			sourceContainer := fmt.Sprintf("%s/%s/%s", volumesPathContainer, name, hostPath)
			fmt.Println("Creating volume at:", sourceContainer)
			if err := os.MkdirAll(sourceContainer, 0755); err != nil {
				return "", fmt.Errorf("failed to create volume directory %s: %w", sourceContainer, err)
			}
			sourceHost = fmt.Sprintf("%s/%s/%s", volumesPathHost, name, hostPath)
		}

		mounts = append(mounts, mount.Mount{Type: mount.TypeBind, Source: sourceHost, Target: containerPath})
	}

	config := &container.Config{Image: image, Env: envVars, ExposedPorts: exposed}
	hostConfig := &container.HostConfig{PortBindings: bindings, Mounts: mounts}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	return resp.ID, nil
}

func (dc *DockerClient) ContainerStatus(name string) (container.ContainerState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cli := dc.cli

	resp, err := cli.ContainerInspect(ctx, name)
	if err != nil {
		return container.StateDead, fmt.Errorf("failed to inspect container %s: %w", name, err)
	}
	return resp.State.Status, nil
}

func (dc *DockerClient) ContainerHealth(name string) (container.HealthStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cli := dc.cli

	resp, err := cli.ContainerInspect(ctx, name)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container %s: %w", name, err)
	}
	return resp.State.Health.Status, nil
}

func (dc *DockerClient) GetContainerLogs(name string, n string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cli := dc.cli

	reader, err := cli.ContainerLogs(ctx, name, container.LogsOptions{ShowStdout: true, ShowStderr: true, Tail: n})
	if err != nil {
		return "", fmt.Errorf("failed to get logs for container %s: %w", name, err)
	}
	defer reader.Close()

	logs, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read logs for container %s: %w", name, err)
	}

	return string(logs), nil
}
