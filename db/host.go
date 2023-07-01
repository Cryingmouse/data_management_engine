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
	IP          string `gorm:"unique;column:ip"`
	Name        string `gorm:"column:name"`
	Username    string `gorm:"column:username"`
	Password    string `gorm:"type:password;column:password"`
	StorageType string `gorm:"column:storage_type"`
}

func (h *Host) Get(engine *DatabaseEngine) (err error) {
	if err = engine.DB.Where(h).First(h).Error; err != nil {
		return err
	}

	h.Password, err = utils.Decrypt(h.Password, context.SecurityKey)

	return err
}

func (h *Host) Save(engine *DatabaseEngine) error {
	host := Host{
		IP:          h.IP,
		Name:        h.Name,
		Username:    h.Username,
		StorageType: h.StorageType,
	}

	encrypted_password, err := utils.Encrypt(h.Password, context.SecurityKey)
	if err != nil {
		return err
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

func (hl *HostList) Get(engine *DatabaseEngine, storageType string) (err error) {
	conds := Host{
		StorageType: storageType,
	}

	if err = engine.DB.Find(&hl.Hosts, conds).Error; err != nil {
		return err
	}

	for i := range hl.Hosts {
		decryptedPassword, err := utils.Decrypt(hl.Hosts[i].Password, context.SecurityKey)
		if err != nil {
			return err
		}

		hl.Hosts[i].Password = decryptedPassword
	}

	return nil
}

func (hl *HostList) Save(engine *DatabaseEngine) (err error) {
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

func (hl *HostList) Delete(engine *DatabaseEngine, storageType string, names, ips []string) error {
	var host Host

	query := engine.DB.Where("1 = 1")
	if storageType != "" {
		query = engine.DB.Where("storage_type = ?", storageType)
	}
	if names != nil {
		query.Where("name IN [?]", names)
	}

	if ips != nil {
		query.Where("ip IN [?]", ips)
	}

	return query.Unscoped().Delete(&host).Error
}
