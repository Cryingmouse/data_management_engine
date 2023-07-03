package db

import (
	"errors"
	"fmt"

	"github.com/cryingmouse/data_management_engine/context"
	"github.com/cryingmouse/data_management_engine/utils"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"unique;column:name"`
	Password string `gorm:"type:password;column:password"`
}

func (u *User) Get(engine *DatabaseEngine) (err error) {
	if err = engine.DB.Where(u).First(u).Error; err != nil {
		return err
	}

	u.Password, err = utils.Decrypt(u.Password, context.SecurityKey)

	return err
}

func (u *User) Save(engine *DatabaseEngine) error {
	user := User{
		Name: u.Name,
	}

	encrypted_password, err := utils.Encrypt(u.Password, context.SecurityKey)
	if err != nil {
		return err
	}
	user.Password = string(encrypted_password)

	return engine.DB.Save(&user).Error
}

func (u *User) Delete(engine *DatabaseEngine) error {
	return engine.DB.Unscoped().Delete(&u, u).Error
}

type UserList struct {
	Users []User
}

func (ul *UserList) Get(engine *DatabaseEngine, filter *context.QueryFilter) (err error) {
	model := User{}

	if filter.Pagination != nil {
		return fmt.Errorf("invalid filter: there is pagination in the filter")
	}

	if _, err := Query(engine, model, filter, &ul.Users); err != nil {
		return fmt.Errorf("failed to query the users by the filter %v in database: %w", filter, err)
	}

	for _, user := range ul.Users {
		user.Password, err = utils.Decrypt(user.Password, context.SecurityKey)
	}

	return
}

type PaginationUser struct {
	Users      []User
	TotalCount int64
}

func (ul *UserList) Pagination(engine *DatabaseEngine, filter *context.QueryFilter) (response *PaginationUser, err error) {
	var totalCount int64
	model := User{}

	if filter.Pagination == nil {
		return response, fmt.Errorf("invalid filter: missing pagination in the filter")
	}

	totalCount, err = Query(engine, model, filter, &ul.Users)
	if err != nil {
		return response, fmt.Errorf("failed to query the users by the filter %v in database: %w", filter, err)
	}

	for _, user := range ul.Users {
		user.Password, err = utils.Decrypt(user.Password, context.SecurityKey)
	}

	response = &PaginationUser{
		Users:      ul.Users,
		TotalCount: totalCount,
	}

	return response, err
}

func (ul *UserList) Save(engine *DatabaseEngine) (err error) {
	if len(ul.Users) == 0 {
		return errors.New("UserList is empty")
	}

	for i, user := range ul.Users {
		// Encrypt the password
		ul.Users[i].Password, err = utils.Encrypt(user.Password, context.SecurityKey)
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

func (ul *UserList) Delete(engine *DatabaseEngine, filter *context.QueryFilter) (err error) {
	var users []User

	err = Delete(engine, filter, users)
	if err != nil {
		return fmt.Errorf("failed to delete users by the filter %v in database: %w", filter, err)
	}

	return
}
