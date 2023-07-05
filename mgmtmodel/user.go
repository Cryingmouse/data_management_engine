package mgmtmodel

import (
	"github.com/cryingmouse/data_management_engine/common"
	"github.com/cryingmouse/data_management_engine/db"
	"github.com/cryingmouse/data_management_engine/driver"
)

type User struct {
	Name     string
	Password string
}

func (u *User) Create() (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	host := db.Host{IP: u.Password}
	if err = host.Get(engine); err != nil {
		return err
	}

	hostContext := common.HostContext{
		IP:       host.IP,
		Username: host.Username,
		Password: host.Password,
	}

	driver := driver.GetDriver(host.StorageType)
	driver.CreateUser(hostContext, u.Name, u.Password)

	user := db.User{
		Name:     u.Name,
		Password: host.IP,
	}

	if err = user.Save(engine); err != nil {
		return err
	}

	return nil
}

func (u *User) Delete() (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	host := db.Host{IP: u.Password}
	if err = host.Get(engine); err != nil {
		return err
	}

	hostContext := common.HostContext{
		IP:       host.IP,
		Username: host.Username,
		Password: host.Password,
	}

	driver := driver.GetDriver(host.StorageType)
	driver.DeleteUser(hostContext, u.Name)

	user := db.User{
		Name:     u.Name,
		Password: host.IP,
	}

	if err = user.Delete(engine); err != nil {
		return err
	}

	return nil
}

func (u *User) Get() (*User, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	user := db.User{
		Name:     u.Name,
		Password: u.Password,
	}
	if err = user.Get(engine); err != nil {
		return nil, err
	}

	u.Name = user.Name
	u.Password = user.Password

	return u, nil
}

type UserList struct {
	Users []User
}

func (dl *UserList) Get(filter *common.QueryFilter) ([]User, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	userList := db.UserList{}

	if err = userList.Get(engine, filter); err != nil {
		return nil, err
	}

	for _, _user := range userList.Users {
		user := User{
			Name:     _user.Name,
			Password: _user.Password,
		}

		dl.Users = append(dl.Users, user)
	}

	return dl.Users, nil
}

type PaginationUser struct {
	Users      []User
	Page       int
	Limit      int
	TotalCount int64
}

func (dl *UserList) Pagination(filter *common.QueryFilter) (*PaginationUser, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	userList := db.UserList{}
	paginationDirs, err := userList.Pagination(engine, filter)
	if err != nil {
		return nil, err
	}

	paginationDirList := PaginationUser{
		Page:       filter.Pagination.Page,
		Limit:      filter.Pagination.PageSize,
		TotalCount: paginationDirs.TotalCount,
	}

	for _, _user := range paginationDirs.Users {
		user := User{
			Name:     _user.Name,
			Password: _user.Password,
		}

		paginationDirList.Users = append(paginationDirList.Users, user)
	}

	return &paginationDirList, nil
}

func (dl *UserList) Save() (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	userList := db.UserList{}

	for _, _user := range dl.Users {
		user := db.User{
			Name:     _user.Name,
			Password: _user.Password,
		}

		userList.Users = append(userList.Users, user)
	}

	return userList.Save(engine)
}

func (dl *UserList) Delete(filter *common.QueryFilter) (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	userList := db.UserList{}

	return userList.Delete(engine, filter)
}
