package webservice

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cryingmouse/data_management_engine/agent"
	"github.com/cryingmouse/data_management_engine/context"
	"github.com/cryingmouse/data_management_engine/mgmtmodel"
	"github.com/cryingmouse/data_management_engine/utils"
	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type PaginationUserInfo struct {
	Users      []UserInfo `json:"users"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	TotalCount int64      `json:"total_count"`
}

func createUserHandler(c *gin.Context) {
	var userInfo UserInfo
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userModel := mgmtmodel.User{
		Name:     userInfo.Name,
		Password: userInfo.Password,
	}

	if err := userModel.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to create the user with the parameters: host_ip=%s,name=%s", userInfo.Password, userInfo.Name),
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Create the user '%s' on host '%s' successfully.", userInfo.Name, userInfo.Password)})
}

func deleteUserHandler(c *gin.Context) {
	var userInfo UserInfo
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userModel := mgmtmodel.User{
		Name:     userInfo.Name,
		Password: userInfo.Password,
	}

	if err := userModel.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to delete the user '%s' on host '%s'.", userInfo.Name, userInfo.Password),
			"error":   err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Delete the user '%s' on host '%s' successfully.", userInfo.Name, userInfo.Password)})
}

func getUserHandler(c *gin.Context) {
	dirName := c.Query("name")
	hostIp := c.Query("host_ip")
	fields := c.Query("fields")
	nameKeyword := c.Query("q")

	page, limit, err := validatePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request.", "error": err.Error()})
		return
	}

	if dirName == "" {
		userListModel := mgmtmodel.UserList{}
		if page == 0 && limit == 0 {
			// Query users without pagination.
			filter := context.QueryFilter{
				Fields: utils.SplitToList(fields),
				Keyword: map[string]string{
					"name": nameKeyword,
				},
				Conditions: struct {
					Password string
				}{
					Password: hostIp,
				},
			}
			users, err := userListModel.Get(&filter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("Failed to get the users with the parameters: host_ip=%s", hostIp),
					"error":   err.Error(),
				})
				return
			}

			userInfoList := []UserInfo{}
			for _, user := range users {
				userInfoList = append(userInfoList, UserInfo{
					Name:     user.Name,
					Password: user.Password,
				})
			}

			c.JSON(http.StatusOK, gin.H{"message": "Get the users successfully.", "users": userInfoList})
			return

		} else {
			// Query users with pagination.
			filter := context.QueryFilter{
				Fields: strings.Split(fields, ","),
				Keyword: map[string]string{
					"name": nameKeyword,
				},
				Pagination: &context.Pagination{
					Page:     page,
					PageSize: limit,
				},
				Conditions: struct {
					Password string
				}{
					Password: hostIp,
				},
			}
			paginationDirs, err := userListModel.Pagination(&filter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("Failed to get the users with the parameters: host_ip=%s,page=%d,limit=%d", hostIp, page, limit),
					"error":   err.Error(),
				})
				return
			}

			paginationDirList := PaginationUserInfo{
				Page:       page,
				Limit:      limit,
				TotalCount: paginationDirs.TotalCount,
			}

			for _, _user := range paginationDirs.Users {
				user := UserInfo{
					Name:     _user.Name,
					Password: _user.Password,
				}

				paginationDirList.Users = append(paginationDirList.Users, user)
			}

			c.JSON(http.StatusOK, gin.H{"message": "Get the users successfully.", "pagination": paginationDirList})
			return
		}
	} else {
		userModel := mgmtmodel.User{
			Name:     dirName,
			Password: hostIp,
		}

		user, err := userModel.Get()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Failed to get the users with parameters: name=%s,host_ip=%s", dirName, hostIp),
				"error":   err.Error(),
			})
			return
		}

		// Convert to UserInfo as REST API response.
		userInfo := UserInfo{
			Name:     user.Name,
			Password: user.Password,
		}

		c.JSON(http.StatusOK, gin.H{"message": "Get the user successfully.", "user": userInfo})
	}
}

func createUserOnAgentHandler(c *gin.Context) {
	type AgentUserInfo struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	var userInfo AgentUserInfo
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hostContext := context.HostContext{
		Username: c.Request.Header.Get("X-agent-username"),
		Password: c.Request.Header.Get("X-agent-password"),
	}

	agent := agent.GetAgent()
	agent.CreateLocalUser(hostContext, userInfo.Name, userInfo.Password)

	c.JSON(http.StatusOK, gin.H{"message": "Create user on agent successfully."})
}

func deleteUserOnAgentHandler(c *gin.Context) {
	type AgentUserInfo struct {
		Name string `json:"name"`
	}

	var userInfo AgentUserInfo
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hostContext := context.HostContext{
		Username: c.Request.Header.Get("X-agent-username"),
		Password: c.Request.Header.Get("X-agent-password"),
	}

	agent := agent.GetAgent()
	agent.DeleteLocalUser(hostContext, userInfo.Name)

	c.JSON(http.StatusOK, gin.H{"message": "Delete user on agent successfully.", "user": userInfo.Name})
}

func getUserOnAgentHandler(c *gin.Context) {
	username := c.Query("name")

	hostContext := context.HostContext{
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
