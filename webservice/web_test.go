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
	shareAgnet.POST("/directories/create", CreateDirectoryOnAgentHandler)
	shareAgnet.POST("/directories/delete", DeleteDirectoryOnAgentHandler)
	shareAgnet.POST("/directories/batch-create", CreateDirectoriesOnAgentHandler)
	shareAgnet.POST("/directories/batch-delete", DeleteDirectoriesOnAgentHandler)
	shareAgnet.GET("/directories/detail", GetDirectoryDetailOnAgentHandler)

	shareAgnet.POST("/shares/create", CreateShareOnAgentHandler)
	shareAgnet.POST("/shares/delete", DeleteShareOnAgentHandler)
	shareAgnet.POST("/shares/mount", MountShareOnAgentHandler)
	shareAgnet.POST("/shares/unmount", UnmountShareOnAgentHandler)
	shareAgnet.GET("/shares/detail", GetShareOnAgentHandler)

	shareAgnet.POST("/users/create", CreateLocalUserOnAgentHandler)
	shareAgnet.POST("/users/delete", DeleteLocalUserOnAgentHandler)
	shareAgnet.GET("/users/detail", GetLocalUserOnAgentHandler)

	shareAgnet.GET("/system-info", GetSystemInfoOnAgentHandler)

	// 执行测试
	exitCode := m.Run()

	// 退出测试
	os.Exit(exitCode)
}
