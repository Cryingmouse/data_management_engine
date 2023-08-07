package mgmtmodel

import (
	"context"
	"errors"
	"strings"

	"github.com/cryingmouse/data_management_engine/common"
	"github.com/cryingmouse/data_management_engine/db"
	"github.com/cryingmouse/data_management_engine/driver"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/sync/errgroup"
)

type Host struct {
	ID             uint   `json:"id,omitempty"`
	ComputerName   string `json:"computer_name,omitempty"`
	IP             string `json:"ip,omitempty"`
	Username       string `json:"username,omitempty"`
	Password       string `json:"password,omitempty"`
	StorageType    string `json:"storage_type,omitempty"`
	Caption        string `json:"caption,omitempty"`
	OSArchitecture string `json:"os_arch,omitempty"`
	OSVersion      string `json:"os_verion,omitempty"`
	BuildNumber    string `json:"build_number,omitempty"`
	Connected      bool   `json:"connected,omitempty"`

	Directories []Directory `json:"directories,omitempty"`
}

func (h *Host) Register(ctx context.Context) error {
	systemInfo, err := h.GetSystemInfo(ctx)
	if err != nil {
		return err
	}

	// Update mgmtmodel with system information
	h.ComputerName = systemInfo.ComputerName
	h.Caption = systemInfo.Caption
	h.OSArchitecture = systemInfo.OSArchitecture
	h.OSVersion = systemInfo.OSVersion
	h.BuildNumber = systemInfo.BuildNumber
	h.Connected = true

	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	var host db.Host

	common.DeepCopy(h, &host)

	err = host.Save(engine)
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		switch sqliteErr.ExtendedCode {
		// Map SQLite ErrNo to specific error scenarios
		case sqlite3.ErrConstraintUnique: // SQLite constraint violation
			error := common.ErrHostAlreadyRegistered
			error.Params = []string{host.IP}
			return error
		}
	}

	return nil
}

func (h *Host) Unregister(ctx context.Context) error {
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
	if directories, err := directoryList.Get(ctx, &filter); err != nil || len(directories) != 0 {
		return err
	}

	// Delete host from database.
	host := db.Host{IP: h.IP}
	return host.Delete(engine)
}

func (h *Host) Get(ctx context.Context) (*Host, error) {
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

	common.DeepCopy(host, h)

	return h, nil
}

type HostList struct {
	Hosts []Host
}

func (hl *HostList) getIPList() []string {
	ipList := make([]string, len(hl.Hosts))
	for i, host := range hl.Hosts {
		ipList[i] = host.IP
	}
	return ipList
}

func (hl *HostList) Register(ctx context.Context) error {
	g, _ := errgroup.WithContext(context.Background())

	results := make([]common.SystemInfo, len(hl.Hosts))
	var resultErr error

	for i, h := range hl.Hosts {
		index := i // 避免闭包问题
		host := h  // 避免闭包问题
		g.Go(func() error {
			systemInfo, err := host.GetSystemInfo(ctx)
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

	if err := common.DeepCopy(results, &hl.Hosts); err != nil {
		return err
	}

	dbHostList := db.HostList{}

	if err := common.DeepCopy(hl.Hosts, &dbHostList.Hosts); err != nil {
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
				error := common.ErrHostAlreadyRegistered
				ipList := hl.getIPList()
				error.Params = []string{strings.Join(ipList, ",")}
				return error
			}
		}

		return err
	}
}

func (hl *HostList) Unregister(ctx context.Context) error {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	// TODO: Return error if there is any related directories.

	hostList := db.HostList{}

	common.DeepCopy(hl.Hosts, &hostList.Hosts)

	return hostList.Delete(engine, nil)
}

func (hl *HostList) Update(ctx context.Context) error {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	hostList := db.HostList{}
	if err := hostList.Get(engine, &common.QueryFilter{}); err != nil {
		return err
	}

	g, _ := errgroup.WithContext(context.Background())

	for _, h := range hostList.Hosts {
		dbHost := h // 避免闭包问题
		g.Go(func() error {
			var host Host
			common.DeepCopy(dbHost, &host)

			if _, err := host.GetSystemInfo(ctx); err != nil {
				dbHost.Connected = false
				dbHost.Save(engine)
				return err
			}
			dbHost.Connected = true
			dbHost.Save(engine)

			return nil
		})
	}

	return g.Wait()
}

func (hl *HostList) Get(ctx context.Context, filter *common.QueryFilter) ([]Host, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	hostList := db.HostList{}
	if err := hostList.Get(engine, filter); err != nil {
		return nil, err
	}

	common.DeepCopy(hostList.Hosts, &hl.Hosts)

	return hl.Hosts, nil
}

type PaginationHost struct {
	Hosts      []Host
	Page       int
	Limit      int
	TotalCount int64
}

func (hl *HostList) Pagination(ctx context.Context, filter *common.QueryFilter) (*PaginationHost, error) {
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

	common.DeepCopy(paginationHosts.Hosts, &paginationHostList.Hosts)

	return &paginationHostList, nil
}

func (h *Host) GetSystemInfo(ctx context.Context) (systemInfo common.SystemInfo, err error) {
	hostContext := common.HostContext{
		IP:       h.IP,
		Username: h.Username,
		Password: h.Password,
	}

	ctx = context.WithValue(ctx, common.HostContextkey("hostContext"), hostContext)

	driver := driver.GetDriver(h.StorageType)

	return driver.GetSystemInfo(ctx)
}
