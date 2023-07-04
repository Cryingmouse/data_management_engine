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
	IP             string `gorm:"unique;column:ip"`
	Name           string `gorm:"column:name"`
	Username       string `gorm:"column:username"`
	Password       string `gorm:"type:password;column:password"`
	StorageType    string `gorm:"column:storage_type"`
	Caption        string `gorm:"column:os_type"`
	OSArchitecture string `gorm:"column:os_arch"`
	Version        string `gorm:"column:version"`
	BuildNumber    string `gorm:"column:build_number"`
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
		IP:             h.IP,
		Name:           h.Name,
		Username:       h.Username,
		StorageType:    h.StorageType,
		Caption:        h.Caption,
		OSArchitecture: h.OSArchitecture,
		Version:        h.Version,
		BuildNumber:    h.BuildNumber,
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

func (hl *HostList) Get(engine *DatabaseEngine, filter *context.QueryFilter) (err error) {
	model := Host{}

	if filter.Pagination != nil {
		return fmt.Errorf("invalid filter: there is pagination in the filter")
	}

	if _, err := Query(engine, model, filter, &hl.Hosts); err != nil {
		return fmt.Errorf("failed to query the hosts by the filter %v in database: %w", filter, err)
	}

	for _, host := range hl.Hosts {
		host.Password, err = utils.Decrypt(host.Password, context.SecurityKey)
	}

	return
}

type PaginationHost struct {
	Hosts      []Host
	TotalCount int64
}

func (hl *HostList) Pagination(engine *DatabaseEngine, filter *context.QueryFilter) (response *PaginationHost, err error) {
	var totalCount int64
	model := Host{}

	if filter.Pagination == nil {
		return response, fmt.Errorf("invalid filter: missing pagination in the filter")
	}

	totalCount, err = Query(engine, model, filter, &hl.Hosts)
	if err != nil {
		return response, fmt.Errorf("failed to query the hosts by the filter %v in database: %w", filter, err)
	}

	for _, host := range hl.Hosts {
		host.Password, err = utils.Decrypt(host.Password, context.SecurityKey)
	}

	response = &PaginationHost{
		Hosts:      hl.Hosts,
		TotalCount: totalCount,
	}

	return response, err
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
