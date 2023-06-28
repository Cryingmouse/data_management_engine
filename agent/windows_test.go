package agent

import (
	"testing"

	"github.com/cryingmouse/data_management_engine/context"
)

func Test_Windows_CreateShare(t *testing.T) {

	agent := GetAgent()

	hostContext := context.HostContext{
		IP:       "127.09.0.1",
		Username: "admin",
		Password: "password",
	}

	shareName, err := agent.CreateShare(hostContext, "NewCIFS_1", "C:\\test\\folder-2")

	// Assert
	if err != nil {
		t.Errorf("Failed to create share %s. Error: %v", shareName, err)
	}
}
