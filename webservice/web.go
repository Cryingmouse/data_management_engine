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
	config, err := common.GetConfig()
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Use(cors.Default())

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validatePassword", PasswordValidator)
		v.RegisterValidation("validateStorageType", StorageTypeValidator)
	}

	// Router 'portal' for Portal
	portal := router.Group("/api")
	portal.Use(TraceMiddleware(), LoggingMiddleware())

	// Router 'agent' for Agent
	agent := router.Group("/agent")
	agent.Use(LoggingMiddleware())

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
	// Portal API about swagger-ui
	portal.Static("/docs", "./docs/swagger-ui/dist")

	// Agent API about host
	agent.GET("/system-info", GetSystemInfoOnAgentHandler)
	// Agent API about directory
	agent.GET("/directories/detail", getDirectoryDetailsOnAgentHandler)
	agent.POST("/directories/create", createDirectoryOnAgentHandler)
	agent.POST("/directories/batch-create", createDirectoriesOnAgentHandler)
	agent.POST("/directories/delete", deleteDirectoryOnAgentHandler)
	agent.POST("/directories/batch-delete", deleteDirectoriesOnAgentHandler)
	// Agent API about local user
	agent.POST("/users/create", createLocalUserOnAgentHandler)
	agent.POST("/users/delete", deleteLocalUserOnAgentHandler)
	agent.GET("/users/detail", getLocalUserOnAgentHandler)

	addr := fmt.Sprintf(":%s", config.Webservice.Port)
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
	traceID, _ := c.Request.Context().Value(common.TraceIDKey("TraceID")).(string)
	return context.WithValue(context.Background(), common.TraceIDKey("TraceID"), traceID)
}
