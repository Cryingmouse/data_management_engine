package agent

import (
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

func (agent *WindowsAgent) DeleteDirectory(hostContext context.HostContext, name string) (dirPath string, err error) {
	dirPath = fmt.Sprintf("%s\\%s", "c:\\test", name)

	err = os.Remove(dirPath)
	if err != nil {
		return "", err
	}

	return dirPath, nil
}

func (agent *WindowsAgent) CreateShare(hostContext context.HostContext, name, directory_name string) (shareName string, err error) {
	// TODO: Check if the root path and directory name is valid

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
		return "", err
	}

	// Print the command output
	fmt.Println(string(output))

	return "", nil
}
