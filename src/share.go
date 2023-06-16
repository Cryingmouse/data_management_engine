package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func protectedHandler(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{"message": "Access granted", "username": username})
}
