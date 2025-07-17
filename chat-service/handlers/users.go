// chat-service/handlers/users.go
package handlers

import (
	"chat-service/db"
	"chat-service/middleware"
	"chat-service/models"
	"encoding/json"
	"net/http"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := `SELECT id, name, email FROM chat_users WHERE id != $1 ORDER BY name`
	rows, err := db.DB.Query(query, user.ID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.ChatUser
	for rows.Next() {
		var u models.ChatUser
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			continue
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	searchTerm := r.URL.Query().Get("q")
	if searchTerm == "" {
		http.Error(w, "Search term is required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT id, name, email 
		FROM chat_users 
		WHERE id != $1 AND (name ILIKE $2 OR email ILIKE $2)
		ORDER BY name
		LIMIT 20
	`
	rows, err := db.DB.Query(query, user.ID, "%"+searchTerm+"%")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.ChatUser
	for rows.Next() {
		var u models.ChatUser
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			continue
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetChatList(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := `
		SELECT DISTINCT 
			CASE 
				WHEN m.sender_id = $1 THEN m.receiver_id 
				ELSE m.sender_id 
			END as other_user_id,
			cu.name,
			cu.email,
			MAX(m.created_at) as last_message_time,
			COUNT(CASE WHEN m.receiver_id = $1 AND m.is_read = false THEN 1 END) as unread_count
		FROM messages m
		JOIN chat_users cu ON cu.id = CASE 
			WHEN m.sender_id = $1 THEN m.receiver_id 
			ELSE m.sender_id 
		END
		WHERE m.sender_id = $1 OR m.receiver_id = $1
		GROUP BY other_user_id, cu.name, cu.email
		ORDER BY last_message_time DESC
	`

	rows, err := db.DB.Query(query, user.ID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type ChatListItem struct {
		UserID          int    `json:"user_id"`
		Name            string `json:"name"`
		Email           string `json:"email"`
		LastMessageTime string `json:"last_message_time"`
		UnreadCount     int    `json:"unread_count"`
	}

	var chatList []ChatListItem
	for rows.Next() {
		var item ChatListItem
		if err := rows.Scan(&item.UserID, &item.Name, &item.Email, &item.LastMessageTime, &item.UnreadCount); err != nil {
			continue
		}
		chatList = append(chatList, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatList)
}
