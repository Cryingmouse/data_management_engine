package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cryingmouse/data_management_engine/common"
)

const (
	baseURL = "http://%s:8080/"
)

type RestClient struct {
	client      *http.Client
	hostContext common.HostContext
	prefixURL   string
}

// GetRestClient returns a new instance of the RestClient.
func GetRestClient(hostContext common.HostContext, prefixURL string) *RestClient {
	return &RestClient{
		client:      &http.Client{},
		hostContext: hostContext,
		prefixURL:   prefixURL,
	}
}

// Get performs an HTTP GET request.
func (c *RestClient) Get(url string, contentType string) (*http.Response, error) {
	fullURL := fmt.Sprintf(baseURL+c.prefixURL+"/%s", c.hostContext.IP, url)

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-agent-username", c.hostContext.Username)
	req.Header.Set("X-agent-password", c.hostContext.Password)
	req.Header.Set("Content-Type", contentType)

	return c.client.Do(req)
}

// Post performs an HTTP POST request.
func (c *RestClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	fullURL := fmt.Sprintf(baseURL+c.prefixURL+"/%s", c.hostContext.IP, url)

	req, err := http.NewRequest(http.MethodPost, fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-agent-username", c.hostContext.Username)
	req.Header.Set("X-agent-password", c.hostContext.Password)
	req.Header.Set("Content-Type", contentType)

	return c.client.Do(req)
}

// GetResponseBody reads the response body and unmarshals it into the provided result.
func (c *RestClient) GetResponseBody(response *http.Response, result interface{}) error {
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(responseBody, result)
}
