package webservice

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cryingmouse/data_management_engine/agent"
	"github.com/cryingmouse/data_management_engine/common"
	"github.com/cryingmouse/data_management_engine/mgmtmodel"

	"github.com/gin-gonic/gin"
)

type LocalUserResponse struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type PaginationLocalUserResponse struct {
	Users      []LocalUserResponse `json:"users"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalCount int64               `json:"total_count"`
}

func createUserHandler(c *gin.Context) {
	request := []struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
		HostName string `json:"host_name" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var userListModel mgmtmodel.LocalUserList

	common.CopyStructList(request, &userListModel.Users)

	if err := userListModel.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create the users.",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Create the users successfully."})
}

func deleteUserHandler(c *gin.Context) {
	type Request struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userModel := mgmtmodel.LocalUser{
		Name:     request.Name,
		Password: request.Password,
	}

	if err := userModel.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to delete the user '%s' on host '%s'.", request.Name, request.Password),
			"error":   err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Delete the user '%s' on host '%s' successfully.", request.Name, request.Password)})
}

func getUserHandler(c *gin.Context) {
	userName := c.Query("name")
	// isLockout := c.Query("is_lockout")
	// computerName := c.Query("host_name")
	fields := c.Query("fields")

	page, limit, err := validatePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request.", "error": err.Error()})
		return
	}

	if userName == "" {
		userListModel := mgmtmodel.LocalUserList{}
		if page == 0 && limit == 0 {
			// Query users without pagination.
			filter := common.QueryFilter{
				Fields: common.SplitToList(fields),
				// Conditions: struct {
				// 	Password string
				// }{
				// 	Password: hostIp,
				// },
			}
			users, err := userListModel.Get(&filter)
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
					Name:     user.Name,
					Password: user.Password,
				})
			}

			c.JSON(http.StatusOK, gin.H{"message": "Get the users successfully.", "users": userInfoList})
			return
		} else {
			// Query users with pagination.
			filter := common.QueryFilter{
				Fields: strings.Split(fields, ","),
				Pagination: &common.Pagination{
					Page:     page,
					PageSize: limit,
				},
				// Conditions: struct {
				// 	Password string
				// }{
				// 	Password: hostIp,
				// },
			}
			paginationDirs, err := userListModel.Pagination(&filter)
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
					Name:     _user.Name,
					Password: _user.Password,
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

		user, err := userModel.Get()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to get the users",
				"error":   err.Error(),
			})
			return
		}

		// Convert to UserInfo as REST API response.
		userInfo := LocalUserResponse{
			Name:     user.Name,
			Password: user.Password,
		}

		c.JSON(http.StatusOK, gin.H{"message": "Get the user successfully.", "user": userInfo})
	}
}

func createLocalUserOnAgentHandler(c *gin.Context) {
	type Request struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hostContext := common.HostContext{
		Username: c.Request.Header.Get("X-agent-username"),
		Password: c.Request.Header.Get("X-agent-password"),
	}

	agent := agent.GetAgent()
	agent.CreateLocalUser(hostContext, request.Name, request.Password)

	c.JSON(http.StatusOK, gin.H{"message": "Create user on agent successfully."})
}

func deleteLocalUserOnAgentHandler(c *gin.Context) {
	type Request struct {
		Name string `json:"name"`
	}

	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hostContext := common.HostContext{
		Username: c.Request.Header.Get("X-agent-username"),
		Password: c.Request.Header.Get("X-agent-password"),
	}

	agent := agent.GetAgent()
	agent.DeleteLocalUser(hostContext, request.Name)

	c.JSON(http.StatusOK, gin.H{"message": "Delete user on agent successfully.", "user": request.Name})
}

func getLocalUserOnAgentHandler(c *gin.Context) {
	username := c.Query("name")

	hostContext := common.HostContext{
		Username: c.Request.Header.Get("X-agent-username"),
		Password: c.Request.Header.Get("X-agent-password"),
	}

	agent := agent.GetAgent()

	if username != "" {
		user, _ := agent.GetLocalUser(hostContext, username)
		c.JSON(http.StatusOK, gin.H{"message": "Get the user on agent successfully.", "user": user})

	} else {
		users, _ := agent.GetLocalUsers(hostContext)
		c.JSON(http.StatusOK, gin.H{"message": "Get the users on agent successfully.", "users": users})
	}
}
