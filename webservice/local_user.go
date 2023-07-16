package webservice

import (
	"net/http"
	"strconv"

	"github.com/cryingmouse/data_management_engine/agent"
	"github.com/cryingmouse/data_management_engine/common"
	"github.com/cryingmouse/data_management_engine/mgmtmodel"

	"github.com/gin-gonic/gin"
)

type LocalUserResponse struct {
	HostIP               string `json:"host_ip,omitempty"`
	Name                 string `json:"name,omitempty"`
	UID                  string `json:"id,omitempty"`
	FullName             string `json:"full_name,omitempty"`
	Description          string `json:"description,omitempty"`
	Status               string `json:"status,omitempty"`
	IsDisabled           bool   `json:"disabled"`
	IsPasswordRequired   bool   `json:"is_password_required"`
	IsPasswordExpired    bool   `json:"is_password_expired"`
	IsPasswordChangeable bool   `json:"is_password_changeable"`
	IsLockout            bool   `json:"is_lockout"`
}

type PaginationLocalUserResponse struct {
	LocalUsers []LocalUserResponse `json:"users"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalCount int64               `json:"total_count"`
}

type requestLocalUser struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,validatePassword"`
	HostIP   string `json:"host_ip" binding:"required"`
}

func CreateLocalUserHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request requestLocalUser
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	var localUserModel mgmtmodel.LocalUser
	common.CopyStructList(request, &localUserModel)

	if err := localUserModel.Create(ctx); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create the local user", err.Error())
		return
	}

	localUserResponse := LocalUserResponse{}
	common.CopyStructList(localUserModel, &localUserResponse)

	c.JSON(http.StatusOK, localUserResponse)
}

func CreateLocalUsersHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request []requestLocalUser
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	var localUserListModel mgmtmodel.LocalUserList
	common.CopyStructList(request, &localUserListModel.LocalUsers)

	if err := localUserListModel.Create(ctx); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create the local users", err.Error())
		return
	}

	localUserResponseList := make([]LocalUserResponse, len(localUserListModel.LocalUsers))
	common.CopyStructList(localUserListModel.LocalUsers, &localUserResponseList)

	c.JSON(http.StatusOK, localUserResponseList)
}

func DeleteLocalUserHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request requestLocalUser
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	var localUserModel mgmtmodel.LocalUser
	common.CopyStructList(request, &localUserModel)

	if err := localUserModel.Delete(ctx); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete the local user", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func DeleteLocalUsersHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request []requestLocalUser
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	var localUserListModel mgmtmodel.LocalUserList
	common.CopyStructList(request, &localUserListModel.LocalUsers)

	if err := localUserListModel.Delete(ctx, nil); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete the local users", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func GetlocalUsersHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	userName := c.Query("name")
	isLockout := c.Query("is_lockout")
	hostIP := c.Query("host_ip")
	fields := c.Query("fields")

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	if userName == "" {
		localUserListModel := mgmtmodel.LocalUserList{}

		filter := common.QueryFilter{
			Fields: common.SplitToList(fields),
			Conditions: struct {
				HostIP    string
				Name      string
				IsLockout string
			}{
				HostIP:    hostIP,
				Name:      userName,
				IsLockout: isLockout,
			},
		}

		if page == 0 && limit == 0 {
			// Query local users without pagination.
			localUsers, err := localUserListModel.Get(ctx, &filter)
			if err != nil {
				ErrorResponse(c, http.StatusInternalServerError, "Failed to get the local users", err.Error())
				return
			}

			localUserList := make([]LocalUserResponse, len(localUsers))
			common.CopyStructList(localUsers, &localUserList)

			c.JSON(http.StatusOK, localUserList)
		} else {
			// Query directories with pagination.
			filter.Pagination = &common.Pagination{
				Page:     page,
				PageSize: limit,
			}

			paginationLocalUsers, err := localUserListModel.Pagination(ctx, &filter)
			if err != nil {
				ErrorResponse(c, http.StatusInternalServerError, "Failed to get the local users", err.Error())
				return
			}

			paginationlocalUserList := PaginationLocalUserResponse{
				Page:       page,
				Limit:      limit,
				TotalCount: paginationLocalUsers.TotalCount,
			}

			common.CopyStructList(paginationLocalUsers.LocalUsers, &paginationlocalUserList.LocalUsers)

			c.JSON(http.StatusOK, paginationlocalUserList)
		}
	} else {
		userModel := mgmtmodel.LocalUser{
			Name:   userName,
			HostIP: hostIP,
		}

		localUser, err := userModel.Get(ctx)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "Failed to get the local user", err.Error())
			return
		}

		localUserInfo := DirectoryResponse{}
		common.CopyStructList(localUser, &localUserInfo)

		directoryInfoList := []DirectoryResponse{localUserInfo}

		c.JSON(http.StatusOK, directoryInfoList)
	}
}

func ManageLocalUserHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request requestLocalUser
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	var localUserModel mgmtmodel.LocalUser
	common.CopyStructList(request, &localUserModel)

	if err := localUserModel.Manage(ctx); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to manage the local user", err.Error())
		return
	}

	localUserResponse := LocalUserResponse{}
	common.CopyStructList(localUserModel, &localUserResponse)

	c.JSON(http.StatusOK, localUserResponse)
}

func ManageLocalUsersHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request []requestLocalUser
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	var localUserListModel mgmtmodel.LocalUserList
	common.CopyStructList(request, &localUserListModel.LocalUsers)

	if err := localUserListModel.Manage(ctx); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create the local users", err.Error())
		return
	}

	localUserResponseList := make([]LocalUserResponse, len(localUserListModel.LocalUsers))
	common.CopyStructList(localUserListModel.LocalUsers, &localUserResponseList)

	c.JSON(http.StatusOK, localUserResponseList)
}

func UnmanageLocalUserHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request requestLocalUser
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	var localUserModel mgmtmodel.LocalUser
	common.CopyStructList(request, &localUserModel)

	if err := localUserModel.Unmanage(ctx); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete the local user", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func UnmanageLocalUsersHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request []requestLocalUser
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	var localUserListModel mgmtmodel.LocalUserList
	common.CopyStructList(request, &localUserListModel.LocalUsers)

	if err := localUserListModel.Unmanage(ctx, nil); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete the local users", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func createLocalUserOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,validatePassword"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	hostContext := common.HostContext{
		Username: c.Request.Header.Get("X-agent-username"),
		Password: c.Request.Header.Get("X-agent-password"),
	}

	agent := agent.GetAgent()
	err := agent.CreateLocalUser(ctx, hostContext, request.Name, request.Password)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create the local user", err.Error())
		return
	}

	localUserDetail, err := agent.GetLocalUserDetail(ctx, hostContext, request.Name)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to get the local user detail", err.Error())
		return
	}

	c.JSON(http.StatusOK, localUserDetail)
}

func deleteLocalUserOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	hostContext := common.HostContext{
		Username: c.Request.Header.Get("X-agent-username"),
		Password: c.Request.Header.Get("X-agent-password"),
	}

	agent := agent.GetAgent()
	if err := agent.DeleteLocalUser(ctx, hostContext, request.Name); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete the local user", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func getLocalUserOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	name := c.Query("name")
	names := common.SplitToList(name)

	hostContext := common.HostContext{
		Username: c.Request.Header.Get("X-agent-username"),
		Password: c.Request.Header.Get("X-agent-password"),
	}

	agent := agent.GetAgent()
	if len(names) == 0 {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", "")
	} else if len(names) == 1 {
		localUserDetail, err := agent.GetLocalUserDetail(ctx, hostContext, name)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "Failed to get the local user detail", err.Error())
			return
		}

		c.JSON(http.StatusOK, localUserDetail)
	} else {
		localUsersDetail, err := agent.GetLocalUsersDetail(ctx, hostContext, names)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "Failed to get the local users detail", err.Error())
			return
		}

		c.JSON(http.StatusOK, localUsersDetail)
	}
}
