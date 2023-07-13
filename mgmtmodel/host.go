package mgmtmodel

import (
	"context"
	"errors"
	"fmt"

	"github.com/cryingmouse/data_management_engine/common"
	"github.com/cryingmouse/data_management_engine/db"
	"github.com/cryingmouse/data_management_engine/driver"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/sync/errgroup"
)

type Host struct {
	ComputerName   string
	IP             string
	Username       string
	Password       string
	StorageType    string
	Caption        string
	OSArchitecture string
	OSVersion      string
	BuildNumber    string
}

func (h *Host) Register() error {
	systemInfo, err := h.getSystemInfo()
	if err != nil {
		return err
	}

	// Update mgmtmodel with system information
	h.ComputerName = systemInfo.ComputerName
	h.Caption = systemInfo.Caption
	h.OSArchitecture = systemInfo.OSArchitecture
	h.OSVersion = systemInfo.OSVersion
	h.BuildNumber = systemInfo.BuildNumber

	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	var host db.Host

	common.CopyStructList(h, &host)

	err = host.Save(engine)
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		switch sqliteErr.ExtendedCode {
		// Map SQLite ErrNo to specific error scenarios
		case sqlite3.ErrConstraintUnique: // SQLite constraint violation
			return fmt.Errorf("the host %v has already been registered", host.IP)
		}
	}

	return err
}

func (h *Host) Unregister() error {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	// Retur error if there is any directory on the host which will be unregistered.
	directoryList := DirectoryList{}
	filter := common.QueryFilter{
		Conditions: struct {
			HostIP string
		}{
			HostIP: h.IP,
		},
	}
	if directories, err := directoryList.Get(&filter); err != nil || len(directories) != 0 {
		return err
	}

	// Delete host from database.
	host := db.Host{IP: h.IP}
	return host.Delete(engine)
}

func (h *Host) Get() (*Host, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	host := db.Host{
		ComputerName: h.ComputerName,
		IP:           h.IP,
	}
	if err = host.Get(engine); err != nil {
		return nil, err
	}

	common.CopyStructList(host, h)

	return h, nil
}

type HostList struct {
	Hosts []Host
}

func (hl *HostList) Register() error {
	g, _ := errgroup.WithContext(context.Background())

	results := make([]common.SystemInfo, len(hl.Hosts))
	var resultErr error

	for i, h := range hl.Hosts {
		index := i // 避免闭包问题
		host := h  // 避免闭包问题
		g.Go(func() error {
			systemInfo, err := host.getSystemInfo()
			if err != nil {
				resultErr = errors.Join(resultErr, err)
				return err
			}
			results[index] = systemInfo // 保存协程的返回值

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		if resultErr != nil {
			return resultErr
		}
		return err
	}

	if err := common.CopyStructList(results, &hl.Hosts); err != nil {
		return err
	}

	dbHostList := db.HostList{}

	if err := common.CopyStructList(hl.Hosts, &dbHostList.Hosts); err != nil {
		return err
	}

	if engine, err := db.GetDatabaseEngine(); err != nil {
		return err
	} else {
		err := dbHostList.Save(engine)
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			switch sqliteErr.ExtendedCode {
			// Map SQLite ErrNo to specific error scenarios
			case sqlite3.ErrConstraintUnique: // SQLite constraint violation
				return fmt.Errorf("some hosts have already been registered")
			}
		}

		return err
	}
}

func (hl *HostList) Unregister() error {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	// TODO: Return error if there is any related directories.

	hostList := db.HostList{}

	common.CopyStructList(hl.Hosts, &hostList.Hosts)

	return hostList.Delete(engine)
}

func (hl *HostList) Get(filter *common.QueryFilter) ([]Host, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	hostList := db.HostList{}
	if err := hostList.Get(engine, filter); err != nil {
		return nil, err
	}

	common.CopyStructList(hostList.Hosts, &hl.Hosts)

	return hl.Hosts, nil
}

type PaginationHost struct {
	Hosts      []Host
	Page       int
	Limit      int
	TotalCount int64
}

func (hl *HostList) Pagination(filter *common.QueryFilter) (*PaginationHost, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	hostList := db.HostList{}
	paginationHosts, err := hostList.Pagination(engine, filter)
	if err != nil {
		return nil, err
	}

	paginationHostList := PaginationHost{
		Page:       filter.Pagination.Page,
		Limit:      filter.Pagination.PageSize,
		TotalCount: paginationHosts.TotalCount,
	}

	common.CopyStructList(paginationHosts.Hosts, &paginationHostList.Hosts)

	return &paginationHostList, nil
}

func (h *Host) getSystemInfo() (systemInfo common.SystemInfo, err error) {
	hostContext := common.HostContext{
		IP:       h.IP,
		Username: h.Username,
		Password: h.Password,
	}
	driver := driver.GetDriver(h.StorageType)

	return driver.GetSystemInfo(hostContext)
}
