package db

import (
	"github.com/cryingmouse/data_management_engine/context"
	"github.com/cryingmouse/data_management_engine/utils"
	"gorm.io/gorm"
)

type Host struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	Ip          string `gorm:"unique"`
	Name        string
	Username    string
	Password    string
	StorageType string
}

func (hostModel *Host) Get(engine *DatabaseEngine, name, ip string) (*Host, error) {
	host := &Host{}
	query := engine.DB

	switch {
	case name != "" && ip != "":
		query = query.Where("name = ? AND ip = ?", name, ip)
	case ip == "":
		query = query.Where("name = ?", name)
	case name == "":
		query = query.Where("ip = ?", ip)
	}

	result := query.First(host)
	if result.Error != nil {
		return nil, result.Error
	}

	password, err := utils.Decrypt(host.Password, context.SecurityKey)
	if err != nil {
		return nil, err
	}
	host.Password = password

	return host, nil
}

func (hostModel *Host) Save(engine *DatabaseEngine) error {
	encrypted_password, err := utils.Encrypt(hostModel.Password, context.SecurityKey)
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
