package server

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type user struct {
	ID       *string `json:"id"`
	UserName string  `json:"username"`
	Password string  `json:"password"`
}

func LoginUser(c *gin.Context) {
	var User user

	if err := c.BindJSON(&User); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incomplete information",
		})
	}

	// Check if the user exists

	// Add json web token

	c.JSON(http.StatusOK, gin.H{
		"message": "Login User",
	})
}

func CreateUser(c *gin.Context) {
	var User user

	if err := c.BindJSON(&User); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incomplete information",
		})
	}

	User.Password = fmt.Sprintf("%x", hashData(User.Password))

	// Check if the username is available

	// Some type of database logic

	// Add json web token

	c.JSON(http.StatusOK, gin.H{
		"message": "Create User",
	})
}

func hashData(data string) [32]byte {
	sum := sha256.Sum256([]byte(data))
	return sum
}
