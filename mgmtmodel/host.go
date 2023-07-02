package mgmtmodel

import (
	"github.com/cryingmouse/data_management_engine/context"
	"github.com/cryingmouse/data_management_engine/db"
)

type Host struct {
	Name        string
	IP          string
	Username    string
	Password    string
	StorageType string
}

func (h *Host) Register() error {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	host := db.Host{
		Name:        h.Name,
		IP:          h.IP,
		Username:    h.Username,
		Password:    h.Password,
		StorageType: h.StorageType,
	}

	return host.Save(engine)
}

func (h *Host) Unregister() error {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	directoryList := db.DirectoryList{}
	directoryList.Delete(engine, nil)

	host := db.Host{
		IP:   h.IP,
		Name: h.Name,
	}

	return host.Delete(engine)
}

func (h *Host) Get() (*Host, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	host := db.Host{
		Name: h.Name,
		IP:   h.IP,
	}
	if err = host.Get(engine); err != nil {
		return nil, err
	}

	h.IP = host.IP
	h.Name = host.Name
	h.Username = host.Username
	h.Password = host.Password
	h.StorageType = host.StorageType

	return h, nil
}

type HostList struct {
	Hosts []Host
}

func (hl *HostList) Get(filter *context.QueryFilter) ([]Host, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	hostList := db.HostList{}
	if err := hostList.Get(engine, filter); err != nil {
		return nil, err
	}

	for _, _host := range hostList.Hosts {
		host := Host{
			IP:          _host.IP,
			Name:        _host.Name,
			Username:    _host.Username,
			Password:    _host.Password,
			StorageType: _host.StorageType,
		}

		hl.Hosts = append(hl.Hosts, host)
	}

	return hl.Hosts, nil
}

type PaginationHost struct {
	Hosts      []Host
	Page       int
	Limit      int
	TotalCount int64
}

func (hl *HostList) Pagination(filter *context.QueryFilter) (*PaginationHost, error) {
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

	for _, _host := range paginationHosts.Hosts {
		host := Host{
			Name:        _host.Name,
			IP:          _host.IP,
			Username:    _host.Username,
			Password:    _host.Password,
			StorageType: _host.StorageType,
		}

		paginationHostList.Hosts = append(paginationHostList.Hosts, host)
	}

	return &paginationHostList, nil
}

func (hl *HostList) Register() error {
	return nil
}

func (hl *HostList) Unregister() error {
	return nil
}
