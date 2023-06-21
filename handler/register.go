package handler

import (
	"fmt"
	"net/http"

	"github.com/cryingmouse/data_management_engine/db"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"
)

type HostRegisterInfo struct {
	Name        string `json:"name"`
	Ip          string `json:"ip"`
	Username    string `json:"user_name"`
	Password    string `json:"password"`
	StorageType string `json:"storage_type"`
}

type HostUnregisterInfo struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

type HostInfoWithoutPassword struct {
	Name        string `json:"name"`
	Ip          string `json:"ip"`
	Username    string `json:"user_name"`
	StorageType string `json:"storage_type"`
}

func HostRegistrationHandler(c *gin.Context) {
	var registerInfo HostRegisterInfo
	if err := c.ShouldBindJSON(&registerInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hostInfo := HostInfoWithoutPassword{
		Name:        registerInfo.Name,
		Ip:          registerInfo.Ip,
		Username:    registerInfo.Username,
		StorageType: registerInfo.StorageType,
	}

	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	hostModel := db.Host{
		Name:        registerInfo.Name,
		Ip:          registerInfo.Ip,
		Username:    registerInfo.Username,
		Password:    registerInfo.Password,
		StorageType: registerInfo.StorageType,
	}

	if err = hostModel.Save(engine); err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			// Map SQLite ErrNo to specific error scenarios
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique: // SQLite constraint violation
				c.JSON(http.StatusInternalServerError, gin.H{"Message": "The host information has already been registered.", "HostRegisterInfo": hostInfo})
				return
			default:
				fmt.Println("Error")
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Register the host information successfully.", "HostRegisterInfo": hostInfo})
}

func GetRegisteredHostsHandler(c *gin.Context) {
	hostName := c.Query("name")
	hostIp := c.Query("ip")

	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	if hostName == "" && hostIp == "" {
		hostListModel := db.HostList{}

		hosts, err := hostListModel.Get(engine)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Message": "Failed to get the registered hosts."})
			return
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
		hostModel := db.Host{}
		host, err := hostModel.Get(engine, hostName, hostIp)
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

func HostUnregistrationHandler(c *gin.Context) {
	var unregister_info HostUnregisterInfo
	if err := c.ShouldBindJSON(&unregister_info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	hostInfoModel := db.Host{}
	host, err := hostInfoModel.Get(engine, unregister_info.Name, unregister_info.Ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Failed to get the registered host."})
		return
	}

	if err := host.Delete(engine); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Failed to delete the registered host."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Unregister the host successfully.", "HostRegisterInfo": unregister_info})
}
