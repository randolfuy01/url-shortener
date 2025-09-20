package main

import (
	"github.com/gin-gonic/gin"
	"url-shortener.com/m/server"
)

func main() {
	router := gin.Default()
	router.POST("/createUser", server.CreateUser)
	router.POST("/loginUser", server.LoginUser)

	router.Run("localhost:8080")
}
