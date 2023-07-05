package client

import (
	"fmt"
	"io"
	"net/http"

	"github.com/cryingmouse/data_management_engine/common"
)

type RestfulAPI interface {
	Get(url string, contentType string) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (*http.Response, error)
}

type RestClient struct {
	client      *http.Client
	hostContext common.HostContext
	prefixURL   string
}

func GetRestClient(hostContext common.HostContext, prefixURL string) RestfulAPI {
	return &RestClient{
		client:      &http.Client{},
		hostContext: hostContext,
		prefixURL:   prefixURL,
	}
}

func (c *RestClient) Get(url string, contentType string) (*http.Response, error) {
	// Append the base URL to the input URL.
	fullURL := fmt.Sprintf("http://%s:8080/%s/%s", c.hostContext.IP, c.prefixURL, url)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-agent-username", c.hostContext.Username)
	req.Header.Add("X-agent-password", c.hostContext.Password)

	req.Header.Set("Content-Type", contentType)
	return c.client.Do(req)
}

func (c *RestClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	// Append the base URL to the input URL.
	fullURL := fmt.Sprintf("http://%s:8080/%s/%s", c.hostContext.IP, c.prefixURL, url)

	req, err := http.NewRequest("POST", fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-agent-username", c.hostContext.Username)
	req.Header.Add("X-agent-password", c.hostContext.Password)

	req.Header.Set("Content-Type", contentType)
	return c.client.Do(req)
}
