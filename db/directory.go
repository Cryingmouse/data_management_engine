package db

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/cryingmouse/data_management_engine/common"
)

type Directory struct {
	gorm.Model
	Name   string `gorm:"uniqueIndex:idx_name_host_ip;column:name"`
	HostIP string `gorm:"uniqueIndex:idx_name_host_ip;column:host_ip"`
}

func (d *Directory) Get(engine *DatabaseEngine) error {
	// The query information should be in the instance of Directory struct pointer 'd'
	return engine.DB.Where(d).First(d).Error
}

func (d *Directory) Save(engine *DatabaseEngine) error {
	return engine.DB.Save(d).Error
}

func (d *Directory) Delete(engine *DatabaseEngine) error {
	return engine.DB.Unscoped().Delete(&d, d).Error
}

type DirectoryList struct {
	Directories []Directory
}

func (dl *DirectoryList) Get(engine *DatabaseEngine, filter *common.QueryFilter) (err error) {
	model := Directory{}

	if filter.Pagination != nil {
		return fmt.Errorf("invalid filter: there is pagination in the filter")
	}

	if _, err := Query(engine, model, filter, &dl.Directories); err != nil {
		return fmt.Errorf("failed to query the directories by the filter %v in database: %w", filter, err)
	}

	return
}

type PaginationDirectory struct {
	Directories []Directory
	TotalCount  int64
}

func (dl *DirectoryList) Pagination(engine *DatabaseEngine, filter *common.QueryFilter) (response *PaginationDirectory, err error) {
	var totalCount int64
	model := Directory{}

	if filter.Pagination == nil {
		return response, fmt.Errorf("invalid filter: missing pagination in the filter")
	}

	totalCount, err = Query(engine, model, filter, &dl.Directories)
	if err != nil {
		return response, fmt.Errorf("failed to query the directories by the filter %v in database: %w", filter, err)
	}

	response = &PaginationDirectory{
		Directories: dl.Directories,
		TotalCount:  totalCount,
	}

	return response, err
}

func (dl *DirectoryList) Save(engine *DatabaseEngine) (err error) {
	if len(dl.Directories) == 0 {
		return errors.New("directories are empty")
	}

	err = engine.DB.CreateInBatches(dl.Directories, len(dl.Directories)).Error
	if err != nil {
		return fmt.Errorf("failed to save the directories in database: %w", err)
	}
	return
}

func (dl *DirectoryList) Delete(engine *DatabaseEngine, filter *common.QueryFilter) (err error) {
	var directories []Directory

	if filter != nil {
		err = Delete(engine, filter, directories)
		if err != nil {
			return fmt.Errorf("failed to delete directories by the filter %v in database: %w", filter, err)
		}
	} else {
		if err := engine.DB.Delete(dl.Directories).Error; err != nil {
			return err
		}
	}

	return
}
