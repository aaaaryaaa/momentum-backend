// user_handlers.go
package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getCurrentUser(c *gin.Context) {
	user := c.MustGet("user").(*User)
	c.JSON(http.StatusOK, user)
}

func updateCurrentUser(c *gin.Context) {
	user := c.MustGet("user").(*User)

	var updateData struct {
		Name      string `json:"name"`
		Username  string `json:"username"`
		Bio       string `json:"bio"`
		Avatar    string `json:"avatar"`
		IsPrivate bool   `json:"is_private"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update user fields
	if updateData.Name != "" {
		user.Name = updateData.Name
	}
	if updateData.Username != "" {
		user.Username = updateData.Username
	}
	if updateData.Bio != "" {
		user.Bio = updateData.Bio
	}
	if updateData.Avatar != "" {
		user.Avatar = updateData.Avatar
	}
	user.IsPrivate = updateData.IsPrivate

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func searchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var users []User
	err := db.Where("name ILIKE ? OR username ILIKE ? OR email ILIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%").
		Where("is_active = ?", true).
		Limit(limit).Offset(offset).Find(&users).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"page":  page,
		"limit": limit,
	})
}

func getUserByID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user User
	err = db.Where("id = ? AND is_active = ?", userID, true).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func getUserFollowers(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var follows []Follow
	err = db.Preload("Follower").
		Where("following_id = ? AND status = ?", userID, "accepted").
		Limit(limit).Offset(offset).Find(&follows).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get followers"})
		return
	}

	followers := make([]User, len(follows))
	for i, follow := range follows {
		followers[i] = follow.Follower
	}

	c.JSON(http.StatusOK, gin.H{
		"followers": followers,
		"page":      page,
		"limit":     limit,
	})
}

func getUserFollowing(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var follows []Follow
	err = db.Preload("Following").
		Where("follower_id = ? AND status = ?", userID, "accepted").
		Limit(limit).Offset(offset).Find(&follows).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get following"})
		return
	}

	following := make([]User, len(follows))
	for i, follow := range follows {
		following[i] = follow.Following
	}

	c.JSON(http.StatusOK, gin.H{
		"following": following,
		"page":      page,
		"limit":     limit,
	})
}
