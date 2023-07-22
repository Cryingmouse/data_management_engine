package webservice

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
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

	shareRouter = gin.Default()
	shareAgnet = shareRouter.Group("/agent")

	// 定义测试路由
	shareAgnet.POST("/directories/create", createDirectoryOnAgentHandler)
	shareAgnet.POST("/directories/delete", deleteDirectoryOnAgentHandler)
	shareAgnet.POST("/directories/batch-create", createDirectoriesOnAgentHandler)
	shareAgnet.POST("/directories/batch-delete", deleteDirectoriesOnAgentHandler)
	shareAgnet.GET("/directories/detail", getDirectoryDetailOnAgentHandler)

	shareAgnet.GET("/system-info", GetSystemInfoOnAgentHandler)

	// 执行测试
	exitCode := m.Run()

	// 退出测试
	os.Exit(exitCode)
}
