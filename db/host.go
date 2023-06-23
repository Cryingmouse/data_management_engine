package db

import (
	"errors"
	"fmt"

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
	Password    string `gorm:"type:password"`
	StorageType string
}

func (hostModel *Host) Get(engine *DatabaseEngine, name, ip string) (*Host, error) {
	host := Host{
		Name: name,
		Ip:   ip,
	}

	err := engine.Get(&host).Error
	if err != nil {
		return nil, err
	}

	if host.Password, err = utils.Decrypt(host.Password, context.SecurityKey); err != nil {
		return nil, err
	}

	return &host, nil
}

func (h *Host) Save(engine *DatabaseEngine) error {
	host := Host{
		Ip:          h.Ip,
		Name:        h.Name,
		Username:    h.Username,
		StorageType: h.StorageType,
	}

	encrypted_password, err := utils.Encrypt(h.Password, context.SecurityKey)
	if err != nil {
		panic(err)
	}
	host.Password = string(encrypted_password)

	return engine.DB.Save(&host).Error
}

func (h *Host) Delete(engine *DatabaseEngine) error {
	return engine.DB.Unscoped().Delete(&h, h).Error
}

type HostList struct {
	Hosts []Host
}

func (hl *HostList) Get(engine *DatabaseEngine, storageType string) ([]Host, error) {
	var hosts []Host
	conds := Host{
		StorageType: storageType,
	}

	err := engine.DB.Find(&hosts, conds).Error
	if err != nil {
		return nil, err
	}
	return hosts, nil
}

func (hl *HostList) Create(engine *DatabaseEngine) (err error) {
	if len(hl.Hosts) == 0 {
		return errors.New("HostList is empty")
	}

	for i, host := range hl.Hosts {
		// Encrypt the password
		hl.Hosts[i].Password, err = utils.Encrypt(host.Password, context.SecurityKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt password for host %v: %w", host.Name, err)
		}
	}

	err = engine.DB.CreateInBatches(hl.Hosts, len(hl.Hosts)).Error
	if err != nil {
		return fmt.Errorf("failed to create hosts in database: %w", err)
	}
	return nil
}

func (hl *HostList) Delete(engine *DatabaseEngine, storageType string, names, ips []string) (err error) {
	var hosts []Host

	query := engine.DB.Where("storage_type = ?", storageType)
	if names != nil {
		query.Where("name IN [?]", names)
	}

	if ips != nil {
		query.Where("ip IN [?]", ips)
	}

	err = query.Unscoped().Delete(&hosts).Error
	if err != nil {
		return fmt.Errorf("failed to delete hosts in database: %w", err)
	}
	return nil
}
