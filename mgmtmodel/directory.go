package mgmtmodel

import (
	"context"
	"errors"

	"github.com/cryingmouse/data_management_engine/common"
	"github.com/cryingmouse/data_management_engine/db"
	"github.com/cryingmouse/data_management_engine/driver"
	"golang.org/x/sync/errgroup"
)

type Directory struct {
	Name           string
	HostIP         string
	CreationTime   string
	LastAccessTime string
	LastWriteTime  string
	Exist          bool
	FullPath       string
	ParentFullPath string
}

func (d *Directory) Create() (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	// Get the right driver and call driver to create directory.
	host := db.Host{IP: d.HostIP}
	if err = host.Get(engine); err != nil {
		return err
	}
	driver := driver.GetDriver(host.StorageType)

	hostContext := common.HostContext{
		IP:       host.IP,
		Username: host.Username,
		Password: host.Password,
	}
	directoryDetails, err := driver.CreateDirectory(hostContext, d.Name)
	if err != nil {
		return err
	}

	// Save the details of the directory into database.
	directory := db.Directory{
		Name:           d.Name,
		HostIP:         host.IP,
		CreationTime:   directoryDetails.CreationTime,
		LastAccessTime: directoryDetails.LastAccessTime,
		LastWriteTime:  directoryDetails.LastWriteTime,
		Exist:          directoryDetails.Exist,
		FullPath:       directoryDetails.FullPath,
		ParentFullPath: directoryDetails.ParentFullPath,
	}

	common.CopyStructList(directory, d)

	return directory.Save(engine)
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
	if err := driver.DeleteDirectory(hostContext, d.Name); err != nil {
		return err
	}

	directory := db.Directory{
		Name:   d.Name,
		HostIP: host.IP,
	}
	return directory.Delete(engine)
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

	common.CopyStructList(directoryList.Directories, &dl.Directories)

	return dl.Directories, nil
}

func (dl *DirectoryList) Create() error {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	input := make([]interface{}, len(dl.Directories))
	for index, directory := range dl.Directories {
		input[index] = directory
	}

	g, _ := errgroup.WithContext(context.Background())

	results := make([]common.DirectoryDetail, len(dl.Directories))
	var resultErr error

	for _, d := range dl.Directories {
		directory := d // 避免闭包问题
		g.Go(func() error {
			host := db.Host{IP: directory.HostIP}
			if err = host.Get(engine); err != nil {
				resultErr = errors.Join(resultErr, err)
				return err
			}

			hostContext := common.HostContext{
				IP:       host.IP,
				Username: host.Username,
				Password: host.Password,
			}

			driver := driver.GetDriver(host.StorageType)
			directoryDetail, err := driver.CreateDirectory(hostContext, directory.Name)
			if err != nil {
				resultErr = errors.Join(resultErr, err)
				return err
			}

			results = append(results, directoryDetail) // 保存协程的返回值

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		if resultErr != nil {
			return resultErr
		}
		return err
	}

	if err := common.CopyStructList(results, &dl.Directories); err != nil {
		return err
	}

	// Save to database.
	dbDirectoryList := db.DirectoryList{}

	if err := common.CopyStructList(dl.Directories, &dbDirectoryList.Directories); err != nil {
		return err
	}

	return dbDirectoryList.Save(engine)
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

	g, _ := errgroup.WithContext(context.Background())

	results := make([]Directory, len(dl.Directories))
	var resultErr error

	for i, d := range dl.Directories {
		index := i     // 避免闭包问题
		directory := d // 避免闭包问题
		g.Go(func() error {
			host := db.Host{IP: directory.HostIP}
			if err = host.Get(engine); err != nil {
				resultErr = errors.Join(resultErr, err)
				return err
			}

			hostContext := common.HostContext{
				IP:       host.IP,
				Username: host.Username,
				Password: host.Password,
			}

			driver := driver.GetDriver(host.StorageType)
			if err := driver.DeleteDirectory(hostContext, directory.Name); err != nil {
				resultErr = errors.Join(resultErr, err)
				return err
			}

			results[index] = directory // 保存协程的返回值

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		if resultErr != nil {
			return resultErr
		}
		return err
	}

	if err := common.CopyStructList(results, &dl.Directories); err != nil {
		return err
	}

	dbDirectoryList := db.DirectoryList{}

	if err := common.CopyStructList(dl.Directories, &dbDirectoryList.Directories); err != nil {
		return err
	}

	directoryList := db.DirectoryList{}

	return directoryList.Delete(engine, filter)
}
