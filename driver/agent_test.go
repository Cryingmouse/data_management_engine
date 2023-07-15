package driver

import (
	"context"
	"reflect"
	"testing"

	"github.com/cryingmouse/data_management_engine/common"
)

func TestAgentDriver_GetDirectoriesDetail(t *testing.T) {
	type args struct {
		hostContext common.HostContext
		names       []string
		ctx         context.Context
	}
	tests := []struct {
		name       string
		d          *AgentDriver
		args       args
		wantDetail []common.DirectoryDetail
		wantErr    bool
	}{
		{
			name: "Get directories detail success",
			d:    &AgentDriver{},
			args: args{
				hostContext: common.HostContext{
					IP: "localhost",
				},
				names: []string{"Directory-1", "Directory-2"},
				ctx:   nil,
			},
			wantDetail: []common.DirectoryDetail{
				{
					Name:           "Directory-1",
					CreationTime:   "Thursday, July 13, 2023 4:13:06 PM",
					LastAccessTime: "Thursday, July 13, 2023 4:13:16 PM",
					LastWriteTime:  "Thursday, July 13, 2023 4:13:06 PM",
					Exist:          true,
					FullPath:       "C:\\test\\Directory-1",
					ParentFullPath: "C:\\test",
				},
				{
					Name:           "Directory-2",
					CreationTime:   "Thursday, July 13, 2023 4:13:15 PM",
					LastAccessTime: "Thursday, July 13, 2023 4:13:16 PM",
					LastWriteTime:  "Thursday, July 13, 2023 4:13:15 PM",
					Exist:          true,
					FullPath:       "C:\\test\\Directory-2",
					ParentFullPath: "C:\\test",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDetail, err := tt.d.GetDirectoriesDetail(tt.args.ctx, tt.args.names)
			if (err != nil) != tt.wantErr {
				t.Errorf("AgentDriver.GetDirectoriesDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDetail, tt.wantDetail) {
				t.Errorf("AgentDriver.GetDirectoriesDetail() = %v, want %v", gotDetail, tt.wantDetail)
			}
		})
	}
}

func TestAgentDriver_GetDirectoryDetail(t *testing.T) {
	type args struct {
		hostContext common.HostContext
		name        string
		ctx         context.Context
	}
	tests := []struct {
		name       string
		d          *AgentDriver
		args       args
		wantDetail common.DirectoryDetail
		wantErr    bool
	}{
		{
			name: "Get directory detail success",
			d:    &AgentDriver{},
			args: args{
				hostContext: common.HostContext{
					IP: "127.0.0.1",
				},
				name: "Directory-1",
				ctx:  nil,
			},
			wantDetail: common.DirectoryDetail{
				Name:           "Directory-1",
				CreationTime:   "Thursday, July 13, 2023 4:13:06 PM",
				LastAccessTime: "Thursday, July 13, 2023 4:13:16 PM",
				LastWriteTime:  "Thursday, July 13, 2023 4:13:06 PM",
				Exist:          true,
				FullPath:       "C:\\test\\Directory-1",
				ParentFullPath: "C:\\test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDetail, err := tt.d.GetDirectoryDetail(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("AgentDriver.GetDirectoryDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDetail, tt.wantDetail) {
				t.Errorf("AgentDriver.GetDirectoryDetail() = %v, want %v", gotDetail, tt.wantDetail)
			}
		})
	}
}
