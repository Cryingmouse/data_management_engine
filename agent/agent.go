package agent

import (
	"runtime"

	"github.com/cryingmouse/data_management_engine/common"
)

type Agent interface {
	// The area method returns the area of the shape.
	CreateDirectory(hostContext common.HostContext, name string) (dirPath string, err error)
	DeleteDirectory(hostContext common.HostContext, name string) (err error)
	CreateShare(hostContext common.HostContext, name, directory_name string) (err error)
	CreateLocalUser(hostContext common.HostContext, username, password string) (err error)
	DeleteLocalUser(hostContext common.HostContext, username string) (err error)
	GetLocalUsers(hostContext common.HostContext) (users []User, err error)
	GetLocalUser(hostContext common.HostContext, username string) (user User, err error)
	GetSystemInfo(hostContext common.HostContext) (system common.SystemInfo, err error)
}

func GetAgent() Agent {
	agents := map[string]Agent{
		"windows": &WindowsAgent{},
		"linux":   &LinuxAgent{},
	}

	agent, ok := agents[runtime.GOOS]
	if !ok {
		return nil
	}

	return agent
}
