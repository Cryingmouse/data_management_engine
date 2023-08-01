package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cryingmouse/data_management_engine/common"
	log "github.com/sirupsen/logrus"
)

type RestClient struct {
	client      *http.Client
	baseURL     string
	hostContext common.HostContext
	ContentType string
	authEnabled bool
	tokenKey    string
	AuthToken   string
	TraceID     string
}

// GetRestClient returns a new instance of the RestClient.
func GetRestClient(scheme string, hostContext common.HostContext, port int, prefixURL, tokenKey, traceID string, authEnabled bool) *RestClient {
	return &RestClient{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		baseURL:     fmt.Sprintf("%s://%s:%d/%s", scheme, hostContext.IP, port, prefixURL),
		hostContext: hostContext,
		ContentType: "application/json",
		authEnabled: authEnabled,
		tokenKey:    tokenKey,
		AuthToken:   "",
		TraceID:     traceID,
	}
}

// getAuthorizationHeader returns the Authorization header value based on the current authentication state.
func (c *RestClient) getAuthorizationHeader() string {
	if c.AuthToken != "" {
		return "Bearer " + c.AuthToken
	}
	// Fallback to Basic Authentication if token is missing
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(c.hostContext.Username+":"+c.hostContext.Password))
}

func (c *RestClient) refreshAuthToken(response *http.Response) error {
	if c.tokenKey != "" {
		token := response.Header.Get(c.tokenKey)
		if token != "" {
			c.AuthToken = token
		} else {
			common.Logger.WithFields(log.Fields{"token": c.tokenKey}).Error("Failed to get the token from response.")
			return fmt.Errorf("failed to get the token from response. Token key: %s", c.tokenKey)
		}
	} else {
		common.Logger.Error("No token key in RestClient.")
		return fmt.Errorf("no token key in RestClient")
	}

	return nil
}

// Get performs an HTTP GET request.
func (c *RestClient) Get(url string) (*http.Response, error) {
	fullURL := fmt.Sprintf(c.baseURL+"/%s", url)

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}

	// Set the Authorization header
	req.Header.Set("Content-Type", c.ContentType)
	req.Header.Set("X-Trace-ID", c.TraceID)

	if c.authEnabled {
		req.Header.Set("Authorization", c.getAuthorizationHeader())
	}

	resp, err := c.client.Do(req)
	if resp != nil && resp.StatusCode == http.StatusUnauthorized {
		// Try using Basic Authentication if token returns a 401 status code
		c.AuthToken = "" // Reset the AuthToken to trigger Basic Authentication
		return c.Get(url)
	} else if c.authEnabled {
		c.refreshAuthToken(resp)
	}

	return resp, err
}

// Post performs an HTTP POST request.
func (c *RestClient) Post(url string, body io.Reader) (*http.Response, error) {
	fullURL := fmt.Sprintf(c.baseURL+"/%s", url)

	req, err := http.NewRequest(http.MethodPost, fullURL, body)
	if err != nil {
		return nil, err
	}

	// Set the Authorization header
	req.Header.Set("Content-Type", c.ContentType)
	req.Header.Set("X-Trace-ID", c.TraceID)

	if c.authEnabled {
		req.Header.Set("Authorization", c.getAuthorizationHeader())
	}

	resp, err := c.client.Do(req)
	if resp != nil && resp.StatusCode == http.StatusUnauthorized {
		// Try using Basic Authentication if token returns a 401 status code
		c.AuthToken = "" // Reset the AuthToken to trigger Basic Authentication
		return c.Post(url, body)
	} else if c.authEnabled {
		c.refreshAuthToken(resp)
	}

	return resp, err
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
