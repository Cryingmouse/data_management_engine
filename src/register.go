package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

	if err := c.ShouldBindJSON(&register_info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

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

	db, err := gorm.Open(sqlite.Open("./sqlite3.db"), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Failed to open database."})
		return
	}

	type HostInfo struct {
		Name        string
		Ip          string `gorm:"unique"`
		Username    string
		Password    string
		StorageType string
	}

	encrypted_password, err := bcrypt.GenerateFromPassword([]byte(register_info.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	encrypted_password_str := string(encrypted_password)

	// Create a new variable of the updated struct type
	host_info := HostInfo{
		Name:        register_info.Name,
		Ip:          register_info.Ip,
		Username:    register_info.Username,
		Password:    encrypted_password_str,
		StorageType: register_info.StorageType,
	}

	// Migrate the schema (create the 'users' table if it doesn't exist)
	db.AutoMigrate(&HostInfo{})

	// Insert the user into the database
	if err = db.Create(&host_info).Error; err != nil {
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
