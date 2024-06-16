package transmission

import (
	"log/slog"
	"strconv"
)

func (s *Service) Host() string {
	if s.labels["dockarr.override.host"] != "" {
		return s.labels["dockarr.override.host"]
	} else {
		return s.hostname
	}
}

func (s *Service) Port() float64 {
	if s.labels["dockarr.override.port"] != "" {
		i, err := strconv.Atoi(s.labels["dockarr.override.port"])
		if err != nil {
			slog.Warn("Failed to parse port override.", "label", s.labels["dockarr.override.port"])
		} else {
			return float64(i)
		}
	}
	return s.config.RpcPort
}

func (s *Service) SSL() bool {
	if s.labels["dockarr.override.ssl"] != "" {
		i, err := strconv.ParseBool(s.labels["dockarr.override.ssl"])
		if err != nil {
			slog.Warn("Failed to parse ssl override.", "label", s.labels["dockarr.override.ssl"])
		} else {
			return i
		}
	}
	return false
}

func (s *Service) UrlBase() string {
	if s.labels["dockarr.override.urlbase"] != "" {
		return s.labels["dockarr.override.urlbase"]
	}
	return s.config.RpcUrl
}

func (s *Service) Username() string {
	if s.env["USER"] != "" {
		return s.env["USER"]
	} else {
		return ""
	}
}

func (s *Service) Password() string {
	if s.env["PASS"] != "" {
		return s.env["PASS"]
	} else {
		return ""
	}
}
