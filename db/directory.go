package db

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/cryingmouse/data_management_engine/common"
)

type Directory struct {
	gorm.Model
	Name           string `gorm:"uniqueIndex:idx_name_host_ip;column:name"`
	CreationTime   string `gorm:"column:creation_time"`
	LastAccessTime string `gorm:"column:last_access_time"`
	LastWriteTime  string `gorm:"column:last_write_time"`
	Exist          bool   `gorm:"column:exist"`
	FullPath       string `gorm:"column:full_path"`
	ParentFullPath string `gorm:"column:parent_full_path"`
	HostIP         string `gorm:"uniqueIndex:idx_name_host_ip;column:host_ip"`
	// Host           Host   `gorm:"foreignKey:HostIP"`
}

func (d *Directory) Get(engine *DatabaseEngine) error {
	return engine.DB.Where(d).First(d).Error
}

func (d *Directory) Save(engine *DatabaseEngine) error {
	return engine.DB.Save(d).Error
}

func (d *Directory) Delete(engine *DatabaseEngine) error {
	return engine.DB.Unscoped().Where(d).Delete(d).Error
}

type DirectoryList struct {
	Directories []Directory
}

func (dl *DirectoryList) Get(engine *DatabaseEngine, filter *common.QueryFilter) error {
	model := Directory{}

	if filter.Pagination != nil {
		return fmt.Errorf("invalid filter: pagination is not supported")
	}

	if _, err := Query(engine, model, filter, &dl.Directories); err != nil {
		return fmt.Errorf("failed to query the directories by the filter %v in database: %w", filter, err)
	}

	return nil
}

type PaginationDirectory struct {
	Directories []Directory
	TotalCount  int64
}

func (dl *DirectoryList) Pagination(engine *DatabaseEngine, filter *common.QueryFilter) (*PaginationDirectory, error) {
	var totalCount int64
	model := Directory{}

	if filter.Pagination == nil {
		return nil, fmt.Errorf("invalid filter: missing pagination")
	}

	totalCount, err := Query(engine, model, filter, &dl.Directories)
	if err != nil {
		return nil, fmt.Errorf("failed to query directories by the filter %v in the database: %w", filter, err)
	}

	response := &PaginationDirectory{
		Directories: dl.Directories,
		TotalCount:  totalCount,
	}

	return response, nil
}

func (dl *DirectoryList) Save(engine *DatabaseEngine) error {
	if len(dl.Directories) == 0 {
		return errors.New("no directories to save")
	}

	return engine.DB.CreateInBatches(dl.Directories, len(dl.Directories)).Error
}

func (dl *DirectoryList) Delete(engine *DatabaseEngine, filter *common.QueryFilter) error {
	var directories []Directory
	if filter != nil {
		return Delete(engine, filter, &dl.Directories)
	}

	query := engine.DB.Unscoped()

	conditions := common.StructListToMapList(dl.Directories)
	for _, condition := range conditions {
		query = query.Or(condition)
	}

	return query.Find(&directories).Delete(&directories).Error
}
