package agent

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/cryingmouse/data_management_engine/common"
)

const testHostIP = "192.168.0.166"
const testShareName = "test_cifs_share"
const testDirectoryName = "test_directory"
const testDirectoryName1 = "test_directory_1"
const testDirectoryName2 = "test_directory_2"
const testLocalUserName = "test_account"
const testLocalUserPassword = "Passw0rd!"
const testDeviceName = "Y:"

var testSharePath = fmt.Sprintf("\\\\%s\\%s", testHostIP, testShareName)

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

	if err := common.InitializeConfig("config.ini"); err != nil {
		panic(err)
	}

	// 执行测试
	exitCode := m.Run()

	// 退出测试
	os.Exit(exitCode)
}

func setupCreateDirectory(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.CreateDirectory(ctx, testDirectoryName)
}

func setupCreateDirectories(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.CreateDirectories(ctx, []string{testDirectoryName1, testDirectoryName2})
}

func teardownDeleteDirectory(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.DeleteDirectory(ctx, testDirectoryName)
}

func teardownDeleteDirectories(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.DeleteDirectories(ctx, []string{testDirectoryName1, testDirectoryName2})
}

func setupCreateCIFSShare(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.CreateCIFSShare(ctx, testShareName, testDirectoryName, "this is a test cifs share", []string{testLocalUserName})
}

func teardownDeleteCIFSShare(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.DeleteCIFSShare(ctx, testShareName)
}

func setupMountCIFSShare(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.MountCIFSShare(ctx, testDeviceName, testSharePath, testLocalUserName, testLocalUserPassword)
}

func teardownUnmountCIFSShare(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.UnmountCIFSShare(ctx, testDeviceName)
}

func setupCreateLocalUser(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.CreateLocalUser(ctx, testLocalUserName, testLocalUserPassword)
}

func teardownDeleteLocalUser(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.DeleteLocalUser(ctx, testLocalUserName)
}

func TestWindowsAgent_CreateDirectory(t *testing.T) {
	defer teardownDeleteDirectory(t)

	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name        string
		agent       *WindowsAgent
		args        args
		wantDirPath string
		wantErr     bool
	}{
		{
			name:  "test_create_directory",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				name: testDirectoryName,
			},
			wantErr:     false,
			wantDirPath: fmt.Sprintf("%s\\%s", common.Config.Agent.WindowsRootFolder, testDirectoryName),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDirPath, err := tt.agent.CreateDirectory(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.CreateDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDirPath != tt.wantDirPath {
				t.Errorf("WindowsAgent.CreateDirectory() = %v, want %v", gotDirPath, tt.wantDirPath)
			}
		})
	}
}

func TestWindowsAgent_CreateDirectories(t *testing.T) {
	defer teardownDeleteDirectories(t)

	type args struct {
		ctx   context.Context
		names []string
	}
	tests := []struct {
		name         string
		agent        *WindowsAgent
		args         args
		wantDirPaths []string
		wantErr      bool
	}{
		{
			name:  "test_create_directories",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				names: []string{testDirectoryName1, testDirectoryName2},
			},
			wantErr: false,
			wantDirPaths: []string{
				fmt.Sprintf("%s\\%s", common.Config.Agent.WindowsRootFolder, testDirectoryName1),
				fmt.Sprintf("%s\\%s", common.Config.Agent.WindowsRootFolder, testDirectoryName2),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDirPaths, err := tt.agent.CreateDirectories(tt.args.ctx, tt.args.names)
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.CreateDirectories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDirPaths, tt.wantDirPaths) {
				t.Errorf("WindowsAgent.CreateDirectories() = %v, want %v", gotDirPaths, tt.wantDirPaths)
			}
		})
	}
}

