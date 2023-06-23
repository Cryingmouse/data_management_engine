package webservice

import (
	"net/http"

	"github.com/cryingmouse/data_management_engine/mgmtmodel"
	"github.com/cryingmouse/data_management_engine/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/mattn/go-sqlite3"
)

type HostInfoWithoutPassword struct {
	Name        string `json:"name" binding:"required"`
	Ip          string `json:"ip" binding:"required"`
	Username    string `json:"user_name" binding:"required"`
	StorageType string `json:"storage_type" binding:"required,validateStorageType"`
}

type HostRegisterInfo struct {
	Name        string `json:"name" binding:"required"`
	Ip          string `json:"ip" binding:"required"`
	Username    string `json:"user_name" binding:"required"`
	Password    string `json:"password" binding:"required,validatePassword"`
	StorageType string `json:"storage_type" binding:"required,validateStorageType"`
}

type HostUnregisterInfo struct {
	Name string `json:"name"`
	Ip   string `json:"ip" binding:"required"`
}

func hostRegistrationHandler(c *gin.Context) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validateStorageType", storageTypeValidator)
	}

	var registerInfo HostRegisterInfo
	if err := c.ShouldBindJSON(&registerInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid request"})
		return
	}

	hostInfo := HostInfoWithoutPassword{
		Name:        registerInfo.Name,
		Ip:          registerInfo.Ip,
		Username:    registerInfo.Username,
		StorageType: registerInfo.StorageType,
	}

	hostModel := mgmtmodel.Host{
		Name:        registerInfo.Name,
		Ip:          registerInfo.Ip,
		Username:    registerInfo.Username,
		Password:    registerInfo.Password,
		StorageType: registerInfo.StorageType,
	}

	if err := hostModel.Register(); err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			// Map SQLite ErrNo to specific error scenarios
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique: // SQLite constraint violation
				c.JSON(http.StatusInternalServerError, gin.H{"Message": "The host information has already been registered.", "HostRegisterInfo": hostInfo, "Error": err})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"Message": "Failed to register the host.", "HostRegisterInfo": hostInfo, "Error": err})
			}
		}

	}

	c.JSON(http.StatusOK, gin.H{"Message": "Register the host information successfully.", "HostRegisterInfo": hostInfo})
}

func getRegisteredHostsHandler(c *gin.Context) {
	hostName := c.Query("name")
	hostIp := c.Query("ip")
	storageType := c.Query("storage_type")

	if hostName == "" && hostIp == "" {
		hostListModel := mgmtmodel.HostList{}
		hosts, err := hostListModel.Get(storageType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Message": "Failed to get the registered host.", "Error": err})
		}

		var hostInfoList []HostInfoWithoutPassword
		for _, host := range hosts {
			host := HostInfoWithoutPassword{
				Ip:          host.Ip,
				Name:        host.Name,
				Username:    host.Username,
				StorageType: host.StorageType,
			}

			hostInfoList = append(hostInfoList, host)
		}

		c.JSON(http.StatusOK, gin.H{"Message": "Get the registered hosts successfully.", "RegisteredHosts": hostInfoList})
		return
	} else {
		hostModel := mgmtmodel.Host{
			Ip:   hostIp,
			Name: hostName,
		}

		host, err := hostModel.Get()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Message": "Failed to get the registered host."})
			return
		}

		hostInfo := HostInfoWithoutPassword{
			Ip:          host.Ip,
			Name:        host.Name,
			Username:    host.Username,
			StorageType: host.StorageType,
		}

		c.JSON(http.StatusOK, gin.H{"Message": "Get the registered host successfully.", "RegisteredHosts": hostInfo})
	}
}

func hostUnregistrationHandler(c *gin.Context) {
	var unregister_info HostUnregisterInfo
	if err := c.ShouldBindJSON(&unregister_info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hostModel := mgmtmodel.Host{
		Name: unregister_info.Name,
		Ip:   unregister_info.Ip,
	}

	if err := hostModel.Unregister(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Failed to delete the registered host.", "HostUnregisterInfo": unregister_info, "Error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Unregister the host successfully.", "HostUnregisterInfo": unregister_info})
}

func storageTypeValidator(fl validator.FieldLevel) bool {
	storageType := fl.Field().String()

	storageTypeList := []string{"agent", "ontap", "magnascale"}

	return utils.In(storageType, storageTypeList)
}
