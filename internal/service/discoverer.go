package service

import (
	"log/slog"

	"github.com/huseyz/dockarr/internal/dockerclient"
)

type ServiceDiscoverer struct {
	dockerClient *dockerclient.Client
}

func NewServiceDiscoverer(dockerClient *dockerclient.Client) *ServiceDiscoverer {
	return &ServiceDiscoverer{
		dockerClient: dockerClient,
	}
}

func (s *ServiceDiscoverer) DiscoverServices() ([]*Service, error) {
	containers, err := s.dockerClient.GetContainersByLabel("dockarr.discover")
	if err != nil {
		return nil, err
	}

	services := []*Service{}

	for _, container := range containers {
		serviceType, err := detectServiceType(container.Image)
		if err != nil {
			slog.Warn("Error detecting service type from image. Ignoring...", "image", container.Image)
			continue
		}

		service := &Service{
			Type:      serviceType,
			Container: &container,
		}
		services = append(services, service)
	}

	return services, err
}
