// admin_handlers.go
package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// func setupAdminRoutes(api *gin.RouterGroup) {
// 	admin := api.Group("/admin")

// 	// Add users (for easy testing/setup)
// 	admin.POST("/users", addUser)

// 	// Get all users
// 	admin.GET("/users", getAllUsers)

// 	// Get service stats
// 	admin.GET("/stats", getServiceStats)
// }

func addUser(c *gin.Context) {
	var reqData struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Bio      string `json:"bio"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if reqData.ID == 0 || reqData.Name == "" || reqData.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID, name, and email are required"})
		return
	}

	// Set default username if not provided
	if reqData.Username == "" {
		reqData.Username = strings.Split(reqData.Email, "@")[0]
	}

	user := User{
		ID:       reqData.ID,
		Name:     reqData.Name,
		Email:    reqData.Email,
		Username: reqData.Username,
		Bio:      reqData.Bio,
		Avatar:   reqData.Avatar,
		IsActive: true,
		LastSeen: time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func getAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	var users []User
	var total int64

	db.Model(&User{}).Count(&total)
	err := db.Where("is_active = ?", true).
		Limit(limit).Offset(offset).Find(&users).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func getServiceStats(c *gin.Context) {
	var stats struct {
		TotalUsers         int64 `json:"total_users"`
		ActiveUsers        int64 `json:"active_users"`
		TotalConversations int64 `json:"total_conversations"`
		TotalMessages      int64 `json:"total_messages"`
		TotalFollows       int64 `json:"total_follows"`
	}

	db.Model(&User{}).Count(&stats.TotalUsers)
	db.Model(&User{}).Where("is_active = ?", true).Count(&stats.ActiveUsers)
	db.Model(&Conversation{}).Count(&stats.TotalConversations)
	db.Model(&Message{}).Count(&stats.TotalMessages)
	db.Model(&Follow{}).Where("status = ?", "accepted").Count(&stats.TotalFollows)

	c.JSON(http.StatusOK, stats)
}
