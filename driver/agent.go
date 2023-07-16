package driver

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/cryingmouse/data_management_engine/client"
	"github.com/cryingmouse/data_management_engine/common"
)

type AgentDriver struct {
}

func (d *AgentDriver) CreateDirectory(ctx context.Context, name string) (directoryDetails common.DirectoryDetail, err error) {
	hostContext := ctx.Value(common.HostContextkey("hostContext")).(common.HostContext)
	traceID := ctx.Value(common.TraceIDKey("TraceID")).(string)

	restClient := client.GetRestClient(hostContext, traceID, "agent")

	// Create the request body as a string
	request_body := fmt.Sprintf(`{"name": "%s"}`, name)

	// Convert the string to an io.Reader
	reader := strings.NewReader(request_body)

	response, err := restClient.Post("directories/create", "application/json", reader)
	if err != nil {
		directoryDetails.Name = name
		directoryDetails.Exist = false
		return directoryDetails, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		directoryDetails.Name = name
		directoryDetails.Exist = false
		return directoryDetails, fmt.Errorf("Failed")
	}

	restClient.GetResponseBody(response, &directoryDetails)

	return directoryDetails, err
}

func (d *AgentDriver) DeleteDirectory(ctx context.Context, name string) (err error) {
	hostContext := ctx.Value(common.HostContextkey("hostContext")).(common.HostContext)
	traceID := ctx.Value(common.TraceIDKey("TraceID")).(string)

	restClient := client.GetRestClient(hostContext, traceID, "agent")

	// Create the request body as a string
	body := fmt.Sprintf(`{"name": "%s"}`, name)

	// Convert the string to an io.Reader
	reader := strings.NewReader(body)

	response, err := restClient.Post("directories/delete", "application/json", reader)
	if err != nil {
		return err
	} else if response.StatusCode != http.StatusOK {
		var result common.FailedRESTResponse
		restClient.GetResponseBody(response, &result)

		return fmt.Errorf(result.Error)
	}

	return nil
}

func (d *AgentDriver) GetDirectoryDetail(ctx context.Context, name string) (detail common.DirectoryDetail, err error) {
	hostContext := ctx.Value(common.HostContextkey("hostContext")).(common.HostContext)
	traceID := ctx.Value(common.TraceIDKey("TraceID")).(string)

	restClient := client.GetRestClient(hostContext, traceID, "agent")

	url := fmt.Sprintf("directories/detail?name=%s", name)

	response, err := restClient.Get(url, "application/json")
	if err != nil {
		detail.Name = name
		detail.Exist = false
		return detail, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		detail.Name = name
		detail.Exist = false
		return detail, fmt.Errorf("Failed")
	}

	restClient.GetResponseBody(response, &detail)

	return detail, err
}

func (d *AgentDriver) GetDirectoriesDetail(ctx context.Context, names []string) (detail []common.DirectoryDetail, err error) {
	hostContext := ctx.Value(common.HostContextkey("hostContext")).(common.HostContext)
	traceID := ctx.Value(common.TraceIDKey("TraceID")).(string)

	restClient := client.GetRestClient(hostContext, traceID, "agent")

	url := fmt.Sprintf("directories/detail?name=%s", strings.Join(names, ","))

	response, err := restClient.Get(url, "application/json")
	if err != nil {
		return detail, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return detail, fmt.Errorf("Failed")
	}

	restClient.GetResponseBody(response, &detail)

	return detail, err
}

func (d *AgentDriver) CreateShare(ctx context.Context, name string) (resp *http.Response, err error) {
	// TODO: Check if the root path and directory name is valid

	// Create a new folder called `newFolderName` in the current working directory.

	return nil, nil
}

func (d *AgentDriver) DeleteShare(ctx context.Context, name string) (resp *http.Response, err error) {
	// TODO: Check if the root path and directory name is valid

	// Delete a new folder called `newFolderName` in the current working directory.

	return nil, nil
}

