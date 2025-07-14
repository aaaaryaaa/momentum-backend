package handlers

import (
	"encoding/json"
	"net/http"
	"post-service/db"
	"post-service/middleware"
	"post-service/models"
	"strconv"

	"github.com/gorilla/mux"
)

// GET /thread/{id} - Get all posts in a thread (including the parent)
func GetThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query(`
		SELECT id, user_id, video_id, caption, category, thread_id, created_at
		FROM posts
		WHERE id = $1 OR thread_id = $1
		ORDER BY created_at ASC
	`, threadID)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(&p.ID, &p.UserID, &p.VideoID, &p.Caption, &p.Category, &p.ThreadID, &p.CreatedAt)
		if err == nil {
			posts = append(posts, p)
		}
	}

	json.NewEncoder(w).Encode(posts)
}

// GET /mythreads - Get user's own thread-starting posts
func GetMyThreads(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.DB.Query(`
		SELECT 
			p.id, p.user_id, p.video_id, p.caption, p.category, p.thread_id, p.created_at,
			COALESCE(n.nudges, 0) AS nudges,
			COALESCE(r.reactions, '{}') AS reactions
		FROM posts p
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS nudges
			FROM nudges
			GROUP BY post_id
		) n ON p.id = n.post_id
		LEFT JOIN (
			SELECT post_id, json_object_agg(emoji, count) AS reactions
			FROM (
				SELECT post_id, emoji, COUNT(*) AS count
				FROM reactions
				GROUP BY post_id, emoji
			) sub
			GROUP BY post_id
		) r ON p.id = r.post_id
		WHERE p.user_id = $1 AND p.thread_id IS NULL
		ORDER BY p.created_at DESC
	`, user.ID)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var threads []models.Post
	for rows.Next() {
		var p models.Post
		var reactionsJSON []byte
		err := rows.Scan(&p.ID, &p.UserID, &p.VideoID, &p.Caption, &p.Category, &p.ThreadID, &p.CreatedAt, &p.Nudges, &reactionsJSON)
		if err != nil {
			continue
		}
		if len(reactionsJSON) > 0 {
			_ = json.Unmarshal(reactionsJSON, &p.Reactions)
		} else {
			p.Reactions = map[string]int{}
		}
		threads = append(threads, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(threads)
}
