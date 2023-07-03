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

	shareName := "NewCIFS_1"

	err := agent.CreateShare(hostContext, shareName, "C:\\test\\folder-2")
	// Assert
	if err != nil {
		t.Errorf("Failed to create share %s. Error: %v", shareName, err)
	}
}

func Test_Windows_GetLocalUsers(t *testing.T) {

	agent := GetAgent()

	hostContext := context.HostContext{
		IP:       "127.09.0.1",
		Username: "admin",
		Password: "password",
	}

	users, err := agent.GetLocalUsers(hostContext)
	if err != nil {
		t.Errorf("Failed to get local users: %s. Error: %v", users, err)
	}
}

func Test_Windows_CreateLocalUsers(t *testing.T) {

	agent := GetAgent()

	hostContext := context.HostContext{
		IP:       "127.09.0.1",
		Username: "admin",
		Password: "password",
	}

	username := "test"
	password := "test"

	err := agent.CreateLocalUser(hostContext, username, password)
	if err != nil {
		t.Errorf("Failed to create local users: %s. Error: %v", username, err)
	}
}

func Test_Windows_DeleteLocalUsers(t *testing.T) {

	agent := GetAgent()

	hostContext := context.HostContext{
		IP:       "127.09.0.1",
		Username: "admin",
		Password: "password",
	}

	username := "test"

	err := agent.DeleteLocalUser(hostContext, username)
	if err != nil {
		t.Errorf("Failed to create local users: %s. Error: %v", username, err)
	}
}
