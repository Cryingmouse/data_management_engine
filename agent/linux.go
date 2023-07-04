package agent

import (
	"fmt"
	"os"

	"github.com/cryingmouse/data_management_engine/context"
)

type LinuxAgent struct {
}

func (agent *LinuxAgent) CreateDirectory(hostContext context.HostContext, name string) (dirPath string, err error) {
	dirPath = fmt.Sprintf("%s\\%s", "c:\\test", name)

	err = os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return dirPath, nil
}

func (agent *LinuxAgent) DeleteDirectory(hostContext context.HostContext, name string) (err error) {
	dirPath := fmt.Sprintf("%s\\%s", "c:\\test", name)

	return os.Remove(dirPath)
}

func (agent *LinuxAgent) CreateShare(hostContext context.HostContext, name, directory_name string) (err error) {
	return err
}

func (agent *LinuxAgent) CreateLocalUser(hostContext context.HostContext, username, password string) (err error) {
	return err
}

func (agent *LinuxAgent) DeleteLocalUser(hostContext context.HostContext, username string) (err error) {
	return err
}

func (agent *LinuxAgent) GetLocalUsers(hostContext context.HostContext) (users []User, err error) {
	return nil, nil
}

func (agent *LinuxAgent) GetLocalUser(hostContext context.HostContext, username string) (user User, err error) {
	return user, nil
}

func (agent *LinuxAgent) GetSystemInfo(hostContext context.HostContext) (system context.SystemInfo, err error) {
	return system, nil
}
