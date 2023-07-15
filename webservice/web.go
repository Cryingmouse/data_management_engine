package webservice

import (
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
	portal.POST("/directories/create", createDirectoryHandler)
	portal.POST("/directories/batch-create", createDirectoriesHandler)
	portal.POST("/directories/delete", deleteDirectoryHandler)
	portal.POST("/directories/batch-delete", deleteDirectoriesHandler)
	portal.GET("/directories", getDirectoriesHandler)
	// Portal API about local user
	portal.POST("/users/create", createLocalUserHandler)
	portal.POST("/users/batch-create", createLocalUsersHandler)
	portal.POST("/users/delete", deleteLocalUserHandler)
	portal.POST("/users/batch-delete", deleteLocalUsersHandler)
	portal.GET("/users", getlocalUsersHandler)
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
	agent.GET("/users", getLocalUserOnAgentHandler)

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
