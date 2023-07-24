package webservice

import (
	"context"
	"fmt"

	"github.com/cryingmouse/data_management_engine/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func Start() {
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(TraceMiddleware(), LoggingMiddleware())

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validatePassword", PasswordValidator)
		v.RegisterValidation("validateStorageType", StorageTypeValidator)
	}

	// Router 'portal' for Portal
	portal := router.Group("/api")

	// Router 'agent' for Agent
	agent := router.Group("/agent")

	// 登录路由，验证用户凭证并生成JWT令牌
	// router.POST("/login", getTokenHandler)

	// Portal API about host
	portal.POST("/hosts/register", RegisterHostHandler)
	portal.POST("/hosts/batch-register", RegisterHostsHandler)
	portal.POST("/hosts/unregister", UnregisterHostHandler)
	portal.POST("/hosts/batch-unregister", UnregisterHostsHandler)
	portal.GET("/hosts", GetRegisteredHostsHandler)
	// Portal API about directory
	portal.POST("/directories/create", CreateDirectoryHandler)
	portal.POST("/directories/batch-create", CreateDirectoriesHandler)
	portal.POST("/directories/delete", DeleteDirectoryHandler)
	portal.POST("/directories/batch-delete", DeleteDirectoriesHandler)
	portal.GET("/directories", GetDirectoriesHandler)
	// Portal API about local user
	portal.POST("/users/create", CreateLocalUserHandler)
	portal.POST("/users/batch-create", CreateLocalUsersHandler)
	portal.POST("/users/delete", DeleteLocalUserHandler)
	portal.POST("/users/batch-delete", DeleteLocalUsersHandler)
	portal.POST("/users/manage", ManageLocalUserHandler)
	portal.POST("/users/batch-manage", ManageLocalUsersHandler)
	portal.POST("/users/unmanage", UnmanageLocalUserHandler)
	portal.POST("/users/batch-unmanage", UnmanageLocalUsersHandler)
	portal.GET("/users", GetlocalUsersHandler)
	// Portal API about share
	portal.POST("/shares/create", CreateShareHandler)
	portal.POST("/shares/delete", DeleteShareHandler)
	portal.POST("/shares/mount", MountCIFSShareHandler)
	portal.POST("/shares/unmount", UnmountShareHandler)
	portal.GET("/shares", GetSharesHandler)

	// Portal API about swagger-ui
	portal.Static("/docs", "./docs/swagger-ui/dist")

	// Agent API about host
	agent.GET("/system-info", GetSystemInfoOnAgentHandler)
	// Agent API about directory
	agent.GET("/directories/detail", GetDirectoryDetailOnAgentHandler)
	agent.POST("/directories/create", CreateDirectoryOnAgentHandler)
	agent.POST("/directories/batch-create", CreateDirectoriesOnAgentHandler)
	agent.POST("/directories/delete", DeleteDirectoryOnAgentHandler)
	agent.POST("/directories/batch-delete", DeleteDirectoriesOnAgentHandler)
	// Agent API about share
	agent.POST("/shares/create", CreateShareOnAgentHandler)
	agent.POST("/shares/delete", DeleteShareOnAgentHandler)
	agent.POST("/shares/mount", MountShareOnAgentHandler)
	agent.POST("/shares/unmount", UnmountShareOnAgentHandler)
	agent.GET("/shares/detail", GetShareOnAgentHandler)
	// Agent API about local user
	agent.POST("/users/create", CreateLocalUserOnAgentHandler)
	agent.POST("/users/delete", DeleteLocalUserOnAgentHandler)
	agent.GET("/users/detail", GetLocalUserOnAgentHandler)

	addr := fmt.Sprintf(":%d", common.Config.WebService.Port)
	router.Run(addr)

}

func ErrorResponse(c *gin.Context, statusCode int, message string, errMessage string) {
	response := gin.H{"message": message}
	if errMessage != "" {
		response["error"] = errMessage
	}

	c.JSON(statusCode, response)
}

func SetTraceIDInContext(c *gin.Context) context.Context {
	traceID := c.Request.Header.Get("X-Trace-ID")
	return context.WithValue(context.Background(), common.TraceIDKey("TraceID"), traceID)
}
