package mgmtmodel

import (
	"context"
	"fmt"
	"strings"

	"github.com/cryingmouse/data_management_engine/common"
	"github.com/cryingmouse/data_management_engine/db"
	"github.com/cryingmouse/data_management_engine/driver"
)

type CIFSShare struct {
	Name            string
	HostIP          string
	SharePath       string
	DirectoryName   string
	Description     string
	MountPoint      string
	AccessUserNames []string
}

func (c *CIFSShare) Create(ctx context.Context) (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	// Get the right driver and call driver to create share.
	host := db.Host{IP: c.HostIP}
	if err = host.Get(engine); err != nil {
		return err
	}
	driver := driver.GetDriver(host.StorageType)

	hostContext := common.HostContext{
		IP:       host.IP,
		Username: host.Username,
		Password: host.Password,
	}
	ctx = context.WithValue(ctx, common.HostContextkey("hostContext"), hostContext)
	if err = driver.CreateCIFSShare(ctx, c.Name, c.DirectoryName, c.Description, c.AccessUserNames); err != nil {
		return err
	}

	c.SharePath = buildCIFSSharePath(host.IP, c.Name)

	share := db.CIFSShare{}

	if err = common.DeepCopy(c, &share); err != nil {
		return err
	}
	share.AccessUserNames = strings.Join(c.AccessUserNames, ",")

	return share.Save(engine)
}

func (c *CIFSShare) Delete(ctx context.Context) (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	host := db.Host{IP: c.HostIP}
	if err = host.Get(engine); err != nil {
		return err
	}

	hostContext := common.HostContext{
		IP:       host.IP,
		Username: host.Username,
		Password: host.Password,
	}
	ctx = context.WithValue(ctx, common.HostContextkey("hostContext"), hostContext)

	driver := driver.GetDriver(host.StorageType)
	if err := driver.DeleteCIFSShare(ctx, c.Name); err != nil {
		return err
	}

	c.SharePath = buildCIFSSharePath(host.IP, c.Name)
	share := db.CIFSShare{
		HostIP: c.HostIP,
		Name:   c.Name,
		Path:   c.SharePath,
	}
	return share.Delete(engine)
}

func (c *CIFSShare) Get(ctx context.Context) (*CIFSShare, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	share := db.CIFSShare{
		Name:   c.Name,
		HostIP: c.HostIP,
	}
	if err = share.Get(engine); err != nil {
		return nil, err
	}

	common.DeepCopy(share, c)
	c.AccessUserNames = strings.Split(share.AccessUserNames, ",")

	return c, nil
}

func (c *CIFSShare) Mount(ctx context.Context, userName, password string) (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	// Get the right driver and call driver to create share.
	host := db.Host{IP: c.HostIP}
	if err = host.Get(engine); err != nil {
		return err
	}
	driver := driver.GetDriver(host.StorageType)

	hostContext := common.HostContext{
		IP:       host.IP,
		Username: host.Username,
		Password: host.Password,
	}
	ctx = context.WithValue(ctx, common.HostContextkey("hostContext"), hostContext)
	if err = driver.MountCIFSShare(ctx, c.MountPoint, c.SharePath, userName, password); err != nil {
		return err
	}

	share := db.CIFSShare{
		HostIP: c.HostIP,
		Path:   c.SharePath,
	}
	share.Get(engine)

	share.MountPoint = c.MountPoint

	return share.Save(engine)
}

func (c *CIFSShare) Unmount(ctx context.Context) (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	// Get the right driver and call driver to create share.
	host := db.Host{IP: c.HostIP}
	if err = host.Get(engine); err != nil {
		return err
	}
	driver := driver.GetDriver(host.StorageType)

	hostContext := common.HostContext{
		IP:       host.IP,
		Username: host.Username,
		Password: host.Password,
	}
	ctx = context.WithValue(ctx, common.HostContextkey("hostContext"), hostContext)

	if err = driver.UnmountCIFSShare(ctx, c.MountPoint); err != nil {
		return err
	}

	share := db.CIFSShare{
		HostIP:     c.HostIP,
		MountPoint: c.MountPoint,
	}
	share.Get(engine)

	share.MountPoint = ""

	return share.Save(engine)
}

func buildCIFSSharePath(ip string, shareName string) string {
	return fmt.Sprintf("\\\\%s\\%s", ip, shareName)
}

type CIFSShareList struct {
	Shares []CIFSShare
}

func (cl *CIFSShareList) Get(ctx context.Context, filter *common.QueryFilter) ([]CIFSShare, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	shareList := db.CIFSShareList{}

	if err = shareList.Get(engine, filter); err != nil {
		return nil, err
	}

	common.DeepCopy(shareList.Shares, &cl.Shares)

	for index, share := range shareList.Shares {
		cl.Shares[index].AccessUserNames = strings.Split(share.AccessUserNames, ",")
	}

	return cl.Shares, nil
}

type PaginationShare struct {
	Shares     []CIFSShare
	Page       int
	Limit      int
	TotalCount int64
}

func (cl *CIFSShareList) Pagination(ctx context.Context, filter *common.QueryFilter) (*PaginationShare, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	shareList := db.CIFSShareList{}
	paginationShares, err := shareList.Pagination(engine, filter)
	if err != nil {
		return nil, err
	}

	paginationShareList := PaginationShare{
		Page:       filter.Pagination.Page,
		Limit:      filter.Pagination.PageSize,
		TotalCount: paginationShares.TotalCount,
	}

	for _, _share := range paginationShares.Shares {
		share := CIFSShare{
			Name:            _share.Name,
			HostIP:          _share.HostIP,
			SharePath:       _share.Path,
			DirectoryName:   _share.DirectoryName,
			Description:     _share.Description,
			AccessUserNames: strings.Split(_share.AccessUserNames, ","),
		}

		paginationShareList.Shares = append(paginationShareList.Shares, share)
	}

	return &paginationShareList, nil
}
