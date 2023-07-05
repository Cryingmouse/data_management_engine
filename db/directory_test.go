package db

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cryingmouse/data_management_engine/common"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setup() {
	engine, _ := GetDatabaseEngine()

	// Define the number of Directory objects you want to create.
	numDirectories := 10

	// Generate and append Directory objects to the slice.
	for _, hostIP := range []string{"127.0.0.1", "127.0.0.2"} {
		dl := DirectoryList{}
		directories := []Directory{}

		for i := 0; i < numDirectories; i++ {
			directory := Directory{
				Name: fmt.Sprintf("Directory %d", i),
				// HostIP: fmt.Sprintf("192.168.1.%d", i),
				HostIP: hostIP,
			}
			directories = append(directories, directory)
		}
		dl.Directories = directories

		dl.Save(engine)
	}
}

func cleanup() {
	engine, _ := GetDatabaseEngine()

	// Create an empty slice of Directory objects.
	filter := common.QueryFilter{}

	dl := DirectoryList{}
	dl.Delete(engine, &filter)
}

func TestMain(m *testing.M) {
	// 获取当前文件所在的目录
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	// 切换到项目的根目录
	projectPath := filepath.Join(dir, "../") // 假设项目的根目录在当前目录的上一级目录
	err := os.Chdir(projectPath)
	if err != nil {
		panic(err)
	}

	setup()
	// 在所有测试运行之后执行清理操作
	defer cleanup()

	// 执行测试
	exitCode := m.Run()

	// 退出测试
	os.Exit(exitCode)
}

func Test_Directory_Save_Get_Delete(t *testing.T) {
	engine, _ := GetDatabaseEngine()

	// Create a new Directory object.
	directory := Directory{
		Name:   "My Directory",
		HostIP: "192.168.1.1",
	}
	// Save the Directory object to the database.
	err := directory.Save(engine)
	assert.NoError(t, err)

	// Get the Directory object from the database.
	retrievedDirectory := Directory{
		Name: "My Directory",
	}
	err = retrievedDirectory.Get(engine)
	assert.NoError(t, err)

	assert.Equal(t, directory.Name, retrievedDirectory.Name)

	assert.Equal(t, directory.HostIP, retrievedDirectory.HostIP)

	deletedDirectory := &Directory{
		Name: "My Directory",
	}

	err = deletedDirectory.Delete(engine)
	assert.NoError(t, err)

	err = deletedDirectory.Get(engine)
	assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

func Test_DirectoryList_Get_With_Fields_Keyword_Condition(t *testing.T) {
	engine, err := GetDatabaseEngine()
	if err != nil {
		panic(err)
	}
	// Create a new Directory object.
	dl := DirectoryList{}

	filter := common.QueryFilter{
		Fields: []string{"Name"},
		Keyword: map[string]string{
			"name": "Directory 1",
		},
		Conditions: Directory{
			HostIP: "127.0.0.1",
		},
	}

	// Get the Directory object from the database.
	err = dl.Get(engine, &filter)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(dl.Directories))

	for _, directory := range dl.Directories {
		assert.NotEmpty(t, directory.Name)
		assert.Empty(t, directory.HostIP)
	}
}

func Test_DirectoryList_Get_With_Fields_Condition(t *testing.T) {
	engine, err := GetDatabaseEngine()
	if err != nil {
		panic(err)
	}
	// Create a new Directory object.
	dl := DirectoryList{}

	filter := common.QueryFilter{
		Fields: []string{"Name"},
		Conditions: Directory{
			HostIP: "127.0.0.1",
		},
	}

	// Get the Directory object from the database.
	err = dl.Get(engine, &filter)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(dl.Directories))

	for _, directory := range dl.Directories {
		assert.NotEmpty(t, directory.Name)
		assert.Empty(t, directory.HostIP)
	}
}

func Test_DirectoryList_Get_With_Fields_Keyword(t *testing.T) {
	engine, err := GetDatabaseEngine()
	if err != nil {
		panic(err)
	}
	// Create a new Directory object.
	dl := DirectoryList{}

	filter := common.QueryFilter{
		Fields: []string{"Name"},
		Keyword: map[string]string{
			"name": "Directory 1",
		},
	}

	// Get the Directory object from the database.
	err = dl.Get(engine, &filter)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(dl.Directories))

	for _, directory := range dl.Directories {
		assert.NotEmpty(t, directory.Name)
		assert.Empty(t, directory.HostIP)
	}
}

