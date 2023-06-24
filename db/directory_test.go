package db

import (
	"testing"

	"gorm.io/gorm"
)

func TestDirectory_Save_Get_Delete(t *testing.T) {
	engine, _ := GetDatabaseEngine()

	// Create a new Directory object.
	directory := Directory{
		Name:   "My Directory",
		HostIp: "192.168.1.1",
	}
	// Save the Directory object to the database.
	err := directory.Save(engine)
	// Check if the error is nil.
	if err != nil {
		t.Errorf("Error saving Directory: %v", err)
	}

	// Get the Directory object from the database.
	retrievedDirectory := Directory{
		Name: "My Directory",
	}
	err = retrievedDirectory.Get(engine)
	// Check if the error is nil and the Directory objects are equal.
	if err != nil {
		t.Errorf("Error getting Directory: %v", err)
	}

	if retrievedDirectory.Name != directory.Name {
		t.Errorf("Expected directory name to be '%s', got '%s'", directory.Name, retrievedDirectory.Name)
	}

	if retrievedDirectory.HostIp != directory.HostIp {
		t.Errorf("Expected directory hostIp to be '%s', got '%s'", directory.HostIp, retrievedDirectory.HostIp)
	}

	deletedDirectory := &Directory{
		Name: "My Directory",
	}

	err = deletedDirectory.Delete(engine)
	if err != nil {
		t.Errorf("Error deleting directory: %v", err)
	}

	err = deletedDirectory.Get(engine)
	if err != gorm.ErrRecordNotFound {
		t.Errorf("Failed to delete directory: %v", err)
	}

}

func TestDirectoryListGet(t *testing.T) {

	engine, err := GetDatabaseEngine()
	if err != nil {
		panic(err)
	}
	// Create a new Directory object.
	directoryList := DirectoryList{}

	// Get the Directory object from the database.
	result, err := directoryList.Get(engine, "127.0.0.1")

	// Check if the error is nil and the Directory objects are equal.
	if err != nil {
		t.Errorf("Error getting Directory: %v", err)
	} else {
		t.Log(result)
	}
}
