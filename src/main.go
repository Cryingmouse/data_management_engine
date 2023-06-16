package main

import (
	"github.com/gin-gonic/gin"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	router := gin.Default()

	// 登录路由，验证用户凭证并生成JWT令牌
	router.POST("/login", getTokenHandler)

	// 登录路由，验证用户凭证并生成JWT令牌
	router.POST("/api/add-host", hostRegistrationHandler)

	// 受保护的路由，需要在请求头中带有有效的JWT令牌才能访问
	router.GET("/protected", authMiddleware(), protectedHandler)

	router.Run(":8080")

}
