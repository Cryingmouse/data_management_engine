package driver

import (
	"context"
	"net/http"

	"github.com/cryingmouse/data_management_engine/common"
)

type Driver interface {
	// The area method returns the area of the shape.
	GetDirectoryDetail(ctx context.Context, name string) (directoryDetail common.DirectoryDetail, err error)
	CreateDirectory(ctx context.Context, name string) (directoryDetails common.DirectoryDetail, err error)
	DeleteDirectory(ctx context.Context, name string) (err error)
	CreateShare(ctx context.Context, name string) (resp *http.Response, err error)
	DeleteShare(ctx context.Context, name string) (resp *http.Response, err error)
	CreateLocalUser(ctx context.Context, name, password string) (resp *http.Response, err error)
	DeleteUser(ctx context.Context, name string) (resp *http.Response, err error)
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
