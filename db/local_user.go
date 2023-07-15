package db

import (
	"errors"
	"fmt"

	"github.com/cryingmouse/data_management_engine/common"

	"gorm.io/gorm"
)

type LocalUser struct {
	gorm.Model
	Name     string `gorm:"unique;column:name"`
	Password string `gorm:"type:password;column:password"`
	HostName string `gorm:"column:host_name"`
}

func (u *LocalUser) Get(engine *DatabaseEngine) (err error) {
	if err = engine.DB.Where(u).First(u).Error; err != nil {
		return err
	}

	u.Password, err = common.Decrypt(u.Password, common.SecurityKey)

	return err
}

func (u *LocalUser) Save(engine *DatabaseEngine) error {
	encrypted_password, err := common.Encrypt(u.Password, common.SecurityKey)
	if err != nil {
		return err
	}

	u.Password = string(encrypted_password)

	return engine.DB.Save(u).Error
}

func (u *LocalUser) Delete(engine *DatabaseEngine) error {
	err := engine.DB.Unscoped().Where(u).Delete(u).Error
	return err
}

type LocalUserList struct {
	Users []LocalUser
}

func (ul *LocalUserList) Get(engine *DatabaseEngine, filter *common.QueryFilter) (err error) {
	model := LocalUser{}

	if filter.Pagination != nil {
		return fmt.Errorf("invalid filter: there is pagination in the filter")
	}

	if _, err := Query(engine, model, filter, &ul.Users); err != nil {
		return fmt.Errorf("failed to query the users by the filter %v in database: %w", filter, err)
	}

	for _, user := range ul.Users {
		user.Password, err = common.Decrypt(user.Password, common.SecurityKey)
	}

	return
}

type PaginationUser struct {
	Users      []LocalUser
	TotalCount int64
}

func (ul *LocalUserList) Pagination(engine *DatabaseEngine, filter *common.QueryFilter) (response *PaginationUser, err error) {
	var totalCount int64
	model := LocalUser{}

	if filter.Pagination == nil {
		return response, fmt.Errorf("invalid filter: missing pagination in the filter")
	}

	totalCount, err = Query(engine, model, filter, &ul.Users)
	if err != nil {
		return response, fmt.Errorf("failed to query the users by the filter %v in database: %w", filter, err)
	}

	for _, user := range ul.Users {
		user.Password, err = common.Decrypt(user.Password, common.SecurityKey)
	}

	response = &PaginationUser{
		Users:      ul.Users,
		TotalCount: totalCount,
	}

	return response, err
}

func (ul *LocalUserList) Save(engine *DatabaseEngine) (err error) {
	if len(ul.Users) == 0 {
		return errors.New("UserList is empty")
	}

	for i, user := range ul.Users {
		// Encrypt the password
		ul.Users[i].Password, err = common.Encrypt(user.Password, common.SecurityKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt password for user %v: %w", user.Name, err)
		}
	}

	err = engine.DB.CreateInBatches(ul.Users, len(ul.Users)).Error
	if err != nil {
		return fmt.Errorf("failed to create users in database: %w", err)
	}
	return nil
}

func (ul *LocalUserList) Delete(engine *DatabaseEngine, filter *common.QueryFilter) (err error) {
	var users []LocalUser

	if filter != nil {
		err = Delete(engine, filter, users)
		if err != nil {
			return fmt.Errorf("failed to delete local users by the filter %v in database: %w", filter, err)
		}
	} else {
		if err := engine.DB.Delete(ul.Users).Error; err != nil {
			return err
		}
	}

	return
}
