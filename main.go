package main

import (
	"db"

	"github.com/gin-gonic/gin"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	engine.Migrate()

	router := gin.Default()

	// 登录路由，验证用户凭证并生成JWT令牌
	router.POST("/login", getTokenHandler)

	router.POST("/api/register-host", hostRegistrationHandler)

	router.GET("/api/registered-hosts", getRegisteredHostsHandler)

	// router.POST("/api/unregister-host", hostUnregistrationHandler)

	router.Run(":8080")

}
