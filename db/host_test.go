package db

import (
	"testing"

	"gorm.io/gorm"
)

func compareHosts(expected []Host, actual []Host) bool {
	if len(expected) != len(actual) {
		return false
	}

	for i := range expected {
		// Compare the desired fields of each struct
		if expected[i].Name != actual[i].Name ||
			expected[i].Ip != actual[i].Ip ||
			expected[i].Username != actual[i].Username ||
			expected[i].StorageType != actual[i].StorageType {
			return false
		}
		// Add more fields for comparison as needed
	}

	return true
}

func Test_Host_Save_Get_Delete(t *testing.T) {
	engine, _ := GetDatabaseEngine()
	host := Host{
		Name:        "test_host_1",
		Ip:          "127.0.0.1",
		Username:    "test_user",
		Password:    "test_password",
		StorageType: "local",
	}

	err := host.Save(engine)
	if err != nil {
		t.Errorf("Error saving host: %v", err)
	}

	retrievedHost := &Host{
		Ip: "127.0.0.1",
	}

	err = retrievedHost.Get(engine)
	if err != nil {
		t.Errorf("Error getting host: %v", err)
	}

	if retrievedHost.Name != host.Name {
		t.Errorf("Expected host name to be '%s', got '%s'", host.Name, retrievedHost.Name)
	}

	if retrievedHost.Ip != host.Ip {
		t.Errorf("Expected host IP to be '%s', got '%s'", host.Ip, retrievedHost.Ip)
	}

	if retrievedHost.Password != host.Password {
		t.Errorf("Expected password to be '%s', got '%s'", host.Password, retrievedHost.Password)
	}

	deletedHost := &Host{
		Name:        "test_host_1",
		Ip:          "127.0.0.1",
		Username:    "test_user",
		StorageType: "local",
	}

	err = deletedHost.Delete(engine)
	if err != nil {
		t.Errorf("Error deleting host: %v", err)
	}

	err = deletedHost.Get(engine)
	if err != gorm.ErrRecordNotFound {
		t.Errorf("Failed to delete host: %v", err)
	}
}

func TestHostList_Create_Get(t *testing.T) {
	engine, _ := GetDatabaseEngine()

	expectedHosts := []Host{
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

	expectedLocalHosts := []Host{
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
	}

	hostList := HostList{
		Hosts: expectedHosts,
	}

	err := hostList.Save(engine)
	if err != nil {
		t.Errorf("Error saving host: %v", err)
	}

	hostList = HostList{}
	if err = hostList.Get(engine, "local"); err != nil {
		t.Errorf("Error get hosts: %v", err)
	}

	if !compareHosts(hostList.Hosts, expectedLocalHosts) {
		t.Errorf("Expected hosts to be '%v', got '%v'", expectedLocalHosts, hostList.Hosts)
	}

	hostList = HostList{}
	if err = hostList.Get(engine, ""); err != nil {
		t.Errorf("Error get hosts: %v", err)
	}

	if !compareHosts(hostList.Hosts, expectedHosts) {
		t.Errorf("Expected hosts to be '%v', got '%v'", expectedHosts, hostList.Hosts)
	}

	if err = hostList.Delete(engine, "", nil, nil); err != nil {
		t.Errorf("Failed to delete directory: %v", err)
	}
}
