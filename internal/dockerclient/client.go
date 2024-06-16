package dockerclient

import (
	"archive/tar"
	"bytes"
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

func (c *Client) GetContainersByLabel(label string) ([]types.Container, error) {
	opts := container.ListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", label),
		),
	}
	return c.dockerClient.ContainerList(context.Background(), opts)
}

func (c *Client) GetFileContent(containerID string, filePath string) (*bytes.Buffer, error) {
	reader, _, err := c.dockerClient.CopyFromContainer(context.Background(), containerID, filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	tr := tar.NewReader(reader)
	tr.Next() // Ignore the tar header

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(tr)
	return buffer, nil
}

func (c *Client) GetEnvironmentVariables(containerID string) (map[string]string, error) {
	result := map[string]string{}
	config, err := c.dockerClient.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return nil, err
	}
	for _, pair := range config.Config.Env {
		splitted := strings.Split(pair, "=")
		result[splitted[0]] = splitted[1]
	}
	return result, nil
}
