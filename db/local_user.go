package db

import (
	"errors"
	"fmt"

	"github.com/cryingmouse/data_management_engine/common"

	"gorm.io/gorm"
)

type LocalUser struct {
	gorm.Model
	HostIP               string `gorm:"uniqueIndex:idx_name_host_ip;column:host_ip"`
	Name                 string `gorm:"uniqueIndex:idx_name_host_ip;column:name"`
	UID                  string `gorm:"column:user_id"`
	Password             string `gorm:"type:password;column:password"`
	Fullname             string `gorm:"column:full_name"`
	Status               string `gorm:"column:status"`
	Description          string `gorm:"column:description"`
	IsDisabled           bool   `gorm:"column:is_disabled"`
	IsPasswordExpired    bool   `gorm:"column:is_password_expired"`
	IsPasswordChangeable bool   `gorm:"column:is_password_changeable"`
	IsPasswordRequired   bool   `gorm:"column:is_password_required"`
	IsLockout            bool   `gorm:"column:is_lockout"`
}

func (u *LocalUser) Get(engine *DatabaseEngine) error {
	err := engine.DB.Where(u).First(u).Error
	if err != nil {
		return err
	}

	u.Password, err = common.Decrypt(u.Password, common.SecurityKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt password: %w", err)
	}

	return nil
}

func (u *LocalUser) Save(engine *DatabaseEngine) error {
	encryptedPassword, err := common.Encrypt(u.Password, common.SecurityKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	localUser := &LocalUser{
		HostIP:               u.HostIP,
		UID:                  u.UID,
		Name:                 u.Name,
		Fullname:             u.Fullname,
		Status:               u.Status,
		Description:          u.Description,
		IsPasswordExpired:    u.IsPasswordExpired,
		IsPasswordChangeable: u.IsPasswordChangeable,
		IsPasswordRequired:   u.IsPasswordRequired,
		IsLockout:            u.IsLockout,
		Password:             string(encryptedPassword),
	}

	return engine.DB.Save(localUser).Error
}

func (u *LocalUser) Delete(engine *DatabaseEngine) error {
	return engine.DB.Unscoped().Where(u).Delete(u).Error
}

type LocalUserList struct {
	LocalUsers []LocalUser
}

func (ul *LocalUserList) Get(engine *DatabaseEngine, filter *common.QueryFilter) (err error) {
	model := LocalUser{}

	if filter.Pagination != nil {
		return errors.New("invalid filter: pagination is not supported")
	}

	if _, err := Query(engine, model, filter, &ul.LocalUsers); err != nil {
		return fmt.Errorf("failed to query hosts from the database: %w", err)
	}

	for i := range ul.LocalUsers {
		if ul.LocalUsers[i].Password != "" {
			var err error
			ul.LocalUsers[i].Password, err = common.Decrypt(ul.LocalUsers[i].Password, common.SecurityKey)
			if err != nil {
				return fmt.Errorf("failed to decrypt password: %w", err)
			}
		}
	}

	return
}

type PaginationLocalUser struct {
	LocalUsers []LocalUser
	TotalCount int64
}

func (ul *LocalUserList) Pagination(engine *DatabaseEngine, filter *common.QueryFilter) (*PaginationLocalUser, error) {
	var totalCount int64
	model := LocalUser{}

	if filter.Pagination == nil {
		return nil, fmt.Errorf("invalid filter: missing pagination")
	}

	totalCount, err := Query(engine, model, filter, &ul.LocalUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to query local users by the filter %v in the database: %w", filter, err)
	}

	response := &PaginationLocalUser{
		LocalUsers: ul.LocalUsers,
		TotalCount: totalCount,
	}

	return response, nil
}

func (ul *LocalUserList) Save(engine *DatabaseEngine) (err error) {
	if len(ul.LocalUsers) == 0 {
		return errors.New("UserList is empty")
	}

	for i, user := range ul.LocalUsers {
		// Encrypt the password
		ul.LocalUsers[i].Password, err = common.Encrypt(user.Password, common.SecurityKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt password for user %v: %w", user.Name, err)
		}
	}

	err = engine.DB.CreateInBatches(ul.LocalUsers, len(ul.LocalUsers)).Error
	if err != nil {
		return fmt.Errorf("failed to create users in database: %w", err)
	}
	return nil
}

func (ul *LocalUserList) Delete(engine *DatabaseEngine, filter *common.QueryFilter) (err error) {
	var users []LocalUser
	if filter != nil {
		return Delete(engine, filter, &ul.LocalUsers)
	}

	query := engine.DB.Unscoped()

	conditions := common.StructListToMapList(ul.LocalUsers)
	for _, condition := range conditions {
		query = query.Or(condition)
	}

	return query.Find(&users).Delete(&users).Error
}
