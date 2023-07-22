package webservice

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupCreateDirectoryOnAgent(t *testing.T) {
	requestBody := bytes.NewBuffer([]byte(`{"name": "test_directory"}`))

	// 创建测试请求
	req, _ := http.NewRequest("POST", "/agent/directories/create", requestBody)

	req.Header.Set("X-Trace-ID", "123456")

	// 创建响应的Recorder
	w := httptest.NewRecorder()

	// 处理测试请求
	shareRouter.ServeHTTP(w, req)
}

func teardownDeleteDirectoryOnAgent(t *testing.T) {
	requestBody := bytes.NewBuffer([]byte(`{"name": "test_directory"}`))

	// 创建测试请求
	req, _ := http.NewRequest("POST", "/agent/directories/delete", requestBody)

	req.Header.Set("X-Trace-ID", "123456")

	// 创建响应的Recorder
	w := httptest.NewRecorder()

	// 处理测试请求
	shareRouter.ServeHTTP(w, req)
}

func setupCreateDirectoriesOnAgent(t *testing.T) {
	requestBody := bytes.NewBuffer([]byte(`[{"name": "test_directory_1"},{"name": "test_directory_2"}]`))

	// 创建测试请求
	req, _ := http.NewRequest("POST", "/agent/directories/batch-create", requestBody)

	req.Header.Set("X-Trace-ID", "123456")

	// 创建响应的Recorder
	w := httptest.NewRecorder()

	// 处理测试请求
	shareRouter.ServeHTTP(w, req)
}

func teardownDeleteDirectoriesOnAgent(t *testing.T) {
	requestBody := bytes.NewBuffer([]byte(`[{"name": "test_directory_1"},{"name": "test_directory_2"}]`))

	// 创建测试请求
	req, _ := http.NewRequest("POST", "/agent/directories/batch-delete", requestBody)

	req.Header.Set("X-Trace-ID", "123456")

	// 创建响应的Recorder
	w := httptest.NewRecorder()

	// 处理测试请求
	shareRouter.ServeHTTP(w, req)
}

func Test_createDirectoryOnAgentHandler(t *testing.T) {
	defer teardownDeleteDirectoryOnAgent(t)

	requestBody := bytes.NewBuffer([]byte(`{"name": "test_directory"}`))

	// 创建测试请求
	req, _ := http.NewRequest("POST", "/agent/directories/create", requestBody)

	req.Header.Set("X-Trace-ID", "123456")

	// 创建响应的Recorder
	w := httptest.NewRecorder()

	// 处理测试请求
	shareRouter.ServeHTTP(w, req)

	// 断言检查状态码是否为200
	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_deleteDirectoryOnAgentHandler(t *testing.T) {
	setupCreateDirectoryOnAgent(t)

	requestBody := bytes.NewBuffer([]byte(`{"name": "test_directory"}`))

	// 创建测试请求
	req, _ := http.NewRequest("POST", "/agent/directories/delete", requestBody)

	req.Header.Set("X-Trace-ID", "123456")

	// 创建响应的Recorder
	w := httptest.NewRecorder()

	// 处理测试请求
	shareRouter.ServeHTTP(w, req)

	// 断言检查状态码是否为200
	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_createDirectoriesOnAgentHandler(t *testing.T) {
	defer teardownDeleteDirectoriesOnAgent(t)

	requestBody := bytes.NewBuffer([]byte(`[{"name": "test_directory_1"},{"name": "test_directory_2"}]`))

	// 创建测试请求
	req, _ := http.NewRequest("POST", "/agent/directories/batch-create", requestBody)

	req.Header.Set("X-Trace-ID", "123456")

	// 创建响应的Recorder
	w := httptest.NewRecorder()

	// 处理测试请求
	shareRouter.ServeHTTP(w, req)

	// 断言检查状态码是否为200
	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_deleteDirectoriesOnAgentHandler(t *testing.T) {
	setupCreateDirectoriesOnAgent(t)

	requestBody := bytes.NewBuffer([]byte(`[{"name": "test_directory_1"},{"name": "test_directory_2"}]`))

	// 创建测试请求
	req, _ := http.NewRequest("POST", "/agent/directories/batch-delete", requestBody)

	req.Header.Set("X-Trace-ID", "123456")

	// 创建响应的Recorder
	w := httptest.NewRecorder()

	// 处理测试请求
	shareRouter.ServeHTTP(w, req)

	// 断言检查状态码是否为200
	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_getDirectoryDetailOnAgentHandler(t *testing.T) {
	setupCreateDirectoriesOnAgent(t)
	defer teardownDeleteDirectoriesOnAgent(t)

	tests := []struct {
		name             string
		request          *http.Request
		wantedStatusCode int
	}{
		{
			name: "test_get_directory_detail",
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/agent/directories/detail?name=test_directory_1", nil)
				return req
			}(),
			wantedStatusCode: http.StatusOK,
		},
		{
			name: "test_get_directory_detail",
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/agent/directories/detail?name=test_directory_1,test_directory_2", nil)
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
