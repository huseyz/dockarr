package arrclient

type Client struct {
	address    string
	apiKey     string
	apiVersion string
}

func NewClient(Address, ApiKey string) *Client {
	return &Client{
		address: Address,
		apiKey:  ApiKey,
	}
}

type ApiVersionResponse struct {
	Current string `json:"current"`
}

type DownloadClientField struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value,omitempty"`
}

func NewDownloadClientField(name string, value interface{}) DownloadClientField {
	return DownloadClientField{
		Name:  name,
		Value: value,
	}
}

type DownloadClient struct {
	Id             int                   `json:"id,omitempty"`
	Name           string                `json:"name"`
	Implementation string                `json:"implementation"`
	ConfigContract string                `json:"configContract"`
	Protocol       string                `json:"protocol"`
	Fields         []DownloadClientField `json:"fields"`
	Enable         bool                  `json:"enable"`
}

func (a *DownloadClient) Equals(b *DownloadClient) bool {
	if a.Name != b.Name ||
		a.Implementation != b.Implementation ||
		a.ConfigContract != b.ConfigContract ||
		a.Protocol != b.Protocol ||
		a.Enable != b.Enable {
		return false
	}
	aFields := fieldsToMap(a.Fields)
	bFields := fieldsToMap(b.Fields)

	return compareFields(aFields, bFields, []string{
		"host", "port", "useSsl", "urlBase", "username", "password",
	})
}

func fieldsToMap(fields []DownloadClientField) map[string]interface{} {
	m := map[string]interface{}{}
	for _, field := range fields {
		m[field.Name] = field.Value
	}
	return m
}

func compareFields(aFields, bFields map[string]interface{}, fields []string) bool {
	for _, field := range fields {
		if aFields[field] != bFields[field] {
			return false
		}
	}
	return true
}
