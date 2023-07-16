package driver

import (
	"context"
	"net/http"

	"github.com/cryingmouse/data_management_engine/common"
)

type Driver interface {
	CreateDirectory(ctx context.Context, name string) (directoryDetails common.DirectoryDetail, err error)

	DeleteDirectory(ctx context.Context, name string) (err error)

	GetDirectoryDetail(ctx context.Context, name string) (detail common.DirectoryDetail, err error)

	GetDirectoriesDetail(ctx context.Context, names []string) (detail []common.DirectoryDetail, err error)

	CreateShare(ctx context.Context, name string) (resp *http.Response, err error)

	DeleteShare(ctx context.Context, name string) (resp *http.Response, err error)

	CreateLocalUser(ctx context.Context, name, password string) (localUserDetail common.LocalUserDetail, err error)

	DeleteLocalUser(ctx context.Context, name string) (err error)

	GetLocalUserDetail(ctx context.Context, name string) (detail common.LocalUserDetail, err error)

	GetLocalUsersDetail(ctx context.Context, names []string) (detail []common.LocalUserDetail, err error)

	GetSystemInfo(ctx context.Context) (systemInfo common.SystemInfo, err error)
}

func GetDriver(storageType string) Driver {
	drivers := map[string]Driver{
		"agent": &AgentDriver{},
	}

	driver, ok := drivers[storageType]
	if !ok {
		return nil
	}

	return driver
}
