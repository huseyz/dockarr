package prowlarrclient

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/huseyz/dockarr/internal/arrclient"
)

func (c *Client) GetApplications() (*[]ProwlarrApplication, error) {
	req, err := c.ArrClient.ApiRequest("GET", "/applications", nil)
	if err != nil {
		return nil, err
	}
	return arrclient.DoRequest[[]ProwlarrApplication](req)
}

func (c *Client) AddApplication(app *ProwlarrApplication) (*ProwlarrApplication, error) {
	return c.application(app, false)
}

func (c *Client) UpdateApplication(app *ProwlarrApplication) (*ProwlarrApplication, error) {
	return c.application(app, true)
}

func (c *Client) DeleteApplication(id int) error {
	req, err := c.ArrClient.ApiRequest("DELETE", "/applications/"+fmt.Sprint(id), nil)
	if err != nil {
		return err
	}
	_, err = arrclient.DoRequest[interface{}](req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) application(app *ProwlarrApplication, update bool) (*ProwlarrApplication, error) {
	jsonBytes, err := json.Marshal(app)
	if err != nil {
		return nil, err
	}

	operation := "POST"
	if update {
		operation = "PUT"
	}

	req, err := c.ArrClient.ApiRequest(operation, "/applications", bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	return arrclient.DoRequest[ProwlarrApplication](req)
}
