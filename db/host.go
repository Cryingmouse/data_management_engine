package db

import (
	"errors"
	"fmt"

	"github.com/cryingmouse/data_management_engine/common"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

type Host struct {
	gorm.Model
	IP             string `gorm:"unique;column:ip"`
	ComputerName   string `gorm:"column:name"`
	Username       string `gorm:"column:username"`
	Password       string `gorm:"type:password;column:password"`
	StorageType    string `gorm:"column:storage_type"`
	Caption        string `gorm:"column:os_type"`
	OSArchitecture string `gorm:"column:os_arch"`
	OSVersion      string `gorm:"column:os_version"`
	BuildNumber    string `gorm:"column:build_number"`
}

// Get retrieves a Host from the database.
func (h *Host) Get(engine *DatabaseEngine) error {
	err := engine.DB.Where(h).First(h).Error
	if err != nil {
		return err
	}

	h.Password, err = common.Decrypt(h.Password, common.SecurityKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt password: %w", err)
	}

	return nil
}

// Save a Host to the database.
func (h *Host) Save(engine *DatabaseEngine) error {
	encryptedPassword, err := common.Encrypt(h.Password, common.SecurityKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	host := &Host{
		IP:             h.IP,
		ComputerName:   h.ComputerName,
		Username:       h.Username,
		StorageType:    h.StorageType,
		Caption:        h.Caption,
		OSArchitecture: h.OSArchitecture,
		OSVersion:      h.OSVersion,
		BuildNumber:    h.BuildNumber,
		Password:       string(encryptedPassword),
	}

	return engine.DB.Create(host).Error
}

// Delete a Host from the database.
func (h *Host) Delete(engine *DatabaseEngine) error {
	return engine.DB.Unscoped().Where(h).Delete(h).Error
}

type HostList struct {
	Hosts []Host
}

// Get a list of Hosts from the database.
func (hl *HostList) Get(engine *DatabaseEngine, filter *common.QueryFilter) error {
	model := Host{}

	if filter.Pagination != nil {
		return errors.New("invalid filter: pagination is not supported")
	}

	if _, err := Query(engine, model, filter, &hl.Hosts); err != nil {
		return fmt.Errorf("failed to query hosts from the database: %w", err)
	}

	for i := range hl.Hosts {
		if hl.Hosts[i].Password != "" {
			var err error
			hl.Hosts[i].Password, err = common.Decrypt(hl.Hosts[i].Password, common.SecurityKey)
			if err != nil {
				return fmt.Errorf("failed to decrypt password: %w", err)
			}
		}
	}

	return nil
}

type PaginationHost struct {
	Hosts      []Host
	TotalCount int64
}

// Pagination retrieves a paginated list of Hosts from the database.
func (hl *HostList) Pagination(engine *DatabaseEngine, filter *common.QueryFilter) (*PaginationHost, error) {
	if filter.Pagination == nil {
		return nil, errors.New("invalid filter: pagination is required")
	}

	model := Host{}
	var totalCount int64
	var err error

	totalCount, err = Query(engine, model, filter, &hl.Hosts)
	if err != nil {
		return nil, fmt.Errorf("failed to query hosts from the database: %w", err)
	}

	for i := range hl.Hosts {
		if hl.Hosts[i].Password != "" {
			var err error
			hl.Hosts[i].Password, err = common.Decrypt(hl.Hosts[i].Password, common.SecurityKey)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt password: %w", err)
			}
		}
	}

	response := &PaginationHost{
		Hosts:      hl.Hosts,
		TotalCount: totalCount,
	}

	return response, nil
}

// Save saves a list of Hosts to the database.
func (hl *HostList) Save(engine *DatabaseEngine) error {
	if len(hl.Hosts) == 0 {
		return errors.New("HostList is empty")
	}

	for i, host := range hl.Hosts {
		encryptedPassword, err := common.Encrypt(host.Password, common.SecurityKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt password for host %v: %w", host.ComputerName, err)
		}

		hl.Hosts[i].Password = string(encryptedPassword)
	}

	return engine.DB.CreateInBatches(hl.Hosts, len(hl.Hosts)).Error
}

// Delete deletes a list of Hosts from the database.
func (hl *HostList) Delete(engine *DatabaseEngine) error {
	var hosts []Host
	ips := funk.Map(hl.Hosts, func(host Host) string {
		return host.IP
	}).([]string)
	if ips != nil {
		return engine.DB.Where("ip IN (?)", ips).Unscoped().Delete(&hosts).Error
	}

	return nil
}
