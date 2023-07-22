package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cryingmouse/data_management_engine/common"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type WindowsAgent struct{}

func (agent *WindowsAgent) CreateDirectory(ctx context.Context, name string) (dirPath string, err error) {
	dirPath = fmt.Sprintf("%s\\%s", common.Config.Agent.WindowsRootFolder, name)

	err = os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return dirPath, nil
}

func (agent *WindowsAgent) CreateDirectories(ctx context.Context, names []string) (dirPaths []string, err error) {
	for _, name := range names {
		dirPath, err := agent.CreateDirectory(ctx, name)
		if err != nil {
			return dirPaths, err
		}

		dirPaths = append(dirPaths, dirPath)
	}

	return dirPaths, err
}

func (agent *WindowsAgent) DeleteDirectory(ctx context.Context, name string) (err error) {
	dirPath := fmt.Sprintf("%s\\%s", common.Config.Agent.WindowsRootFolder, name)

	return os.Remove(dirPath)
}

func (agent *WindowsAgent) GetDirectoryList(ctx context.Context) (err error) {
	cmdlet := "Get-ChildItem"

	// Define the arguments
	args := []string{
		common.Config.Agent.WindowsRootFolder,
	}

	// Execute the PowerShell command
	cmd := exec.Command("powershell.exe", append([]string{"-Command", cmdlet}, args...)...)

	// Capture the command output
	_, err = cmd.CombinedOutput()

	return err
}

func (agent *WindowsAgent) DeleteDirectories(ctx context.Context, names []string) (err error) {
	for _, name := range names {
		if err = agent.DeleteDirectory(ctx, name); err != nil {
			return err
		}
	}

	return err
}

func (agent *WindowsAgent) GetDirectoryDetail(ctx context.Context, name string) (detail common.DirectoryDetail, err error) {
	script := "./agent/windows/Get-DirectoryDetail.ps1"

	dirPath := fmt.Sprintf("%s\\%s", common.Config.Agent.WindowsRootFolder, name)

	output, err := execPowerShellCmdlet(script, "-DirectoryPaths", dirPath)
	if err != nil {
		return detail, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(output, &result)
	if err != nil {
		return detail, err
	}

	detail = common.DirectoryDetail{
		Name:           result["Name"].(string),
		FullPath:       result["FullPath"].(string),
		CreationTime:   result["CreationTime"].(string),
		LastWriteTime:  result["LastWriteTime"].(string),
		LastAccessTime: result["LastAccessTime"].(string),
		Exist:          result["Exist"].(bool),
		ParentFullPath: result["ParentFullPath"].(string),
	}

	return detail, err
}

func (agent *WindowsAgent) GetDirectoriesDetail(ctx context.Context, names []string) (detail []common.DirectoryDetail, err error) {
	script := "./agent/windows/Get-DirectoryDetail.ps1"

	dirPaths := make([]string, len(names))
	for i, name := range names {
		dirPaths[i] = fmt.Sprintf("%s\\%s", common.Config.Agent.WindowsRootFolder, name)
	}

	output, err := execPowerShellCmdlet(script, "-DirectoryPaths", strings.Join(dirPaths, ","))
	if err != nil {
		return detail, err
	}

	var result []map[string]interface{}
	err = json.Unmarshal(output, &result)
	if err != nil {
		// In case that only one file path to query.
		var result map[string]interface{}
		err = json.Unmarshal(output, &result)
		if err != nil {
			return detail, err
		}
		directory := common.DirectoryDetail{
			Name:           result["Name"].(string),
			FullPath:       result["FullPath"].(string),
			CreationTime:   result["CreationTime"].(string),
			LastWriteTime:  result["LastWriteTime"].(string),
			LastAccessTime: result["LastAccessTime"].(string),
			Exist:          result["Exist"].(bool),
			ParentFullPath: result["ParentFullPath"].(string),
		}
		detail = append(detail, directory)
	} else {

		for _, item := range result {
			directory := common.DirectoryDetail{
				Name:           item["Name"].(string),
				FullPath:       item["FullPath"].(string),
				CreationTime:   item["CreationTime"].(string),
				LastWriteTime:  item["LastWriteTime"].(string),
				LastAccessTime: item["LastAccessTime"].(string),
				Exist:          item["Exist"].(bool),
				ParentFullPath: item["ParentFullPath"].(string),
			}
			detail = append(detail, directory)
		}
	}

	return detail, err
}

func (agent *WindowsAgent) CreateCIFSShare(ctx context.Context, name, directoryName, description string, usernames []string) (err error) {
	cmdlet := "New-SmbShare"

	directoryPath := fmt.Sprintf("%s\\%s", common.Config.Agent.WindowsRootFolder, directoryName)

	// Define the arguments
	args := []string{
		"-Name", name,
		"-Path", common.AddQuotes(directoryPath),
		"-Description", common.AddQuotes(description),
		"-FullAccess", strings.Join(usernames, ", "),
		"-FolderEnumerationMode", "Unrestricted",
	}

	// Execute the PowerShell command
	cmd := exec.Command("powershell.exe", append([]string{"-Command", cmdlet}, args...)...)

	// Capture the command output
	_, err = cmd.CombinedOutput()

	return err
}

func (agent *WindowsAgent) DeleteCIFSShare(ctx context.Context, name string) (err error) {
	cmdlet := "Remove-SmbShare"

	// Define the arguments
	args := []string{
		"-Name", name,
		"-Force",
	}

	// Execute the PowerShell command
	cmd := exec.Command("powershell.exe", append([]string{"-Command", cmdlet}, args...)...)

	// Capture the command output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing PowerShell command:", err.Error())
	}

	fmt.Println(string(output))
	return err
}

