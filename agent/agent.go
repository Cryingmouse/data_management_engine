package agent

import (
	"runtime"

	"github.com/cryingmouse/data_management_engine/context"
)

type Agent interface {
	// The area method returns the area of the shape.
	CreateDirectory(hostContext context.HostContext, name string) (dirPath string, err error)
	CreateShare(hostContext context.HostContext, name string) (shareName string, err error)
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
