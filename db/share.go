package db

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/cryingmouse/data_management_engine/common"
)

type CIFSShare struct {
	gorm.Model
	Name            string `gorm:"column:name"`
	HostIP          string `gorm:"column:host_ip"`
	Path            string `gorm:"uniqueIndex:idx_cifs_share_unique;column:path"`
	DirectoryName   string `gorm:"column:directory_name"`
	Description     string `gorm:"column:description"`
	AccessUserNames string `gorm:"column:access_usernames"`
}

func (c *CIFSShare) Get(engine *DatabaseEngine) error {
	return engine.DB.Where(c).First(c).Error
}

func (c *CIFSShare) Save(engine *DatabaseEngine) error {
	return engine.DB.Save(c).Error
}

func (c *CIFSShare) Delete(engine *DatabaseEngine) error {
	return engine.DB.Unscoped().Where(c).Delete(c).Error
}

type CIFSShareList struct {
	Shares []CIFSShare
}

func (cl *CIFSShareList) Get(engine *DatabaseEngine, filter *common.QueryFilter) error {
	model := CIFSShare{}

	if filter.Pagination != nil {
		return fmt.Errorf("invalid filter: pagination is not supported")
	}

	if _, err := Query(engine, model, filter, &cl.Shares); err != nil {
		return fmt.Errorf("failed to query the shares by the filter %v in database: %w", filter, err)
	}

	return nil
}

type PaginationCIFSShare struct {
	Shares     []CIFSShare
	TotalCount int64
}

func (cl *CIFSShareList) Pagination(engine *DatabaseEngine, filter *common.QueryFilter) (paginationShare PaginationCIFSShare, err error) {

	model := CIFSShare{}

	if filter.Pagination == nil {
		return paginationShare, fmt.Errorf("invalid filter: missing pagination")
	}

	var totalCount int64
	totalCount, err = Query(engine, model, filter, &cl.Shares)
	if err != nil {
		return paginationShare, fmt.Errorf("failed to query shares by the filter %v in the database: %w", filter, err)
	}

	paginationShare.Shares = cl.Shares
	paginationShare.TotalCount = totalCount

	return paginationShare, nil
}
