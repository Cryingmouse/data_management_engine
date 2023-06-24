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

func (h *Host) Get(engine *DatabaseEngine) (err error) {
	if err = engine.Get(h).Error; err != nil {
		return err
	}

	h.Password, err = utils.Decrypt(h.Password, context.SecurityKey)

	return err
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
