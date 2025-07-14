// // handlers/post.go
// package handlers

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"net/http"
// 	"post-service/db"
// 	"post-service/middleware"
// 	"post-service/models"
// 	"strconv"
// 	"time"

// 	"github.com/gorilla/mux"
// )

// // POST /posts - Create a new post
// func CreatePost(w http.ResponseWriter, r *http.Request) {
// 	user := middleware.GetUserFromContext(r.Context())
// 	if user == nil {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	var req models.CreatePostRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	// var username string
// 	// err := db.DB.QueryRow("SELECT name FROM users WHERE id=$1", user.ID).Scan(&username)
// 	// if err != nil {
// 	// 	http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
// 	// 	return
// 	// }

// 	_, err := db.DB.Exec(`
// 		INSERT INTO posts(user_id, video_id, caption, category, thread_id, created_at)
// 		VALUES ($1, $2, $3, $4, $5, $6)
// 	`, user.ID, req.VideoID, req.Caption, req.Category, req.ThreadID, time.Now())

// 	if err != nil {
// 		http.Error(w, "DB insert failed", http.StatusInternalServerError)
// 		return
// 	}

// 	// After successful post insert:
// 	today := time.Now().UTC().Truncate(24 * time.Hour)

// 	var lastPosted time.Time
// 	var streak int

// 	err = db.DB.QueryRow("SELECT last_posted, streak FROM user_streaks WHERE user_id = $1", user.ID).Scan(&lastPosted, &streak)

// 	// if err == sql.ErrNoRows {
// 	// 	_, _ = db.DB.Exec("INSERT INTO user_streaks(user_id, last_posted, streak) VALUES($1, $2, 1)", user.ID, today)
// 	// } else if err == nil {
// 	// 	yesterday := today.AddDate(0, 0, -1)
// 	// 	if lastPosted.Equal(yesterday) {
// 	// 		streak++
// 	// 	} else if !lastPosted.Equal(today) {
// 	// 		streak = 1
// 	// 	}
// 	// 	_, _ = db.DB.Exec("UPDATE user_streaks SET last_posted = $1, streak = $2 WHERE user_id = $3", today, streak, user.ID)
// 	// }
// 	if err == sql.ErrNoRows {
// 		// First post ever
// 		_, _ = db.DB.Exec("INSERT INTO user_streaks(user_id, last_posted, streak) VALUES($1, $2, 1)", user.ID, today)
// 	} else if err == nil {
// 		lastPosted = lastPosted.Truncate(24 * time.Hour)
// 		yesterday := today.AddDate(0, 0, -1)

// 		if lastPosted.Equal(today) {
// 			// Already posted today, no update needed
// 		} else if lastPosted.Equal(yesterday) {
// 			// Continue streak
// 			streak++
// 			_, _ = db.DB.Exec("UPDATE user_streaks SET last_posted = $1, streak = $2 WHERE user_id = $3", today, streak, user.ID)
// 		} else {
// 			// Broke streak, reset to 1
// 			streak = 1
// 			_, _ = db.DB.Exec("UPDATE user_streaks SET last_posted = $1, streak = $2 WHERE user_id = $3", today, streak, user.ID)
// 		}
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]string{"message": "Post created"})
// }

// // GET /feed - Return latest posts (paginated)
// // func GetFeed(w http.ResponseWriter, r *http.Request) {
// // 	rows, err := db.DB.Query(`
// // 		SELECT id, user_id, video_id, caption, category, thread_id, created_at
// // 		FROM posts
// // 		ORDER BY created_at DESC
// // 		LIMIT 20
// // 	`)
// // 	if err != nil {
// // 		http.Error(w, "DB error", http.StatusInternalServerError)
// // 		return
// // 	}
// // 	defer rows.Close()

// // 	var posts []models.Post
// // 	for rows.Next() {
// // 		var p models.Post
// // 		err := rows.Scan(&p.ID, &p.UserID, &p.VideoID, &p.Caption, &p.Category, &p.ThreadID, &p.CreatedAt)
// // 		if err == nil {
// // 			posts = append(posts, p)
// // 		}
// // 	}

