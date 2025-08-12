package docker

import (
	"axolotl-cloud/infra/logger"
	"axolotl-cloud/infra/shared"
	"axolotl-cloud/internal/app/model"
	"axolotl-cloud/utils"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	dImage "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
)

func (dc *DockerClient) ContainerExists(ctx context.Context, name string, log *logger.Logger) (bool, error) {
	cli := dc.cli
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return false, fmt.Errorf("failed to list containers: %w", err)
	}
	for _, c := range containers {
		if c.Names[0] == "/"+name {
			log.Info("Container %s exists", name)
			return true, nil
		}
	}
	log.Info("Container %s does not exist", name)
	return false, nil
}

func (dc *DockerClient) StartContainer(ctx context.Context, name string, log *logger.Logger) (string, error) {
	cli := dc.cli

	err := cli.ContainerStart(ctx, name, container.StartOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to start container %s: %w", name, err)
	}

	log.Info("Container %s started successfully", name)
	return name, nil
}

func (dc *DockerClient) StopContainer(ctx context.Context, name string, log *logger.Logger) error {
	cli := dc.cli

	err := cli.ContainerStop(ctx, name, container.StopOptions{})
	if err != nil {
		return fmt.Errorf("failed to stop container %s: %w", name, err)
	}

	log.Info("Container %s stopped successfully", name)
	return nil
}

func (dc *DockerClient) RemoveContainer(ctx context.Context, name string, log *logger.Logger) error {
	cli := dc.cli

	err := cli.ContainerRemove(ctx, name, container.RemoveOptions{Force: true})
	if err != nil {
		return fmt.Errorf("failed to remove container %s: %w", name, err)
	}

	log.Info("Container %s removed successfully", name)
	return nil
}

func (dc *DockerClient) PullImage(ctx context.Context, image string, log *logger.Logger) error {
	cli := dc.cli
	reader, err := cli.ImagePull(ctx, image, dImage.PullOptions{})
	if err != nil {
		// check if image exists locally
		_, err := cli.ImageInspect(ctx, image)
		if err == nil {
			log.Info("Image %s exists locally", image)
			return nil
		}
		return fmt.Errorf("failed to pull image %s: %w", image, err)
	}
	io.Copy(io.Discard, reader)
	reader.Close()
	log.Info("Image %s pulled successfully", image)
	return nil
}

func (dc *DockerClient) CreateContainer(ctx context.Context, name string, image string, ports map[string]string, env map[string]string, volumes map[string]string, networkMode string, log *logger.Logger) (string, error) {
	cli := dc.cli

	if err := dc.PullImage(ctx, image, log); err != nil {
		return "", fmt.Errorf("failed to pull image %s: %w", image, err)
	}

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
			log.Info("Creating volume at: %s", sourceContainer)
			if err := os.MkdirAll(sourceContainer, 0755); err != nil {
				return "", fmt.Errorf("failed to create volume directory %s: %w", sourceContainer, err)
			}
			sourceHost = fmt.Sprintf("%s/%s/%s", volumesPathHost, name, hostPath)
		}

		mounts = append(mounts, mount.Mount{Type: mount.TypeBind, Source: sourceHost, Target: containerPath})
	}

	config := &container.Config{Image: image, Env: envVars, ExposedPorts: exposed}
	hostConfig := &container.HostConfig{
		PortBindings: bindings,
		Mounts:       mounts,
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		Resources: container.Resources{
			CgroupParent: "/docker.slice",
		},
	}
	if networkMode != "" {
		hostConfig.NetworkMode = container.NetworkMode(networkMode)
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}
	log.Info("Container %s created successfully", name)

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

func (dc *DockerClient) ContainerVolumes(ctx context.Context, name string) ([]*model.Volume, error) {
	cli := dc.cli

	if exists, err := dc.ContainerExists(ctx, name, logger.GlobalLogger); err != nil || !exists {
		return []*model.Volume{}, nil
	}

	resp, err := cli.ContainerInspect(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container %s: %w", name, err)
	}

	var volumes []*model.Volume
	for _, mount := range resp.Mounts {
		volumes = append(volumes, &model.Volume{
			Source:      mount.Source,
			Destination: mount.Destination,
			Type:        typeToString(mount.Type),
			Size:        getSize(mount.Source),
		})
	}

	return volumes, nil
}

// Get size in bytes of a file or directory
func getSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Error getting size for %s: %v\n", path, err)
		return 0 // or handle error as needed
	}
	if info.IsDir() {
		size := int64(0)
		files, err := os.ReadDir(path)
		if err != nil {
			return 0 // or handle error as needed
		}
		for _, file := range files {
			size += getSize(fmt.Sprintf("%s/%s", path, file.Name()))
		}
		return size
	}
	return info.Size()
}

func typeToString(t mount.Type) string {
	switch t {
	case mount.TypeBind:
		return "bind"
	case mount.TypeVolume:
		return "volume"
	case mount.TypeTmpfs:
		return "tmpfs"
	case mount.TypeNamedPipe:
		return "npipe"
	case mount.TypeCluster:
		return "cluster"
	case mount.TypeImage:
		return "image"
	default:
		return "unknown"
	}
}
