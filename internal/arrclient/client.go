package arrclient

import (
	"encoding/json"
	"io"
	"net/http"
)

func DoRequest[V any](req *http.Request) (*V, error) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var result V
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) ApiRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	version, err := c.getApiVersion()
	if err != nil {
		return nil, err
	}
	return c.authenticatedReqest(method, "/api/"+version+endpoint, body)
}

func (c *Client) getApiVersion() (string, error) {
	if c.apiVersion != "" {
		return c.apiVersion, nil
	} else {
		req, err := c.authenticatedReqest("GET", "/api", nil)
		if err != nil {
			return "", err
		}
		response, err := DoRequest[ApiVersionResponse](req)
		if err != nil {
			return "", err
		}
		c.apiVersion = response.Current
		return c.apiVersion, nil
	}
}

func (c *Client) authenticatedReqest(method, endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.address+endpoint, body)
	req.Header.Add("X-Api-Key", c.apiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	return req, err
}
