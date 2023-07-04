package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/cryingmouse/data_management_engine/context"
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
	Name string
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

	var result map[string]interface{}
	err = json.Unmarshal(output, &result)
	if err != nil {
	}

	// 使用解析后的map进行操作
	fmt.Println("解析后的map:", result)

	return users, nil
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
	err = cmd.Run()
	if err != nil {
	}

	// 检查错误输出
	if stderr.Len() > 0 {
	}

	return stdout.Bytes(), err
}
