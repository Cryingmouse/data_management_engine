package mgmtmodel

import (
	"context"
	"errors"

	"github.com/cryingmouse/data_management_engine/common"
	"github.com/cryingmouse/data_management_engine/db"
	"github.com/cryingmouse/data_management_engine/driver"
	"golang.org/x/sync/errgroup"
)

type LocalUser struct {
	HostIP               string
	UID                  string
	Name                 string
	Password             string
	Fullname             string
	Description          string
	Status               string
	IsDisabled           bool
	IsPasswordExpired    bool
	IsPasswordChangeable bool
	IsPasswordRequired   bool
	IsLockout            bool
}

func (u *LocalUser) Create(ctx context.Context) (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	// Get the right driver and call driver to create directory.
	host := db.Host{IP: u.HostIP}
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
	localUserDetail, err := driver.CreateLocalUser(ctx, u.Name, u.Password)
	if err != nil {
		return err
	}

	// Save the details of the directory into database.
	LocalUser := db.LocalUser{
		HostIP:               host.IP,
		Name:                 u.Name,
		Password:             u.Password,
		UID:                  localUserDetail.UID,
		Fullname:             localUserDetail.FullName,
		Description:          localUserDetail.Description,
		IsPasswordExpired:    localUserDetail.IsPasswordExpired,
		IsPasswordChangeable: localUserDetail.IsPasswordChangeable,
		IsPasswordRequired:   localUserDetail.IsPasswordRequired,
		IsLockout:            localUserDetail.IsLockout,
	}

	common.CopyStructList(LocalUser, u)

	return LocalUser.Save(engine)
}

func (u *LocalUser) Delete(ctx context.Context) (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	host := db.Host{IP: u.HostIP}
	if err = host.Get(engine); err != nil {
		return err
	}

	hostContext := common.HostContext{
		IP:       host.IP,
		Username: host.Username,
		Password: host.Password,
	}
	ctx = context.WithValue(ctx, common.HostContextkey("hostContext"), hostContext)

	// Create local user on agent host.
	driver := driver.GetDriver(host.StorageType)
	if err := driver.DeleteLocalUser(ctx, u.Name); err != nil {
		return err
	}

	localUser := db.LocalUser{
		Name:   u.Name,
		HostIP: host.IP,
	}
	return localUser.Delete(engine)
}

func (u *LocalUser) Get(ctx context.Context) (*LocalUser, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	localUser := db.LocalUser{
		Name:   u.Name,
		HostIP: u.HostIP,
	}
	if err = localUser.Get(engine); err != nil {
		return nil, err
	}

	common.CopyStructList(localUser, u)

	return u, nil
}

func (u *LocalUser) Manage(ctx context.Context) (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	// Get the right driver and call driver to create directory.
	host := db.Host{IP: u.HostIP}
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
	localUserDetail, err := driver.GetLocalUserDetail(ctx, u.Name)
	if err != nil {
		return err
	}

	// Save the details of the directory into database.
	LocalUser := db.LocalUser{
		HostIP:               host.IP,
		Name:                 u.Name,
		Password:             u.Password,
		UID:                  localUserDetail.UID,
		Fullname:             localUserDetail.FullName,
		Description:          localUserDetail.Description,
		Status:               localUserDetail.Status,
		IsPasswordExpired:    localUserDetail.IsPasswordExpired,
		IsPasswordChangeable: localUserDetail.IsPasswordChangeable,
		IsPasswordRequired:   localUserDetail.IsPasswordRequired,
		IsLockout:            localUserDetail.IsLockout,
	}

	common.CopyStructList(LocalUser, u)

	return LocalUser.Save(engine)
}

func (u *LocalUser) Unmanage(ctx context.Context) (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	localUser := db.LocalUser{
		HostIP: u.HostIP,
		Name:   u.Name,
	}
	return localUser.Delete(engine)
}

type LocalUserList struct {
	LocalUsers []LocalUser
}

