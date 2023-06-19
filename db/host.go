package db

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type HostInfo struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Ip          string `gorm:"unique"`
	Username    string
	Password    string
	StorageType string
}

func (host_info *HostInfo) Save(engine *DatabaseEngine) error {
	encrypted_password, err := bcrypt.GenerateFromPassword([]byte(host_info.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	host_info.Password = string(encrypted_password)

	result := engine.DB.Create(&host_info)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
