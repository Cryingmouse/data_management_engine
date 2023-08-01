package webservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cryingmouse/data_management_engine/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en_US"
	"github.com/go-playground/locales/zh_Hans_CN"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

var (
	I18NBundle      *i18n.Bundle
	LanguageMatcher language.Matcher
	UniTranslator   *ut.UniversalTranslator
	Validate        *validator.Validate
)

func Start() {
	Validate = binding.Validator.Engine().(*validator.Validate)
	Validate.RegisterValidation("validatePassword", PasswordValidator)

	initializeI18N(Validate)

	router := gin.Default()

	router.Use(cors.Default())
	router.Use(TraceMiddleware(), LoggingMiddleware(), TimeoutMiddleware(8*time.Second), I18nMiddleware())

	// Router 'portal' for Portal
	portal := router.Group("/api")

	// Router 'agent' for Agent
	agent := router.Group("/agent")

	// ====================================
	// Portal related APIs
	// ====================================
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

	// ====================================
	// Agent related APIs
	// ====================================
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

	router.Run(fmt.Sprintf(":%d", common.Config.WebService.Port))
}

func SetTraceIDToContext(c *gin.Context) (context.Context, string) {
	traceID := c.Request.Header.Get("X-Trace-ID")

	common.Logger.WithFields(log.Fields{
		"X-Trace-ID": traceID,
	}).Debug("Get trace id from request header.")

	return context.WithValue(context.Background(), common.TraceIDKey("TraceID"), traceID), traceID
}

func GetTraceIDFromContext(ctx context.Context) string {
	return ctx.Value(common.TraceIDKey("TraceID")).(string)
}

func initializeI18N(validate *validator.Validate) {
	I18NBundle = i18n.NewBundle(language.English)
	I18NBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	I18NBundle.LoadMessageFile("locales/en_US.json")
	I18NBundle.LoadMessageFile("locales/zh_CN.json")
	// 创建一个 language.Matcher
	LanguageMatcher = language.NewMatcher(I18NBundle.LanguageTags())

	en := en_US.New()
	zh := zh_Hans_CN.New()
	UniTranslator = ut.New(en, en, zh)

	USTranslator, _ := UniTranslator.GetTranslator("en_US")
	ZHTranslator, _ := UniTranslator.GetTranslator("zh_Hans_CN")

	en_translations.RegisterDefaultTranslations(validate, USTranslator)
	zh_translations.RegisterDefaultTranslations(validate, ZHTranslator)

	validate.RegisterTranslation("validatePassword", USTranslator, func(ut ut.Translator) error {
		return ut.Add("validatePassword", "Invalid Password", false) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("validatePassword", fe.Field())

		return t
	})

	validate.RegisterTranslation("validatePassword", ZHTranslator, func(ut ut.Translator) error {
		return ut.Add("validatePassword", "无效的密码", false) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("validatePassword", fe.Field())

		return t
	})
}

func TranslateValidationError(c *gin.Context, err error) string {
	acceptLanguage := "en_US"

	switch c.GetHeader("Accept-Language") {
	case "en_US":
		acceptLanguage = "en_US"
	case "zh_CN":
		acceptLanguage = "zh_Hans_CN"
	}

	if trans, ok := UniTranslator.GetTranslator(acceptLanguage); ok {
		errs, _ := err.(validator.ValidationErrors)

		return common.ConvertMapToString(errs.Translate(trans))
	}

	return ""
}

func TranslateError(c *gin.Context, err error) string {
	acceptLanguage := "en_US"

	switch c.GetHeader("Accept-Language") {
	case "en_US":
		acceptLanguage = "en_US"
	case "zh_CN":
		acceptLanguage = "zh_Hans_CN"
	}

	if trans, ok := UniTranslator.GetTranslator(acceptLanguage); ok {
		errs, _ := err.(validator.ValidationErrors)

		return common.ConvertMapToString(errs.Translate(trans))
	}

	return ""
}

func GetLocalizer(c *gin.Context) *i18n.Localizer {
	acceptLanguage := c.GetHeader("Accept-Language")
	tag, _, confidence := LanguageMatcher.Match(language.Make(acceptLanguage))
	if confidence.String() == "No" {
		// Failed to match the language, use a default language.
		tag = language.English
	}

	return i18n.NewLocalizer(I18NBundle, tag.String())
}

func ErrorResponse(c *gin.Context, statusCode int, message string, errMessage string) {
	response := gin.H{"message": message}
	if errMessage != "" {
		response["error"] = errMessage
	}

	c.JSON(statusCode, response)
}

func ErrorResponse_1(c *gin.Context, statusCode int, errorCode *common.Error, err error) {
	response := gin.H{"error": err}
	if errorCode != nil {
		response["error_code"] = errorCode
	}

	c.JSON(statusCode, response)
}
