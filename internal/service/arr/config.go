package arr

import (
	"bytes"
	"encoding/xml"
)

type XmlConfig struct {
	XMLName   xml.Name `xml:"Config"`
	ApiKey    string   `xml:"ApiKey"`
	UrlBase   string   `xml:"UrlBase"`
	Port      string   `xml:"Port"`
	SslPort   string   `xml:"SslPort"`
	EnableSsl bool     `xml:"EnableSsl"`
}

func NewXmlConfig(content *bytes.Buffer) (XmlConfig, error) {
	var config XmlConfig
	error := xml.Unmarshal(content.Bytes(), &config)
	return config, error
}

func (config *XmlConfig) getPort() string {
	if config.EnableSsl {
		return config.SslPort
	} else {
		return config.Port
	}
}

func (config *XmlConfig) getProtocol() string {
	if config.EnableSsl {
		return "https://"
	} else {
		return "http://"
	}
}
