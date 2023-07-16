package agent

import (
	"context"
	"runtime"

	"github.com/cryingmouse/data_management_engine/common"
)

type Agent interface {
	// The area method returns the area of the shape.
	GetDirectoryDetail(ctx context.Context, hostContext common.HostContext, path string) (detail common.DirectoryDetail, err error)
	GetDirectoriesDetail(ctx context.Context, hostContext common.HostContext, paths []string) (detail []common.DirectoryDetail, err error)
	CreateDirectory(ctx context.Context, hostContext common.HostContext, name string) (dirPath string, err error)
	CreateDirectories(ctx context.Context, hostContext common.HostContext, names []string) (dirPaths []string, err error)
	DeleteDirectory(ctx context.Context, hostContext common.HostContext, name string) (err error)
	DeleteDirectories(ctx context.Context, hostContext common.HostContext, names []string) (err error)
	CreateShare(ctx context.Context, hostContext common.HostContext, name, directory_name string) (err error)
	CreateLocalUser(ctx context.Context, hostContext common.HostContext, username, password string) (err error)
	DeleteLocalUser(ctx context.Context, hostContext common.HostContext, username string) (err error)
	GetLocalUserDetail(ctx context.Context, hostContext common.HostContext, username string) (detail common.LocalUserDetail, err error)
	GetLocalUsersDetail(ctx context.Context, hostContext common.HostContext, usernames []string) (detail []common.LocalUserDetail, err error)
	GetSystemInfo(ctx context.Context, hostContext common.HostContext) (system common.SystemInfo, err error)
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
