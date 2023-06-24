package webservice

import (
	"regexp"

	"github.com/cryingmouse/data_management_engine/db"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func Start() {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	engine.Migrate()

	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validatePassword", passwordValidator)
	}

	// 登录路由，验证用户凭证并生成JWT令牌
	// router.POST("/login", getTokenHandler)

	router.POST("/api/host/register", hostRegistrationHandler)

	router.GET("/api/hosts", getRegisteredHostsHandler)

	router.POST("/api/host/unregister", hostUnregistrationHandler)

	router.POST("/api/directory/create", createDirectoryHandler)

	router.POST("/api/directory/delete", deleteDirectoryHandler)

	router.GET("/api/directories", getDirectoryHandler)

	router.POST("/agent/directory/create", createDirectoryOnAgentHandler)

	router.POST("/agent/directory/delete", deleteDirectoryOnAgentHandler)

	router.Run(":8080")

}

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) >= 8 && regexp.MustCompile(`[A-Z]+`).MatchString(password) && regexp.MustCompile(`[a-z]+`).MatchString(password) && regexp.MustCompile(`[0-9]+`).MatchString(password) {
		return true
	}
	return false
}
