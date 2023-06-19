package db

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type HostInfo struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex:idx_name_ip"`
	Ip          string `gorm:"uniqueIndex:idx_name_ip"`
	Username    string
	Password    string
	StorageType string
}

func (host_info *HostInfo) GetHosts(engine *DatabaseEngine) ([]HostInfo, error) {
	var hosts []HostInfo
	result := engine.DB.Find(&hosts)
	if result.Error != nil {
		return nil, result.Error
	}
	return hosts, nil
}

func (host_info *HostInfo) Get(engine *DatabaseEngine, name, ip string) (*HostInfo, error) {
	host := &HostInfo{}
	query := engine.DB
	query = query.Where("name = ? AND ip = ?", name, ip)
	result := query.First(host)
	if result.Error != nil {
		return nil, result.Error
	}
	return host, nil
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

func (host_info *HostInfo) Delete(engine *DatabaseEngine, name, ip string) error {
	result := engine.DB.Delete(&host_info)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
