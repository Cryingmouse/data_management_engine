package webservice

import (
	"regexp"

	"github.com/cryingmouse/data_management_engine/common"
	"github.com/go-playground/validator/v10"
)

func PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) >= 8 && regexp.MustCompile(`[A-Z]+`).MatchString(password) && regexp.MustCompile(`[a-z]+`).MatchString(password) && regexp.MustCompile(`[0-9]+`).MatchString(password) {
		return true
	}
	return false
}

func StorageTypeValidator(fl validator.FieldLevel) bool {
	storageType := fl.Field().String()

	storageTypeList := []string{StorageTypeWorkstation, StorageTypeOntap, StorageTypeMagnascale}

	return common.In(storageType, storageTypeList)
}

func validateIPAddress(ip string) (err error) {
	type IPAddress struct {
		IP string `validate:"required,ip"`
	}

	validate := validator.New()

	// 验证 IP 地址
	ipAddress := IPAddress{IP: ip}

	return validate.Struct(ipAddress)
}
