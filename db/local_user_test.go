package db

import (
	"testing"
)

func TestLocalUser_Save(t *testing.T) {
	type args struct {
		engine *DatabaseEngine
	}

	_engine, _ := GetDatabaseEngine()
	tests := []struct {
		name    string
		u       *LocalUser
		args    args
		wantErr bool
	}{
		{
			name: "Test Case 1",
			u: &LocalUser{
				Name:     "jay",
				Password: "1234",
				HostName: "computer1",
			},
			args:    args{engine: _engine},
			wantErr: false, // 设置期望的错误结果
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.Save(tt.args.engine); (err != nil) != tt.wantErr {
				t.Errorf("LocalUser.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocalUser_Delete(t *testing.T) {
	type args struct {
		engine *DatabaseEngine
	}
	_engine, _ := GetDatabaseEngine()
	tests := []struct {
		name    string
		u       LocalUser
		args    args
		wantErr bool
	}{
		{
			name: "Test Case 1",
			u: LocalUser{
				Name:     "jay",
				Password: "1234",
				HostName: "computer1",
			},
			args:    args{engine: _engine},
			wantErr: false, // 设置期望的错误结果
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.Delete(tt.args.engine); (err != nil) != tt.wantErr {
				t.Errorf("LocalUser.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
