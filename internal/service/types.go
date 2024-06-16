package service

import (
	"errors"
	"strings"

	"github.com/docker/docker/api/types"
)

type ServiceType int

const (
	Arr          ServiceType = 0
	Transmission ServiceType = 1
)

type Service struct {
	Container *types.Container
	Type      ServiceType
}

func detectServiceType(image string) (ServiceType, error) {
	if strings.Contains(image, "radarr") ||
		strings.Contains(image, "sonarr") ||
		strings.Contains(image, "prowlarr") ||
		strings.Contains(image, "lidarr") ||
		strings.Contains(image, "readarr") {
		return Arr, nil
	} else if strings.Contains(image, "transmission") {
		return Transmission, nil
	} else {
		return -1, errors.New("could not detect the service type from the image name")
	}
}
