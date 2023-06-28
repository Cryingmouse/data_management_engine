package db

import (
	"encoding/json"
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

func Test_DirectoryList_Get(t *testing.T) {
	engine, err := GetDatabaseEngine()
	if err != nil {
		panic(err)
	}
	// Create a new Directory object.
	directoryList := DirectoryList{}

	// Get the Directory object from the database.
	if err := directoryList.Get(engine, "", "Directory 1"); err != nil {
		t.Errorf("Error getting Directory: %v", err)
	}

	dirs, _ := json.Marshal(directoryList.Directories)

	t.Log(string(dirs))
}

func Test_DirectoryList_QueryByPagination(t *testing.T) {
	engine, err := GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	// // Create an empty slice of Directory objects.
	// directories := []Directory{}

	// // Define the number of Directory objects you want to create.
	// numDirectories := 50

	// // Generate and append Directory objects to the slice.
	// for i := 1; i <= numDirectories; i++ {
	// 	directory := Directory{
	// 		Name:   fmt.Sprintf("Directory %d", i),
	// 		HostIp: fmt.Sprintf("192.168.1.%d", i),
	// 	}
	// 	directories = append(directories, directory)
	// }

	// // Save the Directory objects to the database.
	// for _, directory := range directories {
	// 	err := directory.Save(engine)
	// 	// Check if the error is nil.
	// 	if err != nil {
	// 		t.Errorf("Error saving Directory: %v", err)
	// 	}
	// }

	dl := DirectoryList{}
	pagination_directories, err1 := dl.GetByPagination(engine, []string{}, "", 2, 10)
	if err != nil {
		t.Errorf("Error pagination directories: %v", err1)
	}

	dirs, _ := json.Marshal(pagination_directories)
	t.Log(string(dirs))
}
