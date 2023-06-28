package mgmtmodel

import (
	"github.com/cryingmouse/data_management_engine/db"
)

type Host struct {
	Name        string
	Ip          string
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
		Ip:          h.Ip,
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
	directoryList.Delete(engine, nil, h.Ip)

	host := db.Host{
		Ip:   h.Ip,
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
		Ip:   h.Ip,
	}
	if err = host.Get(engine); err != nil {
		return nil, err
	}

	h.Ip = host.Ip
	h.Name = host.Name
	h.Username = host.Username
	h.Password = host.Password
	h.StorageType = host.StorageType

	return h, nil
}

type HostList struct {
	Hosts []Host
}

func (hl *HostList) Get(storageType string) ([]Host, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	hostList := db.HostList{}
	if err := hostList.Get(engine, storageType); err != nil {
		return nil, err
	}

	for _, _host := range hostList.Hosts {
		host := Host{
			Ip:          _host.Ip,
			Name:        _host.Name,
			Username:    _host.Username,
			Password:    _host.Password,
			StorageType: _host.StorageType,
		}

		hl.Hosts = append(hl.Hosts, host)
	}

	return hl.Hosts, nil
}

func (hl *HostList) Register() error {
	return nil
}

func (hl *HostList) Unregister() error {
	return nil
}
