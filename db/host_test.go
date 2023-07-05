package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func compareHosts(expected []Host, actual []Host) bool {
	if len(expected) != len(actual) {
		return false
	}

	for i := range expected {
		// Compare the desired fields of each struct
		if expected[i].ComputerName != actual[i].ComputerName ||
			expected[i].IP != actual[i].IP ||
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
		ComputerName: "test_host_1",
		IP:           "127.0.0.1",
		Username:     "test_user",
		Password:     "test_password",
		StorageType:  "local",
	}

	err := host.Save(engine)
	assert.NoError(t, err, "Failed to save host: %v", host)

	retrievedHost := &Host{
		IP: "127.0.0.1",
	}

	err = retrievedHost.Get(engine)
	assert.NoError(t, err, "Failed to get host: %v", retrievedHost)
	assert.Equal(t, host.ComputerName, retrievedHost.ComputerName)
	assert.Equal(t, host.IP, retrievedHost.IP)
	assert.Equal(t, host.Password, retrievedHost.Password)

	deletedHost := &Host{
		ComputerName: "test_host_1",
		IP:           "127.0.0.1",
		Username:     "test_user",
		StorageType:  "local",
	}

	err = deletedHost.Delete(engine)
	assert.NoError(t, err, "Failed to delete the host: %v", deletedHost)

	err = deletedHost.Get(engine)
	assert.NoError(t, err, "Failed to get the host: %v", deletedHost)
}

// func TestHostList_Create_Get(t *testing.T) {
// 	engine, _ := GetDatabaseEngine()

// 	expectedHosts := []Host{
// 		{
// 			Name:        "test_host_1",
// 			IP:          "127.0.0.1",
// 			Username:    "test_user1",
// 			Password:    "test_password1",
// 			StorageType: "local",
// 		},
// 		{
// 			Name:        "test_host_2",
// 			IP:          "127.0.0.2",
// 			Username:    "test_user2",
// 			Password:    "test_password2",
// 			StorageType: "local",
// 		},
// 		{
// 			Name:        "test_host_3",
// 			IP:          "127.0.0.3",
// 			Username:    "test_user3",
// 			Password:    "test_password3",
// 			StorageType: "remote",
// 		},
// 	}

// 	expectedLocalHosts := []Host{
// 		{
// 			Name:        "test_host_1",
// 			IP:          "127.0.0.1",
// 			Username:    "test_user1",
// 			Password:    "test_password1",
// 			StorageType: "local",
// 		},
// 		{
// 			Name:        "test_host_2",
// 			IP:          "127.0.0.2",
// 			Username:    "test_user2",
// 			Password:    "test_password2",
// 			StorageType: "local",
// 		},
// 	}

// 	hostList := HostList{
// 		Hosts: expectedHosts,
// 	}

// 	err := hostList.Save(engine)
// 	assert.NoError(t, err, "Failed to save the host list: %v", hostList)

// 	hostList = HostList{}
// 	err = hostList.Get(engine, "local")
// 	assert.NoError(t, err, "Failed to get host list: %v", hostList)

// 	assert.EqualValues(t, expectedLocalHosts, hostList.Hosts)

// 	if !compareHosts(hostList.Hosts, expectedLocalHosts) {
// 		t.Errorf("Expected hosts to be '%v', got '%v'", expectedLocalHosts, hostList.Hosts)
// 	}

// 	hostList = HostList{}
// 	if err = hostList.Get(engine, ""); err != nil {
// 		t.Errorf("Error get hosts: %v", err)
// 	}

// 	if !compareHosts(hostList.Hosts, expectedHosts) {
// 		t.Errorf("Expected hosts to be '%v', got '%v'", expectedHosts, hostList.Hosts)
// 	}

// 	if err = hostList.Delete(engine, "", nil, nil); err != nil {
// 		t.Errorf("Failed to delete directory: %v", err)
// 	}
// }
