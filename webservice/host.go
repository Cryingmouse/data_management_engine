package webservice

import (
	"fmt"
	"net/http"

	"errors"

	"github.com/cryingmouse/data_management_engine/agent"
	"github.com/cryingmouse/data_management_engine/common"
	"github.com/cryingmouse/data_management_engine/mgmtmodel"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/mattn/go-sqlite3"
)

type HostResponse struct {
	IP             string `json:"ip,omitempty"`
	ComputerName   string `json:"name,omitempty"`
	StorageType    string `json:"storage_type,omitempty"`
	Caption        string `json:"os_type,omitempty"`
	OSArchitecture string `json:"os_arch,omitempty"`
	Version        string `json:"os_version,omitempty"`
	BuildNumber    string `json:"build_number,omitempty"`
}

type PaginationHostResponse struct {
	Hosts      []HostResponse `json:"hosts"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalCount int64          `json:"total_count"`
}

func hostRegistrationHandler(c *gin.Context) {
	request := []struct {
		IP          string `json:"ip" binding:"required,validateIP"`
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required,validatePassword"`
		StorageType string `json:"storage_type" binding:"required,validateStorageType"`
	}{}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validateStorageType", storageTypeValidator)
		v.RegisterValidation("validateIP", common.IPValidator)
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	var hostListModel mgmtmodel.HostList
	common.CopyStructList(request, &hostListModel.Hosts)

	if err := hostListModel.Register(); err != nil {
		for err != nil {
			err = errors.Unwrap(err)
			if sqliteErr, ok := err.(sqlite3.Error); ok {
				switch sqliteErr.ExtendedCode {
				// Map SQLite ErrNo to specific error scenarios
				case sqlite3.ErrConstraintUnique: // SQLite constraint violation
					c.JSON(http.StatusBadRequest, gin.H{"message": "The hosts have already been registered.", "error": err.Error()})
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register the hosts.", "error": err.Error()})
				return
			}
		}
	}

	if len(hostListModel.Hosts) == 1 {
		var host HostResponse
		common.CopyStructList(hostListModel.Hosts[0], &host)

		c.JSON(http.StatusOK, host)
	} else {
		var response []HostResponse
		common.CopyStructList(hostListModel.Hosts, &response)

		c.JSON(http.StatusOK, response)
	}
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
			filter := common.QueryFilter{
				Fields: common.SplitToList(fields),
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

			var hostInfoList []HostResponse

			common.CopyStructList(hosts, &hostInfoList)

			c.JSON(http.StatusOK, gin.H{"hosts": hostInfoList})
			return
		} else {
			// Query hosts with pagination.
			filter := common.QueryFilter{
				Fields: common.SplitToList(fields),
				Keyword: map[string]string{
					"name": hostNameKeyword,
				},
				Pagination: &common.Pagination{
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

			paginationHostList := PaginationHostResponse{
				Page:       page,
				Limit:      limit,
				TotalCount: paginationHosts.TotalCount,
			}

			common.CopyStructList(paginationHosts.Hosts, &paginationHostList.Hosts)

			c.JSON(http.StatusOK, paginationHostList)
			return

		}
	} else {
		// Using mgmtmodel.Host, to get the host.
		hostModel := mgmtmodel.Host{
			IP:           hostIp,
			ComputerName: hostName,
		}

		host, err := hostModel.Get()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get the registered host.", "error": err.Error()})
			return
		}

		var hostInfo HostResponse
		common.CopyStructList(host, &hostInfo)

		c.JSON(http.StatusOK, hostInfo)
	}
}

func hostUnregistrationHandler(c *gin.Context) {
	request := []struct {
		IP string `json:"ip" binding:"required,validateIP"`
	}{}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validateIP", common.IPValidator)
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	var hostListModel mgmtmodel.HostList
	common.CopyStructList(request, &hostListModel.Hosts)

	if err := hostListModel.Unregister(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to unregister the host.", "error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func getSystemInfoOnAgentHandler(c *gin.Context) {
	hostContext := common.HostContext{
		Username: c.Request.Header.Get("X-agent-username"),
		Password: c.Request.Header.Get("X-agent-password"),
	}

	agent := agent.GetAgent()

	if systemInfo, err := agent.GetSystemInfo(hostContext); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get system info on agent.", "error": err})

	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Get system info on agent successfully.", "system-info": systemInfo})
	}

}

func storageTypeValidator(fl validator.FieldLevel) bool {
	storageType := fl.Field().String()

	storageTypeList := []string{"agent", "ontap", "magnascale"}

	return common.In(storageType, storageTypeList)
}
