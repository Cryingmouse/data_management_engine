package webservice

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cryingmouse/data_management_engine/common"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var shareRouter *gin.Engine
var shareAgnet *gin.RouterGroup

func TestMain(m *testing.M) {
	// 获取当前文件所在的目录
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	// 切换到项目的根目录
	projectPath := filepath.Join(dir, "../") // 假设项目的根目录在当前目录的上一级目录
	err := os.Chdir(projectPath)
	if err != nil {
		panic(err)
	}

	common.InitializeConfig("config.ini")

	shareRouter = gin.Default()
	shareAgnet = shareRouter.Group("/agent")

	shareRouter.Use(TraceMiddleware(), LoggingMiddleware())

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validatePassword", PasswordValidator)
		v.RegisterValidation("validateStorageType", StorageTypeValidator)
	}

	// 定义测试路由
	shareAgnet.POST("/directories/create", createDirectoryOnAgentHandler)
	shareAgnet.POST("/directories/delete", deleteDirectoryOnAgentHandler)
	shareAgnet.POST("/directories/batch-create", createDirectoriesOnAgentHandler)
	shareAgnet.POST("/directories/batch-delete", deleteDirectoriesOnAgentHandler)
	shareAgnet.GET("/directories/detail", getDirectoryDetailOnAgentHandler)

	shareAgnet.POST("/shares/create", createShareOnAgentHandler)
	shareAgnet.POST("/shares/delete", deleteShareOnAgentHandler)
	shareAgnet.GET("/shares/detail", getShareOnAgentHandler)

	shareAgnet.POST("/users/create", createLocalUserOnAgentHandler)
	shareAgnet.POST("/users/delete", deleteLocalUserOnAgentHandler)
	shareAgnet.GET("/users/detail", getLocalUserOnAgentHandler)

	shareAgnet.GET("/system-info", GetSystemInfoOnAgentHandler)

	// 执行测试
	exitCode := m.Run()

	// 退出测试
	os.Exit(exitCode)
}
