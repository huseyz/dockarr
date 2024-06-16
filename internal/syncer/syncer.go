package syncer

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/huseyz/dockarr/internal/arrclient"
	"github.com/huseyz/dockarr/internal/config"
	"github.com/huseyz/dockarr/internal/prowlarrclient"
	"github.com/huseyz/dockarr/internal/service"
	"github.com/huseyz/dockarr/internal/service/arr"
	"github.com/huseyz/dockarr/internal/service/transmission"
)

func (syncer *Syncer) SyncServices(services []*service.Service) {
	arrServices, downloaders := syncer.classifyServices(services)

	prowlarrService := arrServices[arr.Prowlarr]

	if prowlarrService == nil {
		slog.Warn("No Prowlarr instance is detected.")
		return
	}

	slog.Debug("Syncing applications...")

	prowlarrClient := &prowlarrclient.Client{
		ArrClient: arrclient.NewClient(prowlarrService.Address(), prowlarrService.ApiKey()),
	}

	existingApplications, err := prowlarrClient.GetApplications()
	if err != nil {
		slog.Warn("Failed to get existing applications from Prowlarr.")
		return
	}

	slog.Debug("Fetched existing applications.", "applications", len(*existingApplications))
	syncer.syncApplications(prowlarrClient, arrServices, prowlarrService, *existingApplications)

	slog.Debug("Syncing download clients...")
	syncer.syncDownloadClients(arrServices, downloaders)
}

func (syncer *Syncer) createTransmissionService(s *service.Service) (*transmission.Service, error) {
	jsonConfigContent, jsonConfigError := syncer.dockerClient.GetFileContent(s.Container.ID, "/config/settings.json")
	envVariables, envVariablesError := syncer.dockerClient.GetEnvironmentVariables(s.Container.ID)
	if envVariablesError != nil || jsonConfigError != nil {
		return nil, errors.Join(envVariablesError, jsonConfigError)
	}
	return transmission.NewTransmissionService(s.Container, jsonConfigContent, envVariables)
}

func (syncer *Syncer) createArrService(s *service.Service) (*arr.Service, error) {
	xmlConfigContent, err := syncer.dockerClient.GetFileContent(s.Container.ID, "/config/config.xml")
	if err != nil {
		return nil, err
	}
	return arr.NewArrService(s.Container, xmlConfigContent)
}

func getSyncedApplication(apps []prowlarrclient.ProwlarrApplication) map[string]*prowlarrclient.ProwlarrApplication {
	syncedApps := map[string]*prowlarrclient.ProwlarrApplication{}

	for _, app := range apps {
		if !strings.Contains(app.Name, "Dockarr") {
			continue
		}
		syncedApps[app.Implementation] = &app
	}
	return syncedApps
}

func getSyncedDownloaders(downloaders []arrclient.DownloadClient) map[string]*arrclient.DownloadClient {
	synced := map[string]*arrclient.DownloadClient{}

	for _, downloader := range downloaders {
		if !strings.Contains(downloader.Name, "Dockarr") {
			continue
		}
		synced[downloader.Implementation] = &downloader
	}
	return synced
}

func (syncer *Syncer) classifyServices(services []*service.Service) (map[arr.ServiceType]*arr.Service, map[string]*arrclient.DownloadClient) {
	arrServices := map[arr.ServiceType]*arr.Service{}
	downloaders := map[string]*arrclient.DownloadClient{}

	for _, s := range services {
		switch s.Type {
		case service.Arr:
			arrService, err := syncer.createArrService(s)
			if err != nil {
				slog.Warn("Error creating an *arr service.", err)
			}
			arrServices[arrService.Type] = arrService
		case service.Transmission:
			transmissionService, err := syncer.createTransmissionService(s)
			if err != nil {
				slog.Warn("Error creating Transmission service", err)
			}
			cli := transmissionService.ToDownloadClient()
			downloaders[cli.Implementation] = cli
		}
	}
	return arrServices, downloaders
}

