package main

import (
	"github.com/gin-gonic/gin"
	server "url-shortener.com/m/server"
)

func main() {
	router := gin.Default()
	router.POST("/createUser", server.CreateUser)
	router.POST("/loginUser", server.LoginUser)
	router.POST("/shortenUrl", server.ShortenUrl)
	router.GET("/shortenUrl", server.ShortenUrl)

	router.Run("localhost:8080")
}
