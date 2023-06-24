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

func (d *Directory) Create() (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	host := db.Host{Ip: d.HostIp}
	if err = host.Get(engine); err != nil {
		return err
	}

	hostContext := context.HostContext{
		IP:       host.Ip,
		Username: host.Username,
		Password: host.Password,
	}

	driver := driver.GetDriver(host.StorageType)
	driver.CreateDirectory(hostContext, d.Name)

	directory := db.Directory{
		Name:   d.Name,
		HostIp: host.Ip,
	}

	if err = directory.Save(engine); err != nil {
		return err
	}

	return nil
}

func (d *Directory) Delete() (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	host := db.Host{Ip: d.HostIp}
	if err = host.Get(engine); err != nil {
		return err
	}

	hostContext := context.HostContext{
		IP:       host.Ip,
		Username: host.Username,
		Password: host.Password,
	}

	driver := driver.GetDriver(host.StorageType)
	driver.DeleteDirectory(hostContext, d.Name)

	directory := db.Directory{
		Name:   d.Name,
		HostIp: host.Ip,
	}

	if err = directory.Delete(engine); err != nil {
		return err
	}

	return nil
}

func (d *Directory) Get() (*Directory, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	directory := db.Directory{
		Name:   d.Name,
		HostIp: d.HostIp,
	}
	if err = directory.Get(engine); err != nil {
		return nil, err
	}

	d.Name = directory.Name
	d.HostIp = directory.HostIp

	return d, nil
}

type DirectoryList struct {
	Directories []Directory
}

func (dl *DirectoryList) Get(hostIp string) ([]Directory, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	directoryList := db.DirectoryList{}
	err = directoryList.Get(engine, hostIp)
	if err != nil {
		return nil, err
	}

	for _, _directory := range directoryList.Directories {
		directory := Directory{
			Name:   _directory.Name,
			HostIp: _directory.HostIp,
		}

		dl.Directories = append(dl.Directories, directory)
	}

	return dl.Directories, nil
}
