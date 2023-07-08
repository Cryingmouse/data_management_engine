package webservice

import (
	"fmt"
	"regexp"

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
		v.RegisterValidation("validatePassword", passwordValidator)
		v.RegisterValidation("validateStorageType", storageTypeValidator)
	}

	// Router 'portal' for Portal
	portal := router.Group("/api")
	// Router 'agent' for Agent
	agent := router.Group("/agnet")

	// 登录路由，验证用户凭证并生成JWT令牌
	// router.POST("/login", getTokenHandler)

	// Portal API about host
	portal.POST("/hosts/register", registerHostHandler)
	portal.POST("/hosts/batch-register", registerHostsHandler)
	portal.POST("/hosts/unregister", unregisterHostHandler)
	portal.POST("/hosts/batch-unregister", unregisterHostsHandler)
	portal.GET("/hosts", getRegisteredHostsHandler)
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
	agent.GET("/system-info", getSystemInfoOnAgentHandler)
	// Agent API about directory
	agent.POST("/directories/create", createDirectoryOnAgentHandler)
	agent.POST("/directories/delete", deleteDirectoryOnAgentHandler)
	// Agent API about local user
	agent.POST("/users/create", createLocalUserOnAgentHandler)
	agent.POST("/users/delete", deleteLocalUserOnAgentHandler)
	agent.GET("/users", getLocalUserOnAgentHandler)

	addr := fmt.Sprintf(":%s", config.Webservice.Port)
	router.Run(addr)

}

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) >= 8 && regexp.MustCompile(`[A-Z]+`).MatchString(password) && regexp.MustCompile(`[a-z]+`).MatchString(password) && regexp.MustCompile(`[0-9]+`).MatchString(password) {
		return true
	}
	return false
}

func validatePagination(page, limit int) (err error) {

	// Create a validator instance.
	v := validator.New()

	type Pagination struct {
		Page  int
		Limit int
	}

	// Define validation rules for page and limit.
	pagination := Pagination{Page: page, Limit: limit}

	// Custom validation function to check if both page and limit have values or both are empty.
	v.RegisterValidation("pageLimit", func(fl validator.FieldLevel) bool {
		if pagination.Page >= 0 && pagination.Limit > 0 {
			return true
		}
		if pagination.Page == 0 && pagination.Limit == 0 {
			return true
		}
		return false
	})

	// Perform validation.
	return v.Struct(pagination)
}

func validateIPAddress(ip string) (err error) {
	type IPAddress struct {
		IP string `validate:"required,ip"`
	}

	validate := validator.New()

	// 验证 IP 地址
	ipAddress := IPAddress{IP: ip}

	return validate.Struct(ipAddress)
}
