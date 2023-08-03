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
	log "github.com/sirupsen/logrus"
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
}

type PaginationHostResponse struct {
	Hosts      []HostResponse `json:"hosts"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalCount int64          `json:"total_count"`
}

func RegisterHostHandler(c *gin.Context) {
	ctx, traceID := SetTraceIDToContext(c)

	var request struct {
		IP          string `json:"ip" binding:"required,ip"`
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required,validatePassword"`
		StorageType string `json:"storage_type" binding:"required,oneof=agent ontap magnascale"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		common.Logger.WithFields(log.Fields{
			"TraceID": traceID,
			"error":   err.Error(),
		}).Error("Invalid request.")

		SetErrorToContext(c, "E0100010001", err)
		return
	}

	hostModel := mgmtmodel.Host{}
	common.DeepCopy(request, &hostModel)
	common.Logger.WithFields(log.Fields{
		"TraceID":   traceID,
		"HostModel": common.MaskPassword(hostModel),
	}).Debug("Copy host model.")

	if err := hostModel.Register(ctx); err != nil {
		common.Logger.WithFields(log.Fields{
			"TraceID":   traceID,
			"HostModel": common.MaskPassword(hostModel),
			"error":     err.Error(),
		}).Error("Failed to register the host.")

		SetErrorToContext(c, "", err)
		return
	}

	common.Logger.WithFields(log.Fields{
		"TraceID":   traceID,
		"HostModel": common.MaskPassword(hostModel),
	}).Debug("Rregister the host successfully.")

	hostResponse := HostResponse{}
	common.DeepCopy(hostModel, &hostResponse)

	c.JSON(http.StatusOK, hostResponse)
}

func RegisterHostsHandler(c *gin.Context) {
	ctx, traceID := SetTraceIDToContext(c)

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
		common.Logger.WithFields(log.Fields{
			"TraceID": traceID,
			"error":   err.Error(),
		}).Error("Invalid request.")

		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	hostListModel := mgmtmodel.HostList{}
	common.DeepCopy(request, &hostListModel.Hosts)
	common.Logger.WithFields(log.Fields{
		"TraceID":       traceID,
		"HostListModel": common.MaskPassword(hostListModel),
	}).Debug("Copy host list model.")

	if err := hostListModel.Register(ctx); err != nil {
		common.Logger.WithFields(log.Fields{
			"TraceID":       traceID,
			"HostListModel": hostListModel,
			"error":         err.Error(),
		}).Error("Failed to register the hosts.")

		ErrorResponse(c, http.StatusInternalServerError, "Failed to register the hosts", err.Error())
		return
	}

	common.Logger.WithFields(log.Fields{
		"TraceID":       traceID,
		"HostListModel": common.MaskPassword(hostListModel),
	}).Debug("Register the hosts successfully.")

	hostResponseList := make([]HostResponse, len(hostListModel.Hosts))
	common.DeepCopy(hostListModel.Hosts, &hostResponseList)

	c.JSON(http.StatusOK, hostResponseList)
}

func UnregisterHostHandler(c *gin.Context) {
	ctx, traceID := SetTraceIDToContext(c)

	var request struct {
		IP string `json:"ip" binding:"required,ip"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		common.Logger.WithFields(log.Fields{
			"TraceID": traceID,
			"error":   err.Error(),
		}).Error("Invalid request.")

		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	hostModel := mgmtmodel.Host{}
	common.DeepCopy(request, &hostModel)
	common.Logger.WithFields(log.Fields{
		"TraceID":   traceID,
		"HostModel": hostModel,
	}).Debug("Copy host model.")

	if err := hostModel.Unregister(ctx); err != nil {
		common.Logger.WithFields(log.Fields{
			"TraceID":   traceID,
			"HostModel": hostModel,
			"error":     err.Error(),
		}).Error("Failed to unregister the host.")

		ErrorResponse(c, http.StatusInternalServerError, "Failed to unregister the host", err.Error())
		return
	}

	common.Logger.WithFields(log.Fields{
		"TraceID":   traceID,
		"HostModel": hostModel,
	}).Debug("Unregister the host successfully.")

	c.Status(http.StatusOK)
}

