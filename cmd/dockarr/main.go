package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/huseyz/dockarr/internal/config"
	"github.com/huseyz/dockarr/internal/dockerclient"
	"github.com/huseyz/dockarr/internal/service"
	"github.com/huseyz/dockarr/internal/syncer"
)

func sync(discoverer *service.ServiceDiscoverer, syncer *syncer.Syncer) {
	slog.Debug("Syncing...")

	services, _ := discoverer.DiscoverServices()

	syncer.SyncServices(services)

}

func main() {
	dockarrConfig := config.InitConfig()
	slog.Info("Initializing Dockarr", "Config", fmt.Sprintf("%+v\n", dockarrConfig))

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: dockarrConfig.LogLevel}))
	slog.SetDefault(logger)

	cli, _ := client.NewClientWithOpts(client.FromEnv)

	dockerClient := dockerclient.NewClient(cli)

	discoverer := service.NewServiceDiscoverer(dockerClient)
	syncer := syncer.NewSyncer(dockarrConfig, dockerClient)

	slog.Info("Starting the sync loop.", "sync interval", dockarrConfig.SyncInterval)
	sync(discoverer, syncer)

	ticker := time.NewTicker(dockarrConfig.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		sync(discoverer, syncer)
	}
}