func (syncer *Syncer) syncApplications(prowlarrClient *prowlarrclient.Client, arrServices map[arr.ServiceType]*arr.Service, prowlarrService *arr.Service, existingApplications []prowlarrclient.ProwlarrApplication) {
	syncedApplications := getSyncedApplication(existingApplications)

	for t, arrService := range arrServices {
		if t == arr.Prowlarr {
			continue
		}
		existing := syncedApplications[t.Implementation()]
		if existing != nil {
			slog.Debug("Found existing application.", "type", arrService.Type.Implementation())
			prowlarrApp := arrService.ToProwlarrApplication(prowlarrService)
			if prowlarrApp.Equals(existing) {
				slog.Debug("Existing application has not changed.", "name", existing.Name)
			} else {
				slog.Info("Existing application changed, updating.", "name", existing.Name)
				prowlarrApp.Id = existing.Id
				_, err := prowlarrClient.UpdateApplication(prowlarrApp)
				if err != nil {
					slog.Error("Failed to update existing application.", "name", existing.Name, "err", err)
				}
			}
		} else {
			slog.Info("New application discovered. Adding.", "type", t.Implementation())
			prowlarrApp := arrService.ToProwlarrApplication(prowlarrService)
			_, err := prowlarrClient.AddApplication(prowlarrApp)
			if err != nil {
				slog.Error("Failed to add new application.", "name", prowlarrApp.Name, "err", err)
			}
		}
	}

	for _, syncedApplication := range syncedApplications {
		t := arr.FromImplementation(syncedApplication.Implementation)
		if arrServices[t] == nil {
			slog.Info("Existing application is not discovered.", "name", syncedApplication.Name)
			if syncer.config.DeleteBehaviour == config.Delete {
				slog.Info("Deleting application.", "name", syncedApplication.Name)
				err := prowlarrClient.DeleteApplication(syncedApplication.Id)
				if err != nil {
					slog.Error("Failed to delete application.", "app", t.Implementation())
				}
			} else if syncer.config.DeleteBehaviour == config.Disable && syncedApplication.SyncLevel != "disabled" {
				slog.Info("Disabling application.", "name", syncedApplication.Name)
				syncedApplication.SyncLevel = "disabled"
				_, err := prowlarrClient.UpdateApplication(syncedApplication)
				if err != nil {
					slog.Error("Failed to update application.", "app", t.Implementation())
				}
			}
		}
	}
}

func (syncer *Syncer) syncDownloadClients(arrServices map[arr.ServiceType]*arr.Service, downloaders map[string]*arrclient.DownloadClient) {
	for t, arrService := range arrServices {
		if t == arr.Prowlarr {
			continue
		}
		slog.Debug("Syncing service.", "service", t.Implementation())
		client := arrService.Client()
		existingDownloaders, _ := client.GetDownloaders()
		syncedDownloaders := getSyncedDownloaders(*existingDownloaders)
		for name, downloader := range downloaders {
			existing := syncedDownloaders[name]
			if existing == nil {
				slog.Info("Found new downloader. Adding...", "app", t.Implementation(), "downloader", downloader.Name)
				_, err := client.AddDownloader(downloader)
				if err != nil {
					slog.Error("Failed to add new downloader.", "app", t.Implementation(), "downloader", downloader.Name, "err", err)
				}
			} else {
				slog.Debug("Found existing downloader.", "app", t.Implementation(), "downloader", downloader.Name)
				if existing.Equals(downloader) {
					slog.Debug("Existing downloader did not change.", "app", t.Implementation(), "downloader", existing.Name)
				} else {
					slog.Debug("Existing downloader did change.", "app", t.Implementation(), "downloader", existing.Name)
					downloader.Id = existing.Id
					_, err := client.UpdateDownloader(downloader)
					if err != nil {
						slog.Error("Failed to update downloader.", "app", t.Implementation(), "downloader", existing.Name)
					}
				}
			}
		}

		for name, syncedDownloader := range syncedDownloaders {
			if downloaders[name] == nil {
				slog.Info("Existing downloader is not discovered.", "app", t.Implementation(), "downloader", syncedDownloader.Name)
				if syncer.config.DeleteBehaviour == config.Delete {
					slog.Info("Deleting the downloader.", "app", t.Implementation(), "downloader", syncedDownloader.Name)
					err := client.DeleteDownloader(syncedDownloader.Id)
					if err != nil {
						slog.Error("Failed to delete downloader.", "app", t.Implementation(), "downloader", syncedDownloader.Name)
					}
				} else if syncer.config.DeleteBehaviour == config.Disable {
					slog.Info("Disabling the downloader.", "app", t.Implementation(), "downloader", syncedDownloader.Name)
					syncedDownloader.Enable = false
					_, err := client.UpdateDownloader(syncedDownloader)
					if err != nil {
						slog.Error("Failed to update downloader.", "app", t.Implementation(), "downloader", syncedDownloader.Name)
					}
				}
			}
		}
	}
}
