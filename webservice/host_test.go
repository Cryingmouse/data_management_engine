package webservice

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSystemInfoOnAgentHandler(t *testing.T) {
	// 创建测试请求
	req, _ := http.NewRequest("GET", "/agent/system-info", nil)

	req.Header.Set("X-Trace-ID", "123456")

	// 创建响应的Recorder
	w := httptest.NewRecorder()

	// 处理测试请求
	shareRouter.ServeHTTP(w, req)

	// 断言检查状态码是否为200
	assert.Equal(t, http.StatusOK, w.Code)
}
