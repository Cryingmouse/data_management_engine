package agent

import (
	"context"
	"runtime"

	"github.com/cryingmouse/data_management_engine/common"
)

type Agent interface {
	// The area method returns the area of the shape.
	GetDirectoryDetail(ctx context.Context, path string) (detail common.DirectoryDetail, err error)
	GetDirectoriesDetail(ctx context.Context, paths []string) (detail []common.DirectoryDetail, err error)
	CreateDirectory(ctx context.Context, name string) (dirPath string, err error)
	CreateDirectories(ctx context.Context, names []string) (dirPaths []string, err error)
	DeleteDirectory(ctx context.Context, name string) (err error)
	DeleteDirectories(ctx context.Context, names []string) (err error)
	CreateShare(ctx context.Context, name, directoryName, description string, usernames []string) (err error)
	DeleteShare(ctx context.Context, name string) (err error)
	GetShareDetail(ctx context.Context, name string) (detail common.ShareDetail, err error)
	GetSharesDetail(ctx context.Context, names []string) (detail []common.ShareDetail, err error)
	CreateShareMapping(ctx context.Context, deviceName, sharePath, userName, password string) (err error)
	DeleteShareMapping(ctx context.Context, deviceName string) (err error)
	CreateLocalUser(ctx context.Context, username, password string) (err error)
	DeleteLocalUser(ctx context.Context, username string) (err error)
	GetLocalUserDetail(ctx context.Context, username string) (detail common.LocalUserDetail, err error)
	GetLocalUsersDetail(ctx context.Context, usernames []string) (detail []common.LocalUserDetail, err error)
	GetSystemInfo(ctx context.Context) (system common.SystemInfo, err error)
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
