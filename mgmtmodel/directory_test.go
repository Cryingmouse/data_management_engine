package mgmtmodel

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cryingmouse/data_management_engine/common"
	"github.com/stretchr/testify/assert"
)

func setup() {

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
		dl.Save()
	}
}

func cleanup() {
	// Create an empty slice of Directory objects.
	filter := common.QueryFilter{}

	dl := DirectoryList{}
	dl.Delete(&filter)
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

func Test_DirectoryList_Get_With_Fields_Keyword_Condition(t *testing.T) {
	// Create a new Directory object.
	dl := DirectoryList{}

	filter := common.QueryFilter{
		Fields: []string{"Name", "HostIP"},
		Keyword: map[string]string{
			"name": "Directory 1",
		},
		Conditions: Directory{
			HostIP: "127.0.0.1",
		},
	}

	// Get the Directory object from the database.
	_, err := dl.Get(&filter)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(dl.Directories))

	for _, directory := range dl.Directories {
		assert.NotEmpty(t, directory.Name)
		assert.NotEmpty(t, directory.HostIP)
	}
}
