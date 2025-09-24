package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	migrations "url-shortener.com/m/migrations/driver"
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

	ensureDB()

	var (
		id         int64
		storedHash string
	)

	// Lookup user by username
	row := pool.QueryRow(context.Background(), "SELECT id, password FROM users WHERE name = $1 LIMIT 1", User.UserName)
	if err := row.Scan(&id, &storedHash); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid credentials",
		})
		return
	}

	// Compare password hash
	inputHash := fmt.Sprintf("%x", hashData(User.Password))
	if inputHash != storedHash {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid credentials",
		})
		return
	}

	// TODO: Issue JWT token
	token, err := issueJWT(id, User.UserName, time.Hour*24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to issue token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"user_id":  id,
		"username": User.UserName,
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

	ensureDB()

	// Check if the username is available
	var exists int
	row := pool.QueryRow(context.Background(), "SELECT 1 FROM users WHERE name = $1 LIMIT 1", User.UserName)
	if err := row.Scan(&exists); err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Username already taken",
		})
		return
	}

	// Create the user via sqlc
	queries := migrations.New(pool)
	created, err := queries.CreateUser(context.Background(), migrations.CreateUserParams{
		Name:     User.UserName,
		Password: User.Password,
	})
	if err != nil {
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
	// For now this handler only documents the expected flow for both POST (create)
	// and GET (expand/resolve) without touching database logic.

	// POST /shortenUrl
	// Expected JSON payload:
	// {
	//   "url": "https://example.com/some/long/path?with=query",
	//   "custom_alias": "optional-custom",
	//   "expiry_seconds": 604800
	// }
	// Steps:
	// 1) Parse and validate payload (required url, optional alias, optional ttl)
	// 2) Validate URL format using pkg.Format_validation with a URL regex
	// 3) Normalize URL (ensure scheme, trim spaces)
	// 4) (Optional) Authenticate user via JWT in Authorization header Bearer <token>
	// 5) (Optional) Apply rate limiting per user/IP
	// 6) If custom_alias provided: validate allowed charset/length; check availability
	// 7) If no alias: generate short code (e.g., base62) with collision checks
	// 8) Persist mapping: short code -> original URL (+ owner, expiry, created_at)
	// 9) Return 201 with JSON { short_url, code, expiry, original_url }

	// GET /shortenUrl?code=<shortCode>
	// Steps:
	// 1) Read query param "code"; validate format
	// 2) Lookup original URL by code
	// 3) Check expiry/disabled flags
	// 4) (Optional) Record click stats: timestamp, referer, user-agent, ip, geo
	// 5) Redirect with 301/302 to original URL

	if c.Request.Method == http.MethodPost {
		// Lightweight scaffold for payload parsing without implementation
		var payload struct {
			URL           string  `json:"url"`
			CustomAlias   *string `json:"custom_alias"`
			ExpirySeconds *int64  `json:"expiry_seconds"`
		}
		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid payload",
			})
			return
		}

		// NOTE: Implementation intentionally omitted; see steps above.
		c.JSON(http.StatusNotImplemented, gin.H{
			"message":        "Shorten URL not implemented",
			"received_url":   payload.URL,
			"custom_alias":   payload.CustomAlias,
			"expiry_seconds": payload.ExpirySeconds,
		})
		return
	}

	if c.Request.Method == http.MethodGet {
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Missing code parameter",
			})
			return
		}

		// NOTE: Implementation intentionally omitted; see steps above.
		c.JSON(http.StatusNotImplemented, gin.H{
			"message": "Expand URL not implemented",
			"code":    code,
		})
		return
	}

	// For other methods, respond with 405
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"message": "Method not allowed",
	})
}
