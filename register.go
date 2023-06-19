package main

import (
	"fmt"
	"net/http"

	"db"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"
)

type HostRegisterInfo struct {
	Name        string `json:"name"`
	Ip          string `json:"ip"`
	Username    string `json:"user_name"`
	Password    string `json:"password"`
	StorageType string `json:"storage_type"`
}

func hostRegistrationHandler(c *gin.Context) {
	var register_info HostRegisterInfo

	type HostInfoWithoutPassword struct {
		Name        string
		Ip          string
		Username    string
		StorageType string
	}

	host_info_without_password := HostInfoWithoutPassword{
		Name:        register_info.Name,
		Ip:          register_info.Ip,
		Username:    register_info.Username,
		StorageType: register_info.StorageType,
	}

	if err := c.ShouldBindJSON(&register_info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	engine.Migrate()

	model, err := engine.GetModel("host_info")
	if err != nil {
		panic(err)
	}

	fmt.Println(model)

	if err = model.(*db.HostInfo).Save(engine); err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			// Map SQLite ErrNo to specific error scenarios
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique: // SQLite constraint violation
				c.JSON(http.StatusInternalServerError, gin.H{"Message": "The host information has already been registered.", "HostRegisterInfo": host_info_without_password})
				return
			default:
				fmt.Println("Error")
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Register the host information successfully.", "HostRegisterInfo": host_info_without_password})
}