func Test_DirectoryList_Get_With_Keyword_Condition(t *testing.T) {
	engine, err := GetDatabaseEngine()
	if err != nil {
		panic(err)
	}
	// Create a new Directory object.
	dl := DirectoryList{}

	filter := common.QueryFilter{
		Keyword: map[string]string{
			"name": "Directory 1",
		},
		Conditions: Directory{
			HostIP: "127.0.0.1",
		},
	}

	// Get the Directory object from the database.
	err = dl.Get(engine, &filter)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(dl.Directories))

	for _, directory := range dl.Directories {
		assert.NotEmpty(t, directory.Name)
		assert.NotEmpty(t, directory.HostIP)
	}
}

func Test_DirectoryList_Get_With_Condition(t *testing.T) {
	engine, err := GetDatabaseEngine()
	if err != nil {
		panic(err)
	}
	// Create a new Directory object.
	dl := DirectoryList{}

	filter := common.QueryFilter{
		Conditions: Directory{
			HostIP: "127.0.0.1",
		},
	}

	// Get the Directory object from the database.
	err = dl.Get(engine, &filter)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(dl.Directories))

	for _, directory := range dl.Directories {
		assert.NotEmpty(t, directory.Name)
		assert.NotEmpty(t, directory.HostIP)
	}
}

func Test_DirectoryList_Get_With_Fields(t *testing.T) {
	engine, err := GetDatabaseEngine()
	if err != nil {
		panic(err)
	}
	// Create a new Directory object.
	dl := DirectoryList{}

	filter := common.QueryFilter{
		Fields: []string{"Name"},
	}

	// Get the Directory object from the database.
	err = dl.Get(engine, &filter)
	assert.NoError(t, err)
	assert.Equal(t, 20, len(dl.Directories))

	for _, directory := range dl.Directories {
		assert.NotEmpty(t, directory.Name)
		assert.Empty(t, directory.HostIP)
	}
}

func Test_DirectoryList_Get_With_Keyword(t *testing.T) {
	engine, err := GetDatabaseEngine()
	if err != nil {
		panic(err)
	}
	// Create a new Directory object.
	dl := DirectoryList{}

	filter := common.QueryFilter{
		Keyword: map[string]string{
			"name": "Directory 1",
		},
	}

	// Get the Directory object from the database.
	err = dl.Get(engine, &filter)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(dl.Directories))

	for _, directory := range dl.Directories {
		assert.NotEmpty(t, directory.Name)
		assert.NotEmpty(t, directory.HostIP)
	}
}

func Test_DirectoryList_Get(t *testing.T) {
	engine, err := GetDatabaseEngine()
	if err != nil {
		panic(err)
	}
	// Create a new Directory object.
	dl := DirectoryList{}

	filter := common.QueryFilter{}

	// Get the Directory object from the database.
	err = dl.Get(engine, &filter)
	assert.NoError(t, err)
	assert.Equal(t, 20, len(dl.Directories))

	for _, directory := range dl.Directories {
		assert.NotEmpty(t, directory.Name)
		assert.NotEmpty(t, directory.HostIP)
	}
}

func Test_DirectoryList_Delete_With_Condition(t *testing.T) {
	engine, err := GetDatabaseEngine()
	assert.NoError(t, err)

	// Create an empty slice of Directory objects.
	filter := common.QueryFilter{
		Conditions: Directory{
			HostIP: "127.0.0.1",
		},
	}

	dl := DirectoryList{}
	dl.Delete(engine, &filter)

	filter = common.QueryFilter{
		Fields: []string{"HostIP", "Name"},
	}

	err = dl.Get(engine, &filter)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(dl.Directories))
}

func Test_DirectoryList_Pagination_With_Fields_Keyword_Condition(t *testing.T) {
	engine, err := GetDatabaseEngine()
	assert.NoError(t, err)

	filter := common.QueryFilter{
		Fields: []string{"HostIP", "Name"},
		Keyword: map[string]string{
			"name": "Directory 1",
		},
		Pagination: &common.Pagination{
			Page:     1,
			PageSize: 2,
		},
		Conditions: Directory{
			HostIP: "127.0.0.2",
		},
	}

	dl := DirectoryList{}
	pagination_directories, err := dl.Pagination(engine, &filter)
	assert.NoError(t, err)

	var expectedTotalCount int64 = 1
	assert.Equal(t, expectedTotalCount, pagination_directories.TotalCount)

	assert.Equal(t, 1, len(pagination_directories.Directories))
}
