package db

import (
	"errors"
	"fmt"

	"github.com/cryingmouse/data_management_engine/common"

	"gorm.io/gorm"
)

type LocalUser struct {
	gorm.Model
	HostIP               string `gorm:"uniqueIndex:idx_local_user_unique;column:host_ip"`
	Name                 string `gorm:"uniqueIndex:idx_local_user_unique;column:name"`
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

func (u *LocalUser) Get(engine *DatabaseEngine) (err error) {
	err = engine.DB.Where(u).First(u).Error
	if err != nil {
		return err
	}

	u.Password, err = common.Decrypt(u.Password, common.SecurityKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt password: %w", err)
	}

	return nil
}

func (u *LocalUser) Save(engine *DatabaseEngine) (err error) {
	u.Password, err = common.Encrypt(u.Password, common.SecurityKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	return engine.DB.Save(u).Error
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

	return nil
}

type PaginationLocalUser struct {
	LocalUsers []LocalUser
	TotalCount int64
}

func (ul *LocalUserList) Pagination(engine *DatabaseEngine, filter *common.QueryFilter) (paginationLocalUser PaginationLocalUser, err error) {
	if filter.Pagination == nil {
		return paginationLocalUser, fmt.Errorf("invalid filter: missing pagination")
	}

	model := LocalUser{}
	var totalCount int64
	totalCount, err = Query(engine, model, filter, &ul.LocalUsers)
	if err != nil {
		return paginationLocalUser, fmt.Errorf("failed to query local users by the filter %v in the database: %w", filter, err)
	}

	for i := range ul.LocalUsers {
		if ul.LocalUsers[i].Password != "" {
			ul.LocalUsers[i].Password, err = common.Decrypt(ul.LocalUsers[i].Password, common.SecurityKey)
			if err != nil {
				return paginationLocalUser, fmt.Errorf("failed to decrypt password: %w", err)
			}
		}
	}

	paginationLocalUser.LocalUsers = ul.LocalUsers
	paginationLocalUser.TotalCount = totalCount

	return paginationLocalUser, nil
}

func (ul *LocalUserList) Save(engine *DatabaseEngine) (err error) {
	if len(ul.LocalUsers) == 0 {
		return errors.New("UserList is empty")
	}

	for i := range ul.LocalUsers {
		// Encrypt the password
		ul.LocalUsers[i].Password, err = common.Encrypt(ul.LocalUsers[i].Password, common.SecurityKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt password for user %v: %w", ul.LocalUsers[i].Name, err)
		}
	}

	err = engine.DB.CreateInBatches(ul.LocalUsers, len(ul.LocalUsers)).Error
	if err != nil {
		return fmt.Errorf("failed to save the users in database: %w", err)
	}
	return nil
}

func (ul *LocalUserList) Delete(engine *DatabaseEngine, filter *common.QueryFilter) error {
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
