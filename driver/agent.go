package driver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cryingmouse/data_management_engine/client"
	"github.com/cryingmouse/data_management_engine/context"
)

type AgentDriver struct {
}

func (d *AgentDriver) CreateDirectory(hostContext context.HostContext, name string) (resp *http.Response, err error) {
	restClient := client.GetRestClient(hostContext, "agent")

	// Create the request body as a string
	body := fmt.Sprintf(`{"name": "%s"}`, name)

	// Convert the string to an io.Reader
	reader := strings.NewReader(body)

	return restClient.Post("directory/create", "application/json", reader)
}

func (d *AgentDriver) CreateShare(hostContext context.HostContext, name string) (resp *http.Response, err error) {
	// TODO: Check if the root path and directory name is valid

	// Create a new folder called `newFolderName` in the current working directory.

	return nil, nil
}