func (agent *WindowsAgent) GetCIFSShareDetail(ctx context.Context, name string) (detail common.ShareDetail, err error) {
	script := "./agent/windows/Get-ShareDetail.ps1"
	output, err := execPowerShellCmdlet(script, "-ShareNames", name)
	if err != nil {
		return detail, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(output, &result)
	if err != nil {
		return detail, err
	}

	detail = common.ShareDetail{
		Name:          result["Name"].(string),
		DirectoryPath: result["DirectoryPath"].(string),
		Description:   result["Description"].(string),
		State: func() string {
			if result["ShareState"].(float64) == 1 {
				return "online"
			}
			return "offline"
		}(),
	}

	return detail, err
}

func (agent *WindowsAgent) GetCIFSSharesDetail(ctx context.Context, names []string) (detail []common.ShareDetail, err error) {
	script := "./agent/windows/Get-ShareDetail.ps1"
	output, err := execPowerShellCmdlet(script, "-ShareNames", strings.Join(names, ","))
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	err = json.Unmarshal(output, &result)
	if err != nil {
		// In case that only one file path to query.
		var result map[string]interface{}
		err = json.Unmarshal(output, &result)
		if err != nil {
			return detail, err
		}
		share := common.ShareDetail{
			Name:          result["Name"].(string),
			DirectoryPath: result["DirectoryPath"].(string),
			Description:   result["Description"].(string),
			State: func() string {
				if result["ShareState"].(float64) == 1 {
					return "online"
				}
				return "offline"
			}(),
		}
		detail = append(detail, share)
	}

	for _, item := range result {
		share := common.ShareDetail{
			Name:          item["Name"].(string),
			DirectoryPath: item["DirectoryPath"].(string),
			Description:   item["Description"].(string),
			State: func() string {
				if item["ShareState"].(float64) == 1 {
					return "online"
				}
				return "offline"
			}(),
		}
		detail = append(detail, share)
	}

	return detail, err
}

func (agent *WindowsAgent) MountCIFSShare(ctx context.Context, deviceName, sharePath, userName, password string) (err error) {
	cmdlet := "net"

	// Define the arguments
	args := []string{
		"use",
		deviceName,
		sharePath,
		password,
		"/user:" + userName,
	}

	// Execute the PowerShell command
	cmd := exec.Command("powershell.exe", append([]string{"-Command", cmdlet}, args...)...)

	// Capture the command output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing PowerShell command:", err.Error())
	}

	fmt.Println(string(output))
	return err
}

