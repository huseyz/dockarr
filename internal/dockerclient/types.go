package dockerclient

import "github.com/docker/docker/client"

type Client struct {
	dockerClient *client.Client
}

func NewClient(dockerClient *client.Client) *Client {
	return &Client{
		dockerClient: dockerClient,
	}
}
