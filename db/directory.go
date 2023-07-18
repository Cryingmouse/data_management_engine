package db

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/cryingmouse/data_management_engine/common"
)

type Directory struct {
	gorm.Model
	HostIP         string `gorm:"uniqueIndex:idx_directory_unique;column:host_ip"`
	Name           string `gorm:"uniqueIndex:idx_directory_unique;column:name"`
	CreationTime   string `gorm:"column:creation_time"`
	LastAccessTime string `gorm:"column:last_access_time"`
	LastWriteTime  string `gorm:"column:last_write_time"`
	Exist          bool   `gorm:"column:exist"`
	FullPath       string `gorm:"column:full_path"`
	ParentFullPath string `gorm:"column:parent_full_path"`
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

func (dl *DirectoryList) Pagination(engine *DatabaseEngine, filter *common.QueryFilter) (paginationDirectory PaginationDirectory, err error) {

	model := Directory{}

	if filter.Pagination == nil {
		return paginationDirectory, fmt.Errorf("invalid filter: missing pagination")
	}

	var totalCount int64
	totalCount, err = Query(engine, model, filter, &dl.Directories)
	if err != nil {
		return paginationDirectory, fmt.Errorf("failed to query directories by the filter %v in the database: %w", filter, err)
	}

	paginationDirectory.Directories = dl.Directories
	paginationDirectory.TotalCount = totalCount

	return paginationDirectory, nil
}

func (dl *DirectoryList) Save(engine *DatabaseEngine) (err error) {
	if len(dl.Directories) == 0 {
		return errors.New("no directories to save")
	}

	err = engine.DB.CreateInBatches(dl.Directories, len(dl.Directories)).Error
	if err != nil {
		return fmt.Errorf("failed to save the directories in database: %w", err)
	}
	return nil
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
