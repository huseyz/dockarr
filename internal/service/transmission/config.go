package transmission

import (
	"bytes"
	"encoding/json"
)

type JsonConfig struct {
	RpcPort float64 `json:"rpc-port"`
	RpcUrl  string  `json:"rpc-url"`
}

func NewJsonConfig(content *bytes.Buffer) (JsonConfig, error) {
	var config JsonConfig
	err := json.Unmarshal(content.Bytes(), &config)
	return config, err
}
