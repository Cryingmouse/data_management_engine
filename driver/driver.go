package driver

import (
	"net/http"

	"github.com/cryingmouse/data_management_engine/context"
)

type Driver interface {
	// The area method returns the area of the shape.
	CreateDirectory(hostContext context.HostContext, name string) (resp *http.Response, err error)
	DeleteDirectory(hostContext context.HostContext, name string) (resp *http.Response, err error)
	CreateShare(hostContext context.HostContext, name string) (resp *http.Response, err error)
	DeleteShare(hostContext context.HostContext, name string) (resp *http.Response, err error)
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
