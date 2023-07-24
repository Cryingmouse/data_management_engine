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

const (
	StorageTypeAgent      = "agent"
	StorageTypeOntap      = "ontap"
	StorageTypeMagnascale = "magnascale"
)

type HostResponse struct {
	IP             string `json:"ip,omitempty"`
	ComputerName   string `json:"name,omitempty"`
	StorageType    string `json:"storage_type,omitempty"`
	Caption        string `json:"os_type,omitempty"`
	OSArchitecture string `json:"os_arch,omitempty"`
	OSVersion      string `json:"os_version,omitempty"`
	BuildNumber    string `json:"build_number,omitempty"`
	Username       string `json:"username,omitempty"`
	Password       string `json:"password,omitempty"`
}

type PaginationHostResponse struct {
	Hosts      []HostResponse `json:"hosts"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalCount int64          `json:"total_count"`
}

func RegisterHostHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request struct {
		IP          string `json:"ip" binding:"required,ip"`
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required,validatePassword"`
		StorageType string `json:"storage_type" binding:"required,oneof=agent ontap magnascale"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	hostModel := mgmtmodel.Host{}
	common.DeepCopy(request, &hostModel)

	if err := hostModel.Register(ctx); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to register the host", err.Error())
		return
	}

	hostResponse := HostResponse{}
	common.DeepCopy(hostModel, &hostResponse)

	c.JSON(http.StatusOK, hostResponse)
}

func RegisterHostsHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request []struct {
		IP          string `json:"ip" binding:"required,ip"`
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required,validatePassword"`
		StorageType string `json:"storage_type" binding:"required,oneof=agent ontap magnascale"`
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validateStorageType", StorageTypeValidator)
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	hostListModel := mgmtmodel.HostList{}
	common.DeepCopy(request, &hostListModel.Hosts)

	if err := hostListModel.Register(ctx); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to register the hosts", err.Error())
		return
	}

	hostResponseList := make([]HostResponse, len(hostListModel.Hosts))
	common.DeepCopy(hostListModel.Hosts, &hostResponseList)

	c.JSON(http.StatusOK, hostResponseList)
}

func UnregisterHostHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request struct {
		IP string `json:"ip" binding:"required,ip"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	hostModel := mgmtmodel.Host{}
	common.DeepCopy(request, &hostModel)

	if err := hostModel.Unregister(ctx); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to unregister the host", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func UnregisterHostsHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request []struct {
		IP string `json:"ip" binding:"required,ip"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	hostListModel := mgmtmodel.HostList{}
	common.DeepCopy(request, &hostListModel.Hosts)

	if err := hostListModel.Unregister(ctx); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to unregister the hosts", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func GetRegisteredHostsHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	hostName := c.Query("name")
	hostIP := c.Query("ip")
	fields := c.Query("fields")
	storageType := c.DefaultQuery("storage_type", StorageTypeAgent)
	nameKeyword := c.Query("name-like")
	osTypeKeyword := c.Query("os_type-like")
	page, errPage := strconv.Atoi(c.Query("page"))
	limit, errLimit := strconv.Atoi(c.Query("limit"))

	if (errPage != nil && errLimit == nil) || (errPage == nil && errLimit != nil) || (hostIP != "" && validateIPAddress(hostIP) != nil) {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", "")
		return
	}

	if hostIP == "" && hostName == "" {
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
			hosts, err := hostListModel.Get(ctx, &filter)
			if err != nil {
				ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to get the hosts with the parameters: storage_type=%s", storageType), err.Error())
				return
			}

			hostResponseList := make([]HostResponse, len(hosts))
			common.DeepCopy(hosts, &hostResponseList)

			c.JSON(http.StatusOK, hostResponseList)
		} else {
			filter.Pagination = &common.Pagination{
				Page:     page,
				PageSize: limit,
			}

			paginationHosts, err := hostListModel.Pagination(ctx, &filter)
			if err != nil {
				ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to get the hosts with the parameters: storage_type=%s, page=%d, limit=%d", storageType, page, limit), err.Error())
				return
			}

			paginationHostResponse := PaginationHostResponse{
				Page:       page,
				Limit:      limit,
				TotalCount: paginationHosts.TotalCount,
			}
			common.DeepCopy(paginationHosts.Hosts, &paginationHostResponse.Hosts)

			c.JSON(http.StatusOK, paginationHostResponse)
		}
	} else {
		hostModel := mgmtmodel.Host{
			IP:           hostIP,
			ComputerName: hostName,
		}

		host, err := hostModel.Get(ctx)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "Failed to get the registered host", err.Error())
			return
		}

		hostResponse := HostResponse{}
		common.DeepCopy(host, &hostResponse)

		hostResponseList := []HostResponse{hostResponse}

		c.JSON(http.StatusOK, hostResponseList)
	}
}

func GetSystemInfoOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	agent := agent.GetAgent()

	if systemInfo, err := agent.GetSystemInfo(ctx); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to get system info on agent", err.Error())
	} else {
		c.JSON(http.StatusOK, systemInfo)
	}
}
