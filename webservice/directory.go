package webservice

import (
	"net/http"

	"github.com/cryingmouse/data_management_engine/agent"
	"github.com/cryingmouse/data_management_engine/context"
	"github.com/cryingmouse/data_management_engine/mgmtmodel"
	"github.com/gin-gonic/gin"
)

type DirectoryInfo struct {
	Name   string `json:"name" binding:"required"`
	HostIp string `json:"host_ip" binding:"required"`
}

func createDirectoryHandler(c *gin.Context) {
	var directoryInfo DirectoryInfo
	if err := c.ShouldBindJSON(&directoryInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	directoryModel := mgmtmodel.Directory{
		Name:   directoryInfo.Name,
		HostIp: directoryInfo.HostIp,
	}

	if err := directoryModel.Create(); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Create the directory successfully."})
}

func deleteDirectoryHandler(c *gin.Context) {
	var directoryInfo DirectoryInfo
	if err := c.ShouldBindJSON(&directoryInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	directoryModel := mgmtmodel.Directory{
		Name:   directoryInfo.Name,
		HostIp: directoryInfo.HostIp,
	}

	_, err := directoryModel.Delete()
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Create the directory successfully."})
}

func getDirectoryHandler(c *gin.Context) {
	dirName := c.Query("name")
	hostIp := c.Query("ip")

	if dirName == "" && hostIp == "" {
		directoryListModel := mgmtmodel.DirectoryList{}
		directories, err := directoryListModel.Get(hostIp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Message": "Failed to get the directories.", "Error": err})
		}

		c.JSON(http.StatusOK, gin.H{"Message": "Get the directories successfully.", "Directories": directories})
		return

	} else {
		directoryModel := mgmtmodel.Directory{
			Name:   dirName,
			HostIp: hostIp,
		}

		directory, err := directoryModel.Get()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Message": "Failed to get the directory."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"Message": "Get the registered host successfully.", "RegisteredHosts": directory})

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

	c.JSON(http.StatusOK, gin.H{"Message": "Create directory on agent successfully.", "DirectoryPath": dirPath})
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

	c.JSON(http.StatusOK, gin.H{"Message": "Delete directory on agent successfully.", "DirectoryPath": dirPath})
}
