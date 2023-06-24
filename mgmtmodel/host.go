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

	dl := db.DirectoryList{}
	dl.Delete(engine, nil, h.Ip)

	host := db.Host{
		Ip:   h.Ip,
		Name: h.Name,
	}

	return host.Delete(engine)
}

func (h *Host) Get() (*Host, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	dbHost := db.Host{
		Name: h.Name,
		Ip:   h.Ip,
	}
	if err = dbHost.Get(engine); err != nil {
		return nil, err
	}

	h.Ip = dbHost.Ip
	h.Name = dbHost.Name
	h.Username = dbHost.Username
	h.Password = dbHost.Password
	h.StorageType = dbHost.StorageType

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

	dbHostList := db.HostList{}
	if err := dbHostList.Get(engine, storageType); err != nil {
		return nil, err
	}

	for _, _host := range dbHostList.Hosts {
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
