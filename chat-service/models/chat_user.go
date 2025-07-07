// chat-service/models/chat_user.go
package models

type ChatUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
