// chat-service/handlers/sync.go
package handlers

import (
	"chat-service/db"
	"chat-service/middleware"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type SyncUserRequest struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// SyncUser - endpoint called by auth service to sync user data
func SyncUser(w http.ResponseWriter, r *http.Request) {
	var req SyncUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Insert or update user in chat_users table
	_, err := db.DB.Exec(`
		INSERT INTO chat_users (id, name, email, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (id) 
		DO UPDATE SET 
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			updated_at = NOW()
	`, req.ID, req.Name, req.Email)

	if err != nil {
		log.Printf("Failed to sync user: %v", err)
		http.Error(w, "Failed to sync user", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully synced user: %d - %s", req.ID, req.Name)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User synced successfully",
	})
}

// BulkSyncUsers - endpoint to sync all users from auth service
func BulkSyncUsers(w http.ResponseWriter, r *http.Request) {
	var users []SyncUserRequest
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Begin transaction
	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Insert or update all users
	for _, user := range users {
		_, err := tx.Exec(`
			INSERT INTO chat_users (id, name, email, created_at, updated_at)
			VALUES ($1, $2, $3, NOW(), NOW())
			ON CONFLICT (id) 
			DO UPDATE SET 
				name = EXCLUDED.name,
				email = EXCLUDED.email,
				updated_at = NOW()
		`, user.ID, user.Name, user.Email)

		if err != nil {
			log.Printf("Failed to sync user %d: %v", user.ID, err)
			http.Error(w, "Failed to sync users", http.StatusInternalServerError)
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully synced %d users", len(users))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": fmt.Sprintf("Successfully synced %d users", len(users)),
		"count":   len(users),
	})
}

// GetFollowStatus - get follow status between current user and another user
func GetFollowStatus(w http.ResponseWriter, r *http.Request) {
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

	type FollowStatus struct {
		IFollowThem  bool `json:"i_follow_them"`
		TheyFollowMe bool `json:"they_follow_me"`
		MutualFollow bool `json:"mutual_follow"`
	}

	var status FollowStatus

	// Check if current user follows the other user
	err := db.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $1 AND following_id = $2)
	`, user.ID, reqBody.OtherUserID).Scan(&status.IFollowThem)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Check if other user follows current user
	err = db.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $1 AND following_id = $2)
	`, reqBody.OtherUserID, user.ID).Scan(&status.TheyFollowMe)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	status.MutualFollow = status.IFollowThem && status.TheyFollowMe

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// UnfollowUser - unfollow a user
func UnfollowUser(w http.ResponseWriter, r *http.Request) {
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
		DELETE FROM follows 
		WHERE follower_id = $1 AND following_id = $2
	`, user.ID, reqBody.OtherUserID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Unfollowed successfully",
	})
}
