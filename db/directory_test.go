package db

import (
	"testing"
)

func TestDirectorySave(t *testing.T) {
	// Create a new Directory object.
	directory := Directory{
		Name:   "My Directory",
		HostIp: "192.168.1.1",
	}

	// Save the Directory object to the database.
	err := directory.Save(nil)

	// Check if the error is nil.
	if err != nil {
		t.Errorf("Error saving Directory: %v", err)
	}
}

func TestDirectoryGet(t *testing.T) {

	engine, err := GetDatabaseEngine()
	if err != nil {
		panic(err)
	}
	// Create a new Directory object.
	directory := Directory{
		Name:   "My Directory",
		HostIp: "192.168.1.1",
	}

	// Save the Directory object to the database.
	err = directory.Save(engine)

	// Check if the error is nil.
	if err != nil {
		t.Errorf("Error saving Directory: %v", err)
	}

	// Get the Directory object from the database.
	retrievedDirectory, err := directory.Get(engine, "My Directory", "192.168.1.1")

	// Check if the error is nil and the Directory objects are equal.
	if err != nil {
		t.Errorf("Error getting Directory: %v", err)
	} else if *retrievedDirectory != directory {
		t.Errorf("Retrieved Directory is not equal to saved Directory")
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
