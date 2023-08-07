package db

import (
	"errors"
	"fmt"

	"github.com/cryingmouse/data_management_engine/common"
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
	Connected      bool   `json:"connected,omitempty"`

	// Association for the Host's Directories using foreign key
	Directories []Directory `gorm:"foreignKey:HostIP;references:IP"`
}

// Get retrieves a Host from the database.
func (h *Host) Get(engine *DatabaseEngine) error {
	err := engine.DB.Where(h).Preload("Directories").First(h).Error
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
func (h *Host) Save(engine *DatabaseEngine) (err error) {
	h.Password, err = common.Encrypt(h.Password, common.SecurityKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	return engine.DB.Save(h).Error
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
func (hl *HostList) Pagination(engine *DatabaseEngine, filter *common.QueryFilter) (paginationHost PaginationHost, err error) {
	if filter.Pagination == nil {
		return paginationHost, errors.New("invalid filter: pagination is required")
	}

	model := Host{}
	var totalCount int64
	totalCount, err = Query(engine, model, filter, &hl.Hosts)
	if err != nil {
		return paginationHost, fmt.Errorf("failed to query hosts from the database: %w", err)
	}

	for i := range hl.Hosts {
		if hl.Hosts[i].Password != "" {
			hl.Hosts[i].Password, err = common.Decrypt(hl.Hosts[i].Password, common.SecurityKey)
			if err != nil {
				return paginationHost, fmt.Errorf("failed to decrypt password: %w", err)
			}
		}
	}

	paginationHost.Hosts = hl.Hosts
	paginationHost.TotalCount = totalCount

	return paginationHost, nil
}

// Save saves a list of Hosts to the database.
func (hl *HostList) Save(engine *DatabaseEngine) (err error) {
	if len(hl.Hosts) == 0 {
		return errors.New("UserList is empty")
	}

	for i, host := range hl.Hosts {
		hl.Hosts[i].Password, err = common.Encrypt(host.Password, common.SecurityKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt password for host %v: %w", host.ComputerName, err)
		}
	}

	err = engine.DB.CreateInBatches(hl.Hosts, len(hl.Hosts)).Error
	if err != nil {
		return fmt.Errorf("failed to save the hosts in database: %w", err)
	}
	return nil
}

// Delete deletes a list of Hosts from the database.
func (hl *HostList) Delete(engine *DatabaseEngine, filter *common.QueryFilter) error {
	var hosts []Host
	if filter != nil {
		return Delete(engine, filter, &hl.Hosts)
	}

	query := engine.DB.Unscoped()

	conditions := common.StructListToMapList(hl.Hosts)
	for _, condition := range conditions {
		query = query.Or(condition)
	}

	return query.Find(&hosts).Delete(&hosts).Error
}
