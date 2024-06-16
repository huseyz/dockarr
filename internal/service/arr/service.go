package arr

import "fmt"

func (s *Service) Address() string {
	if s.labels["dockarr.override.address"] != "" {
		return s.labels["dockarr.override.address"]
	}
	return fmt.Sprintf("%s%s:%s%s", s.xmlConfig.getProtocol(), s.hostname, s.xmlConfig.getPort(), s.xmlConfig.UrlBase)

}

func (s *Service) SyncLevel() string {
	override := s.syncLevelOverride()
	if override == "" {
		return "fullSync"
	} else {
		return override
	}
}

func (s *Service) ApiKey() string {
	return s.xmlConfig.ApiKey
}

func (s *Service) syncLevelOverride() string {
	override := s.labels["dockarr.override.syncLevel"]
	if override == "addOnly" || override == "disabled" || override == "fullSync" {
		return override
	} else {
		return ""
	}
}
