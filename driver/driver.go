package driver

import (
	"net/http"

	"github.com/cryingmouse/data_management_engine/common"
)

type Driver interface {
	// The area method returns the area of the shape.
	CreateDirectory(hostContext common.HostContext, name string) (resp *http.Response, err error)
	DeleteDirectory(hostContext common.HostContext, name string) (resp *http.Response, err error)
	CreateShare(hostContext common.HostContext, name string) (resp *http.Response, err error)
	DeleteShare(hostContext common.HostContext, name string) (resp *http.Response, err error)
	CreateLocalUser(hostContext common.HostContext, name, password string) (resp *http.Response, err error)
	DeleteUser(hostContext common.HostContext, name string) (resp *http.Response, err error)
	GetSystemInfo(hostContext common.HostContext) (resp *http.Response, err error)
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
