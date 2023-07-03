package agent

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	// Execute the PowerShell command
	cmd := exec.Command("powershell", "Get-LocalUser | Select-Object -Property Name")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing PowerShell command:", err.Error())
		return nil, err
	}

	local_users := parseLocalUserOutput(string(output))
	for _, user := range local_users {
		users = append(users, User{Name: user})
	}

	return users, nil
}

func parseLocalUserOutput(output string) []string {
	lines := strings.Split(output, "\r\n")
	users := make([]string, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Name") || strings.HasPrefix(line, "----") {
			continue
		}
		users = append(users, line)
	}

	return users
}
