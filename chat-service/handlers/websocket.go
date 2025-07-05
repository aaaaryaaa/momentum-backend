// chat-service/handlers/websocket.go
package handlers

import (
	"chat-service/db"
	"chat-service/middleware"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Adjust for CORS/security in production
}

func getRecentMessages(userID int) ([]map[string]interface{}, error) {
	query := `
		SELECT sender_id, receiver_id, content, created_at
		FROM messages
		WHERE (sender_id = $1 OR receiver_id = $1)
		ORDER BY created_at DESC
		LIMIT 50
	`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []map[string]interface{}
	for rows.Next() {
		var senderID, receiverID int
		var content string
		var createdAt time.Time

		if err := rows.Scan(&senderID, &receiverID, &content, &createdAt); err != nil {
			return nil, err
		}

		messages = append([]map[string]interface{}{{
			"sender_id":   senderID,
			"receiver_id": receiverID,
			"content":     content,
			"created_at":  createdAt,
		}}, messages...) // reverse to chronological order
	}

	return messages, nil
}

func getMessagesByReadStatus(userID int) (map[string][]map[string]interface{}, error) {
	query := `
		SELECT sender_id, receiver_id, content, created_at, is_read
		FROM messages
		WHERE receiver_id = $1
		ORDER BY created_at ASC
	`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := map[string][]map[string]interface{}{
		"read":   {},
		"unread": {},
	}

	for rows.Next() {
		var senderID, receiverID int
		var content string
		var createdAt time.Time
		var isRead bool

		if err := rows.Scan(&senderID, &receiverID, &content, &createdAt, &isRead); err != nil {
			return nil, err
		}

		msg := map[string]interface{}{
			"sender_id":   senderID,
			"receiver_id": receiverID,
			"content":     content,
			"created_at":  createdAt,
		}

		if isRead {
			messages["read"] = append(messages["read"], msg)
		} else {
			messages["unread"] = append(messages["unread"], msg)
		}
	}

	return messages, nil
}

func handleSendMessage(senderID, receiverID int, content string, conn *websocket.Conn) {
	// ✅ Check mutual follow
	var mutual bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM follows f1
			JOIN follows f2 ON f1.follower_id = f2.following_id AND f1.following_id = f2.follower_id
			WHERE f1.follower_id = $1 AND f1.following_id = $2
		)
	`
	err := db.DB.QueryRow(query, senderID, receiverID).Scan(&mutual)
	if err != nil || !mutual {
		conn.WriteJSON(map[string]string{
			"error": "You cannot message users who aren't mutual followers.",
		})
		return
	}

	// 💾 Store as unread
	_, err = db.DB.Exec(`
		INSERT INTO messages (sender_id, receiver_id, content, is_read)
		VALUES ($1, $2, $3, false)
	`, senderID, receiverID, content)
	if err != nil {
		conn.WriteJSON(map[string]string{"error": "Failed to save message"})
		return
	}

	// 📡 If receiver online, deliver + mark as read
	if receiverConn, ok := connections[receiverID]; ok {
		receiverConn.WriteJSON(map[string]interface{}{
			"type":    "new_message",
			"from":    senderID,
			"content": content,
		})

		_, err = db.DB.Exec(`
			UPDATE messages
			SET is_read = true
			WHERE sender_id = $1 AND receiver_id = $2 AND content = $3 AND is_read = false
		`, senderID, receiverID, content)
		if err != nil {
			log.Println("Failed to mark message as read:", err)
		}
	}
}

func handleLoadChat(userID, otherID int, before string, conn *websocket.Conn) {
	query := `
		SELECT sender_id, receiver_id, content, created_at, is_read
		FROM messages
		WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
	`
	args := []interface{}{userID, otherID}

	if before != "" {
		query += " AND created_at < $3"
		t, err := time.Parse(time.RFC3339, before)
		if err != nil {
			conn.WriteJSON(map[string]string{"error": "Invalid timestamp"})
			return
		}
		args = append(args, t)
	}

	query += " ORDER BY created_at DESC LIMIT 20"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		log.Println("DB error:", err)
		conn.WriteJSON(map[string]string{"error": "Failed to fetch chat"})
		return
	}
	defer rows.Close()

	readMessages := []map[string]interface{}{}
	unreadMessages := []map[string]interface{}{}

	for rows.Next() {
		var msg struct {
			SenderID   int
			ReceiverID int
			Content    string
			CreatedAt  time.Time
			IsRead     bool
		}
		rows.Scan(&msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.CreatedAt, &msg.IsRead)

		message := map[string]interface{}{
			"sender_id":   msg.SenderID,
			"receiver_id": msg.ReceiverID,
			"content":     msg.Content,
			"created_at":  msg.CreatedAt,
			"is_read":     msg.IsRead,
		}

		// If this message is sent *to* current user and unread, separate it
		if msg.ReceiverID == userID && !msg.IsRead {
			unreadMessages = append(unreadMessages, message)
		} else {
			readMessages = append(readMessages, message)
		}
	}

	conn.WriteJSON(map[string]interface{}{
		"type":   "chat_with_user",
		"with":   otherID,
		"read":   readMessages,
		"unread": unreadMessages,
	})

	// ✅ Mark messages from `otherID` to `userID` as read
	_, err = db.DB.Exec(`
		UPDATE messages
		SET is_read = true
		WHERE sender_id = $1 AND receiver_id = $2 AND is_read = false
	`, otherID, userID)

	if err != nil {
		log.Println("Failed to mark messages as read after load_chat:", err)
	}
}

var connections = make(map[int]*websocket.Conn) // map[userID] => conn

type IncomingMessage struct {
	Type       string `json:"type"`
	ReceiverID int    `json:"receiver_id,omitempty"`
	Content    string `json:"content,omitempty"`
	WithUserID int    `json:"with,omitempty"`
	Before     string `json:"before,omitempty"`
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	connections[user.ID] = conn
	defer delete(connections, user.ID)

	log.Printf("User %d connected via WebSocket\n", user.ID)

	// 🔥 Initial chat history (read/unread split)
	if history, err := getMessagesByReadStatus(user.ID); err == nil {
		conn.WriteJSON(map[string]interface{}{
			"type":   "chat_history",
			"read":   history["read"],
			"unread": history["unread"],
		})

		// // ✅ Mark all unread messages as read
		// _, err = db.DB.Exec(`
		// 	UPDATE messages
		// 	SET is_read = true
		// 	WHERE receiver_id = $1 AND is_read = false
		// `, user.ID)
		// if err != nil {
		// 	log.Println("Failed to mark unread messages as read:", err)
		// }
	} else {
		log.Println("Failed to fetch chat history:", err)
	}

	// 🔁 Main loop for handling incoming WebSocket messages
	for {
		var msg IncomingMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		switch msg.Type {
		case "message":
			handleSendMessage(user.ID, msg.ReceiverID, msg.Content, conn)

		case "load_chat":
			handleLoadChat(user.ID, msg.WithUserID, msg.Before, conn)

		default:
			conn.WriteJSON(map[string]string{"error": "Unknown message type"})
		}
	}
}

// Authenticated test route
func PingWithAuth(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"user_id": strconv.Itoa(user.ID),
		"message": "Authenticated!",
	})
}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var reqBody struct {
		OtherUserID int `json:"other_user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil || reqBody.OtherUserID == 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err := db.DB.Exec(`
		INSERT INTO follows (follower_id, following_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, user.ID, reqBody.OtherUserID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Follow successful",
	})
}

func CheckMutualFollow(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var reqBody struct {
		OtherUserID int `json:"other_user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil || reqBody.OtherUserID == 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	query := `
	SELECT EXISTS (
		SELECT 1 FROM follows f1
		JOIN follows f2 ON f1.follower_id = f2.following_id AND f1.following_id = f2.follower_id
		WHERE f1.follower_id = $1 AND f1.following_id = $2
	)
	`
	var mutual bool
	err := db.DB.QueryRow(query, user.ID, reqBody.OtherUserID).Scan(&mutual)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"mutual": mutual})
}
