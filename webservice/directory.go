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

type DirectoryResponse struct {
	Name           string `json:"name,omitempty"`
	HostIP         string `json:"host_ip,omitempty"`
	CreationTime   string `json:"creation_time,omitempty"`
	LastAccessTime string `json:"last_access_time,omitempty"`
	LastWriteTime  string `json:"last_write_time,omitempty"`
	Exist          bool   `json:"exist,omitempty"`
	FullPath       string `json:"full_path,omitempty"`
	ParentFullPath string `json:"parent_full_path,omitempty"`
}

type PaginationDirectoryResponse struct {
	Directories []DirectoryResponse `json:"directories"`
	Page        int                 `json:"page"`
	Limit       int                 `json:"limit"`
	TotalCount  int64               `json:"total_count"`
}

func createDirectoryHandler(c *gin.Context) {
	request := struct {
		Name   string `json:"name" binding:"required"`
		HostIP string `json:"host_ip" binding:"required,ip"`
	}{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var directoryModel mgmtmodel.Directory
	common.CopyStructList(request, &directoryModel)
	if err := directoryModel.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create the directories.",
			"error":   err.Error(),
		})
		return
	}

	var response DirectoryResponse
	common.CopyStructList(directoryModel, &response)

	c.JSON(http.StatusOK, response)
}

func createDirectoriesHandler(c *gin.Context) {
	request := []struct {
		Name   string `json:"name" binding:"required"`
		HostIP string `json:"host_ip" binding:"required,ip"`
	}{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var directoryListModel mgmtmodel.DirectoryList
	common.CopyStructList(request, &directoryListModel.Directories)
	if err := directoryListModel.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create the directories.",
			"error":   err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func deleteDirectoryHandler(c *gin.Context) {
	request := struct {
		Name   string `json:"name" binding:"required"`
		HostIP string `json:"host_ip" binding:"required,ip"`
	}{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var directoryModel mgmtmodel.Directory
	common.CopyStructList(request, &directoryModel)
	if err := directoryModel.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete the directories.",
			"error":   err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func deleteDirectoriesHandler(c *gin.Context) {
	request := []struct {
		Name   string `json:"name" binding:"required"`
		HostIP string `json:"host_ip" binding:"required,ip"`
	}{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var directoryListModel mgmtmodel.DirectoryList
	common.CopyStructList(request, &directoryListModel.Directories)
	if err := directoryListModel.Delete(nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create the directories.",
			"error":   err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func getDirectoriesHandler(c *gin.Context) {
	dirName := c.Query("name")
	hostIP := c.Query("host_ip")
	fields := c.Query("fields")
	nameKeyword := c.Query("q")
	page, err_page := strconv.Atoi(c.Query("page"))
	limit, err_limit := strconv.Atoi(c.Query("limit"))

	if (err_page != nil && err_limit == nil) || (err_page == nil && err_limit != nil) || (hostIP != "" && validateIPAddress(hostIP) != nil) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request."})
		return
	}

	if dirName == "" || hostIP == "" {
		directoryListModel := mgmtmodel.DirectoryList{}

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
				Name:   dirName,
			},
		}

		if page == 0 && limit == 0 {
			// Query directories without pagination.
			directories, err := directoryListModel.Get(&filter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Failed to get the directories.",
					"error":   err.Error(),
				})
				return
			}

			var directoryList []DirectoryResponse

			common.CopyStructList(directories, &directoryList)

			c.JSON(http.StatusOK, directoryList)
		} else {
			// Query directories with pagination.
			filter.Pagination = &common.Pagination{
				Page:     page,
				PageSize: limit,
			}

			paginationDirs, err := directoryListModel.Pagination(&filter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Failed to get the directories.",
					"error":   err.Error(),
				})
				return
			}

			paginationDirList := PaginationDirectoryResponse{
				Page:       page,
				Limit:      limit,
				TotalCount: paginationDirs.TotalCount,
			}

			common.CopyStructList(paginationDirs.Directories, &paginationDirList.Directories)

			c.JSON(http.StatusOK, paginationDirList)
		}
	} else {
		directoryModel := mgmtmodel.Directory{
			Name:   dirName,
			HostIP: hostIP,
		}

		directory, err := directoryModel.Get()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to get the directories.",
				"error":   err.Error(),
			})
			return
		}

		// Convert to DirectoryResponse as REST API response.
		var directoryInfo DirectoryResponse

		common.CopyStructList(directory, &directoryInfo)

		directoryInfoList := []DirectoryResponse{}
		directoryInfoList = append(directoryInfoList, directoryInfo)

		c.JSON(http.StatusOK, directoryInfoList)
	}
}

func createDirectoryOnAgentHandler(c *gin.Context) {
	request := struct {
		Name string `json:"name"`
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
	dirPath, err := agent.CreateDirectory(hostContext, request.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	DirectoryDetails, err := agent.GetDirectoryDetail(hostContext, dirPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, DirectoryDetails)
}

func deleteDirectoryOnAgentHandler(c *gin.Context) {
	request := struct {
		Name string `json:"name"`
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
	if err := agent.DeleteDirectory(hostContext, request.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func createDirectoriesOnAgentHandler(c *gin.Context) {
	request := []struct {
		Name string `json:"name"`
	}{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	names := []string{}
	for _, item := range request {
		names = append(names, item.Name)
	}

	hostContext := common.HostContext{
		Username: c.Request.Header.Get("X-agent-username"),
		Password: c.Request.Header.Get("X-agent-password"),
	}

	agent := agent.GetAgent()
	dirPaths, err := agent.CreateDirectories(hostContext, names)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	DirectoryDetails, err := agent.GetDirectoriesDetail(hostContext, dirPaths)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, DirectoryDetails)
}

func deleteDirectoriesOnAgentHandler(c *gin.Context) {
	request := struct {
		Name string `json:"name"`
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
	if err := agent.DeleteDirectory(hostContext, request.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func getDirectoryDetailsOnAgentHandler(c *gin.Context) {
	name := c.Query("name")
	names := common.SplitToList(name)

	hostContext := common.HostContext{
		Username: c.Request.Header.Get("X-agent-username"),
		Password: c.Request.Header.Get("X-agent-password"),
	}

	agent := agent.GetAgent()
	if len(names) == 0 {
		c.JSON(http.StatusInternalServerError, fmt.Errorf(".nvalid request"))
	} else if len(names) == 1 {
		DirectoryDetail, err := agent.GetDirectoryDetail(hostContext, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, DirectoryDetail)
	} else {
		DirectoriesDetail, err := agent.GetDirectoriesDetail(hostContext, names)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, DirectoriesDetail)
	}
}
