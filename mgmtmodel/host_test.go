package mgmtmodel

import (
	"testing"

	"github.com/cryingmouse/data_management_engine/db"
)

func TestHost_Register(t *testing.T) {
	// Arrange
	host := Host{
		Name:        "test_host",
		Ip:          "192.168.1.1",
		Username:    "admin",
		Password:    "password",
		StorageType: "s3",
	}

	// Act
	err := host.Register()

	// Assert
	if err != nil {
		t.Errorf("Failed to register host: %v", err)
	}
}

func TestHost_Unregister(t *testing.T) {
	// Arrange
	host := Host{
		Name:        "test_host",
		Ip:          "192.168.1.1",
		Username:    "admin",
		Password:    "password",
		StorageType: "s3",
	}
	if err := host.Register(); err != nil {
		t.Errorf("Failed to register host: %v", err)
	}

	engine, _ := db.GetDatabaseEngine()
	directory1 := db.Directory{
		Name:   "Folder-1",
		HostIp: "192.168.1.1",
	}
	directory2 := db.Directory{
		Name:   "Folder-2",
		HostIp: "192.168.1.1",
	}
	directory1.Save(engine)
	directory2.Save(engine)

	// Act
	if err := host.Unregister(); err != nil {
		t.Errorf("Failed to unregister host: %v", err)
	}
}

func TestHost_Get(t *testing.T) {
	// Arrange
	host := Host{
		Name:        "test_host",
		Ip:          "192.168.1.1",
		Username:    "admin",
		Password:    "password",
		StorageType: "s3",
	}

	// Act
	_, err := host.Get()

	// Assert
	if err != nil {
		t.Errorf("Failed to get host: %v", err)
	}
}

func TestHostList_Get(t *testing.T) {
	// Arrange
	hostList := HostList{}

	// Act
	_, err := hostList.Get("s3")

	// Assert
	if err != nil {
		t.Errorf("Failed to get host list: %v", err)
	}
}
