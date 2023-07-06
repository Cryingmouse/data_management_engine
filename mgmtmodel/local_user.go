package mgmtmodel

import (
	"github.com/cryingmouse/data_management_engine/common"
	"github.com/cryingmouse/data_management_engine/db"
	"github.com/cryingmouse/data_management_engine/driver"
)

type LocalUser struct {
	Name     string
	Password string
	HostName string
}

func (u *LocalUser) Create() (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	// Get host IP, usrename and password by host name.
	host := db.Host{ComputerName: u.HostName}
	if err = host.Get(engine); err != nil {
		return err
	}

	hostContext := common.HostContext{
		IP:       host.IP,
		Username: host.Username,
		Password: host.Password,
	}

	// Create local user on agent host.
	driver := driver.GetDriver(host.StorageType)
	driver.CreateLocalUser(hostContext, u.Name, u.Password)

	// Save to database.
	var user db.LocalUser
	common.CopyStructList(u, &user)
	return user.Save(engine)
}

func (u *LocalUser) Delete() (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	// Get host IP, usrename and password by host name.
	host := db.Host{ComputerName: u.HostName}
	if err = host.Get(engine); err != nil {
		return err
	}

	hostContext := common.HostContext{
		IP:       host.IP,
		Username: host.Username,
		Password: host.Password,
	}

	// Create local user on agent host.
	driver := driver.GetDriver(host.StorageType)
	driver.DeleteUser(hostContext, u.Name)

	// Delete from database.
	var user db.LocalUser
	common.CopyStructList(u, &user)
	return user.Delete(engine)
}

func (u *LocalUser) Get() (*LocalUser, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	var user db.LocalUser
	common.CopyStructList(u, &user)

	if err = user.Get(engine); err != nil {
		return nil, err
	}

	common.CopyStructList(user, &u)

	return u, nil
}

type LocalUserList struct {
	Users []LocalUser
}

func (ul *LocalUserList) Create() error {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	// TODO: Need to create the users on the host by cogouine
	// TODO: Save to database

	userList := db.LocalUserList{}

	common.CopyStructList(ul.Users, &userList.Users)

	return userList.Save(engine)
}

func (dl *LocalUserList) Delete(filter *common.QueryFilter) (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	// TODO: Need to delete the users on the host by cogouine, using filter
	// TODO: Delete from database

	userList := db.LocalUserList{}

	return userList.Delete(engine, filter)
}

func (ul *LocalUserList) Get(filter *common.QueryFilter) ([]LocalUser, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	userList := db.LocalUserList{}

	if err = userList.Get(engine, filter); err != nil {
		return nil, err
	}

	common.CopyStructList(userList.Users, &ul.Users)

	return ul.Users, nil
}

type PaginationUser struct {
	Users      []LocalUser
	Page       int
	Limit      int
	TotalCount int64
}

func (dl *LocalUserList) Pagination(filter *common.QueryFilter) (*PaginationUser, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	userList := db.LocalUserList{}
	paginationUsers, err := userList.Pagination(engine, filter)
	if err != nil {
		return nil, err
	}

	paginationUserList := PaginationUser{
		Page:       filter.Pagination.Page,
		Limit:      filter.Pagination.PageSize,
		TotalCount: paginationUsers.TotalCount,
	}

	common.CopyStructList(paginationUsers.Users, &paginationUserList.Users)

	return &paginationUserList, nil
}
