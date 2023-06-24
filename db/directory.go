package db

import (
	"errors"
	"fmt"

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

func (dl *DirectoryList) Get(engine *DatabaseEngine, hostIp string) error {
	conds := map[string]interface{}{
		"host_ip": hostIp,
	}

	return engine.DB.Find(&dl.Directories, conds).Error
}

func (dl *DirectoryList) Create(engine *DatabaseEngine) (err error) {
	if len(dl.Directories) == 0 {
		return errors.New("HostList is empty")
	}

	err = engine.DB.CreateInBatches(dl.Directories, len(dl.Directories)).Error
	if err != nil {
		return fmt.Errorf("failed to create directories in database: %w", err)
	}

	return nil
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

	err = query.Unscoped().Delete(&directories).Error
	if err != nil {
		return fmt.Errorf("failed to delete hosts in database: %w", err)
	}

	return nil
}
