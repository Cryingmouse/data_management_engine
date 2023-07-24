package webservice

import (
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

type requestDirectory struct {
	Name   string `json:"name" binding:"required"`
	HostIP string `json:"host_ip" binding:"required,ip"`
}

func CreateDirectoryHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request requestDirectory
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	directoryModel := mgmtmodel.Directory{}
	common.DeepCopy(request, &directoryModel)

	if err := directoryModel.Create(ctx); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create the directory", err.Error())
		return
	}

	directoryResponse := DirectoryResponse{}
	common.DeepCopy(directoryModel, &directoryResponse)

	c.JSON(http.StatusOK, directoryResponse)
}

func CreateDirectoriesHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request []requestDirectory
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	directoryListModel := mgmtmodel.DirectoryList{}
	common.DeepCopy(request, &directoryListModel.Directories)

	if err := directoryListModel.Create(ctx); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create the directories", err.Error())
		return
	}

	directoryResponseList := make([]DirectoryResponse, len(directoryListModel.Directories))
	common.DeepCopy(directoryListModel.Directories, &directoryResponseList)

	c.JSON(http.StatusOK, directoryResponseList)
}

func DeleteDirectoryHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request requestDirectory
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	directoryModel := mgmtmodel.Directory{}
	common.DeepCopy(request, &directoryModel)

	if err := directoryModel.Delete(ctx); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete the directory", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func DeleteDirectoriesHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request []requestDirectory
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	directoryListModel := mgmtmodel.DirectoryList{}
	common.DeepCopy(request, &directoryListModel.Directories)

	if err := directoryListModel.Delete(ctx, nil); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete the directories", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func GetDirectoriesHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	dirName := c.Query("name")
	hostIP := c.Query("host_ip")
	fields := c.Query("fields")
	nameKeyword := c.Query("q")
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	if hostIP != "" && validateIPAddress(hostIP) != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", "")
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
			directories, err := directoryListModel.Get(ctx, &filter)
			if err != nil {
				ErrorResponse(c, http.StatusInternalServerError, "Failed to get the directories", err.Error())
				return
			}

			directoryList := make([]DirectoryResponse, len(directories))
			common.DeepCopy(directories, &directoryList)

			c.JSON(http.StatusOK, directoryList)
		} else {
			// Query directories with pagination.
			filter.Pagination = &common.Pagination{
				Page:     page,
				PageSize: limit,
			}

			paginationDirs, err := directoryListModel.Pagination(ctx, &filter)
			if err != nil {
				ErrorResponse(c, http.StatusInternalServerError, "Failed to get the directories", err.Error())
				return
			}

			paginationDirList := PaginationDirectoryResponse{
				Page:       page,
				Limit:      limit,
				TotalCount: paginationDirs.TotalCount,
			}

			common.DeepCopy(paginationDirs.Directories, &paginationDirList.Directories)

			c.JSON(http.StatusOK, paginationDirList)
		}
	} else {
		directoryModel := mgmtmodel.Directory{
			Name:   dirName,
			HostIP: hostIP,
		}

		directory, err := directoryModel.Get(ctx)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "Failed to get the directory", err.Error())
			return
		}

		directoryInfo := DirectoryResponse{}
		common.DeepCopy(directory, &directoryInfo)

		directoryInfoList := []DirectoryResponse{directoryInfo}

		c.JSON(http.StatusOK, directoryInfoList)
	}
}

func CreateDirectoryOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	agent := agent.GetAgent()
	_, err := agent.CreateDirectory(ctx, request.Name)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create the directory", err.Error())
		return
	}

	directoryDetails, err := agent.GetDirectoryDetail(ctx, request.Name)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to get the directory details", err.Error())
		return
	}

	c.JSON(http.StatusOK, directoryDetails)
}

func DeleteDirectoryOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	agent := agent.GetAgent()
	if err := agent.DeleteDirectory(ctx, request.Name); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete the directory", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func CreateDirectoriesOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request []struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	names := make([]string, len(request))
	for i, item := range request {
		names[i] = item.Name
	}

	agent := agent.GetAgent()
	_, err := agent.CreateDirectories(ctx, names)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create the directories", err.Error())
		return
	}

	DirectoryDetails, err := agent.GetDirectoriesDetail(ctx, names)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to get the directories details", err.Error())
		return
	}

	c.JSON(http.StatusOK, DirectoryDetails)
}

func DeleteDirectoriesOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request []struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	names := make([]string, len(request))
	for i, item := range request {
		names[i] = item.Name
	}

	agent := agent.GetAgent()
	if err := agent.DeleteDirectories(ctx, names); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete the directories", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func GetDirectoryDetailOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	name := c.Query("name")
	names := common.SplitToList(name)

	agent := agent.GetAgent()
	if len(names) == 0 {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", "")
	} else if len(names) == 1 {
		directoryDetail, err := agent.GetDirectoryDetail(ctx, name)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "Failed to get the directory details", err.Error())
			return
		}

		c.JSON(http.StatusOK, directoryDetail)
	} else {
		directoriesDetail, err := agent.GetDirectoriesDetail(ctx, names)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "Failed to get the directories details", err.Error())
			return
		}

		c.JSON(http.StatusOK, directoriesDetail)
	}
}
