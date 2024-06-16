package syncer

import (
	"github.com/huseyz/dockarr/internal/config"
	"github.com/huseyz/dockarr/internal/dockerclient"
)

type Syncer struct {
	config       *config.DockarrConfig
	dockerClient *dockerclient.Client
}

func NewSyncer(config *config.DockarrConfig, dockerClient *dockerclient.Client) *Syncer {
	return &Syncer{
		config:       config,
		dockerClient: dockerClient,
	}
}
