package webservice

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cryingmouse/data_management_engine/agent"
	"github.com/cryingmouse/data_management_engine/context"
	"github.com/cryingmouse/data_management_engine/mgmtmodel"
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
			"Message": fmt.Sprintf("Failed to create the directory with the parameters: host_ip=%s,name=%s", directoryInfo.HostIP, directoryInfo.Name),
			"Error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": fmt.Sprintf("Create the directory '%s' on host '%s' successfully.", directoryInfo.Name, directoryInfo.HostIP)})
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
			"Message": fmt.Sprintf("Failed to delete the directory '%s' on host '%s'.", directoryInfo.Name, directoryInfo.HostIP),
			"Error":   err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{"Message": fmt.Sprintf("Delete the directory '%s' on host '%s' successfully.", directoryInfo.Name, directoryInfo.HostIP)})
}

// GetUsers returns a list of users
// swagger:route GET /users users listUsers
// Returns a list of users
// responses:
//
//	200: []userResponse
func getDirectoryHandler(c *gin.Context) {
	dirName := c.Query("name")
	hostIp := c.Query("host_ip")
	fields := c.Query("fields")
	nameKeyword := c.Query("name_key_word")

	page, limit, err := validatePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid request.", "Error": err.Error()})
		return
	}

	if dirName == "" {
		directoryListModel := mgmtmodel.DirectoryList{}
		if page == 0 && limit == 0 {
			// Query directories without pagination.
			filter := context.QueryFilter{
				Fields: strings.Split(fields, ","),
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
					"Message": fmt.Sprintf("Failed to get the directories with the parameters: host_ip=%s", hostIp),
					"Error":   err.Error(),
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

			c.JSON(http.StatusOK, gin.H{"Message": "Get the directories successfully.", "Directories": directoryInfoList})
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
			paginationDirs, err := directoryListModel.GetByPagination(&filter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"Message": fmt.Sprintf("Failed to get the directories with the parameters: host_ip=%s,page=%d,limit=%d", hostIp, page, limit),
					"Error":   err.Error(),
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

			c.JSON(http.StatusOK, gin.H{"Message": "Get the directories successfully.", "PaginationDirectories": paginationDirList})
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
				"Message": fmt.Sprintf("Failed to get the directories with parameters: name=%s,host_ip=%s", dirName, hostIp),
				"Error":   err.Error(),
			})
			return
		}

		// Convert to DirectoryInfo as REST API response.
		directoryInfo := DirectoryInfo{
			Name:   directory.Name,
			HostIP: directory.HostIP,
		}

		c.JSON(http.StatusOK, gin.H{"Message": "Get the registered host successfully.", "RegisteredHosts": directoryInfo})
	}
}

func searchDirectoryHandler(c *gin.Context) {
	nameKeyword := c.Query("name")
	fields := c.Query("fields")

	page, limit, err := validatePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid request.", "Error": err.Error()})
		return
	}

	directoryListModel := mgmtmodel.DirectoryList{}
	if page == 0 && limit == 0 {
		// Query directories without pagination.
		filter := context.QueryFilter{
			Fields: strings.Split(fields, ","),
			Keyword: map[string]string{
				"name": nameKeyword,
			},
		}
		directories, err := directoryListModel.Get(&filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Message": fmt.Sprintf("Failed to search the directories with the keyword: name=%s", nameKeyword),
				"Error":   err.Error(),
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

		c.JSON(http.StatusOK, gin.H{"Message": "Get the directories successfully.", "Directories": directoryInfoList})
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
		}
		paginationDirs, err := directoryListModel.GetByPagination(&filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Message": fmt.Sprintf("Failed to search the directories with the parameters: page=%d,limit=%d", page, limit),
				"Error":   err.Error(),
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

		c.JSON(http.StatusOK, gin.H{"Message": "Get the directories successfully.", "PaginationDirectories": paginationDirList})
		return
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

	c.JSON(http.StatusOK, gin.H{"Message": "Create directory on agent successfully.", "Directory": dirPath})
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
	dirPath, _ := agent.DeleteDirectory(hostContext, directoryInfo.Name)

	c.JSON(http.StatusOK, gin.H{"Message": "Delete directory on agent successfully.", "Directory": dirPath})
}
