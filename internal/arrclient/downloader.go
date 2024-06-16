package arrclient

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (c *Client) GetDownloaders() (*[]DownloadClient, error) {
	req, err := c.ApiRequest("GET", "/downloadclient", nil)
	if err != nil {
		return nil, err
	}
	return DoRequest[[]DownloadClient](req)
}

func (c *Client) AddDownloader(downloader *DownloadClient) (*DownloadClient, error) {
	return c.upsert(downloader, false)
}

func (c *Client) UpdateDownloader(downloader *DownloadClient) (*DownloadClient, error) {
	return c.upsert(downloader, true)
}

func (c *Client) DeleteDownloader(id int) error {
	req, err := c.ApiRequest("DELETE", "/applications/"+fmt.Sprint(id), nil)
	if err != nil {
		return err
	}
	_, err = DoRequest[interface{}](req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) upsert(downloader *DownloadClient, update bool) (*DownloadClient, error) {
	jsonBytes, err := json.Marshal(downloader)
	if err != nil {
		return nil, err
	}

	operation := "POST"
	if update {
		operation = "PUT"
	}

	req, err := c.ApiRequest(operation, "/downloadclient", bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	return DoRequest[DownloadClient](req)

}
