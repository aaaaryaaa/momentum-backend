// message_handlers.go
package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func sendMessage(c *gin.Context) {
	currentUser := c.MustGet("user").(*User)

	var reqData struct {
		ConversationID int    `json:"conversation_id"`
		Content        string `json:"content"`
		Type           string `json:"type"`
		MediaURL       string `json:"media_url"`
		ReplyToID      *int   `json:"reply_to_id"`
	}

	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate message type
	if reqData.Type == "" {
		reqData.Type = "text"
	}

	// Validate content
	if reqData.Content == "" && reqData.MediaURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message content or media URL required"})
		return
	}

	// Check if user is participant
	var conversation Conversation
	err := db.Joins("JOIN conversation_participants ON conversations.id = conversation_participants.conversation_id").
		Where("conversations.id = ? AND conversation_participants.user_id = ?", reqData.ConversationID, currentUser.ID).
		First(&conversation).Error

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found or access denied"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify conversation access"})
		return
	}

	// Create message
	message := Message{
		ConversationID: reqData.ConversationID,
		SenderID:       currentUser.ID,
		Content:        reqData.Content,
		Type:           reqData.Type,
		MediaURL:       reqData.MediaURL,
		ReplyToID:      reqData.ReplyToID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	// Update conversation last message
	conversation.LastMessageID = &message.ID
	conversation.UpdatedAt = time.Now()
	db.Save(&conversation)

	// Create message status for all participants
	var participants []User
	db.Model(&conversation).Association("Participants").Find(&participants)

	for _, participant := range participants {
		status := "sent"
		if participant.ID != currentUser.ID {
			status = "delivered"
		}

		messageStatus := MessageStatus{
			MessageID: message.ID,
			UserID:    participant.ID,
			Status:    status,
			CreatedAt: time.Now(),
		}
		db.Create(&messageStatus)
	}

	// Load message with sender
	db.Preload("Sender").Preload("ReplyTo").Preload("ReplyTo.Sender").First(&message, message.ID)

	c.JSON(http.StatusCreated, message)
}

func getMessage(c *gin.Context) {
	currentUser := c.MustGet("user").(*User)
	messageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var message Message
	err = db.Preload("Sender").
		Preload("ReplyTo").
		Preload("ReplyTo.Sender").
		Preload("Conversation").
		Joins("JOIN conversations ON messages.conversation_id = conversations.id").
		Joins("JOIN conversation_participants ON conversations.id = conversation_participants.conversation_id").
		Where("messages.id = ? AND conversation_participants.user_id = ?", messageID, currentUser.ID).
		First(&message).Error

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found or access denied"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get message"})
		return
	}

	c.JSON(http.StatusOK, message)
}

func editMessage(c *gin.Context) {
	currentUser := c.MustGet("user").(*User)
	messageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var reqData struct {
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if reqData.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content required"})
		return
	}

	// Get message and check ownership
	var message Message
	err = db.Where("id = ? AND sender_id = ?", messageID, currentUser.ID).First(&message).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found or access denied"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get message"})
		return
	}

	// Update message
	message.Content = reqData.Content
	message.IsEdited = true
	message.UpdatedAt = time.Now()

	if err := db.Save(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message"})
		return
	}

	// Load message with sender
	db.Preload("Sender").Preload("ReplyTo").Preload("ReplyTo.Sender").First(&message, message.ID)

	c.JSON(http.StatusOK, message)
}

func deleteMessage(c *gin.Context) {
	currentUser := c.MustGet("user").(*User)
	messageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	// Get message and check ownership
	var message Message
	err = db.Where("id = ? AND sender_id = ?", messageID, currentUser.ID).First(&message).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found or access denied"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get message"})
		return
	}

	// Delete message (soft delete by updating content)
	message.Content = "This message was deleted"
	message.Type = "deleted"
	message.MediaURL = ""
	message.UpdatedAt = time.Now()

	if err := db.Save(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}

func markMessageAsRead(c *gin.Context) {
	currentUser := c.MustGet("user").(*User)
	messageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	// Check if user can access this message
	var message Message
	err = db.Joins("JOIN conversations ON messages.conversation_id = conversations.id").
		Joins("JOIN conversation_participants ON conversations.id = conversation_participants.conversation_id").
		Where("messages.id = ? AND conversation_participants.user_id = ?", messageID, currentUser.ID).
		First(&message).Error

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found or access denied"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify message access"})
		return
	}

	// Update message status to read
	err = db.Model(&MessageStatus{}).
		Where("message_id = ? AND user_id = ?", messageID, currentUser.ID).
		Update("status", "read").Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark message as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message marked as read"})
}
