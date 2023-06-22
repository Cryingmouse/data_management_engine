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
	// HostName might not be necessary here, but this is for hostModel.Get()
	HostName string `json:"host_name"`
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

	_, err := directoryModel.Create()
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Create the directory successfully."})
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
