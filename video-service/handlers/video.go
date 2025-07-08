// package handlers

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"time"

// 	"video-service/db"
// 	"video-service/middleware"
// 	"video-service/models"
// 	"video-service/utils"
// )

// // func GenerateUploadURL(w http.ResponseWriter, r *http.Request) {
// // 	user := middleware.GetUserFromContext(r.Context())
// // 	if user == nil {
// // 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// // 		return
// // 	}

// // 	var reqBody struct {
// // 		Filename string `json:"filename"`
// // 	}
// // 	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil || reqBody.Filename == "" {
// // 		http.Error(w, "Missing filename", http.StatusBadRequest)
// // 		return
// // 	}

// // 	uploadURL, err := utils.GeneratePresignedURL(reqBody.Filename)
// // 	if err != nil {
// // 		http.Error(w, "Failed to generate URL", http.StatusInternalServerError)
// // 		return
// // 	}

// // 	staticURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", os.Getenv("S3_BUCKET"), os.Getenv("AWS_REGION"), reqBody.Filename)

// // 	_, err = db.DB.Exec(
// // 		"INSERT INTO videos(user_id, video_name, s3_url, timestamp) VALUES($1, $2, $3, $4)",
// // 		user.ID, reqBody.Filename, staticURL, time.Now(),
// // 	)
// // 	if err != nil {
// // 		http.Error(w, "Failed to save metadata", http.StatusInternalServerError)
// // 		return
// // 	}

// //		json.NewEncoder(w).Encode(map[string]string{
// //			"upload_url": uploadURL,
// //		})
// //	}
// func GenerateUploadURL(w http.ResponseWriter, r *http.Request) {
// 	user := middleware.GetUserFromContext(r.Context())
// 	if user == nil {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	var reqBody struct {
// 		Filename string `json:"filename"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil || reqBody.Filename == "" {
// 		http.Error(w, "Missing filename", http.StatusBadRequest)
// 		return
// 	}

// 	// Generate presigned PUT URL using user ID and filename
// 	uploadURL, objectKey, err := utils.GeneratePresignedPutURL(strconv.Itoa(user.ID), reqBody.Filename)
// 	if err != nil {
// 		http.Error(w, "Failed to generate URL", http.StatusInternalServerError)
// 		return
// 	}

// 	// Public static S3 URL (assuming file is uploaded with "public-read" ACL)
// 	staticURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", os.Getenv("S3_BUCKET"), os.Getenv("AWS_REGION"), objectKey)

// 	// Save metadata in DB
// 	_, err = db.DB.Exec(
// 		"INSERT INTO videos(user_id, video_name, s3_url, timestamp) VALUES($1, $2, $3, $4)",
// 		user.ID, objectKey, staticURL, time.Now(),
// 	)
// 	if err != nil {
// 		http.Error(w, "Failed to save metadata", http.StatusInternalServerError)
// 		return
// 	}

// 	// Return upload URL and public URL to client
// 	// json.NewEncoder(w).Encode(map[string]string{
// 	// 	"upload_url": uploadURL,
// 	// 	"public_url": staticURL,
// 	// })
// 	w.Header().Set("Content-Type", "application/json")
// 	enc := json.NewEncoder(w)
// 	enc.SetEscapeHTML(false) // ← 🔥 FIX
// 	enc.Encode(map[string]string{
// 		"upload_url": uploadURL,
// 	})
// }

// func GetUserVideos(w http.ResponseWriter, r *http.Request) {
// 	user := middleware.GetUserFromContext(r.Context())
// 	if user == nil {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	rows, err := db.DB.Query("SELECT id, user_id, video_name, s3_url, timestamp FROM videos WHERE user_id = $1", user.ID)
// 	if err != nil {
// 		http.Error(w, "DB error", http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var videos []models.Video
// 	for rows.Next() {
// 		var vid models.Video
// 		err := rows.Scan(&vid.ID, &vid.UserID, &vid.VideoName, &vid.S3URL, &vid.Timestamp)
// 		if err != nil {
// 			continue
// 		}

