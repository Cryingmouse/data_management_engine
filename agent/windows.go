package agent

import (
	"fmt"
	"os"

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

func (agent *WindowsAgent) CreateShare(hostContext context.HostContext, name string) (shareName string, err error) {
	// TODO: Check if the root path and directory name is valid

	// Create a new folder called `newFolderName` in the current working directory.

	return "", nil
}
