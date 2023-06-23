package db

import (
	"fmt"

	"gorm.io/gorm"
)

type Directory struct {
	gorm.Model
	Name   string `gorm:"uniqueIndex:idx_name_hostip"`
	HostIp string `gorm:"uniqueIndex:idx_name_hostip"`
}

func (model *Directory) Save(engine *DatabaseEngine) error {
	result := engine.DB.Save(&model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (model *Directory) Get(engine *DatabaseEngine, name, hostIp string) (directory *Directory, err error) {
	directory = &Directory{}
	query := engine.DB

	query = query.Where("name = ? AND host_ip = ?", name, hostIp)

	result := query.First(directory)
	if result.Error != nil {
		return nil, result.Error
	}

	return directory, nil
}

type DirectoryList struct{}

func (dl *DirectoryList) Get(engine *DatabaseEngine, hostIp string) ([]Directory, error) {
	var directories []Directory

	if hostIp != "" {
		conds := map[string]interface{}{
			"host_ip": hostIp,
		}
		result := engine.DB.Find(&directories, conds)
		if result.Error != nil {
			return nil, result.Error
		}
	} else {
		result := engine.DB.Find(&directories)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	return directories, nil
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
