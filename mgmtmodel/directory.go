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

	dbHost := db.Host{Ip: d.HostIp}
	if err = dbHost.Get(engine); err != nil {
		panic(err)
	}

	hostContext := context.HostContext{
		IP:       dbHost.Ip,
		Username: dbHost.Username,
		Password: dbHost.Password,
	}

	driver := driver.GetDriver(dbHost.StorageType)
	driver.CreateDirectory(hostContext, d.Name)

	directoryModel := db.Directory{
		Name:   d.Name,
		HostIp: dbHost.Ip,
	}

	if err = directoryModel.Save(engine); err != nil {
		return nil, err
	}

	return nil, nil
}

func (d *Directory) Delete() (*Directory, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	dbHost := db.Host{Ip: d.HostIp}
	if err = dbHost.Get(engine); err != nil {
		panic(err)
	}

	hostContext := context.HostContext{
		IP:       dbHost.Ip,
		Username: dbHost.Username,
		Password: dbHost.Password,
	}

	driver := driver.GetDriver(dbHost.StorageType)
	driver.DeleteDirectory(hostContext, d.Name)

	directoryModel := db.Directory{
		Name:   d.Name,
		HostIp: dbHost.Ip,
	}

	if err = directoryModel.Delete(engine); err != nil {
		return nil, err
	}

	return nil, nil
}

func (d *Directory) Get() (*Directory, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	dbDirectory := db.Directory{
		Name:   d.Name,
		HostIp: d.HostIp,
	}
	if err = dbDirectory.Get(engine); err != nil {
		return nil, err
	}

	d.Name = dbDirectory.Name
	d.HostIp = dbDirectory.HostIp

	return d, nil
}

type DirectoryList struct {
	Directories []Directory
}

func (dl *DirectoryList) Get(hostIp string) ([]Directory, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	directoryListModel := db.DirectoryList{}
	directories, err := directoryListModel.Get(engine, hostIp)
	if err != nil {
		return nil, err
	}

	for _, _directory := range directories {
		directory := Directory{
			Name:   _directory.Name,
			HostIp: _directory.HostIp,
		}

		dl.Directories = append(dl.Directories, directory)
	}

	return dl.Directories, nil
}
