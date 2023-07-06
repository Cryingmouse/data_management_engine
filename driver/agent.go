package driver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cryingmouse/data_management_engine/client"
	"github.com/cryingmouse/data_management_engine/common"
)

type AgentDriver struct {
}

func (d *AgentDriver) CreateDirectory(hostContext common.HostContext, name string) (resp *http.Response, err error) {
	restClient := client.GetRestClient(hostContext, "agent")

	// Create the request body as a string
	body := fmt.Sprintf(`{"name": "%s"}`, name)

	// Convert the string to an io.Reader
	reader := strings.NewReader(body)

	return restClient.Post("directory/create", "application/json", reader)
}

func (d *AgentDriver) DeleteDirectory(hostContext common.HostContext, name string) (resp *http.Response, err error) {
	restClient := client.GetRestClient(hostContext, "agent")

	// Create the request body as a string
	body := fmt.Sprintf(`{"name": "%s"}`, name)

	// Convert the string to an io.Reader
	reader := strings.NewReader(body)

	return restClient.Post("directory/delete", "application/json", reader)
}

func (d *AgentDriver) CreateShare(hostContext common.HostContext, name string) (resp *http.Response, err error) {
	// TODO: Check if the root path and directory name is valid

	// Create a new folder called `newFolderName` in the current working directory.

	return nil, nil
}

func (d *AgentDriver) DeleteShare(hostContext common.HostContext, name string) (resp *http.Response, err error) {
	// TODO: Check if the root path and directory name is valid

	// Delete a new folder called `newFolderName` in the current working directory.

	return nil, nil
}

func (d *AgentDriver) CreateLocalUser(hostContext common.HostContext, name, password string) (resp *http.Response, err error) {
	restClient := client.GetRestClient(hostContext, "agent")

	// Create the request body as a string
	body := fmt.Sprintf(`{"name": "%s", "password": "%s"}`, name, password)

	// Convert the string to an io.Reader
	reader := strings.NewReader(body)

	return restClient.Post("user/create", "application/json", reader)
}

func (d *AgentDriver) DeleteUser(hostContext common.HostContext, name string) (resp *http.Response, err error) {
	restClient := client.GetRestClient(hostContext, "agent")

	// Create the request body as a string
	body := fmt.Sprintf(`{"name": "%s"}`, name)

	// Convert the string to an io.Reader
	reader := strings.NewReader(body)

	return restClient.Post("user/delete", "application/json", reader)
}

func (d *AgentDriver) GetSystemInfo(hostContext common.HostContext) (resp *http.Response, err error) {
	restClient := client.GetRestClient(hostContext, "agent")

	resp, err = restClient.Get("system-info", "application/json")
	return resp, err
}
