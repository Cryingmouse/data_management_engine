package webservice

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cryingmouse/data_management_engine/agent"
	"github.com/cryingmouse/data_management_engine/common"
	"github.com/cryingmouse/data_management_engine/mgmtmodel"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type HostResponse struct {
	IP             string `json:"ip,omitempty"`
	ComputerName   string `json:"name,omitempty"`
	StorageType    string `json:"storage_type,omitempty"`
	Caption        string `json:"os_type,omitempty"`
	OSArchitecture string `json:"os_arch,omitempty"`
	OSVersion      string `json:"os_version,omitempty"`
	BuildNumber    string `json:"build_number,omitempty"`
}

type PaginationHostResponse struct {
	Hosts      []HostResponse `json:"hosts"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalCount int64          `json:"total_count"`
}

func registerHostHandler(c *gin.Context) {
	request := struct {
		IP          string `json:"ip" binding:"required,ip"`
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required,validatePassword"`
		StorageType string `json:"storage_type" binding:"required,validateStorageType"`
	}{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	var hostModel mgmtmodel.Host
	common.CopyStructList(request, &hostModel)

	if err := hostModel.Register(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register the hosts.", "error": err.Error()})
		return
	}

	var response HostResponse
	common.CopyStructList(hostModel, &response)

	c.JSON(http.StatusOK, response)
}

func registerHostsHandler(c *gin.Context) {
	request := []struct {
		IP          string `json:"ip" binding:"required,ip"`
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required,validatePassword"`
		StorageType string `json:"storage_type" binding:"required,validateStorageType"`
	}{}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validateStorageType", storageTypeValidator)
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	var hostListModel mgmtmodel.HostList
	common.CopyStructList(request, &hostListModel.Hosts)

	if err := hostListModel.Register(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register the hosts.", "error": err.Error()})
		return
	}

	var response []HostResponse
	common.CopyStructList(hostListModel.Hosts, &response)

	c.JSON(http.StatusOK, response)
}

func unregisterHostHandler(c *gin.Context) {
	request := struct {
		IP string `json:"ip" binding:"required,ip"`
	}{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	var hostModel mgmtmodel.Host
	common.CopyStructList(request, &hostModel)

	if err := hostModel.Unregister(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to unregister the host.", "error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func unregisterHostsHandler(c *gin.Context) {
	request := []struct {
		IP string `json:"ip" binding:"required,ip"`
	}{}

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

func getRegisteredHostsHandler(c *gin.Context) {
	hostName := c.Query("name")
	hostIP := c.Query("ip")
	fields := c.Query("fields")
	storageType := c.DefaultQuery("storage_type", "agent")
	nameKeyword := c.Query("name-like")
	osTypeKeyword := c.Query("os_type-like")
	page, err_page := strconv.Atoi(c.Query("page"))
	limit, err_limit := strconv.Atoi(c.Query("limit"))

	if (err_page != nil && err_limit == nil) || (err_page == nil && err_limit != nil) || (hostIP != "" && validateIPAddress(hostIP) != nil) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request."})
		return
	}

	if hostIP == "" && hostName == "" {
		// Using mgmtmodel.HostList, to get the list of the host with filter.
		hostListModel := mgmtmodel.HostList{}

		filter := common.QueryFilter{
			Fields: common.SplitToList(fields),
			Keyword: map[string]string{
				"name":    nameKeyword,
				"os_type": osTypeKeyword,
			},
			Conditions: struct {
				StorageType string
			}{
				StorageType: storageType,
			},
		}

		if page == 0 && limit == 0 {
			// Query hosts without pagination.
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

			c.JSON(http.StatusOK, hostInfoList)
		} else {
			// Add the pagination into filter, and then query hosts with pagination.
			filter.Pagination = &common.Pagination{
				Page:     page,
				PageSize: limit,
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
		}
	} else {
		// Using mgmtmodel.Host to get the host
		hostModel := mgmtmodel.Host{
			IP:           hostIP,
			ComputerName: hostName,
		}

		host, err := hostModel.Get()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get the registered host.", "error": err.Error()})
			return
		}

		var hostInfo HostResponse
		common.CopyStructList(host, &hostInfo)

		hostInfoList := []HostResponse{}
		hostInfoList = append(hostInfoList, hostInfo)

		c.JSON(http.StatusOK, hostInfoList)
	}
}

func getSystemInfoOnAgentHandler(c *gin.Context) {
	hostContext := common.HostContext{
		Username: c.Request.Header.Get("X-agent-username"),
		Password: c.Request.Header.Get("X-agent-password"),
	}

	agent := agent.GetAgent()

	if systemInfo, err := agent.GetSystemInfo(hostContext); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, systemInfo)
	}
}

func storageTypeValidator(fl validator.FieldLevel) bool {
	storageType := fl.Field().String()

	storageTypeList := []string{"agent", "ontap", "magnascale"}

	return common.In(storageType, storageTypeList)
}
