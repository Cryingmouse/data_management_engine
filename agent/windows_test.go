package agent

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

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

func setupCreateShare(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.CreateShare(ctx, "test_share", "C:\\test\\directory-1", "this is a test share", []string{"JayXu"})
}

func teardownDeleteShare(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.DeleteShare(ctx, "test_share")
}

func setupCreateShareMapping(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.CreateShareMapping(ctx, "Y:", "\\\\10.128.15.56\\test_share", "JayXu", "Qaviq2ew!")
}

func teardownDeleteShareMapping(t *testing.T) {
	windowsAgent := GetAgent().(*WindowsAgent)
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), "123456")
	windowsAgent.DeleteShareMapping(ctx, "Y:")
}

func TestWindowsAgent_CreateShare(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

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
				ctx:           ctx,
				name:          "test_share",
				directoryName: "C:\\test\\directory-1",
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	setupCreateShare(t)

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
				ctx:  ctx,
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

func TestWindowsAgent_CreateShareMapping(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	setupCreateShare(t)
	defer teardownDeleteShareMapping(t)
	defer teardownDeleteShare(t)

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
				ctx:        ctx,
				deviceName: "Y:",
				sharePath:  "\\\\10.128.15.56\\test_share",
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	setupCreateShare(t)
	setupCreateShareMapping(t)
	defer teardownDeleteShare(t)

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
				ctx:        ctx,
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

func TestWindowsAgent_GetShareDetail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	setupCreateShare(t)
	defer teardownDeleteShare(t)

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
				ctx:  ctx,
				name: "test_share",
			},
			wantErr: false,
			wantDetail: common.ShareDetail{
				Name:          "test_share",
				Description:   "this is a test share",
				DirectoryPath: "C:\\test\\directory-1",
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	setupCreateShare(t)
	defer teardownDeleteShare(t)

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
				ctx:   ctx,
				names: []string{"test_share", "Users"},
			},
			wantErr: false,
			wantDetail: []common.ShareDetail{
				{
					Name:          "test_share",
					Description:   "this is a test share",
					DirectoryPath: "C:\\test\\directory-1",
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

// func TestWindowsAgent_GetDirectoryDetail(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		path string
// 	}
// 	tests := []struct {
// 		name       string
// 		agent      *WindowsAgent
// 		args       args
// 		wantDetail common.DirectoryDetail
// 		wantErr    bool
// 	}{
// 		{
// 			name:  "test_get_directory_detail",
// 			agent: GetAgent().(*WindowsAgent),
// 			args: args{
// 				path: "c:\\test\\directory-1",
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			_, err := tt.agent.GetDirectoryDetail(tt.args.ctx, tt.args.path)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("WindowsAgent.GetDirectoryDetail() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }
