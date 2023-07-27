package driver

import (
	"context"
	"encoding/json"
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

	restClient := client.GetRestClient("http", hostContext, 8080, "agent", "", traceID, false)

	// Create the request body as a string
	request_body := fmt.Sprintf(`{"name": "%s"}`, name)

	// Convert the string to an io.Reader
	reader := strings.NewReader(request_body)

	response, err := restClient.Post("directories/create", reader)
	if err != nil {
		return directoryDetails, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return directoryDetails, fmt.Errorf("failed to create the directory on agent")
	}

	restClient.GetResponseBody(response, &directoryDetails)

	return directoryDetails, err
}

func (d *AgentDriver) DeleteDirectory(ctx context.Context, name string) (err error) {
	hostContext := ctx.Value(common.HostContextkey("hostContext")).(common.HostContext)
	traceID := ctx.Value(common.TraceIDKey("TraceID")).(string)

	restClient := client.GetRestClient("http", hostContext, 8080, "agent", "", traceID, false)

	// Create the request body as a string
	body := fmt.Sprintf(`{"name": "%s"}`, name)

	// Convert the string to an io.Reader
	reader := strings.NewReader(body)

	response, err := restClient.Post("directories/delete", reader)
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

	restClient := client.GetRestClient("http", hostContext, 8080, "agent", "", traceID, false)

	url := fmt.Sprintf("directories/detail?name=%s", name)

	response, err := restClient.Get(url)
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

	restClient := client.GetRestClient("http", hostContext, 8080, "agent", "", traceID, false)

	url := fmt.Sprintf("directories/detail?name=%s", strings.Join(names, ","))

	response, err := restClient.Get(url)
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

func (d *AgentDriver) CreateCIFSShare(ctx context.Context, name, directory_name, description string, usernames []string) (err error) {
	hostContext := ctx.Value(common.HostContextkey("hostContext")).(common.HostContext)
	traceID := ctx.Value(common.TraceIDKey("TraceID")).(string)

	restClient := client.GetRestClient("http", hostContext, 8080, "agent", "", traceID, false)

	body := struct {
		ShareName     string   `json:"share_name"`
		DirectoryName string   `json:"directory_name"`
		Description   string   `json:"description"`
		Usernames     []string `json:"usernames"`
	}{
		ShareName:     name,
		DirectoryName: directory_name,
		Description:   description,
		Usernames:     usernames,
	}
	request_body, err := json.Marshal(body)
	if err != nil {
		return
	}

	// Convert the string to an io.Reader
	reader := strings.NewReader(string(request_body))

	response, err := restClient.Post("shares/create", reader)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create a cifs share")
	}

	return nil
}

func (d *AgentDriver) DeleteCIFSShare(ctx context.Context, name string) (err error) {
	hostContext := ctx.Value(common.HostContextkey("hostContext")).(common.HostContext)
	traceID := ctx.Value(common.TraceIDKey("TraceID")).(string)

	restClient := client.GetRestClient("http", hostContext, 8080, "agent", "", traceID, false)

	body := struct {
		ShareName string `json:"share_name"`
	}{
		ShareName: name,
	}
	request_body, err := json.Marshal(body)
	if err != nil {
		return
	}

	// Convert the string to an io.Reader
	reader := strings.NewReader(string(request_body))

	response, err := restClient.Post("shares/delete", reader)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete a cifs share")
	}

	return nil
}

func (d *AgentDriver) MountCIFSShare(ctx context.Context, mountPoint, sharePath, userName, password string) (err error) {
	hostContext := ctx.Value(common.HostContextkey("hostContext")).(common.HostContext)
	traceID := ctx.Value(common.TraceIDKey("TraceID")).(string)

	restClient := client.GetRestClient("http", hostContext, 8080, "agent", "", traceID, false)

	encryptedPassword, _ := common.Encrypt(password, common.SecurityKey)

	body := struct {
		MountPoint string `json:"mount_point"`
		SharePath  string `json:"share_path"`
		UserName   string `json:"username"`
		Password   string `json:"password"`
	}{
		MountPoint: mountPoint,
		SharePath:  sharePath,
		UserName:   userName,
		Password:   encryptedPassword,
	}
	request_body, err := json.Marshal(body)
	if err != nil {
		return
	}

	// Convert the string to an io.Reader
	reader := strings.NewReader(string(request_body))

	response, err := restClient.Post("shares/mount", reader)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to mount a cifs share")
	}

	return nil
}

func (d *AgentDriver) UnmountCIFSShare(ctx context.Context, mountPoint string) (err error) {
	hostContext := ctx.Value(common.HostContextkey("hostContext")).(common.HostContext)
	traceID := ctx.Value(common.TraceIDKey("TraceID")).(string)

	restClient := client.GetRestClient("http", hostContext, 8080, "agent", "", traceID, false)

	body := struct {
		MountPoint string `json:"mount_point"`
	}{
		MountPoint: mountPoint,
	}
	request_body, err := json.Marshal(body)
	if err != nil {
		return
	}

	// Convert the string to an io.Reader
	reader := strings.NewReader(string(request_body))

	response, err := restClient.Post("shares/unmount", reader)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to unmount a cifs share")
	}

	return nil
}

func (d *AgentDriver) CreateLocalUser(ctx context.Context, name, password string) (localUserDetail common.LocalUserDetail, err error) {
	hostContext := ctx.Value(common.HostContextkey("hostContext")).(common.HostContext)
	traceID := ctx.Value(common.TraceIDKey("TraceID")).(string)

	restClient := client.GetRestClient("http", hostContext, 8080, "agent", "", traceID, false)

	// Create the request body as a string
	request_body := fmt.Sprintf(`{"name": "%s", "password": "%s"}`, name, password)

	// Convert the string to an io.Reader
	reader := strings.NewReader(request_body)

	response, err := restClient.Post("users/create", reader)
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

	restClient := client.GetRestClient("http", hostContext, 8080, "agent", "", traceID, false)

	// Create the request body as a string
	body := fmt.Sprintf(`{"name": "%s"}`, name)

	// Convert the string to an io.Reader
	reader := strings.NewReader(body)

	response, err := restClient.Post("users/delete", reader)
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

	restClient := client.GetRestClient("http", hostContext, 8080, "agent", "", traceID, false)

	escapedName := url.QueryEscape(name)
	escapedName = strings.ReplaceAll(escapedName, "+", "%20")

	url := fmt.Sprintf("users/detail?name=%s", escapedName)

	response, err := restClient.Get(url)
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

	restClient := client.GetRestClient("http", hostContext, 8080, "agent", "", traceID, false)

	escapedNames := make([]string, 0, len(names))
	for _, name := range names {
		escapedName := url.QueryEscape(name)
		escapedName = strings.ReplaceAll(escapedName, "+", "%20")
		escapedNames = append(escapedNames, escapedName)
	}

	url := fmt.Sprintf("users/detail?name=%s", strings.Join(escapedNames, ","))

	response, err := restClient.Get(url)
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

	restClient := client.GetRestClient("http", hostContext, 8080, "agent", "", traceID, false)

	response, err := restClient.Get("system-info")
	if err != nil {
		return systemInfo, err
	}

	err = restClient.GetResponseBody(response, &systemInfo)

	return systemInfo, err
}
