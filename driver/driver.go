package driver

import (
	"net/http"

	"github.com/cryingmouse/data_management_engine/common"
)

type Driver interface {
	// The area method returns the area of the shape.
	GetDirectoryDetail(hostContext common.HostContext, name string) (directoryDetail common.DirectoryDetail, err error)
	CreateDirectory(hostContext common.HostContext, name string) (directoryDetails common.DirectoryDetail, err error)
	DeleteDirectory(hostContext common.HostContext, name string) (err error)
	CreateShare(hostContext common.HostContext, name string) (resp *http.Response, err error)
	DeleteShare(hostContext common.HostContext, name string) (resp *http.Response, err error)
	CreateLocalUser(hostContext common.HostContext, name, password string) (resp *http.Response, err error)
	DeleteUser(hostContext common.HostContext, name string) (resp *http.Response, err error)
	GetSystemInfo(hostContext common.HostContext) (systemInfo common.SystemInfo, err error)
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
