package webservice

import (
	"net/http"

	"github.com/cryingmouse/data_management_engine/agent"
	"github.com/cryingmouse/data_management_engine/context"
	"github.com/gin-gonic/gin"
)

func getSystemInfoOnAgentHandler(c *gin.Context) {
	hostContext := context.HostContext{
		Username: c.Request.Header.Get("X-agent-username"),
		Password: c.Request.Header.Get("X-agent-password"),
	}

	agent := agent.GetAgent()

	if systemInfo, err := agent.GetSystemInfo(hostContext); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get system info on agent.", "error": err})

	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Get system info on agent successfully.", "system-info": systemInfo})
	}

}
