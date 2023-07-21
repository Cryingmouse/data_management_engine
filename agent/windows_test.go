package agent

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/cryingmouse/data_management_engine/common"
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

func setupCreateDirectory(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.CreateDirectory(ctx, "test_directory")
}

func setupCreateDirectories(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.CreateDirectories(ctx, []string{"test_directory_1", "test_directory_2"})
}

func teardownDeleteDirectory(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.DeleteDirectory(ctx, "test_directory")
}

func teardownDeleteDirectories(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.DeleteDirectories(ctx, []string{"test_directory_1", "test_directory_2"})
}

func setupCreateShare(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.CreateShare(ctx, "test_share", "C:\\test\\test_directory", "this is a test share", []string{"JayXu"})
}

func teardownDeleteShare(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.DeleteShare(ctx, "test_share")
}

func setupCreateShareMapping(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.CreateShareMapping(ctx, "Y:", "\\\\192.168.0.166\\test_share", "JayXu", "Qaviq2ew!")
}

func teardownDeleteShareMapping(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.DeleteShareMapping(ctx, "Y:")
}

func setupCreateLocalUser(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.CreateLocalUser(ctx, "test_account", "test_account")
}

func teardownDeleteLocalUser(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.DeleteLocalUser(ctx, "test_account")
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
				name: "test_directory",
			},
			wantErr:     false,
			wantDirPath: "c:\\test\\test_directory",
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
				names: []string{"test_directory_1", "test_directory_2"},
			},
			wantErr:      false,
			wantDirPaths: []string{"c:\\test\\test_directory_1", "c:\\test\\test_directory_2"},
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
				name: "test_directory",
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
				names: []string{"test_directory_1", "test_directory_2"},
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
		path string
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
				path: "c:\\test\\test_directory",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.agent.GetDirectoryDetail(tt.args.ctx, tt.args.path)
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
		paths []string
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
				paths: []string{"c:\\test\\test_directory_1", "c:\\test\\test_directory_2"},
			},
			wantErr: false,
		},
		{
			name:  "test_get_directories_detail",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				paths: []string{"c:\\test\\test_directory_1"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.agent.GetDirectoriesDetail(tt.args.ctx, tt.args.paths)
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.GetDirectoriesDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestWindowsAgent_CreateShare(t *testing.T) {
	setupCreateDirectory(t)
	defer teardownDeleteShare(t)

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
			name:  "test_create_share",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				name:          "test_share",
				directoryName: "C:\\test\\test_directory",
				description:   "this is a test share",
				usernames: []string{
					"JayXu",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.agent.CreateShare(tt.args.ctx, tt.args.name, tt.args.directoryName, tt.args.description, tt.args.usernames); (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.CreateShare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWindowsAgent_DeleteShare(t *testing.T) {
	setupCreateDirectory(t)
	setupCreateShare(t)

	defer teardownDeleteDirectory(t)

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
			name:  "test_delete_share",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				name: "test_share",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.agent.DeleteShare(tt.args.ctx, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.DeleteShare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWindowsAgent_GetShareDetail(t *testing.T) {
	setupCreateDirectory(t)
	setupCreateShare(t)
	defer teardownDeleteShare(t)
	defer teardownDeleteDirectory(t)

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
			name:  "test_get_share_detail",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				name: "test_share",
			},
			wantErr: false,
			wantDetail: common.ShareDetail{
				Name:          "test_share",
				Description:   "this is a test share",
				DirectoryPath: "C:\\test\\test_directory",
				State:         "online",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDetail, err := tt.agent.GetShareDetail(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.GetShareDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDetail, tt.wantDetail) {
				t.Errorf("WindowsAgent.GetShareDetail() = %v, want %v", gotDetail, tt.wantDetail)
			}
		})
	}
}

func TestWindowsAgent_GetSharesDetail(t *testing.T) {
	setupCreateDirectory(t)
	setupCreateShare(t)
	defer teardownDeleteShare(t)
	defer teardownDeleteDirectory(t)

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
			name:  "test_get_shares_detail",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				names: []string{"test_share", "Users"},
			},
			wantErr: false,
			wantDetail: []common.ShareDetail{
				{
					Name:          "test_share",
					Description:   "this is a test share",
					DirectoryPath: "C:\\test\\test_directory",
					State:         "online",
				},
				{
					Name:          "Users",
					Description:   "",
					DirectoryPath: "C:\\Users",
					State:         "online",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDetail, err := tt.agent.GetSharesDetail(tt.args.ctx, tt.args.names)
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.GetSharesDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDetail, tt.wantDetail) {
				t.Errorf("WindowsAgent.GetSharesDetail() = %v, want %v", gotDetail, tt.wantDetail)
			}
		})
	}
}

func TestWindowsAgent_CreateShareMapping(t *testing.T) {
	setupCreateDirectory(t)
	setupCreateShare(t)
	defer teardownDeleteShareMapping(t)
	defer teardownDeleteShare(t)
	defer teardownDeleteDirectory(t)

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
			name:  "test_create_share_mapping",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				deviceName: "Y:",
				sharePath:  "\\\\192.168.0.166\\test_share",
				userName:   "JayXu",
				password:   "Qaviq2ew!",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.agent.CreateShareMapping(tt.args.ctx, tt.args.deviceName, tt.args.sharePath, tt.args.userName, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.CreateShareMapping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWindowsAgent_DeleteShareMapping(t *testing.T) {
	setupCreateDirectory(t)
	setupCreateShare(t)
	setupCreateShareMapping(t)
	defer teardownDeleteShare(t)
	defer teardownDeleteDirectory(t)

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
			name:  "test_delete_share_mapping",
			agent: GetAgent().(*WindowsAgent),
			args: args{
				deviceName: "Y:",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.agent.DeleteShareMapping(tt.args.ctx, tt.args.deviceName); (err != nil) != tt.wantErr {
				t.Errorf("WindowsAgent.DeleteShareMapping() error = %v, wantErr %v", err, tt.wantErr)
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
				name:     "test_account",
				password: "test_account",
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
				name: "test_account",
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
				name: "test_account",
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
				names: []string{"test_account"},
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
