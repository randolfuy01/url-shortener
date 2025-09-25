package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"log"
	server "url-shortener.com/m/server"
)

func main() {
	_ = godotenv.Load()
	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	
	server.GetDB()
	defer server.CloseDB()
	router := gin.Default()
	router.POST("/createUser", server.CreateUser)
	router.POST("/loginUser", server.LoginUser)
	router.POST("/shortenUrl", server.ShortenUrl)
	router.GET("/shortenUrl", server.ShortenUrl)

	router.Run("localhost:8080")
}
