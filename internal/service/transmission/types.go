package transmission

import (
	"bytes"

	"github.com/docker/docker/api/types"
	"github.com/huseyz/dockarr/internal/arrclient"
)

type Service struct {
	hostname string
	labels   map[string]string
	env      map[string]string
	config   *JsonConfig
}

func (t *Service) Protocol() string {
	return "torrent"
}

func (t *Service) ConfigContract() string {
	return t.Implementation() + "Settings"
}

func (t *Service) Implementation() string {
	return "Transmission"
}

func NewTransmissionService(container *types.Container, jsonConfigContent *bytes.Buffer, env map[string]string) (*Service, error) {
	jsonConfig, err := NewJsonConfig(jsonConfigContent)
	if err != nil {
		return nil, err
	}
	return &Service{
		hostname: container.Names[0][1:],
		labels:   container.Labels,
		env:      env,
		config:   &jsonConfig,
	}, nil
}

func (t *Service) ToDownloadClient() *arrclient.DownloadClient {
	return &arrclient.DownloadClient{
		Name:           "[Dockarr] " + t.Implementation(),
		Implementation: t.Implementation(),
		ConfigContract: t.ConfigContract(),
		Protocol:       t.Protocol(),
		Enable:         true,
		Fields: []arrclient.DownloadClientField{
			arrclient.NewDownloadClientField("host", t.Host()),
			arrclient.NewDownloadClientField("port", t.Port()),
			arrclient.NewDownloadClientField("useSsl", t.SSL()),
			arrclient.NewDownloadClientField("urlBase", t.UrlBase()),
			arrclient.NewDownloadClientField("username", t.Username()),
			arrclient.NewDownloadClientField("password", t.Password()),
		},
	}
}
