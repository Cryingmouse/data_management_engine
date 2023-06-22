package mgmtmodel

import (
	"github.com/cryingmouse/data_management_engine/context"
	"github.com/cryingmouse/data_management_engine/db"
	"github.com/cryingmouse/data_management_engine/driver"
)

type Directory struct {
	Name   string
	HostIp string
}

func (d *Directory) Create() (*Directory, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	hostModel := db.Host{}
	host, err := hostModel.Get(engine, "", d.HostIp)
	if err != nil {
		panic(err)
	}

	hostContext := context.HostContext{
		IP:       host.Ip,
		Username: host.Username,
		Password: host.Password,
	}

	driver := driver.GetDriver(host.StorageType)
	driver.CreateDirectory(hostContext, d.Name)

	directoryModel := db.Directory{
		Name:   d.Name,
		HostIp: host.Ip,
	}

	if err = directoryModel.Save(engine); err != nil {
		return nil, err
	}

	return nil, nil
}
