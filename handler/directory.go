package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cryingmouse/data_management_engine/context"
	"github.com/cryingmouse/data_management_engine/db"
	"github.com/cryingmouse/data_management_engine/driver"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"
)

type DirectoryInfo struct {
	Name   string `json:"name"`
	HostIp string `json:"host_ip"`
	// HostName might not be necessary here, but this is for hostModel.Get()
	HostName string `json:"host_name"`
}

func CreateDirectoryHandler(c *gin.Context) {
	var directoryInfo DirectoryInfo
	if err := c.ShouldBindJSON(&directoryInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	hostModel := db.Host{}
	host, err := hostModel.Get(engine, directoryInfo.HostName, directoryInfo.HostIp)
	if err != nil {
		panic(err)
	}

	hostContext := context.HostContext{
		IP:       host.Ip,
		Username: host.Username,
		Password: host.Password,
	}

	driver := driver.GetDriver(host.StorageType)
	driver.CreateDirectory(hostContext, directoryInfo.Name)

	directoryModel := db.Directory{
		Name:   directoryInfo.Name,
		HostIp: host.Ip,
	}

	if err = directoryModel.Save(engine); err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			// Map SQLite ErrNo to specific error scenarios
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique: // SQLite constraint violation
				c.JSON(http.StatusInternalServerError, gin.H{"Message": "Failed to create the directory."})
				return
			default:
				fmt.Println("Error")
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Create the directory successfully."})
}

type AgentDirectoryInfo struct {
	Name string `json:"name"`
}

func CreateDirectoryOnAgentHandler(c *gin.Context) {
	var directoryInfo AgentDirectoryInfo
	if err := c.ShouldBindJSON(&directoryInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := os.Mkdir(fmt.Sprintf("%s\\%s", "c:\\test\\", directoryInfo.Name), os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Create directory on agent successfully."})
}
