package agent

import (
	"context"
	"fmt"
	"os"

	"github.com/cryingmouse/data_management_engine/common"
)

type LinuxAgent struct {
}

func (agent *LinuxAgent) CreateDirectory(ctx context.Context, name string) (dirPath string, err error) {
	dirPath = fmt.Sprintf("%s\\%s", "C:\\test", name)

	err = os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return dirPath, nil
}

func (agent *LinuxAgent) CreateDirectories(ctx context.Context, names []string) (dirPaths []string, err error) {
	for index, name := range names {
		dirPath, err := agent.CreateDirectory(ctx, name)
		if err != nil {
			return dirPaths, err
		}

		dirPaths[index] = dirPath
	}

	return dirPaths, err
}

func (agent *LinuxAgent) DeleteDirectory(ctx context.Context, name string) (err error) {
	dirPath := fmt.Sprintf("%s\\%s", "C:\\test", name)

	return os.Remove(dirPath)
}

func (agent *LinuxAgent) DeleteDirectories(ctx context.Context, names []string) (err error) {
	for _, name := range names {
		if err = agent.DeleteDirectory(ctx, name); err != nil {
			return err
		}

	}

	return err
}

func (agent *LinuxAgent) GetDirectoryDetail(ctx context.Context, path string) (detail common.DirectoryDetail, err error) {
	return detail, nil
}

func (agent *LinuxAgent) GetDirectoriesDetail(ctx context.Context, paths []string) (detail []common.DirectoryDetail, err error) {
	return detail, nil
}

func (agent *LinuxAgent) CreateCIFSShare(ctx context.Context, name, directoryName, description string, usernames []string) (err error) {
	return err
}

func (agent *LinuxAgent) DeleteCIFSShare(ctx context.Context, name string) (err error) {
	return err
}

func (agent *LinuxAgent) GetCIFSShareDetail(ctx context.Context, name string) (detail common.ShareDetail, err error) {
	return detail, err
}

func (agent *LinuxAgent) GetCIFSSharesDetail(ctx context.Context, names []string) (detail []common.ShareDetail, err error) {
	return detail, err
}

func (agent *LinuxAgent) MountCIFSShare(ctx context.Context, deviceName, sharePath, userName, password string) (err error) {
	return err
}

func (agent *LinuxAgent) UnmountCIFSShare(ctx context.Context, deviceName string) (err error) {
	return err
}

func (agent *LinuxAgent) CreateLocalUser(ctx context.Context, username, password string) (err error) {
	return err
}

func (agent *LinuxAgent) DeleteLocalUser(ctx context.Context, username string) (err error) {
	return err
}

func (agent *LinuxAgent) GetLocalUserDetail(ctx context.Context, username string) (detail common.LocalUserDetail, err error) {
	return detail, nil
}

func (agent *LinuxAgent) GetLocalUsersDetail(ctx context.Context, usernames []string) (detail []common.LocalUserDetail, err error) {
	return detail, nil
}

func (agent *LinuxAgent) GetSystemInfo(ctx context.Context) (system common.SystemInfo, err error) {
	return system, nil
}
