// conversation_handlers.go
package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getConversations(c *gin.Context) {
	currentUser := c.MustGet("user").(*User)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	var conversations []Conversation
	err := db.Preload("Participants").
		Preload("LastMessage").
		Preload("LastMessage.Sender").
		Joins("JOIN conversation_participants ON conversations.id = conversation_participants.conversation_id").
		Where("conversation_participants.user_id = ?", currentUser.ID).
		Order("conversations.updated_at DESC").
		Limit(limit).Offset(offset).Find(&conversations).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conversations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conversations": conversations,
		"page":          page,
		"limit":         limit,
	})
}

func createConversation(c *gin.Context) {
	currentUser := c.MustGet("user").(*User)

	var reqData struct {
		Type           string `json:"type"`            // "direct" or "group"
		Name           string `json:"name"`            // Required for group chats
		ParticipantIDs []int  `json:"participant_ids"` // User IDs to include
	}

	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default to direct conversation
	if reqData.Type == "" {
		reqData.Type = "direct"
	}

	// Validate conversation type
	if reqData.Type != "direct" && reqData.Type != "group" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation type"})
		return
	}

	// For group conversations, name is required
	if reqData.Type == "group" && reqData.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Group conversation name is required"})
		return
	}

	// For direct conversations, exactly one other participant is required
	if reqData.Type == "direct" && len(reqData.ParticipantIDs) != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Direct conversation requires exactly one other participant"})
		return
	}

	// Check if direct conversation already exists
	if reqData.Type == "direct" {
		var existingConversation Conversation
		err := db.Joins("JOIN conversation_participants cp1 ON conversations.id = cp1.conversation_id").
			Joins("JOIN conversation_participants cp2 ON conversations.id = cp2.conversation_id").
			Where("conversations.type = ? AND cp1.user_id = ? AND cp2.user_id = ?",
				"direct", currentUser.ID, reqData.ParticipantIDs[0]).
			First(&existingConversation).Error

		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Direct conversation already exists"})
			return
		}
	}

	// Get participants (excluding current user if already in list)
	var participants []User
	participantIDs := reqData.ParticipantIDs

	// Add current user to participants if not already included
	found := false
	for _, id := range participantIDs {
		if id == currentUser.ID {
			found = true
			break
		}
	}
	if !found {
		participantIDs = append(participantIDs, currentUser.ID)
	}

	err := db.Where("id IN ? AND is_active = ?", participantIDs, true).Find(&participants).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find participants"})
		return
	}

	if len(participants) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least 2 participants required"})
		return
	}

	// Create conversation
	conversation := Conversation{
		Type:      reqData.Type,
		Name:      reqData.Name,
		CreatedBy: currentUser.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&conversation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
		return
	}

	// Add participants
	if err := db.Model(&conversation).Association("Participants").Append(participants); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add participants"})
		return
	}

	// Load conversation with participants
	db.Preload("Participants").First(&conversation, conversation.ID)

	c.JSON(http.StatusCreated, conversation)
}

func getConversation(c *gin.Context) {
	currentUser := c.MustGet("user").(*User)
	conversationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	var conversation Conversation
	err = db.Preload("Participants").
		Preload("LastMessage").
		Preload("LastMessage.Sender").
		Joins("JOIN conversation_participants ON conversations.id = conversation_participants.conversation_id").
		Where("conversations.id = ? AND conversation_participants.user_id = ?", conversationID, currentUser.ID).
		First(&conversation).Error

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found or access denied"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conversation"})
		return
	}

	c.JSON(http.StatusOK, conversation)
}

func getConversationMessages(c *gin.Context) {
	currentUser := c.MustGet("user").(*User)
	conversationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	// Check if user is participant
	var conversation Conversation
	err = db.Joins("JOIN conversation_participants ON conversations.id = conversation_participants.conversation_id").
		Where("conversations.id = ? AND conversation_participants.user_id = ?", conversationID, currentUser.ID).
		First(&conversation).Error

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found or access denied"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify conversation access"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	var messages []Message
	err = db.Preload("Sender").
		Preload("ReplyTo").
		Preload("ReplyTo.Sender").
		Where("conversation_id = ?", conversationID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&messages).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"page":     page,
		"limit":    limit,
	})
}

func addParticipants(c *gin.Context) {
	currentUser := c.MustGet("user").(*User)
	conversationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	var reqData struct {
		UserIDs []int `json:"user_ids"`
	}

	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user is participant and conversation is group
	var conversation Conversation
	err = db.Joins("JOIN conversation_participants ON conversations.id = conversation_participants.conversation_id").
		Where("conversations.id = ? AND conversation_participants.user_id = ? AND conversations.type = ?",
			conversationID, currentUser.ID, "group").
		First(&conversation).Error

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group conversation not found or access denied"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify conversation access"})
		return
	}

	// Get users to add
	var users []User
	err = db.Where("id IN ? AND is_active = ?", reqData.UserIDs, true).Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find users"})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid users found"})
		return
	}

	// Add participants
	if err := db.Model(&conversation).Association("Participants").Append(users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add participants"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Participants added successfully"})
}

func removeParticipant(c *gin.Context) {
	currentUser := c.MustGet("user").(*User)
	conversationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user is participant and conversation is group
	var conversation Conversation
	err = db.Joins("JOIN conversation_participants ON conversations.id = conversation_participants.conversation_id").
		Where("conversations.id = ? AND conversation_participants.user_id = ? AND conversations.type = ?",
			conversationID, currentUser.ID, "group").
		First(&conversation).Error

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group conversation not found or access denied"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify conversation access"})
		return
	}

	// Remove participant
	var userToRemove User
	err = db.Where("id = ?", userID).First(&userToRemove).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := db.Model(&conversation).Association("Participants").Delete(&userToRemove); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove participant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Participant removed successfully"})
}