func UnregisterHostsHandler(c *gin.Context) {
	ctx, traceID := SetTraceIDToContext(c)

	var request []struct {
		IP string `json:"ip" binding:"required,ip"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		common.Logger.WithFields(log.Fields{
			"TraceID": traceID,
			"error":   err.Error(),
		}).Error("Invalid request.")
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	hostListModel := mgmtmodel.HostList{}
	common.DeepCopy(request, &hostListModel.Hosts)
	common.Logger.WithFields(log.Fields{
		"TraceID":       traceID,
		"HostListModel": hostListModel,
	}).Debug("Copy host list model.")

	if err := hostListModel.Unregister(ctx); err != nil {
		common.Logger.WithFields(log.Fields{
			"TraceID":       traceID,
			"HostListModel": hostListModel,
			"error":         err.Error(),
		}).Error("Failed to unregister the hosts.")

		ErrorResponse(c, http.StatusInternalServerError, "Failed to unregister the hosts", err.Error())
		return
	}

	common.Logger.WithFields(log.Fields{
		"TraceID":       traceID,
		"HostListModel": hostListModel,
	}).Debug("Unregister the hosts successfully.")

	c.Status(http.StatusOK)
}

func GetRegisteredHostsHandler(c *gin.Context) {
	ctx, traceID := SetTraceIDToContext(c)

	hostName := c.Query("name")
	hostIP := c.Query("ip")
	fields := c.Query("fields")
	storageType := c.DefaultQuery("storage_type", StorageTypeAgent)
	nameKeyword := c.Query("name-like")
	osTypeKeyword := c.Query("os_type-like")
	page, errPage := strconv.Atoi(c.Query("page"))
	limit, errLimit := strconv.Atoi(c.Query("limit"))

	if (errPage != nil && errLimit == nil) || (errPage == nil && errLimit != nil) || (hostIP != "" && validateIPAddress(hostIP) != nil) {
		common.Logger.WithFields(log.Fields{
			"TraceID": traceID,
			"URL":     c.Request.URL,
		}).Error("Invalid request.")

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
				common.Logger.WithFields(log.Fields{
					"TraceID": traceID,
					"filter":  filter,
					"error":   err.Error(),
				}).Error("Failed to get the hosts by the filter.")

				ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to get the hosts with the fitler=%v", filter), err.Error())
				return
			}

			common.Logger.WithFields(log.Fields{
				"TraceID": traceID,
				"Hosts":   common.MaskPassword(hosts),
			}).Debug("Qurey the hosts successfully.")

			hostListResponse := make([]HostResponse, len(hosts))
			common.DeepCopy(hosts, &hostListResponse)
			common.Logger.WithFields(log.Fields{
				"TraceID":          traceID,
				"HostListResponse": common.MaskPassword(hostListResponse),
			}).Debug("Copy the hosts response successfully.")

			c.JSON(http.StatusOK, hostListResponse)
		} else {
			filter.Pagination = &common.Pagination{
				Page:     page,
				PageSize: limit,
			}

			paginationHosts, err := hostListModel.Pagination(ctx, &filter)
			if err != nil {
				common.Logger.WithFields(log.Fields{
					"TraceID": traceID,
					"filter":  filter,
					"error":   err.Error(),
				}).Error("Failed to pagination the hosts by the filter.")

				ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to pagination the hosts with the fitler=%v", filter), err.Error())
				return
			}

			common.Logger.WithFields(log.Fields{
				"TraceID":         traceID,
				"PaginationHosts": common.MaskPassword(paginationHosts),
			}).Debug("Qurey the pagination hosts successfully.")

			paginationHostResponse := PaginationHostResponse{
				Page:       page,
				Limit:      limit,
				TotalCount: paginationHosts.TotalCount,
			}
			common.DeepCopy(paginationHosts.Hosts, &paginationHostResponse.Hosts)

			common.Logger.WithFields(log.Fields{
				"TraceID":                traceID,
				"PaginationHostResponse": common.MaskPassword(paginationHostResponse),
			}).Debug("Copy the pagination hosts response successfully.")

			c.JSON(http.StatusOK, paginationHostResponse)
		}
	} else {
		hostModel := mgmtmodel.Host{
			IP:           hostIP,
			ComputerName: hostName,
		}

		host, err := hostModel.Get(ctx)
		if err != nil {
			common.Logger.WithFields(log.Fields{
				"TraceID":   traceID,
				"HostModel": hostModel,
				"error":     err.Error(),
			}).Error("Failed to get the host.")

			ErrorResponse(c, http.StatusInternalServerError, "Failed to get the registered host", err.Error())
			return
		}

		common.Logger.WithFields(log.Fields{
			"TraceID": traceID,
			"Host":    common.MaskPassword(host),
		}).Debug("Qurey the host successfully.")

		hostResponse := HostResponse{}
		common.DeepCopy(host, &hostResponse)
		common.Logger.WithFields(log.Fields{
			"TraceID":      traceID,
			"HostResponse": common.MaskPassword(hostResponse),
		}).Debug("Copy the host response successfully.")

		hostResponseList := []HostResponse{hostResponse}

		c.JSON(http.StatusOK, hostResponseList)
	}
}

func GetSystemInfoOnAgentHandler(c *gin.Context) {
	ctx, traceID := SetTraceIDToContext(c)

	agent := agent.GetAgent()

	if systemInfo, err := agent.GetSystemInfo(ctx); err != nil {
		common.Logger.WithFields(log.Fields{
			"TraceID": traceID,
			"error":   err.Error(),
		}).Error("Failed to get system info on agent.")

		ErrorResponse(c, http.StatusInternalServerError, "Failed to get system info on agent.", err.Error())
	} else {
		common.Logger.WithFields(log.Fields{
			"TraceID":    traceID,
			"SystemInfo": systemInfo,
		}).Debug("Get system info on agent successfully.")

		c.JSON(http.StatusOK, systemInfo)
	}
}
