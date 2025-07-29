package docker

import (
	"axolotl-cloud/infra/shared"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	dImage "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
)

func (dc *DockerClient) ContainerExists(ctx context.Context, name string) (bool, error) {
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

func (dc *DockerClient) StartContainer(ctx context.Context, name string) (string, error) {
	cli := dc.cli

	err := cli.ContainerStart(ctx, name, container.StartOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to start container %s: %w", name, err)
	}

	return name, nil
}

func (dc *DockerClient) StopContainer(ctx context.Context, name string) error {
	cli := dc.cli

	err := cli.ContainerStop(ctx, name, container.StopOptions{})
	if err != nil {
		return fmt.Errorf("failed to stop container %s: %w", name, err)
	}

	return nil
}

func (dc *DockerClient) RemoveContainer(ctx context.Context, name string) error {
	cli := dc.cli

	err := cli.ContainerRemove(ctx, name, container.RemoveOptions{Force: true})
	if err != nil {
		return fmt.Errorf("failed to remove container %s: %w", name, err)
	}

	return nil
}

func (dc *DockerClient) CreateContainer(ctx context.Context, name string, image string, ports map[string]string, env map[string]string, volumes map[string]string) (string, error) {
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
	volumesPath := shared.GetEnv("VOLUMES_PATH")
	for hostPath, containerPath := range volumes {
		source := fmt.Sprintf("%s/%s/%s", volumesPath, name, hostPath)
		fmt.Println("Creating volume at:", source)
		if err := os.MkdirAll(source, 0755); err != nil {
			return "", fmt.Errorf("failed to create volume directory %s: %w", source, err)
		}
		mounts = append(mounts, mount.Mount{Type: mount.TypeBind, Source: source, Target: containerPath})
	}

	config := &container.Config{Image: image, Env: envVars, ExposedPorts: exposed}
	hostConfig := &container.HostConfig{PortBindings: bindings, Mounts: mounts}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	return resp.ID, nil
}

func (dc *DockerClient) ContainerStatus(ctx context.Context, name string) (container.ContainerState, error) {
	cli := dc.cli

	resp, err := cli.ContainerInspect(ctx, name)
	if err != nil {
		return container.StateDead, fmt.Errorf("failed to inspect container %s: %w", name, err)
	}
	return resp.State.Status, nil
}

func (dc *DockerClient) ContainerHealth(ctx context.Context, name string) (container.HealthStatus, error) {
	cli := dc.cli

	resp, err := cli.ContainerInspect(ctx, name)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container %s: %w", name, err)
	}
	return resp.State.Health.Status, nil
}

func (dc *DockerClient) GetContainerLogs(ctx context.Context, name string, n string) (string, error) {
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
