package agent

import (
	"bytes"
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

type User struct {
	Name                 string `json:"name"`
	ID                   string `json:"id"`
	Fullname             string `json:"fullname"`
	Description          string `json:"description"`
	Status               string `json:"status"`
	IsPasswordExpired    bool   `json:"is_password_expired"`
	IsPasswordChangeable bool   `json:"is_password_changeable"`
	IsPasswordRequired   bool   `json:"is_password_required"`
	IsLockout            bool   `json:"is_lockout"`
	ComputerName         string `json:"host_name"`
}

func (agent *WindowsAgent) GetDirectoryDetail(hostContext common.HostContext, path string) (detail common.DirectoryDetail, err error) {
	script := "./agent/windows/Get-DirectoryDetails.ps1"

	output, err := execPowerShellCmdlet(script, "-DirectoryPaths", path)
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

func (agent *WindowsAgent) GetDirectoriesDetail(hostContext common.HostContext, paths []string) (detail []common.DirectoryDetail, err error) {
	script := "./agent/windows/Get-DirectoryDetails.ps1"

	output, err := execPowerShellCmdlet(script, "-DirectoryPaths", strings.Join(paths, ","))
	if err != nil {
		return detail, err
	}

	var result []map[string]interface{}
	err = json.Unmarshal(output, &result)
	if err != nil {
		return detail, err
	}

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

	return detail, err
}

func (agent *WindowsAgent) CreateDirectory(hostContext common.HostContext, name string) (dirPath string, err error) {
	dirPath = fmt.Sprintf("%s\\%s", "c:\\test", name)

	err = os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return dirPath, nil
}

func (agent *WindowsAgent) CreateDirectories(hostContext common.HostContext, names []string) (dirPaths []string, err error) {
	for _, name := range names {
		dirPath, err := agent.CreateDirectory(hostContext, name)
		if err != nil {
			return dirPaths, err
		}

		dirPaths = append(dirPaths, dirPath)
	}

	return dirPaths, err
}

func (agent *WindowsAgent) DeleteDirectory(hostContext common.HostContext, name string) (err error) {
	dirPath := fmt.Sprintf("%s\\%s", "c:\\test", name)

	return os.Remove(dirPath)
}

func (agent *WindowsAgent) DeleteDirectories(hostContext common.HostContext, names []string) (err error) {
	for _, name := range names {
		if err = agent.DeleteDirectory(hostContext, name); err != nil {
			return err
		}
	}

	return err
}

func (agent *WindowsAgent) CreateShare(hostContext common.HostContext, name, directoryName string) (err error) {
	cmdlet := "New-SmbShare"

	// Define the arguments
	args := []string{
		"-Name", name,
		"-Path", directoryName,
		"-FullAccess", "Everyone",
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

func (agent *WindowsAgent) CreateLocalUser(hostContext common.HostContext, username, password string) (err error) {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("New-LocalUser -Name '%s' -Password (ConvertTo-SecureString -String '%s' -AsPlainText -Force)", username, password))
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to create local user:", err)
	}

	return err
}

func (agent *WindowsAgent) DeleteLocalUser(hostContext common.HostContext, username string) (err error) {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Remove-LocalUser -Name '%s'", username))
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to delete local user:", err)
	}

	return err
}

func (agent *WindowsAgent) GetLocalUsers(hostContext common.HostContext) (users []User, err error) {
	script := "./agent/windows/Get-LocalUserDetails.ps1"
	output, err := execPowerShellCmdlet(script)
	if err != nil {
		return nil, err
	}

	var result map[string]map[string]interface{}
	err = json.Unmarshal(output, &result)
	if err != nil {
		return nil, err
	}

	for _, value := range result {
		user := User{
			Name:                 fmt.Sprintf("%v", value["Name"]),
			Fullname:             fmt.Sprintf("%v", value["FullName"]),
			Description:          fmt.Sprintf("%v", value["Description"]),
			Status:               fmt.Sprintf("%v", value["Status"]),
			IsPasswordExpired:    value["PasswordExpires"].(bool),
			IsPasswordChangeable: value["PasswordChangeable"].(bool),
			IsPasswordRequired:   value["PasswordRequired"].(bool),
			IsLockout:            value["Lockout"].(bool),
			ComputerName:         fmt.Sprintf("%v", value["PSComputerName"]),
			ID:                   fmt.Sprintf("%v", value["SID"]),
		}
		users = append(users, user)
	}

	return users, nil
}

func (agent *WindowsAgent) GetLocalUser(hostContext common.HostContext, username string) (user User, err error) {
	script := "./agent/windows/Get-LocalUserDetails.ps1"
	output, err := execPowerShellCmdlet(script, "-UserName", username)
	if err != nil {
		return user, err
	}

	var result map[string]map[string]interface{}
	err = json.Unmarshal(output, &result)
	if err != nil {
		return user, err
	}

	for _, value := range result {
		user = User{
			Name:                 fmt.Sprintf("%v", value["Name"]),
			Fullname:             fmt.Sprintf("%v", value["FullName"]),
			Description:          fmt.Sprintf("%v", value["Description"]),
			Status:               fmt.Sprintf("%v", value["Status"]),
			IsPasswordExpired:    value["PasswordExpires"].(bool),
			IsPasswordChangeable: value["PasswordChangeable"].(bool),
			IsPasswordRequired:   value["PasswordRequired"].(bool),
			IsLockout:            value["Lockout"].(bool),
			ComputerName:         fmt.Sprintf("%v", value["PSComputerName"]),
			ID:                   fmt.Sprintf("%v", value["SID"]),
		}

		return user, nil
	}

	return user, fmt.Errorf("unable to get the user %s", username)
}

func (agent *WindowsAgent) GetSystemInfo(hostContext common.HostContext) (systemInfo common.SystemInfo, err error) {
	script := "./agent/windows/Get-SystemDetails.ps1"
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
