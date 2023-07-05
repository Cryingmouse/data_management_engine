package mgmtmodel

import (
	"github.com/cryingmouse/data_management_engine/common"
	"github.com/cryingmouse/data_management_engine/db"
	"github.com/cryingmouse/data_management_engine/driver"
)

type Directory struct {
	Name   string
	HostIP string
}

func (d *Directory) Create() (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	host := db.Host{IP: d.HostIP}
	if err = host.Get(engine); err != nil {
		return err
	}

	hostContext := common.HostContext{
		IP:       host.IP,
		Username: host.Username,
		Password: host.Password,
	}

	driver := driver.GetDriver(host.StorageType)
	driver.CreateDirectory(hostContext, d.Name)

	directory := db.Directory{
		Name:   d.Name,
		HostIP: host.IP,
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

	host := db.Host{IP: d.HostIP}
	if err = host.Get(engine); err != nil {
		return err
	}

	hostContext := common.HostContext{
		IP:       host.IP,
		Username: host.Username,
		Password: host.Password,
	}

	driver := driver.GetDriver(host.StorageType)
	driver.DeleteDirectory(hostContext, d.Name)

	directory := db.Directory{
		Name:   d.Name,
		HostIP: host.IP,
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
		HostIP: d.HostIP,
	}
	if err = directory.Get(engine); err != nil {
		return nil, err
	}

	d.Name = directory.Name
	d.HostIP = directory.HostIP

	return d, nil
}

type DirectoryList struct {
	Directories []Directory
}

func (dl *DirectoryList) Get(filter *common.QueryFilter) ([]Directory, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	directoryList := db.DirectoryList{}

	if err = directoryList.Get(engine, filter); err != nil {
		return nil, err
	}

	for _, _directory := range directoryList.Directories {
		directory := Directory{
			Name:   _directory.Name,
			HostIP: _directory.HostIP,
		}

		dl.Directories = append(dl.Directories, directory)
	}

	return dl.Directories, nil
}

type PaginationDirectory struct {
	Directories []Directory
	Page        int
	Limit       int
	TotalCount  int64
}

func (dl *DirectoryList) Pagination(filter *common.QueryFilter) (*PaginationDirectory, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	directoryList := db.DirectoryList{}
	paginationDirs, err := directoryList.Pagination(engine, filter)
	if err != nil {
		return nil, err
	}

	paginationDirList := PaginationDirectory{
		Page:       filter.Pagination.Page,
		Limit:      filter.Pagination.PageSize,
		TotalCount: paginationDirs.TotalCount,
	}

	for _, _directory := range paginationDirs.Directories {
		directory := Directory{
			Name:   _directory.Name,
			HostIP: _directory.HostIP,
		}

		paginationDirList.Directories = append(paginationDirList.Directories, directory)
	}

	return &paginationDirList, nil
}

func (dl *DirectoryList) Save() (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	directoryList := db.DirectoryList{}

	for _, _directory := range dl.Directories {
		directory := db.Directory{
			Name:   _directory.Name,
			HostIP: _directory.HostIP,
		}

		directoryList.Directories = append(directoryList.Directories, directory)
	}

	return directoryList.Save(engine)
}

func (dl *DirectoryList) Delete(filter *common.QueryFilter) (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	directoryList := db.DirectoryList{}

	return directoryList.Delete(engine, filter)
}
