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

type DirectoryInfo struct {
	Name   string `json:"name" binding:"required"`
	HostIP string `json:"host_ip" binding:"required"`
}

type PaginationDirectoryInfo struct {
	Directories []DirectoryInfo `json:"directories"`
	Page        int             `json:"page"`
	Limit       int             `json:"limit"`
	TotalCount  int64           `json:"total_count"`
}

func createDirectoryHandler(c *gin.Context) {
	var directoryInfo DirectoryInfo
	if err := c.ShouldBindJSON(&directoryInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	directoryModel := mgmtmodel.Directory{
		Name:   directoryInfo.Name,
		HostIP: directoryInfo.HostIP,
	}

	if err := directoryModel.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to create the directory with the parameters: host_ip=%s,name=%s", directoryInfo.HostIP, directoryInfo.Name),
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Create the directory '%s' on host '%s' successfully.", directoryInfo.Name, directoryInfo.HostIP)})
}

func deleteDirectoryHandler(c *gin.Context) {
	var directoryInfo DirectoryInfo
	if err := c.ShouldBindJSON(&directoryInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	directoryModel := mgmtmodel.Directory{
		Name:   directoryInfo.Name,
		HostIP: directoryInfo.HostIP,
	}

	if err := directoryModel.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to delete the directory '%s' on host '%s'.", directoryInfo.Name, directoryInfo.HostIP),
			"error":   err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Delete the directory '%s' on host '%s' successfully.", directoryInfo.Name, directoryInfo.HostIP)})
}

func getDirectoryHandler(c *gin.Context) {
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
		directoryListModel := mgmtmodel.DirectoryList{}
		if page == 0 && limit == 0 {
			// Query directories without pagination.
			filter := context.QueryFilter{
				Fields: utils.SplitToList(fields),
				Keyword: map[string]string{
					"name": nameKeyword,
				},
				Conditions: struct {
					HostIP string
				}{
					HostIP: hostIp,
				},
			}
			directories, err := directoryListModel.Get(&filter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("Failed to get the directories with the parameters: host_ip=%s", hostIp),
					"error":   err.Error(),
				})
				return
			}

			directoryInfoList := []DirectoryInfo{}
			for _, directory := range directories {
				directoryInfoList = append(directoryInfoList, DirectoryInfo{
					Name:   directory.Name,
					HostIP: directory.HostIP,
				})
			}

			c.JSON(http.StatusOK, gin.H{"message": "Get the directories successfully.", "directories": directoryInfoList})
			return

		} else {
			// Query directories with pagination.
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
					HostIP string
				}{
					HostIP: hostIp,
				},
			}
			paginationDirs, err := directoryListModel.Pagination(&filter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("Failed to get the directories with the parameters: host_ip=%s,page=%d,limit=%d", hostIp, page, limit),
					"error":   err.Error(),
				})
				return
			}

			paginationDirList := PaginationDirectoryInfo{
				Page:       page,
				Limit:      limit,
				TotalCount: paginationDirs.TotalCount,
			}

			for _, _directory := range paginationDirs.Directories {
				directory := DirectoryInfo{
					Name:   _directory.Name,
					HostIP: _directory.HostIP,
				}

				paginationDirList.Directories = append(paginationDirList.Directories, directory)
			}

			c.JSON(http.StatusOK, gin.H{"message": "Get the directories successfully.", "pagination": paginationDirList})
			return
		}
	} else {
		directoryModel := mgmtmodel.Directory{
			Name:   dirName,
			HostIP: hostIp,
		}

		directory, err := directoryModel.Get()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Failed to get the directories with parameters: name=%s,host_ip=%s", dirName, hostIp),
				"error":   err.Error(),
			})
			return
		}

		// Convert to DirectoryInfo as REST API response.
		directoryInfo := DirectoryInfo{
			Name:   directory.Name,
			HostIP: directory.HostIP,
		}

		c.JSON(http.StatusOK, gin.H{"message": "Get the directory successfully.", "directory": directoryInfo})
	}
}

type AgentDirectoryInfo struct {
	Name string `json:"name"`
}

func createDirectoryOnAgentHandler(c *gin.Context) {
	var directoryInfo AgentDirectoryInfo
	if err := c.ShouldBindJSON(&directoryInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hostContext := context.HostContext{
		Username: c.Request.Header.Get("X-agent-username"),
		Password: c.Request.Header.Get("X-agent-password"),
	}

	agent := agent.GetAgent()
	dirPath, _ := agent.CreateDirectory(hostContext, directoryInfo.Name)

	c.JSON(http.StatusOK, gin.H{"message": "Create directory on agent successfully.", "directory": dirPath})
}

func deleteDirectoryOnAgentHandler(c *gin.Context) {
	var directoryInfo AgentDirectoryInfo
	if err := c.ShouldBindJSON(&directoryInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hostContext := context.HostContext{
		Username: c.Request.Header.Get("X-agent-username"),
		Password: c.Request.Header.Get("X-agent-password"),
	}

	agent := agent.GetAgent()
	agent.DeleteDirectory(hostContext, directoryInfo.Name)

	c.JSON(http.StatusOK, gin.H{"message": "Delete directory on agent successfully.", "directory": directoryInfo.Name})
}