func (d *AgentDriver) CreateLocalUser(ctx context.Context, name, password string) (localUserDetail common.LocalUserDetail, err error) {
	hostContext := ctx.Value(common.HostContextkey("hostContext")).(common.HostContext)
	traceID := ctx.Value(common.TraceIDKey("TraceID")).(string)

	restClient := client.GetRestClient(hostContext, traceID, "agent")

	// Create the request body as a string
	request_body := fmt.Sprintf(`{"name": "%s", "password": "%s"}`, name, password)

	// Convert the string to an io.Reader
	reader := strings.NewReader(request_body)

	response, err := restClient.Post("users/create", "application/json", reader)
	if err != nil {
		localUserDetail.Name = name
		return localUserDetail, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		localUserDetail.Name = name
		return localUserDetail, fmt.Errorf("Failed")
	}

	restClient.GetResponseBody(response, &localUserDetail)

	return localUserDetail, err
}

func (d *AgentDriver) DeleteLocalUser(ctx context.Context, name string) (err error) {
	hostContext := ctx.Value(common.HostContextkey("hostContext")).(common.HostContext)
	traceID := ctx.Value(common.TraceIDKey("TraceID")).(string)

	restClient := client.GetRestClient(hostContext, traceID, "agent")

	// Create the request body as a string
	body := fmt.Sprintf(`{"name": "%s"}`, name)

	// Convert the string to an io.Reader
	reader := strings.NewReader(body)

	response, err := restClient.Post("users/delete", "application/json", reader)
	if err != nil {
		return err
	} else if response.StatusCode != http.StatusOK {
		var result common.FailedRESTResponse
		restClient.GetResponseBody(response, &result)

		return fmt.Errorf(result.Error)
	}

	return nil
}

func (d *AgentDriver) GetLocalUserDetail(ctx context.Context, name string) (detail common.LocalUserDetail, err error) {
	hostContext := ctx.Value(common.HostContextkey("hostContext")).(common.HostContext)
	traceID := ctx.Value(common.TraceIDKey("TraceID")).(string)

	restClient := client.GetRestClient(hostContext, traceID, "agent")

	escapedName := url.QueryEscape(name)
	escapedName = strings.ReplaceAll(escapedName, "+", "%20")

	url := fmt.Sprintf("users/detail?name=%s", escapedName)

	response, err := restClient.Get(url, "application/json")
	if err != nil {
		detail.Name = name
		return detail, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		detail.Name = name
		return detail, fmt.Errorf("Failed")
	}

	restClient.GetResponseBody(response, &detail)

	return detail, err
}

func (d *AgentDriver) GetLocalUsersDetail(ctx context.Context, names []string) (detail []common.LocalUserDetail, err error) {
	hostContext := ctx.Value(common.HostContextkey("hostContext")).(common.HostContext)
	traceID := ctx.Value(common.TraceIDKey("TraceID")).(string)

	restClient := client.GetRestClient(hostContext, traceID, "agent")

	escapedNames := make([]string, 0, len(names))
	for _, name := range names {
		escapedName := url.QueryEscape(name)
		escapedName = strings.ReplaceAll(escapedName, "+", "%20")
		escapedNames = append(escapedNames, escapedName)
	}

	url := fmt.Sprintf("users/detail?name=%s", strings.Join(escapedNames, ","))

	response, err := restClient.Get(url, "application/json")
	if err != nil {
		return detail, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return detail, fmt.Errorf("Failed")
	}

	restClient.GetResponseBody(response, &detail)

	return detail, err
}

func (d *AgentDriver) GetSystemInfo(ctx context.Context) (systemInfo common.SystemInfo, err error) {
	hostContext := ctx.Value(common.HostContextkey("hostContext")).(common.HostContext)
	traceID := ctx.Value(common.TraceIDKey("TraceID")).(string)

	restClient := client.GetRestClient(hostContext, traceID, "agent")

	response, err := restClient.Get("system-info", "application/json")

	err = restClient.GetResponseBody(response, &systemInfo)

	return systemInfo, err
}
