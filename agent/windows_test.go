package agent

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cryingmouse/data_management_engine/context"
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

func Test_Windows_CreateShare(t *testing.T) {

	agent := GetAgent()

	hostContext := context.HostContext{
		IP:       "127.09.0.1",
		Username: "admin",
		Password: "password",
	}

	shareName := "NewCIFS_1"

	err := agent.CreateShare(hostContext, shareName, "C:\\test\\folder-2")
	// Assert
	if err != nil {
		t.Errorf("Failed to create share %s. Error: %v", shareName, err)
	}
}

func Test_Windows_GetLocalUsers(t *testing.T) {

	agent := GetAgent()

	hostContext := context.HostContext{
		IP:       "127.09.0.1",
		Username: "admin",
		Password: "password",
	}

	users, err := agent.GetLocalUsers(hostContext)
	if err != nil {
		t.Errorf("Failed to get local users: %v. Error: %v", users, err)
	}
}

func Test_Windows_GetSystemInfo(t *testing.T) {

	agent := GetAgent()

	hostContext := context.HostContext{
		IP:       "127.09.0.1",
		Username: "admin",
		Password: "password",
	}

	systemInfo, err := agent.GetSystemInfo(hostContext)
	if err != nil {
		t.Errorf("Failed to get system info: %v. Error: %v", systemInfo, err)
	}
}

func Test_Windows_CreateLocalUsers(t *testing.T) {

	agent := GetAgent()

	hostContext := context.HostContext{
		IP:       "127.09.0.1",
		Username: "admin",
		Password: "password",
	}

	username := "test"
	password := "test"

	err := agent.CreateLocalUser(hostContext, username, password)
	if err != nil {
		t.Errorf("Failed to create local users: %s. Error: %v", username, err)
	}
}

func Test_Windows_DeleteLocalUsers(t *testing.T) {

	agent := GetAgent()

	hostContext := context.HostContext{
		IP:       "127.09.0.1",
		Username: "admin",
		Password: "password",
	}

	username := "test"

	err := agent.DeleteLocalUser(hostContext, username)
	if err != nil {
		t.Errorf("Failed to create local users: %s. Error: %v", username, err)
	}
}
