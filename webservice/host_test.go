package webservice

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

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

	// 执行测试
	exitCode := m.Run()

	// 退出测试
	os.Exit(exitCode)
}

func TestGetSystemInfoOnAgentHandler(t *testing.T) {
	// 创建Gin引擎
	router := gin.Default()
	agent := router.Group("/agent")

	// 定义测试路由
	agent.GET("/system-info", GetSystemInfoOnAgentHandler)

	// 创建测试请求
	req, _ := http.NewRequest("GET", "/agent/system-info", nil)

	req.Header.Set("X-Trace-ID", "123456")

	// 创建响应的Recorder
	w := httptest.NewRecorder()

	// 处理测试请求
	router.ServeHTTP(w, req)

	// 断言检查状态码是否为200
	assert.Equal(t, http.StatusOK, w.Code)
}