func (agent *WindowsAgent) UnmountCIFSShare(ctx context.Context, deviceName string) (err error) {
	cmdlet := "net"

	// Define the arguments
	args := []string{
		"use",
		deviceName,
		"/delete",
		"/y",
	}

	// Execute the PowerShell command
	cmd := exec.Command("powershell.exe", append([]string{"-Command", cmdlet}, args...)...)

	// Capture the command output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing PowerShell command:", err.Error())
	}

	fmt.Println(string(output))
	return err
}

func (agent *WindowsAgent) CreateLocalUser(ctx context.Context, name, password string) (err error) {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("New-LocalUser -Name '%s' -Password (ConvertTo-SecureString -String '%s' -AsPlainText -Force)", name, password))
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to create local user:", err)
	}

	return err
}

func (agent *WindowsAgent) DeleteLocalUser(ctx context.Context, name string) (err error) {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Remove-LocalUser -Name '%s'", name))
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to delete local user:", err)
	}

	return err
}

func (agent *WindowsAgent) GetLocalUserDetail(ctx context.Context, name string) (detail common.LocalUserDetail, err error) {
	script := "./agent/windows/Get-LocalUserDetail.ps1"
	output, err := execPowerShellCmdlet(script, "-UserName", name)
	if err != nil {
		return detail, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(output, &result)
	if err != nil {
		return detail, err
	}

	detail = common.LocalUserDetail{
		Name:                 result["Name"].(string),
		UID:                  result["SID"].(string),
		FullName:             result["FullName"].(string),
		Description:          result["Description"].(string),
		Status:               result["Status"].(string),
		IsPasswordExpired:    result["PasswordExpires"].(bool),
		IsPasswordChangeable: result["PasswordChangeable"].(bool),
		IsPasswordRequired:   result["PasswordRequired"].(bool),
		IsLockout:            result["Lockout"].(bool),
		IsDisabled:           result["Disabled"].(bool),
	}

	return detail, err
}

func (agent *WindowsAgent) GetLocalUsersDetail(ctx context.Context, names []string) (detail []common.LocalUserDetail, err error) {
	script := "./agent/windows/Get-LocalUserDetail.ps1"
	output, err := execPowerShellCmdlet(script)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	err = json.Unmarshal(output, &result)
	if err != nil {
		return detail, err
	}

	for _, item := range result {
		localUser := common.LocalUserDetail{
			Name:                 item["Name"].(string),
			UID:                  item["SID"].(string),
			FullName:             item["FullName"].(string),
			Description:          item["Description"].(string),
			Status:               item["Status"].(string),
			IsPasswordExpired:    item["PasswordExpires"].(bool),
			IsPasswordChangeable: item["PasswordChangeable"].(bool),
			IsPasswordRequired:   item["PasswordRequired"].(bool),
			IsLockout:            item["Lockout"].(bool),
			IsDisabled:           item["Disabled"].(bool),
		}
		detail = append(detail, localUser)
	}

	return detail, err
}

func (agent *WindowsAgent) GetSystemInfo(ctx context.Context) (systemInfo common.SystemInfo, err error) {
	script := "./agent/windows/Get-SystemDetail.ps1"
	output, err := execPowerShellCmdlet(script)
	if err != nil {
		return systemInfo, err
	}

	result := make(map[string]string)
	err = json.Unmarshal(output, &result)
	if err != nil {
		return systemInfo, err
	}

	systemInfo = common.SystemInfo{
		ComputerName:   result["ComputerName"],
		Caption:        result["Caption"],
		OSArchitecture: result["OSArchitecture"],
		OSVersion:      result["Version"],
		BuildNumber:    result["BuildNumber"],
	}

	return systemInfo, err
}

func execPowerShellCmdlet(script string, args ...string) (output []byte, err error) {
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", script)
	cmd.Args = append(cmd.Args, args...)
	cmd.Dir, err = os.Getwd()
	if err != nil {
		return nil, err
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err = cmd.Run(); err != nil {
		return nil, err
	}

	if stderr.Len() > 0 {
		return nil, fmt.Errorf(stderr.String())
	}

	outputBytes := stdout.Bytes()
	decoder := simplifiedchinese.GB18030.NewDecoder()
	outputStr, _, err := transform.String(decoder, string(outputBytes))
	if err != nil {
		return nil, err
	}

	return []byte(outputStr), nil
}
