package db

import (
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
