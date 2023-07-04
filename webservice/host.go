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

type HostRegisterResponse struct {
	IP             string `json:"ip"`
	ComputerName   string `json:"name"`
	Username       string `json:"username"`
	StorageType    string `json:"storage_type"`
	Caption        string `json:"os_type"`
	OSArchitecture string `json:"os_arch"`
	Version        string `json:"os_version"`
	BuildNumber    string `json:"build_number"`
}

type PaginationHostInfo struct {
	Hosts      []HostRegisterResponse `json:"hosts"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalCount int64                  `json:"total_count"`
}

func hostRegistrationHandler(c *gin.Context) {
	type HostRegisterRequest struct {
		IP          string `json:"ip" binding:"required"`
		Username    string `json:"user_name" binding:"required"`
		Password    string `json:"password" binding:"required,validatePassword"`
		StorageType string `json:"storage_type" binding:"required,validateStorageType"`
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validateStorageType", storageTypeValidator)
	}

	var registerInfo HostRegisterRequest
	if err := c.ShouldBindJSON(&registerInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	hostModel := mgmtmodel.Host{
		IP:          registerInfo.IP,
		Username:    registerInfo.Username,
		Password:    registerInfo.Password,
		StorageType: registerInfo.StorageType,
	}

	if err := hostModel.Register(); err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			// Map SQLite ErrNo to specific error scenarios
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique: // SQLite constraint violation
				c.JSON(http.StatusBadRequest, gin.H{"message": "The host information has already been registered.", "error": err.Error()})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register the host.", "error": err.Error()})
			}
		}
	}

	host := HostRegisterResponse{
		IP:             hostModel.IP,
		ComputerName:   hostModel.Name,
		Username:       hostModel.Username,
		StorageType:    hostModel.StorageType,
		Caption:        hostModel.Caption,
		OSArchitecture: hostModel.OSArchitecture,
		Version:        hostModel.Version,
		BuildNumber:    hostModel.BuildNumber,
	}

	c.JSON(http.StatusOK, gin.H{"message": "Register the host information successfully.", "host": host})
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

			var hostInfoList []HostRegisterResponse
			for _, host := range hosts {
				host := HostRegisterResponse{
					IP:             host.IP,
					ComputerName:   host.Name,
					Username:       host.Username,
					StorageType:    host.StorageType,
					Caption:        host.Caption,
					OSArchitecture: host.OSArchitecture,
					Version:        host.Version,
					BuildNumber:    host.BuildNumber,
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
				host := HostRegisterResponse{
					IP:             _host.IP,
					ComputerName:   _host.Name,
					Username:       _host.Username,
					StorageType:    _host.StorageType,
					Caption:        _host.Caption,
					OSArchitecture: _host.OSArchitecture,
					Version:        _host.Version,
					BuildNumber:    _host.BuildNumber,
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

		hostInfo := HostRegisterResponse{
			IP:             host.IP,
			ComputerName:   host.Name,
			Username:       host.Username,
			StorageType:    host.StorageType,
			Caption:        host.Caption,
			OSArchitecture: host.OSArchitecture,
			Version:        host.Version,
			BuildNumber:    host.BuildNumber,
		}

		c.JSON(http.StatusOK, gin.H{"message": "Get the registered host successfully.", "host": hostInfo})
	}
}

func hostUnregistrationHandler(c *gin.Context) {
	type HostUnregisterInfo struct {
		IP string `json:"ip" binding:"required"`
	}

	var unregister_info HostUnregisterInfo
	if err := c.ShouldBindJSON(&unregister_info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hostModel := mgmtmodel.Host{
		IP: unregister_info.IP,
	}

	if err := hostModel.Unregister(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to unregister the host.", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unregister the host successfully.", "host": unregister_info.IP})
}

func storageTypeValidator(fl validator.FieldLevel) bool {
	storageType := fl.Field().String()

	storageTypeList := []string{"agent", "ontap", "magnascale"}

	return utils.In(storageType, storageTypeList)
}
