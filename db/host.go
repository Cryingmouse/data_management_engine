package db

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Host struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex:idx_name_ip"`
	Ip          string `gorm:"uniqueIndex:idx_name_ip"`
	Username    string
	Password    string
	StorageType string
}

func (hostModel *Host) Get(engine *DatabaseEngine, name, ip string) (*Host, error) {
	host := &Host{}
	query := engine.DB
	query = query.Where("name = ? AND ip = ?", name, ip)
	result := query.First(host)
	if result.Error != nil {
		return nil, result.Error
	}
	return host, nil
}

func (hostModel *Host) Save(engine *DatabaseEngine) error {
	encrypted_password, err := bcrypt.GenerateFromPassword([]byte(hostModel.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	hostModel.Password = string(encrypted_password)

	result := engine.DB.Save(&hostModel)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (hostInfo *Host) Delete(engine *DatabaseEngine) error {
	result := engine.DB.Unscoped().Delete(&hostInfo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

type HostList struct {
	HostList []Host
}

func (hostListModel *HostList) Get(engine *DatabaseEngine) ([]Host, error) {
	var hosts []Host
	result := engine.DB.Find(&hosts)
	if result.Error != nil {
		return nil, result.Error
	}
	return hosts, nil
}
