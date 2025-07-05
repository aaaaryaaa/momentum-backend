// routes.go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRoutes(r *gin.Engine) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "chat-service",
			"version": "1.0.0",
		})
	})

	// Protected routes
	api := r.Group("/api/v1")
	api.Use(authMiddleware())

	// User routes
	setupUserRoutes(api)

	// Follow routes
	setupFollowRoutes(api)

	// Conversation routes
	setupConversationRoutes(api)

	// Message routes
	setupMessageRoutes(api)

	// Admin routes
	setupAdminRoutes(api)
}

func setupUserRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")

	// Get current user profile
	users.GET("/me", getCurrentUser)

	// Update current user profile
	users.PUT("/me", updateCurrentUser)

	// Search users
	users.GET("/search", searchUsers)

	// Get user by ID
	users.GET("/:id", getUserByID)

	// Get user followers
	users.GET("/:id/followers", getUserFollowers)

	// Get user following
	users.GET("/:id/following", getUserFollowing)
}

func setupFollowRoutes(api *gin.RouterGroup) {
	follows := api.Group("/follows")

	// Follow a user
	follows.POST("/:id", followUser)

	// Unfollow a user
	follows.DELETE("/:id", unfollowUser)

	// Get follow requests
	follows.GET("/requests", getFollowRequests)

	// Accept follow request
	follows.PUT("/requests/:id/accept", acceptFollowRequest)

	// Reject follow request
	follows.PUT("/requests/:id/reject", rejectFollowRequest)
}

func setupConversationRoutes(api *gin.RouterGroup) {
	conversations := api.Group("/conversations")

	// Get all conversations
	conversations.GET("", getConversations)

	// Create new conversation
	conversations.POST("", createConversation)

	// Get conversation by ID
	conversations.GET("/:id", getConversation)

	// Get conversation messages
	conversations.GET("/:id/messages", getConversationMessages)

	// Add participants to conversation
	conversations.POST("/:id/participants", addParticipants)

	// Remove participant from conversation
	conversations.DELETE("/:id/participants/:user_id", removeParticipant)
}

func setupMessageRoutes(api *gin.RouterGroup) {
	messages := api.Group("/messages")

	// Send message
	messages.POST("", sendMessage)

	// Get message by ID
	messages.GET("/:id", getMessage)

	// Edit message
	messages.PUT("/:id", editMessage)

	// Delete message
	messages.DELETE("/:id", deleteMessage)

	// Mark message as read
	messages.PUT("/:id/read", markMessageAsRead)
}

func setupAdminRoutes(api *gin.RouterGroup) {
	admin := api.Group("/admin")

	// Add users (for easy testing/setup)
	admin.POST("/users", addUser)

	// Get all users
	admin.GET("/users", getAllUsers)

	// Get service stats
	admin.GET("/stats", getServiceStats)
}