// 		// Generate presigned GET URL based on video_name (i.e., S3 key)
// 		presignedURL, err := utils.GeneratePresignedGetURL(vid.VideoName)
// 		if err != nil {
// 			continue // skip this video if URL generation fails
// 		}

// 		vid.S3URL = presignedURL
// 		videos = append(videos, vid)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(videos)
// }

// handlers/video.go - Add the delete handler
package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"video-service/db"
	"video-service/middleware"
	"video-service/models"
	"video-service/utils"

	"github.com/gorilla/mux"
)

func GenerateUploadURL(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var reqBody struct {
		Filename string `json:"filename"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil || reqBody.Filename == "" {
		http.Error(w, "Missing filename", http.StatusBadRequest)
		return
	}

	// Generate presigned PUT URL using user ID and filename
	uploadURL, objectKey, err := utils.GeneratePresignedPutURL(strconv.Itoa(user.ID), reqBody.Filename)
	if err != nil {
		http.Error(w, "Failed to generate URL", http.StatusInternalServerError)
		return
	}

	// Public static S3 URL (assuming file is uploaded with "public-read" ACL)
	staticURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", os.Getenv("S3_BUCKET"), os.Getenv("AWS_REGION"), objectKey)

	// Save metadata in DB
	result, err := db.DB.Exec(
		"INSERT INTO videos(user_id, video_name, s3_url, timestamp) VALUES($1, $2, $3, $4) RETURNING id",
		user.ID, objectKey, staticURL, time.Now(),
	)
	if err != nil {
		http.Error(w, "Failed to save metadata", http.StatusInternalServerError)
		return
	}

	// Get the inserted video ID
	videoID, err := result.LastInsertId()
	if err != nil {
		// For PostgreSQL, use a different approach
		var id int
		err = db.DB.QueryRow(
			"INSERT INTO videos(user_id, video_name, s3_url, timestamp) VALUES($1, $2, $3, $4) RETURNING id",
			user.ID, objectKey, staticURL, time.Now(),
		).Scan(&id)
		if err != nil {
			http.Error(w, "Failed to save metadata", http.StatusInternalServerError)
			return
		}
		videoID = int64(id)
	}

	// Return upload URL and video ID to client
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(map[string]interface{}{
		"upload_url": uploadURL,
		"video_id":   videoID,
	})
}

func GetUserVideos(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.DB.Query("SELECT id, user_id, video_name, s3_url, timestamp FROM videos WHERE user_id = $1 ORDER BY timestamp DESC", user.ID)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		var vid models.Video
		err := rows.Scan(&vid.ID, &vid.UserID, &vid.VideoName, &vid.S3URL, &vid.Timestamp)
		if err != nil {
			continue
		}

		// Generate presigned GET URL based on video_name (i.e., S3 key)
		presignedURL, err := utils.GeneratePresignedGetURL(vid.VideoName)
		if err != nil {
			continue // skip this video if URL generation fails
		}

		vid.S3URL = presignedURL
		videos = append(videos, vid)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

func DeleteVideo(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get video ID from URL parameter
	vars := mux.Vars(r)
	videoIDStr := vars["id"]
	videoID, err := strconv.Atoi(videoIDStr)
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	// First, get the video info to check ownership and get S3 key
	var video models.Video
	err = db.DB.QueryRow(
		"SELECT id, user_id, video_name, s3_url FROM videos WHERE id = $1",
		videoID,
	).Scan(&video.ID, &video.UserID, &video.VideoName, &video.S3URL)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Video not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Check if the video belongs to the current user
	if video.UserID != user.ID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Delete from S3
	err = utils.DeleteFromS3(video.VideoName)
	if err != nil {
		// Log the error but don't fail the request
		// The video might already be deleted from S3
		fmt.Printf("Warning: Failed to delete from S3: %v\n", err)
	}

	// Delete from database
	_, err = db.DB.Exec("DELETE FROM videos WHERE id = $1", videoID)
	if err != nil {
		http.Error(w, "Failed to delete video", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Video deleted successfully",
	})
}
