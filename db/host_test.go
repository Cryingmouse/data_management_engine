package db

import (
	"testing"

	"github.com/cryingmouse/data_management_engine/context"
	"github.com/cryingmouse/data_management_engine/utils"
)

func Test_Host_Save_Get_Delete(t *testing.T) {
	engine, _ := GetDatabaseEngine()
	host := Host{
		Name:        "test_host_2",
		Ip:          "127.0.0.2",
		Username:    "test_user",
		Password:    "test_password",
		StorageType: "local",
	}

	err := host.Save(engine)
	if err != nil {
		t.Errorf("Error saving host: %v", err)
	}

	newHost, err := host.Get(engine, host.Name, host.Ip)
	if err != nil {
		t.Errorf("Error getting host: %v", err)
	}

	if newHost.Name != host.Name {
		t.Errorf("Expected host name to be '%s', got '%s'", host.Name, newHost.Name)
	}

	if newHost.Ip != host.Ip {
		t.Errorf("Expected host IP to be '%s', got '%s'", host.Ip, newHost.Ip)
	}

	password, err := utils.Decrypt(newHost.Password, context.SecurityKey)
	if err != nil {
		t.Errorf("Error decrypting password: %v", err)
	}

	if password != host.Password {
		t.Errorf("Expected password to be '%s', got '%s'", host.Password, password)
	}

	newHost.Delete(engine)
}

func TestHostList_Create_Get(t *testing.T) {
	engine, _ := GetDatabaseEngine()

	hosts := []Host{
		{
			Name:        "test_host_1",
			Ip:          "127.0.0.1",
			Username:    "test_user1",
			Password:    "test_password1",
			StorageType: "local",
		},
		{
			Name:        "test_host_2",
			Ip:          "127.0.0.2",
			Username:    "test_user2",
			Password:    "test_password2",
			StorageType: "local",
		},
		{
			Name:        "test_host_3",
			Ip:          "127.0.0.3",
			Username:    "test_user3",
			Password:    "test_password3",
			StorageType: "remote",
		},
	}

	hostList := HostList{
		Hosts: hosts,
	}

	err := hostList.Create(engine)
	if err != nil {
		t.Errorf("Error saving host: %v", err)
	}

	newHostList := HostList{}
	hosts, _ = newHostList.Get(engine, "local")

	newHostList.Delete(engine, "local", nil, nil)

	hosts, _ = newHostList.Get(engine, "")

	t.Log("Test Completed.")
}
