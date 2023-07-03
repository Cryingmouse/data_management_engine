package agent

import (
	"runtime"

	"github.com/cryingmouse/data_management_engine/context"
)

type Agent interface {
	// The area method returns the area of the shape.
	CreateDirectory(hostContext context.HostContext, name string) (dirPath string, err error)
	DeleteDirectory(hostContext context.HostContext, name string) (err error)
	CreateShare(hostContext context.HostContext, name, directory_name string) (err error)
	CreateLocalUser(hostContext context.HostContext, username, password string) (err error)
	DeleteLocalUser(hostContext context.HostContext, username string) (err error)
	GetLocalUsers(hostContext context.HostContext) (users []User, err error)
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
