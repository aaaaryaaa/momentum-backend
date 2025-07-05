// follow_handlers.go
package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func followUser(c *gin.Context) {
	currentUser := c.MustGet("user").(*User)
	followingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user exists
	var targetUser User
	err = db.Where("id = ? AND is_active = ?", followingID, true).First(&targetUser).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}

	// Can't follow yourself
	if currentUser.ID == followingID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
		return
	}

	// Check if already following
	var existingFollow Follow
	err = db.Where("follower_id = ? AND following_id = ?", currentUser.ID, followingID).First(&existingFollow).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Already following this user"})
		return
	}

	// Create follow record
	status := "accepted"
	if targetUser.IsPrivate {
		status = "pending"
	}

	follow := Follow{
		FollowerID:  currentUser.ID,
		FollowingID: followingID,
		Status:      status,
		CreatedAt:   time.Now(),
	}

	if err := db.Create(&follow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Follow request sent",
		"status":  status,
	})
}

func unfollowUser(c *gin.Context) {
	currentUser := c.MustGet("user").(*User)
	followingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = db.Where("follower_id = ? AND following_id = ?", currentUser.ID, followingID).Delete(&Follow{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unfollowed successfully"})
}

func getFollowRequests(c *gin.Context) {
	currentUser := c.MustGet("user").(*User)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var follows []Follow
	err := db.Preload("Follower").
		Where("following_id = ? AND status = ?", currentUser.ID, "pending").
		Limit(limit).Offset(offset).Find(&follows).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get follow requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"requests": follows,
		"page":     page,
		"limit":    limit,
	})
}

func acceptFollowRequest(c *gin.Context) {
	currentUser := c.MustGet("user").(*User)
	followID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid follow ID"})
		return
	}

	var follow Follow
	err = db.Where("id = ? AND following_id = ? AND status = ?", followID, currentUser.ID, "pending").First(&follow).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Follow request not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find follow request"})
		return
	}

	follow.Status = "accepted"
	if err := db.Save(&follow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept follow request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Follow request accepted"})
}

func rejectFollowRequest(c *gin.Context) {
	currentUser := c.MustGet("user").(*User)
	followID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid follow ID"})
		return
	}

	err = db.Where("id = ? AND following_id = ? AND status = ?", followID, currentUser.ID, "pending").Delete(&Follow{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject follow request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Follow request rejected"})
}