func (ul *LocalUserList) Create(ctx context.Context) error {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	input := make([]interface{}, len(ul.LocalUsers))
	for index, localUser := range ul.LocalUsers {
		input[index] = localUser
	}

	g, _ := errgroup.WithContext(context.Background())

	results := make([]common.LocalUserDetail, len(ul.LocalUsers))
	var resultErr error

	for i, u := range ul.LocalUsers {
		index := i
		localUser := u // 避免闭包问题
		g.Go(func() error {
			host := db.Host{IP: localUser.HostIP}
			if err = host.Get(engine); err != nil {
				resultErr = errors.Join(resultErr, err)
				return err
			}

			hostContext := common.HostContext{
				IP:       host.IP,
				Username: host.Username,
				Password: host.Password,
			}
			ctx = context.WithValue(ctx, common.HostContextkey("hostContext"), hostContext)

			driver := driver.GetDriver(host.StorageType)
			localUserDetail, err := driver.GetLocalUserDetail(ctx, localUser.Name)
			if err != nil {
				resultErr = errors.Join(resultErr, err)
				return err
			}

			results[index] = localUserDetail // 保存协程的返回值

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		if resultErr != nil {
			return resultErr
		}
		return err
	}

	if err := common.CopyStructList(results, &ul.LocalUsers); err != nil {
		return err
	}

	// Save to database.
	dbLocalUserList := db.LocalUserList{}

	if err := common.CopyStructList(ul.LocalUsers, &dbLocalUserList.LocalUsers); err != nil {
		return err
	}

	return dbLocalUserList.Save(engine)
}

func (ul *LocalUserList) Delete(ctx context.Context, filter *common.QueryFilter) (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	g, _ := errgroup.WithContext(context.Background())

	results := make([]common.LocalUserDetail, len(ul.LocalUsers))
	var resultErr error

	for i, u := range ul.LocalUsers {
		index := i
		localUser := u // 避免闭包问题
		g.Go(func() error {
			host := db.Host{IP: localUser.HostIP}
			if err = host.Get(engine); err != nil {
				resultErr = errors.Join(resultErr, err)
				return err
			}

			hostContext := common.HostContext{
				IP:       host.IP,
				Username: host.Username,
				Password: host.Password,
			}
			ctx = context.WithValue(ctx, common.HostContextkey("hostContext"), hostContext)

			driver := driver.GetDriver(host.StorageType)
			localUserDetail, err := driver.GetLocalUserDetail(ctx, localUser.Name)
			if err != nil {
				resultErr = errors.Join(resultErr, err)
				return err
			}

			results[index] = localUserDetail // 保存协程的返回值

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		if resultErr != nil {
			return resultErr
		}
		return err
	}

	if err := common.CopyStructList(results, &ul.LocalUsers); err != nil {
		return err
	}

	localUserList := db.LocalUserList{}
	if err := common.CopyStructList(ul.LocalUsers, &localUserList.LocalUsers); err != nil {
		return err
	}

	return localUserList.Delete(engine, filter)
}

func (ul *LocalUserList) Get(ctx context.Context, filter *common.QueryFilter) ([]LocalUser, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	localUserList := db.LocalUserList{}

	if err = localUserList.Get(engine, filter); err != nil {
		return nil, err
	}

	common.CopyStructList(localUserList.LocalUsers, &ul.LocalUsers)

	return ul.LocalUsers, nil
}

func (ul *LocalUserList) Manage(ctx context.Context) error {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	input := make([]interface{}, len(ul.LocalUsers))
	for index, localUser := range ul.LocalUsers {
		input[index] = localUser
	}

	g, _ := errgroup.WithContext(context.Background())

	results := make([]common.LocalUserDetail, len(ul.LocalUsers))
	var resultErr error

	for i, u := range ul.LocalUsers {
		index := i
		localUser := u // 避免闭包问题
		g.Go(func() error {
			host := db.Host{IP: localUser.HostIP}
			if err = host.Get(engine); err != nil {
				resultErr = errors.Join(resultErr, err)
				return err
			}

			hostContext := common.HostContext{
				IP:       host.IP,
				Username: host.Username,
				Password: host.Password,
			}
			ctx = context.WithValue(ctx, common.HostContextkey("hostContext"), hostContext)

			driver := driver.GetDriver(host.StorageType)
			localUserDetail, err := driver.GetLocalUserDetail(ctx, localUser.Name)
			if err != nil {
				resultErr = errors.Join(resultErr, err)
				return err
			}

			results[index] = localUserDetail // 保存协程的返回值

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		if resultErr != nil {
			return resultErr
		}
		return err
	}

	if err := common.CopyStructList(results, &ul.LocalUsers); err != nil {
		return err
	}

	// Save to database.
	dbLocalUserList := db.LocalUserList{}

	if err := common.CopyStructList(ul.LocalUsers, &dbLocalUserList.LocalUsers); err != nil {
		return err
	}

	return dbLocalUserList.Save(engine)
}

func (ul *LocalUserList) Unmanage(ctx context.Context, filter *common.QueryFilter) (err error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	g, _ := errgroup.WithContext(context.Background())

	results := make([]common.LocalUserDetail, len(ul.LocalUsers))
	var resultErr error

	for i, u := range ul.LocalUsers {
		index := i
		localUser := u // 避免闭包问题
		g.Go(func() error {
			host := db.Host{IP: localUser.HostIP}
			if err = host.Get(engine); err != nil {
				resultErr = errors.Join(resultErr, err)
				return err
			}

			hostContext := common.HostContext{
				IP:       host.IP,
				Username: host.Username,
				Password: host.Password,
			}
			ctx = context.WithValue(ctx, common.HostContextkey("hostContext"), hostContext)

			driver := driver.GetDriver(host.StorageType)
			localUserDetail, err := driver.GetLocalUserDetail(ctx, localUser.Name)
			if err != nil {
				resultErr = errors.Join(resultErr, err)
				return err
			}

			results[index] = localUserDetail // 保存协程的返回值

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		if resultErr != nil {
			return resultErr
		}
		return err
	}

	if err := common.CopyStructList(results, &ul.LocalUsers); err != nil {
		return err
	}

	localUserList := db.LocalUserList{}
	if err := common.CopyStructList(ul.LocalUsers, &localUserList.LocalUsers); err != nil {
		return err
	}

	return localUserList.Delete(engine, filter)
}

type PaginationLocalUser struct {
	LocalUsers []LocalUser
	Page       int
	Limit      int
	TotalCount int64
}

func (dl *LocalUserList) Pagination(ctx context.Context, filter *common.QueryFilter) (*PaginationLocalUser, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	localUserList := db.LocalUserList{}
	paginationLocalUsers, err := localUserList.Pagination(engine, filter)
	if err != nil {
		return nil, err
	}

	paginationLocalUserList := PaginationLocalUser{
		Page:       filter.Pagination.Page,
		Limit:      filter.Pagination.PageSize,
		TotalCount: paginationLocalUsers.TotalCount,
	}

	for _, _localUser := range paginationLocalUsers.LocalUsers {
		localUser := LocalUser{
			Name:   _localUser.Name,
			HostIP: _localUser.HostIP,
		}

		paginationLocalUserList.LocalUsers = append(paginationLocalUserList.LocalUsers, localUser)
	}

	return &paginationLocalUserList, nil
}
