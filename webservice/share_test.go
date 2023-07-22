package webservice

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupCreateShareOnAgent(t *testing.T) {
	requestBody := bytes.NewBuffer([]byte(`{"share_name": "test_share", "directory_name": "test_directory", "description": "this is a test share",  "user_names": ["test_account"]}`))

	// 创建测试请求
	req, _ := http.NewRequest("POST", "/agent/shares/create", requestBody)

	req.Header.Set("X-Trace-ID", "123456")

	// 创建响应的Recorder
	w := httptest.NewRecorder()

	// 处理测试请求
	shareRouter.ServeHTTP(w, req)
}

func teardownDeleteShareOnAgent(t *testing.T) {
	requestBody := bytes.NewBuffer([]byte(`{"share_name": "test_share"}`))

	// 创建测试请求
	req, _ := http.NewRequest("POST", "/agent/shares/delete", requestBody)

	req.Header.Set("X-Trace-ID", "123456")

	// 创建响应的Recorder
	w := httptest.NewRecorder()

	// 处理测试请求
	shareRouter.ServeHTTP(w, req)
}

// Create a single share on agent.
func Test_createShareOnAgentHandler(t *testing.T) {
	setupCreateLocalUserOnAgent(t)
	setupCreateDirectoryOnAgent(t)
	defer teardownDeleteShareOnAgent(t)
	defer teardownDeleteDirectoryOnAgent(t)
	defer teardownDeleteLocalUserOnAgent(t)

	requestBody := bytes.NewBuffer([]byte(`{"share_name": "test_share", "directory_name": "test_directory", "description": "this is a test share",  "user_names": ["test_account"]}`))

	// 创建测试请求
	req, _ := http.NewRequest("POST", "/agent/shares/create", requestBody)

	req.Header.Set("X-Trace-ID", "123456")

	// 创建响应的Recorder
	w := httptest.NewRecorder()

	// 处理测试请求
	shareRouter.ServeHTTP(w, req)

	// 断言检查状态码是否为200
	assert.Equal(t, http.StatusOK, w.Code)
}

// Delete a single share on agent.
func Test_deleteShareOnAgentHandler(t *testing.T) {
	setupCreateLocalUserOnAgent(t)
	setupCreateDirectoryOnAgent(t)
	setupCreateShareOnAgent(t)
	defer teardownDeleteDirectoryOnAgent(t)
	defer teardownDeleteLocalUserOnAgent(t)

	requestBody := bytes.NewBuffer([]byte(`{"share_name": "test_share"}`))

	// 创建测试请求
	req, _ := http.NewRequest("POST", "/agent/shares/delete", requestBody)

	req.Header.Set("X-Trace-ID", "123456")

	// 创建响应的Recorder
	w := httptest.NewRecorder()

	// 处理测试请求
	shareRouter.ServeHTTP(w, req)

	// 断言检查状态码是否为200
	assert.Equal(t, http.StatusOK, w.Code)
}

// Get detail of the share/shares on agent
func Test_getShareDetailOnAgentHandler(t *testing.T) {
	setupCreateLocalUserOnAgent(t)
	setupCreateDirectoryOnAgent(t)
	setupCreateShareOnAgent(t)
	defer teardownDeleteShareOnAgent(t)
	defer teardownDeleteDirectoryOnAgent(t)
	defer teardownDeleteLocalUserOnAgent(t)

	tests := []struct {
		name             string
		request          *http.Request
		wantedStatusCode int
	}{
		{
			name: "test_get_share_detail",
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/agent/shares/detail?name=test_share", nil)
				return req
			}(),
			wantedStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.request.Header.Set("X-Trace-ID", "123456")

			// 创建响应的Recorder
			w := httptest.NewRecorder()

			// 处理测试请求
			shareRouter.ServeHTTP(w, tt.request)

			assert.Equal(t, tt.wantedStatusCode, w.Code)
		})
	}
}
