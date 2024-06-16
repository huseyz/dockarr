package arr

import (
	"bytes"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/huseyz/dockarr/internal/arrclient"
	"github.com/huseyz/dockarr/internal/prowlarrclient"
)

type ServiceType int

type Service struct {
	Type      ServiceType
	hostname  string
	labels    map[string]string
	xmlConfig XmlConfig
}

const (
	Prowlarr ServiceType = iota
	Radarr
	Sonarr
	Lidarr
	Readarr
)

func (t ServiceType) ConfigContract() string {
	return t.Implementation() + "Settings"
}

func (t ServiceType) Implementation() string {
	return [...]string{"Radarr", "Sonarr", "Lidarr", "Readarr"}[t-1]
}

func FromImplementation(implementation string) ServiceType {
	switch implementation {
	case "Radarr":
		return Radarr
	case "Sonarr":
		return Sonarr
	case "Lidarr":
		return Lidarr
	case "Readarr":
		return Readarr
	}
	panic("")
}

func NewArrService(container *types.Container, xmlConfigContent *bytes.Buffer) (*Service, error) {
	xmlConfig, err := NewXmlConfig(xmlConfigContent)

	if err != nil {
		return nil, err
	}

	return &Service{
		Type:      serviceTypeFromImage(container.Image),
		hostname:  container.Names[0][1:],
		labels:    container.Labels,
		xmlConfig: xmlConfig,
	}, nil
}

func (s *Service) ToProwlarrApplication(prowlarr *Service) *prowlarrclient.ProwlarrApplication {
	return &prowlarrclient.ProwlarrApplication{
		Name:           "[Dockarr] " + s.Type.Implementation(),
		SyncLevel:      s.SyncLevel(),
		Implementation: s.Type.Implementation(),
		ConfigContract: s.Type.ConfigContract(),
		Fields: []prowlarrclient.ProwlarrApplicationField{
			{
				Name:  "prowlarrUrl",
				Value: prowlarr.Address(),
			},
			{
				Name:  "baseUrl",
				Value: s.Address(),
			},
			{
				Name:  "apiKey",
				Value: s.ApiKey(),
			},
		},
	}
}

func (s *Service) Client() *arrclient.Client {
	return arrclient.NewClient(s.Address(), s.ApiKey())
}

func serviceTypeFromImage(image string) ServiceType {
	if strings.Contains(image, "radarr") {
		return Radarr
	} else if strings.Contains(image, "sonarr") {
		return Sonarr
	} else if strings.Contains(image, "prowlarr") {
		return Prowlarr
	} else if strings.Contains(image, "lidarr") {
		return Lidarr
	} else if strings.Contains(image, "readarr") {
		return Readarr
	} else {
		return -1
	}
}
