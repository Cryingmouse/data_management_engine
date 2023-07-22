package webservice

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupCreateLocalUserOnAgent(t *testing.T) {
	requestBody := bytes.NewBuffer([]byte(`{"name": "test_account", "password": "Passw0rd!"}`))

	// 创建测试请求
	req, _ := http.NewRequest("POST", "/agent/users/create", requestBody)

	req.Header.Set("X-Trace-ID", "123456")

	// 创建响应的Recorder
	w := httptest.NewRecorder()

	// 处理测试请求
	shareRouter.ServeHTTP(w, req)
}

func teardownDeleteLocalUserOnAgent(t *testing.T) {
	requestBody := bytes.NewBuffer([]byte(`{"name": "test_account"}`))

	// 创建测试请求
	req, _ := http.NewRequest("POST", "/agent/users/delete", requestBody)

	req.Header.Set("X-Trace-ID", "123456")

	// 创建响应的Recorder
	w := httptest.NewRecorder()

	// 处理测试请求
	shareRouter.ServeHTTP(w, req)
}

// Create a single share on agent.
func Test_createLocalUserOnAgentHandler(t *testing.T) {
	defer teardownDeleteLocalUserOnAgent(t)

	requestBody := bytes.NewBuffer([]byte(`{"name": "test_account", "password": "Passw0rd!"}`))

	// 创建测试请求
	req, _ := http.NewRequest("POST", "/agent/users/create", requestBody)

	req.Header.Set("X-Trace-ID", "123456")

	// 创建响应的Recorder
	w := httptest.NewRecorder()

	// 处理测试请求
	shareRouter.ServeHTTP(w, req)

	// 断言检查状态码是否为200
	assert.Equal(t, http.StatusOK, w.Code)
}

// Delete a single share on agent.
func Test_deleteLocalUserOnAgentHandler(t *testing.T) {
	setupCreateLocalUserOnAgent(t)

	requestBody := bytes.NewBuffer([]byte(`{"name": "test_account"}`))

	// 创建测试请求
	req, _ := http.NewRequest("POST", "/agent/users/delete", requestBody)

	req.Header.Set("X-Trace-ID", "123456")

	// 创建响应的Recorder
	w := httptest.NewRecorder()

	// 处理测试请求
	shareRouter.ServeHTTP(w, req)

	// 断言检查状态码是否为200
	assert.Equal(t, http.StatusOK, w.Code)
}

// Get detail of the share/shares on agent
func Test_getLocalUserDetailOnAgentHandler(t *testing.T) {
	setupCreateLocalUserOnAgent(t)
	defer teardownDeleteLocalUserOnAgent(t)

	tests := []struct {
		name             string
		request          *http.Request
		wantedStatusCode int
	}{
		{
			name: "test_get_local_user_detail",
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/agent/users/detail?name=test_account", nil)
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
