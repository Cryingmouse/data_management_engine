package webservice

import (
	"net/http"

	"github.com/cryingmouse/data_management_engine/agent"
	"github.com/cryingmouse/data_management_engine/common"

	"github.com/gin-gonic/gin"
)

func createShareOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request struct {
		ShareName     string   `json:"share_name" binding:"required"`
		DirectoryName string   `json:"directory_name" binding:"required"`
		Description   string   `json:"description" binding:"required"`
		UserNames     []string `json:"user_names" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	agent := agent.GetAgent()

	err := agent.CreateCIFSShare(ctx, request.ShareName, request.DirectoryName, request.Description, request.UserNames)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create the local user", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func deleteShareOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	var request struct {
		ShareName string `json:"share_name"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	agent := agent.GetAgent()
	if err := agent.DeleteCIFSShare(ctx, request.ShareName); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete the local user", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func getShareOnAgentHandler(c *gin.Context) {
	ctx := SetTraceIDInContext(c)

	name := c.Query("name")
	names := common.SplitToList(name)

	agent := agent.GetAgent()
	if len(names) == 0 {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", "")
	} else if len(names) == 1 {
		ShareDetail, err := agent.GetCIFSShareDetail(ctx, name)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "Failed to get the local user detail", err.Error())
			return
		}

		c.JSON(http.StatusOK, ShareDetail)
	} else {
		SharesDetail, err := agent.GetCIFSSharesDetail(ctx, names)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "Failed to get the local users detail", err.Error())
			return
		}

		c.JSON(http.StatusOK, SharesDetail)
	}
}
