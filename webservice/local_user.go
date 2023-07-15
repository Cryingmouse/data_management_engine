package webservice

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cryingmouse/data_management_engine/agent"
	"github.com/cryingmouse/data_management_engine/common"
	"github.com/cryingmouse/data_management_engine/mgmtmodel"

	"github.com/gin-gonic/gin"
)

type LocalUserResponse struct {
	Name string `json:"name,omitempty"`
}

type PaginationLocalUserResponse struct {
	Users      []LocalUserResponse `json:"users"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalCount int64               `json:"total_count"`
}

func createLocalUserHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	request := struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,validatePassword"`
		HostName string `json:"host_name" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var userModel mgmtmodel.LocalUser

	common.CopyStructList(request, &userModel)

	if err := userModel.Create(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create the users.",
			"error":   err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func createLocalUsersHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	request := []struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,validatePassword"`
		HostName string `json:"host_name" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var userListModel mgmtmodel.LocalUserList

	common.CopyStructList(request, &userListModel.Users)

	if err := userListModel.Create(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create the users.",
			"error":   err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func deleteLocalUserHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	request := struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,validatePassword"`
	}{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var userModel mgmtmodel.LocalUser

	common.CopyStructList(request, &userModel)

	if err := userModel.Delete(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to delete the user '%s' on host '%s'.", request.Name, request.Password),
			"error":   err.Error(),
		})
	}

	c.Status(http.StatusOK)
}

func deleteLocalUsersHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	request := []struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,validatePassword"`
	}{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var userListModel mgmtmodel.LocalUserList

	common.CopyStructList(request, &userListModel.Users)

	if err := userListModel.Delete(ctx, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete the users.",
			"error":   err.Error(),
		})
	}

	c.Status(http.StatusOK)
}

func getlocalUsersHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	userName := c.Query("name")
	// isLockout := c.Query("is_lockout")
	// computerName := c.Query("host_name")
	fields := c.Query("fields")

	page, err_page := strconv.Atoi(c.Query("page"))
	limit, err_limit := strconv.Atoi(c.Query("limit"))
	if (err_page != nil && err_limit == nil) || (err_page == nil && err_limit != nil) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request."})
		return
	}

	if userName == "" {
		userListModel := mgmtmodel.LocalUserList{}

		filter := common.QueryFilter{
			Fields: common.SplitToList(fields),
			// Conditions: struct {
			// 	Password string
			// }{
			// 	Password: hostIp,
			// },
		}

		if page == 0 && limit == 0 {
			// Query users without pagination.

			users, err := userListModel.Get(ctx, &filter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Failed to get the users",
					"error":   err.Error(),
				})
				return
			}

			userInfoList := []LocalUserResponse{}
			for _, user := range users {
				userInfoList = append(userInfoList, LocalUserResponse{
					Name: user.Name,
				})
			}

			c.JSON(http.StatusOK, gin.H{"message": "Get the users successfully.", "users": userInfoList})
			return
		} else {
			// Query users with pagination.
			filter.Pagination = &common.Pagination{
				Page:     page,
				PageSize: limit,
			}

			paginationDirs, err := userListModel.Pagination(ctx, &filter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Failed to get the users.",
					"error":   err.Error(),
				})
				return
			}

			paginationDirList := PaginationLocalUserResponse{
				Page:       page,
				Limit:      limit,
				TotalCount: paginationDirs.TotalCount,
			}

			for _, _user := range paginationDirs.Users {
				user := LocalUserResponse{
					Name: _user.Name,
				}

				paginationDirList.Users = append(paginationDirList.Users, user)
			}

			c.JSON(http.StatusOK, gin.H{"message": "Get the users successfully.", "pagination": paginationDirList})
			return
		}
	} else {
		userModel := mgmtmodel.LocalUser{
			Name: userName,
			// Password: hostIp,
		}

		user, err := userModel.Get(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to get the users",
				"error":   err.Error(),
			})
			return
		}

		// Convert to UserInfo as REST API response.
		userInfo := LocalUserResponse{
			Name: user.Name,
		}

		c.JSON(http.StatusOK, gin.H{"message": "Get the user successfully.", "user": userInfo})
	}
}

func createLocalUserOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	request := struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,validatePassword"`
	}{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hostContext := common.HostContext{
		Username: c.Request.Header.Get("X-agent-username"),
		Password: c.Request.Header.Get("X-agent-password"),
	}

	agent := agent.GetAgent()
	agent.CreateLocalUser(ctx, hostContext, request.Name, request.Password)

	c.JSON(http.StatusOK, gin.H{"message": "Create user on agent successfully."})
}

func deleteLocalUserOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	request := struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,validatePassword"`
	}{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hostContext := common.HostContext{
		Username: c.Request.Header.Get("X-agent-username"),
		Password: c.Request.Header.Get("X-agent-password"),
	}

	agent := agent.GetAgent()
	agent.DeleteLocalUser(ctx, hostContext, request.Name)

	c.JSON(http.StatusOK, gin.H{"message": "Delete user on agent successfully.", "user": request.Name})
}

func getLocalUserOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	username := c.Query("name")

	hostContext := common.HostContext{
		Username: c.Request.Header.Get("X-agent-username"),
		Password: c.Request.Header.Get("X-agent-password"),
	}

	agent := agent.GetAgent()

	if username != "" {
		user, _ := agent.GetLocalUser(ctx, hostContext, username)
		c.JSON(http.StatusOK, gin.H{"message": "Get the user on agent successfully.", "user": user})

	} else {
		users, _ := agent.GetLocalUsers(ctx, hostContext)
		c.JSON(http.StatusOK, gin.H{"message": "Get the users on agent successfully.", "users": users})
	}
}
