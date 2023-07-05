package driver

import (
	"testing"

	"github.com/cryingmouse/data_management_engine/common"
)

func Test_GetSystemInfo(t *testing.T) {
	driver := &AgentDriver{}

	hostContext := common.HostContext{
		IP:       "127.0.0.2",
		Username: "admin",
		Password: "pasword",
	}

	// Test case 1: Create a new folder successfully.
	response, err := driver.GetSystemInfo(hostContext)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	println(response)
}
