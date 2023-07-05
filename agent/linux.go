package agent

import (
	"fmt"
	"os"

	"github.com/cryingmouse/data_management_engine/common"
)

type LinuxAgent struct {
}

func (agent *LinuxAgent) CreateDirectory(hostContext common.HostContext, name string) (dirPath string, err error) {
	dirPath = fmt.Sprintf("%s\\%s", "c:\\test", name)

	err = os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return dirPath, nil
}

func (agent *LinuxAgent) DeleteDirectory(hostContext common.HostContext, name string) (err error) {
	dirPath := fmt.Sprintf("%s\\%s", "c:\\test", name)

	return os.Remove(dirPath)
}

func (agent *LinuxAgent) CreateShare(hostContext common.HostContext, name, directory_name string) (err error) {
	return err
}

func (agent *LinuxAgent) CreateLocalUser(hostContext common.HostContext, username, password string) (err error) {
	return err
}

func (agent *LinuxAgent) DeleteLocalUser(hostContext common.HostContext, username string) (err error) {
	return err
}

func (agent *LinuxAgent) GetLocalUsers(hostContext common.HostContext) (users []User, err error) {
	return nil, nil
}

func (agent *LinuxAgent) GetLocalUser(hostContext common.HostContext, username string) (user User, err error) {
	return user, nil
}

func (agent *LinuxAgent) GetSystemInfo(hostContext common.HostContext) (system common.SystemInfo, err error) {
	return system, nil
}
