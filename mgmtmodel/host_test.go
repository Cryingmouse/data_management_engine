package mgmtmodel

import (
	"testing"

	"github.com/cryingmouse/data_management_engine/db"
)

func Test_Host_Register(t *testing.T) {
	// Arrange
	host := Host{
		Name:        "test_host",
		IP:          "192.168.1.1",
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

func Test_Host_Unregister(t *testing.T) {
	// Arrange
	host := Host{
		Name:        "test_host",
		IP:          "192.168.1.1",
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
		HostIP: "192.168.1.1",
	}
	directory2 := db.Directory{
		Name:   "Folder-2",
		HostIP: "192.168.1.1",
	}
	directory1.Save(engine)
	directory2.Save(engine)

	// Act
	if err := host.Unregister(); err != nil {
		t.Errorf("Failed to unregister host: %v", err)
	}
}

func Test_Host_Get(t *testing.T) {
	// Arrange
	host := Host{
		Name:        "test_host",
		IP:          "192.168.1.1",
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
