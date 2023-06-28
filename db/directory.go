package db

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type Directory struct {
	gorm.Model
	Name   string `gorm:"uniqueIndex:idx_name_hostip"`
	HostIp string `gorm:"uniqueIndex:idx_name_hostip"`
}

func (d *Directory) Get(engine *DatabaseEngine) error {
	// The query information should be in the instance of Directory struct pointer 'd'
	return engine.Get(d).Error
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

func (dl *DirectoryList) Get(engine *DatabaseEngine, hostIp string, nameKeyword string) error {
	conds := Directory{
		HostIp: hostIp,
	}

	db := engine.DB

	// Add the keyword to the conditions for the fuzzy search
	if nameKeyword != "" {
		db = db.Where("name LIKE ?", "%"+nameKeyword+"%")
	}

	return db.Find(&dl.Directories, conds).Error
}

type PaginationDirectory struct {
	Directories []Directory
	TotalCount  int64
}

func (dl *DirectoryList) GetByPagination(engine *DatabaseEngine, attributes []string, hostIp string, page, pageSize int) (*PaginationDirectory, error) {
	var directories []Directory
	var totalCount int64

	conds := Directory{
		HostIp: hostIp,
	}

	db := engine.DB

	// Build the SELECT statement dynamically based on the input attributes or retrieve all attributes
	if len(attributes) > 0 {
		// Validate attributes exist in the Directory struct
		var validAttributes []string
		directory := Directory{}
		for _, attr := range attributes {
			if _, ok := reflect.TypeOf(directory).FieldByName(attr); ok {
				validAttributes = append(validAttributes, attr)
			}
		}
		if len(validAttributes) == 0 {
			return nil, errors.New("no valid attributes found")
		}
		// Use the provided attributes
		selectStatement := strings.Join(validAttributes, ", ")
		db = db.Select(selectStatement)
	}

	// Perform the query to get paginated directories and total count
	if err := db.Model(&Directory{}).Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if err := db.Offset((page-1)*pageSize).Limit(pageSize).Find(&directories, conds).Error; err != nil {
		return nil, err
	}

	response := &PaginationDirectory{
		Directories: directories,
		TotalCount:  totalCount,
	}

	return response, nil
}

func (dl *DirectoryList) Delete(engine *DatabaseEngine, names []string, hostIp string) (err error) {
	var directories []Directory

	query := engine.DB
	if names != nil {
		query = query.Where("name IN [?]", names)
	}

	if hostIp != "" {
		query = query.Where("host_ip = ?", hostIp)
	}

	if err = query.Unscoped().Delete(&directories).Error; err != nil {
		return fmt.Errorf("failed to delete hosts in database: %w", err)
	}

	return nil
}
