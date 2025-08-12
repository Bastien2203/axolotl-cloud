package docker

import (
	"archive/tar"
	"axolotl-cloud/infra/logger"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/build"
	"github.com/moby/buildkit/session"
)

type buildMsg struct {
	Stream      string          `json:"stream"`
	Error       string          `json:"error"`
	ErrorDetail json.RawMessage `json:"errorDetail"`
	Status      string          `json:"status"`
}

func (dc *DockerClient) BuildImage(
	ctx context.Context,
	contextDir string,
	dockerfile string,
	imageName string,
	log *logger.Logger,
) error {
	log.Info("Building Docker image %s from directory %s", imageName, contextDir)

	buf, err := createTarFromDir(contextDir)
	if err != nil {
		return fmt.Errorf("error creating tar from directory %s: %w", contextDir, err)
	}

	// Create a new session for the build
	sess, err := session.NewSession(ctx, "axolotl-cloud")
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	dialSession := func(ctx context.Context, proto string, meta map[string][]string) (net.Conn, error) {
		return dc.cli.DialHijack(ctx, "/session", proto, meta)
	}

	// Run session in a separate goroutine
	go func() {
		defer sess.Close()
		if err := sess.Run(ctx, dialSession); err != nil {
			log.Error("Session error: %v", err)
		}
	}()

	buildOptions := build.ImageBuildOptions{
		Tags:        []string{imageName},
		Remove:      true,
		ForceRemove: true,
		Version:     build.BuilderBuildKit,
		SessionID:   sess.ID(),
		Dockerfile:  dockerfile,
	}
	log.Info("Building image with options: %+v", buildOptions)

	response, err := dc.cli.ImageBuild(ctx, buf, buildOptions)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	dec := json.NewDecoder(response.Body)
	for {
		var m buildMsg
		if err := dec.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if m.Error != "" {
			return fmt.Errorf("build error: %s", m.Error)
		}
		if m.Stream != "" && m.Stream != "\n" {
			log.Info(strings.TrimSpace(m.Stream))
		}
	}

	return nil
}

func createTarFromDir(dir string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	err := filepath.Walk(dir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir, file)
		if err != nil {
			return err
		}

		if fi.Mode().IsRegular() {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()

			hdr, err := tar.FileInfoHeader(fi, "")
			if err != nil {
				return err
			}
			hdr.Name = relPath

			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		} else if fi.IsDir() {
			hdr, err := tar.FileInfoHeader(fi, "")
			if err != nil {
				return err
			}
			hdr.Name = relPath + "/"
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}