func TestWindowsAgent_DeleteDirectory(t *testing.T) {
	setupCreateDirectory(t)

	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		agent   *WindowsAgent
		args    args
		wantErr bool
	}{
		{
			name:  "test_delete_directory",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				name: testDirectoryName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.agent.DeleteDirectory(tt.args.ctx, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.DeleteDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWindowsAgent_DeleteDirectories(t *testing.T) {
	setupCreateDirectories(t)

	type args struct {
		ctx   context.Context
		names []string
	}
	tests := []struct {
		name    string
		agent   *WindowsAgent
		args    args
		wantErr bool
	}{
		{
			name:  "test_create_directories",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				names: []string{testDirectoryName1, testDirectoryName2},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.agent.DeleteDirectories(tt.args.ctx, tt.args.names); (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.DeleteDirectories() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWindowsAgent_GetDirectoryDetail(t *testing.T) {
	setupCreateDirectory(t)
	defer teardownDeleteDirectory(t)

	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name       string
		agent      *WindowsAgent
		args       args
		wantDetail common.DirectoryDetail
		wantErr    bool
	}{
		{
			name:  "test_get_directory_detail",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				name: testDirectoryName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.agent.GetDirectoryDetail(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.GetDirectoryDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestWindowsAgent_GetDirectoriesDetail(t *testing.T) {
	setupCreateDirectories(t)
	defer teardownDeleteDirectories(t)

	type args struct {
		ctx   context.Context
		names []string
	}
	tests := []struct {
		name       string
		agent      *WindowsAgent
		args       args
		wantDetail []common.DirectoryDetail
		wantErr    bool
	}{
		{
			name:  "test_get_directories_detail",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				names: []string{testDirectoryName1, testDirectoryName2},
			},
			wantErr: false,
		},
		{
			name:  "test_get_directories_detail",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				names: []string{testDirectoryName1},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.agent.GetDirectoriesDetail(tt.args.ctx, tt.args.names)
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.GetDirectoriesDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestWindowsAgent_CreateCIFSShare(t *testing.T) {
	setupCreateLocalUser(t)
	setupCreateDirectory(t)
	defer teardownDeleteCIFSShare(t)
	defer teardownDeleteDirectory(t)
	defer teardownDeleteLocalUser(t)

	type args struct {
		ctx           context.Context
		name          string
		directoryName string
		description   string
		usernames     []string
	}
	tests := []struct {
		name    string
		agent   *WindowsAgent
		args    args
		wantErr bool
	}{
		{
			name:  "test_create_cifs_share",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				name:          testShareName,
				directoryName: testDirectoryName,
				description:   "this is a test cifs share",
				usernames: []string{
					testLocalUserName,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.agent.CreateCIFSShare(tt.args.ctx, tt.args.name, tt.args.directoryName, tt.args.description, tt.args.usernames); (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.CreateCIFSShare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWindowsAgent_DeleteCIFSShare(t *testing.T) {
	setupCreateLocalUser(t)
	setupCreateDirectory(t)
	setupCreateCIFSShare(t)

	defer teardownDeleteDirectory(t)
	defer teardownDeleteLocalUser(t)

	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		agent   *WindowsAgent
		args    args
		wantErr bool
	}{
		{
			name:  "test_delete_cifs_share",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				name: testShareName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.agent.DeleteCIFSShare(tt.args.ctx, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.DeleteCIFSShare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWindowsAgent_GetCIFSShareDetail(t *testing.T) {
	setupCreateLocalUser(t)
	setupCreateDirectory(t)
	setupCreateCIFSShare(t)
	defer teardownDeleteCIFSShare(t)
	defer teardownDeleteDirectory(t)
	defer teardownDeleteLocalUser(t)

	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name       string
		agent      *WindowsAgent
		args       args
		wantDetail common.ShareDetail
		wantErr    bool
	}{
		{
			name:  "test_get_cifs_share_detail",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				name: testShareName,
			},
			wantErr: false,
			wantDetail: common.ShareDetail{
				Name:          testShareName,
				Description:   "this is a test cifs share",
				DirectoryPath: fmt.Sprintf("%s\\%s", common.Config.Agent.WindowsRootFolder, testDirectoryName),
				State:         "online",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDetail, err := tt.agent.GetCIFSShareDetail(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.GetCIFSShareDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDetail, tt.wantDetail) {
				t.Errorf("WindowsAgent.GetCIFSShareDetail() = %v, want %v", gotDetail, tt.wantDetail)
			}
		})
	}
}

func TestWindowsAgent_GetCIFSSharesDetail(t *testing.T) {
	setupCreateLocalUser(t)
	setupCreateDirectory(t)
	setupCreateCIFSShare(t)
	defer teardownDeleteCIFSShare(t)
	defer teardownDeleteDirectory(t)
	defer teardownDeleteLocalUser(t)

	type args struct {
		ctx   context.Context
		names []string
	}
	tests := []struct {
		name       string
		agent      *WindowsAgent
		args       args
		wantDetail []common.ShareDetail
		wantErr    bool
	}{
		{
			name:  "test_get_cifs_shares_detail",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				names: []string{testShareName},
			},
			wantErr: false,
			wantDetail: []common.ShareDetail{
				{
					Name:          testShareName,
					Description:   "this is a test cifs share",
					DirectoryPath: fmt.Sprintf("%s\\%s", common.Config.Agent.WindowsRootFolder, testDirectoryName),
					State:         "online",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDetail, err := tt.agent.GetCIFSSharesDetail(tt.args.ctx, tt.args.names)
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.GetCIFSSharesDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDetail, tt.wantDetail) {
				t.Errorf("WindowsAgent.GetCIFSSharesDetail() = %v, want %v", gotDetail, tt.wantDetail)
			}
		})
	}
}

func TestWindowsAgent_MountCIFSShare(t *testing.T) {
	setupCreateLocalUser(t)
	setupCreateDirectory(t)
	setupCreateCIFSShare(t)
	defer teardownUnmountCIFSShare(t)
	defer teardownDeleteCIFSShare(t)
	defer teardownDeleteDirectory(t)
	defer teardownDeleteLocalUser(t)

	type args struct {
		ctx        context.Context
		deviceName string
		sharePath  string
		userName   string
		password   string
	}
	tests := []struct {
		name    string
		agent   *WindowsAgent
		args    args
		wantErr bool
	}{
		{
			name:  "test_mount_cifs_share",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				deviceName: testDeviceName,
				sharePath:  testSharePath,
				userName:   testLocalUserName,
				password:   testLocalUserPassword,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.agent.MountCIFSShare(tt.args.ctx, tt.args.deviceName, tt.args.sharePath, tt.args.userName, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.MountCIFSShare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWindowsAgent_UnmountCIFSShare(t *testing.T) {
	setupCreateLocalUser(t)
	setupCreateDirectory(t)
	setupCreateCIFSShare(t)
	setupMountCIFSShare(t)
	defer teardownDeleteCIFSShare(t)
	defer teardownDeleteDirectory(t)
	defer teardownDeleteLocalUser(t)

	type args struct {
		ctx        context.Context
		deviceName string
	}
	tests := []struct {
		name    string
		agent   *WindowsAgent
		args    args
		wantErr bool
	}{
		{
			name:  "test_unmount_cifs_share",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				deviceName: testDeviceName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.agent.UnmountCIFSShare(tt.args.ctx, tt.args.deviceName); (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.UnmountCIFSShare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWindowsAgent_CreateLocalUser(t *testing.T) {
	defer teardownDeleteLocalUser(t)

	type args struct {
		ctx      context.Context
		name     string
		password string
	}
	tests := []struct {
		name    string
		agent   *WindowsAgent
		args    args
		wantErr bool
	}{
		{
			name:  "test_create_local_user",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				name:     testLocalUserName,
				password: testLocalUserPassword,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.agent.CreateLocalUser(tt.args.ctx, tt.args.name, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.CreateLocalUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWindowsAgent_DeleteLocalUser(t *testing.T) {
	setupCreateLocalUser(t)

	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		agent   *WindowsAgent
		args    args
		wantErr bool
	}{
		{
			name:  "test_delete_local_user",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				name: testLocalUserName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.agent.DeleteLocalUser(tt.args.ctx, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.DeleteLocalUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWindowsAgent_GetLocalUserDetail(t *testing.T) {
	setupCreateLocalUser(t)
	defer teardownDeleteLocalUser(t)

	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name       string
		agent      *WindowsAgent
		args       args
		wantDetail common.LocalUserDetail
		wantErr    bool
	}{
		{
			name:  "test_get_local_user_detail",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				name: testLocalUserName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.agent.GetLocalUserDetail(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.GetLocalUserDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestWindowsAgent_GetLocalUsersDetail(t *testing.T) {
	setupCreateLocalUser(t)
	defer teardownDeleteLocalUser(t)

	type args struct {
		ctx   context.Context
		names []string
	}
	tests := []struct {
		name       string
		agent      *WindowsAgent
		args       args
		wantDetail []common.LocalUserDetail
		wantErr    bool
	}{
		{
			name:  "test_get_local_users_detail",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				names: []string{testLocalUserName},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.agent.GetLocalUsersDetail(tt.args.ctx, tt.args.names)
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.GetLocalUsersDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestWindowsAgent_GetSystemInfo(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name           string
		agent          *WindowsAgent
		args           args
		wantSystemInfo common.SystemInfo
		wantErr        bool
	}{
		{
			name:    "test_get_system_info",
			agent:   GetAgent().(*WindowsAgent),
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.agent.GetSystemInfo(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.GetSystemInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
