// auth.go
package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func validateJWT(token string) (*AuthUser, error) {
	token = strings.TrimPrefix(token, "Bearer ")

	req, err := http.NewRequest("GET", os.Getenv("AUTH_SERVICE_URL")+"/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.New("auth service verification failed")
	}
	defer resp.Body.Close()

	var user AuthUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		authUser, err := validateJWT(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Get or create user in chat service
		user, err := getOrCreateUser(authUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			c.Abort()
			return
		}

		// Update last seen
		db.Model(&user).Update("last_seen", time.Now())

		c.Set("user", user)
		c.Next()
	}
}

func getOrCreateUser(authUser *AuthUser) (*User, error) {
	var user User
	err := db.Where("id = ?", authUser.ID).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		// Create new user
		user = User{
			ID:       authUser.ID,
			Name:     authUser.Name,
			Email:    authUser.Email,
			Username: strings.Split(authUser.Email, "@")[0], // Default username from email
			IsActive: true,
			LastSeen: time.Now(),
		}

		if err := db.Create(&user).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
