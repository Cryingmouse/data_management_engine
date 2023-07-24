package webservice

import (
	"net/http"
	"strconv"

	"github.com/cryingmouse/data_management_engine/agent"
	"github.com/cryingmouse/data_management_engine/common"
	"github.com/cryingmouse/data_management_engine/mgmtmodel"

	"github.com/gin-gonic/gin"
)

type CIFSShareResponse struct {
	HostIP          string   `json:"host_ip,omitempty"`
	Name            string   `json:"share_name,omitempty"`
	Path            string   `json:"share_path,omitempty"`
	DirectoryName   string   `json:"directory_name,omitempty"`
	Description     string   `json:"description,omitempty"`
	AccessUserNames []string `json:"access_users,omitempty"`
}

type PaginationShareResponse struct {
	Shares     []CIFSShareResponse `json:"shares"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalCount int64               `json:"total_count"`
}

func CreateShareHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request struct {
		HostIP          string   `json:"host_ip" binding:"required"`
		Name            string   `json:"share_name" binding:"required"`
		DirectoryName   string   `json:"directory_name" binding:"required"`
		Description     string   `json:"description" binding:"required"`
		AccessUserNames []string `json:"access_users" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	shareModel := mgmtmodel.CIFSShare{}
	common.DeepCopy(request, &shareModel)

	if err := shareModel.Create(ctx); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create the share", err.Error())
		return
	}

	shareResponse := CIFSShareResponse{}
	common.DeepCopy(shareModel, &shareResponse)

	c.JSON(http.StatusOK, shareResponse)
}

func DeleteShareHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request struct {
		HostIP string `json:"host_ip" binding:"required"`
		Name   string `json:"share_name"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	shareModel := mgmtmodel.CIFSShare{}
	common.DeepCopy(request, &shareModel)

	if err := shareModel.Delete(ctx); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create the share", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func MountShareHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request struct {
		DeviceName string `json:"device_name" binding:"required"`
		SharePath  string `json:"share_name" binding:"required"`
		UserName   string `json:"username" binding:"required"`
		Password   string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	agent := agent.GetAgent()

	err := agent.MountCIFSShare(ctx, request.DeviceName, request.SharePath, request.UserName, request.Password)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to mount the share", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func UnmountShareHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request struct {
		DeviceName string `json:"device_name"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	agent := agent.GetAgent()
	if err := agent.UnmountCIFSShare(ctx, request.DeviceName); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to unmount the share", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func GetSharesHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	shareName := c.Query("name")
	hostIP := c.Query("host_ip")
	fields := c.Query("fields")
	nameKeyword := c.Query("q")
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	if hostIP != "" && validateIPAddress(hostIP) != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", "")
		return
	}

	if shareName == "" || hostIP == "" {
		shareListModel := mgmtmodel.CIFSShareList{}
		filter := common.QueryFilter{
			Fields: common.SplitToList(fields),
			Keyword: map[string]string{
				"name": nameKeyword,
			},
			Conditions: struct {
				HostIP string
				Name   string
			}{
				HostIP: hostIP,
				Name:   shareName,
			},
		}

		if page == 0 && limit == 0 {
			// Query shares without pagination.
			shares, err := shareListModel.Get(ctx, &filter)
			if err != nil {
				ErrorResponse(c, http.StatusInternalServerError, "Failed to get the shares", err.Error())
				return
			}

			shareList := make([]CIFSShareResponse, len(shares))
			common.DeepCopy(shares, &shareList)

			c.JSON(http.StatusOK, shareList)
		} else {
			// Query directories with pagination.
			filter.Pagination = &common.Pagination{
				Page:     page,
				PageSize: limit,
			}

			paginationShares, err := shareListModel.Pagination(ctx, &filter)
			if err != nil {
				ErrorResponse(c, http.StatusInternalServerError, "Failed to get the shares", err.Error())
				return
			}

			paginationShareList := PaginationShareResponse{
				Page:       page,
				Limit:      limit,
				TotalCount: paginationShares.TotalCount,
			}

			common.DeepCopy(paginationShares.Shares, &paginationShareList.Shares)

			c.JSON(http.StatusOK, paginationShareList)
		}
	} else {
		shareModel := mgmtmodel.CIFSShare{
			Name:   shareName,
			HostIP: hostIP,
		}

		share, err := shareModel.Get(ctx)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "Failed to get the share", err.Error())
			return
		}

		shareInfo := CIFSShareResponse{}
		common.DeepCopy(share, &shareInfo)

		shareInfoList := []CIFSShareResponse{shareInfo}

		c.JSON(http.StatusOK, shareInfoList)
	}
}

func CreateShareOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request struct {
		ShareName     string   `json:"share_name" binding:"required"`
		DirectoryName string   `json:"directory_name" binding:"required"`
		Description   string   `json:"description" binding:"required"`
		UserNames     []string `json:"usernames" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	agent := agent.GetAgent()

	err := agent.CreateCIFSShare(ctx, request.ShareName, request.DirectoryName, request.Description, request.UserNames)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create the share", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func DeleteShareOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request struct {
		ShareName string `json:"share_name"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	agent := agent.GetAgent()
	if err := agent.DeleteCIFSShare(ctx, request.ShareName); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete the share", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func MountShareOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request struct {
		DeviceName string `json:"device_name" binding:"required"`
		SharePath  string `json:"share_name" binding:"required"`
		UserName   string `json:"username" binding:"required"`
		Password   string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	agent := agent.GetAgent()

	err := agent.MountCIFSShare(ctx, request.DeviceName, request.SharePath, request.UserName, request.Password)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to mount the share", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func UnmountShareOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request struct {
		DeviceName string `json:"device_name"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	agent := agent.GetAgent()
	if err := agent.UnmountCIFSShare(ctx, request.DeviceName); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to unmount the share", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func GetShareOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	name := c.Query("name")
	names := common.SplitToList(name)

	agent := agent.GetAgent()
	if len(names) == 0 {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", "")
	} else if len(names) == 1 {
		ShareDetail, err := agent.GetCIFSShareDetail(ctx, name)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "Failed to get the share detail", err.Error())
			return
		}

		c.JSON(http.StatusOK, ShareDetail)
	} else {
		SharesDetail, err := agent.GetCIFSSharesDetail(ctx, names)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "Failed to get the shares detail", err.Error())
			return
		}

		c.JSON(http.StatusOK, SharesDetail)
	}
}
