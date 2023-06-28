package webservice

import (
	"regexp"
	"strconv"

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

func validatePagination(c *gin.Context) (page, limit int, err error) {
	page, _ = strconv.Atoi(c.Query("page"))
	limit, _ = strconv.Atoi(c.Query("limit"))

	// Create a validator instance.
	v := validator.New()

	type Pagination struct {
		Page  int `validate:"omitempty,gte=0"`
		Limit int `validate:"omitempty,gte=0"`
	}

	// Define validation rules for page and limit.
	pagination := Pagination{Page: page, Limit: limit}

	// Custom validation function to check if both page and limit have values or both are empty.
	v.RegisterValidation("pageLimit", func(fl validator.FieldLevel) bool {
		if pagination.Page != 0 && pagination.Limit != 0 {
			return true
		}
		if pagination.Page == 0 && pagination.Limit == 0 {
			return true
		}
		return false
	})

	// Perform validation.
	return page, limit, v.Struct(pagination)
}
