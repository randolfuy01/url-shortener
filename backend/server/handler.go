package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	migrations "url-shortener.com/m/migrations/driver"
	pkg "url-shortener.com/m/pkg"
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
		return
	}

	if len(User.UserName) == 0 || len(User.Password) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Username and password are required",
		})
		return
	}

	// Lookup user by username using sqlc
	queries := ProvideQueries()
	dbUser, err := queries.GetUserByName(context.Background(), User.UserName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid credentials",
			})
			return
		}
		fmt.Printf("LoginUser DB error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database error",
		})
		return
	}

	// Compare password hash
	inputHash := fmt.Sprintf("%x", hashData(User.Password))
	if inputHash != dbUser.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid credentials",
		})
		return
	}

	// Issue JWT token
	token, err := issueJWT(dbUser.ID, dbUser.Name, time.Hour*24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to issue token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"user_id":  dbUser.ID,
		"username": dbUser.Name,
		"token":    token,
	})
}

func CreateUser(c *gin.Context) {
	var User user

	if err := c.BindJSON(&User); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incomplete information",
		})
		return
	}

	if len(User.UserName) == 0 || len(User.Password) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Username and password are required",
		})
		return
	}

	User.Password = fmt.Sprintf("%x", hashData(User.Password))

	// Check if the username is available using sqlc
	queries := ProvideQueries()
	_, err := queries.GetUserByName(context.Background(), User.UserName)
	if err == nil {
		// User exists
		c.JSON(http.StatusConflict, gin.H{
			"message": "Username already taken",
		})
		return
	} else if !errors.Is(err, pgx.ErrNoRows) {
		// Unexpected DB error
		fmt.Printf("CreateUser check DB error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database error",
		})
		return
	}

	// Create the user via sqlc
	created, err := queries.CreateUser(context.Background(), migrations.CreateUserParams{
		Name:     User.UserName,
		Password: User.Password,
	})
	if err != nil {
		fmt.Printf("CreateUser DB error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user",
		})
		return
	}

	// Issue JWT token
	token, err := issueJWT(created.ID, created.Name, time.Hour*24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to issue token",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "User created",
		"user_id":  created.ID,
		"username": created.Name,
		"token":    token,
	})
}

// ShortenUrl is a placeholder for URL shortening functionality
func ShortenUrl(c *gin.Context) {
	// POST /shortenUrl -> create short url
	if c.Request.Method == http.MethodPost {
		var payload struct {
			URL    string  `json:"url"`
			UserID int64   `json:"user_id"`
			Alias  *string `json:"alias"`
		}
		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid payload",
			})
			return
		}

		// Basic validations
		payload.URL = strings.TrimSpace(payload.URL)
		if payload.UserID <= 0 || payload.URL == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "user_id and url are required",
			})
			return
		}

		// URL format validation (very permissive http/https pattern)
		ok, err := pkg.Format_validation(payload.URL, `^https?://.+`)
		if err != nil || !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid url format",
			})
			return
		}

		// Generate short code if not provided
		shortCode := ""
		if payload.Alias != nil && strings.TrimSpace(*payload.Alias) != "" {
			shortCode = strings.TrimSpace(*payload.Alias)
		} else {
			// Use MD5 of URL and take first 8 chars for brevity
			encoded, encErr := pkg.Encode(payload.URL, pkg.Encryption_MD5)
			if encErr != nil || len(encoded) < 8 {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "failed to generate short code",
				})
				return
			}
			shortCode = encoded[:8]
		}

		queries := ProvideQueries()

		created, err := queries.CreateUrl(c, migrations.CreateUrlParams{
			UserID:     payload.UserID,
			OrginalUrl: payload.URL,
			ShortUrl:   shortCode,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to create short url",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":     "short url created",
			"id":          created.ID,
			"user_id":     created.UserID,
			"orginal_url": created.OrginalUrl,
			"short_url":   created.ShortUrl,
		})
		return
	}

	// GET /shortenUrl?user_id=<id> -> list user's urls
	if c.Request.Method == http.MethodGet {
		userIDStr := c.Query("user_id")
		if strings.TrimSpace(userIDStr) == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Missing user_id parameter",
			})
			return
		}
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil || userID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid user_id",
			})
			return
		}

		queries := ProvideQueries()
		urls, err := queries.GetUserURLs(c, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to fetch urls",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
			"urls":    urls,
		})
		return
	}

	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"message": "Method not allowed",
	})
}
