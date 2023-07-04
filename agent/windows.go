package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/cryingmouse/data_management_engine/context"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type WindowsAgent struct {
}

func (agent *WindowsAgent) CreateDirectory(hostContext context.HostContext, name string) (dirPath string, err error) {
	dirPath = fmt.Sprintf("%s\\%s", "c:\\test", name)

	err = os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return dirPath, nil
}

func (agent *WindowsAgent) DeleteDirectory(hostContext context.HostContext, name string) (err error) {
	dirPath := fmt.Sprintf("%s\\%s", "c:\\test", name)

	return os.Remove(dirPath)
}

func (agent *WindowsAgent) CreateShare(hostContext context.HostContext, name, directory_name string) (err error) {
	cmdlet := "New-SmbShare"

	// Define the arguments
	args := []string{
		"-Name", name,
		"-Path", directory_name,
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

func (agent *WindowsAgent) CreateLocalUser(hostContext context.HostContext, username, password string) (err error) {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("New-LocalUser -Name '%s' -Password (ConvertTo-SecureString -String '%s' -AsPlainText -Force)", username, password))
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to create local user:", err)
	}

	return err
}

func (agent *WindowsAgent) DeleteLocalUser(hostContext context.HostContext, username string) (err error) {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Remove-LocalUser -Name '%s'", username))
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to create local user:", err)

	}

	return err
}

func (agent *WindowsAgent) GetLocalUsers(hostContext context.HostContext) (users []User, err error) {
	// 设置要执行的脚本和参数
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

		// Print the key and value of the current map entry
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

func (agent *WindowsAgent) GetLocalUser(hostContext context.HostContext, username string) (user User, err error) {
	// 设置要执行的脚本和参数
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

func (agent *WindowsAgent) GetSystemInfo(hostContext context.HostContext) (systemInfo context.SystemInfo, err error) {
	// 设置要执行的脚本和参数
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

	systemInfo = context.SystemInfo{
		ComputerName:   result["ComputerName"],
		Caption:        result["Caption"],
		OSArchitecture: result["OSArchitecture"],
		Version:        result["Version"],
		BuildNumber:    result["BuildNumber"],
	}

	return systemInfo, nil
}

func execPowerShellCmdlet(script string, args ...string) (output []byte, err error) {
	// 创建一个执行命令的Cmd对象
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", script)
	cmd.Args = append(cmd.Args, args...)

	// 设置命令的工作目录
	cmd.Dir, err = os.Getwd()
	if err != nil {
		return nil, err
	}

	// 创建一个缓冲区来收集输出
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// 将缓冲区分配给命令对象
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 执行命令
	if err = cmd.Run(); err != nil {
		return nil, err
	}

	if stderr.Len() > 0 {
		return nil, fmt.Errorf(stderr.String())
	}

	// 将输出按照UTF-8编码转换为字符串
	outputBytes := stdout.Bytes()
	decoder := simplifiedchinese.GB18030.NewDecoder()
	outputStr, _, err := transform.String(decoder, string(outputBytes))
	if err != nil {
		return nil, err
	}

	return []byte(outputStr), nil
}
