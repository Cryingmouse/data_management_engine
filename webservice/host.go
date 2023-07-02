package webservice

import (
	"fmt"
	"net/http"

	"github.com/cryingmouse/data_management_engine/context"
	"github.com/cryingmouse/data_management_engine/mgmtmodel"
	"github.com/cryingmouse/data_management_engine/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/mattn/go-sqlite3"
)

type HostInfoWithoutPassword struct {
	Name        string `json:"name" binding:"required"`
	IP          string `json:"ip" binding:"required"`
	Username    string `json:"user_name" binding:"required"`
	StorageType string `json:"storage_type" binding:"required"`
}

type PaginationHostInfo struct {
	Hosts      []HostInfoWithoutPassword `json:"hosts"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalCount int64                     `json:"total_count"`
}

type HostRegisterInfo struct {
	Name        string `json:"name" binding:"required"`
	Ip          string `json:"ip" binding:"required"`
	Username    string `json:"user_name" binding:"required"`
	Password    string `json:"password" binding:"required,validatePassword"`
	StorageType string `json:"storage_type" binding:"required,validateStorageType"`
}

type HostUnregisterInfo struct {
	Ip   string `json:"ip" binding:"required"`
	Name string `json:"name"`
}

func hostRegistrationHandler(c *gin.Context) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validateStorageType", storageTypeValidator)
	}

	var registerInfo HostRegisterInfo
	if err := c.ShouldBindJSON(&registerInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	hostInfo := HostInfoWithoutPassword{
		Name:        registerInfo.Name,
		IP:          registerInfo.Ip,
		Username:    registerInfo.Username,
		StorageType: registerInfo.StorageType,
	}

	hostModel := mgmtmodel.Host{
		Name:        registerInfo.Name,
		IP:          registerInfo.Ip,
		Username:    registerInfo.Username,
		Password:    registerInfo.Password,
		StorageType: registerInfo.StorageType,
	}

	if err := hostModel.Register(); err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			// Map SQLite ErrNo to specific error scenarios
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique: // SQLite constraint violation
				c.JSON(http.StatusInternalServerError, gin.H{"message": "The host information has already been registered.", "HostRegisterInfo": hostInfo, "error": err.Error()})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register the host.", "HostRegisterInfo": hostInfo, "error": err.Error()})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Register the host information successfully.", "host": hostInfo})
}

func getRegisteredHostsHandler(c *gin.Context) {
	hostName := c.Query("name")
	hostIp := c.Query("ip")
	storageType := c.Query("storage_type")
	fields := c.Query("fields")
	hostNameKeyword := c.Query("q")

	page, limit, err := validatePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request.", "error": err.Error()})
		return
	}

	if hostName == "" && hostIp == "" {
		// Using mgmtmodel.HostList, to get the list of the host.
		hostListModel := mgmtmodel.HostList{}
		if page == 0 && limit == 0 {
			// Query hosts without pagination.
			filter := context.QueryFilter{
				Fields: utils.SplitToList(fields),
				Keyword: map[string]string{
					"name": hostNameKeyword,
				},
				Conditions: struct {
					StorageType string
				}{
					StorageType: storageType,
				},
			}
			hosts, err := hostListModel.Get(&filter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("Failed to get the hosts with the parameters: storage_type=%s", storageType),
					"error":   err.Error(),
				})
				return
			}

			var hostInfoList []HostInfoWithoutPassword
			for _, host := range hosts {
				host := HostInfoWithoutPassword{
					IP:          host.IP,
					Name:        host.Name,
					Username:    host.Username,
					StorageType: host.StorageType,
				}

				hostInfoList = append(hostInfoList, host)
			}

			c.JSON(http.StatusOK, gin.H{"message": "Get the registered hosts successfully.", "hosts": hostInfoList})
			return
		} else {
			// Query hosts with pagination.
			filter := context.QueryFilter{
				Fields: utils.SplitToList(fields),
				Keyword: map[string]string{
					"name": hostNameKeyword,
				},
				Pagination: &context.Pagination{
					Page:     page,
					PageSize: limit,
				},
				Conditions: struct {
					StorageType string
				}{
					StorageType: storageType,
				},
			}
			paginationHosts, err := hostListModel.Pagination(&filter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("Failed to get the host with the parameters: storage_type=%s,page=%d,limit=%d", storageType, page, limit),
					"error":   err.Error(),
				})
				return
			}

			paginationHostList := PaginationHostInfo{
				Page:       page,
				Limit:      limit,
				TotalCount: paginationHosts.TotalCount,
			}

			for _, _host := range paginationHosts.Hosts {
				host := HostInfoWithoutPassword{
					IP:          _host.IP,
					Name:        _host.Name,
					Username:    _host.Username,
					StorageType: _host.StorageType,
				}

				paginationHostList.Hosts = append(paginationHostList.Hosts, host)
			}

			c.JSON(http.StatusOK, gin.H{"message": "Get the hosts successfully.", "pagination": paginationHostList})
			return

		}
	} else {
		// Using mgmtmodel.Host, to get the host.
		hostModel := mgmtmodel.Host{
			IP:   hostIp,
			Name: hostName,
		}

		host, err := hostModel.Get()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get the registered host.", "error": err.Error()})
			return
		}

		hostInfo := HostInfoWithoutPassword{
			IP:          host.IP,
			Name:        host.Name,
			Username:    host.Username,
			StorageType: host.StorageType,
		}

		c.JSON(http.StatusOK, gin.H{"message": "Get the registered host successfully.", "host": hostInfo})
	}
}

func hostUnregistrationHandler(c *gin.Context) {
	var unregister_info HostUnregisterInfo
	if err := c.ShouldBindJSON(&unregister_info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hostModel := mgmtmodel.Host{
		IP:   unregister_info.Ip,
		Name: unregister_info.Name,
	}

	if err := hostModel.Unregister(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete the registered host.", "HostUnregisterInfo": unregister_info, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unregister the host successfully.", "host": unregister_info})
}

func storageTypeValidator(fl validator.FieldLevel) bool {
	storageType := fl.Field().String()

	storageTypeList := []string{"agent", "ontap", "magnascale"}

	return utils.In(storageType, storageTypeList)
}
