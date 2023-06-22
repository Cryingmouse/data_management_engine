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

type HostList struct {
	Hosts []Host
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

	hostModel := db.Host{}
	host, err := hostModel.Get(engine, h.Name, h.Ip)
	if err != nil {
		return err
	}

	return host.Delete(engine)
}

func (h *Host) Get() (*Host, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	hostModel := db.Host{}
	host, err := hostModel.Get(engine, h.Name, h.Ip)
	if err != nil {
		return nil, err
	}

	h.Name = host.Name
	h.Username = host.Username
	h.Password = host.Password
	h.StorageType = host.StorageType

	return h, nil
}

func (hl *HostList) Get() ([]Host, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	hostListModel := db.HostList{}
	hosts, err := hostListModel.Get(engine)
	if err != nil {
		return nil, err
	}

	for _, host := range hosts {
		host := Host{
			Ip:          host.Ip,
			Name:        host.Name,
			Username:    host.Username,
			Password:    host.Password,
			StorageType: host.StorageType,
		}

		hl.Hosts = append(hl.Hosts, host)
	}

	return hl.Hosts, nil
}