// //		w.Header().Set("Content-Type", "application/json")
// //		json.NewEncoder(w).Encode(posts)
// //	}
// func GetFeed(w http.ResponseWriter, r *http.Request) {
// 	rows, err := db.DB.Query(`
// 		SELECT
// 			p.id, p.user_id, p.video_id, p.caption, p.category, p.thread_id, p.created_at,
// 			COALESCE(n.nudges, 0) AS nudges,
// 			COALESCE(r.reactions, '{}') AS reactions
// 		FROM posts p
// 		LEFT JOIN (
// 			SELECT post_id, COUNT(*) AS nudges
// 			FROM nudges
// 			GROUP BY post_id
// 		) n ON p.id = n.post_id
// 		LEFT JOIN (
// 			SELECT post_id, json_object_agg(emoji, count) AS reactions
// 			FROM (
// 				SELECT post_id, emoji, COUNT(*) AS count
// 				FROM reactions
// 				GROUP BY post_id, emoji
// 			) sub
// 			GROUP BY post_id
// 		) r ON p.id = r.post_id
// 		ORDER BY p.created_at DESC
// 		LIMIT 20
// 	`)
// 	if err != nil {
// 		http.Error(w, "DB error", http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var posts []models.Post
// 	for rows.Next() {
// 		var p models.Post
// 		var reactionsJSON []byte
// 		err := rows.Scan(&p.ID, &p.UserID, &p.VideoID, &p.Caption, &p.Category, &p.ThreadID, &p.CreatedAt, &p.Nudges, &reactionsJSON)
// 		if err != nil {
// 			continue
// 		}
// 		if len(reactionsJSON) > 0 {
// 			_ = json.Unmarshal(reactionsJSON, &p.Reactions)
// 		} else {
// 			p.Reactions = map[string]int{}
// 		}
// 		posts = append(posts, p)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(posts)
// }

// // DELETE /posts/{id} - Delete a post by ID
// func DeletePost(w http.ResponseWriter, r *http.Request) {
// 	user := middleware.GetUserFromContext(r.Context())
// 	if user == nil {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	vars := mux.Vars(r)
// 	postID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid post ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Ensure the post belongs to the user
// 	var ownerID int
// 	err = db.DB.QueryRow("SELECT user_id FROM posts WHERE id = $1", postID).Scan(&ownerID)
// 	if err != nil {
// 		http.Error(w, "Post not found", http.StatusNotFound)
// 		return
// 	}
// 	if ownerID != user.ID {
// 		http.Error(w, "Forbidden", http.StatusForbidden)
// 		return
// 	}

// 	_, err = db.DB.Exec("DELETE FROM posts WHERE id = $1", postID)
// 	if err != nil {
// 		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted"})
// }

// handlers/post.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"post-service/db"
	"post-service/middleware"
	"post-service/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// POST /posts - Create a new post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.VideoURL == "" {
		http.Error(w, "Video URL is required", http.StatusBadRequest)
		return
	}

	_, err := db.DB.Exec(`
		INSERT INTO posts(user_id, video_id, video_url, caption, category, thread_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, user.ID, req.VideoID, req.VideoURL, req.Caption, req.Category, req.ThreadID, time.Now())
	if err != nil {
		http.Error(w, "DB insert failed", http.StatusInternalServerError)
		return
	}

	// After successful post insert - update user streak:
	today := time.Now().UTC().Truncate(24 * time.Hour)
	var lastPosted time.Time
	var streak int
	err = db.DB.QueryRow("SELECT last_posted, streak FROM user_streaks WHERE user_id = $1", user.ID).Scan(&lastPosted, &streak)

	if err == sql.ErrNoRows {
		// First post ever
		db.DB.Exec("INSERT INTO user_streaks(user_id, last_posted, streak) VALUES($1, $2, 1)", user.ID, today)
	} else if err == nil {
		lastPosted = lastPosted.Truncate(24 * time.Hour)
		yesterday := today.AddDate(0, 0, -1)
		if lastPosted.Equal(today) {
			// Already posted today, no update needed
		} else if lastPosted.Equal(yesterday) {
			// Continue streak
			streak++
			db.DB.Exec("UPDATE user_streaks SET last_posted = $1, streak = $2 WHERE user_id = $3", today, streak, user.ID)
		} else {
			// Broke streak, reset to 1
			streak = 1
			db.DB.Exec("UPDATE user_streaks SET last_posted = $1, streak = $2 WHERE user_id = $3", today, streak, user.ID)
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post created"})
}

// GET /feed - Return latest posts (paginated)
func GetFeed(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT 
			p.id, p.user_id, p.video_id, p.video_url, p.caption, p.category, p.thread_id, p.created_at,
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
		ORDER BY p.created_at DESC
		LIMIT 20
	`)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		var reactionsJSON []byte
		err := rows.Scan(&p.ID, &p.UserID, &p.VideoID, &p.VideoURL, &p.Caption, &p.Category, &p.ThreadID, &p.CreatedAt, &p.Nudges, &reactionsJSON)
		if err != nil {
			continue
		}

		if len(reactionsJSON) > 0 {
			json.Unmarshal(reactionsJSON, &p.Reactions)
		} else {
			p.Reactions = map[string]int{}
		}

		posts = append(posts, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// DELETE /posts/{id} - Delete a post by ID
func DeletePost(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Ensure the post belongs to the user
	var ownerID int
	err = db.DB.QueryRow("SELECT user_id FROM posts WHERE id = $1", postID).Scan(&ownerID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if ownerID != user.ID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_, err = db.DB.Exec("DELETE FROM posts WHERE id = $1", postID)
	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted"})
}
