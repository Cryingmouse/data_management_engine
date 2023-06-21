package main

import (
	"github.com/cryingmouse/data_management_engine/db"
	"github.com/cryingmouse/data_management_engine/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	engine.Migrate()

	router := gin.Default()

	// 登录路由，验证用户凭证并生成JWT令牌
	// router.POST("/login", getTokenHandler)

	router.POST("/api/hosts/register", handler.HostRegistrationHandler)

	router.GET("/api/hosts", handler.GetRegisteredHostsHandler)

	router.POST("/api/hosts/unregister", handler.HostUnregistrationHandler)

	router.POST("/api/directory/create", handler.CreateDirectoryHandler)

	router.POST("/agent/directory/create", handler.CreateDirectoryOnAgentHandler)

	router.Run(":8080")

}
