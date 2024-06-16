package prowlarrclient

import (
	"github.com/huseyz/dockarr/internal/arrclient"
)

type Client struct {
	ArrClient *arrclient.Client
}

type ProwlarrApplication struct {
	Id             int                        `json:"id,omitempty"`
	Name           string                     `json:"name"`
	SyncLevel      string                     `json:"syncLevel"`
	Implementation string                     `json:"implementation"`
	ConfigContract string                     `json:"configContract"`
	Fields         []ProwlarrApplicationField `json:"fields"`
}

type ProwlarrApplicationField struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func (a *ProwlarrApplication) Equals(b *ProwlarrApplication) bool {
	if a.Name != b.Name ||
		a.SyncLevel != b.SyncLevel ||
		a.ConfigContract != b.ConfigContract ||
		a.Implementation != b.Implementation {
		return false
	}
	aFields := fieldsToMap(a.Fields)
	bFields := fieldsToMap(b.Fields)

	return aFields["prowlarrUrl"] == bFields["prowlarrUrl"] &&
		aFields["baseUrl"] == bFields["baseUrl"]
}

func fieldsToMap(fields []ProwlarrApplicationField) map[string]interface{} {
	m := map[string]interface{}{}
	for _, field := range fields {
		m[field.Name] = field.Value
	}
	return m
}
