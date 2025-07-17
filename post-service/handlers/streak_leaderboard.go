package handlers

import (
	"encoding/json"
	"net/http"
	"post-service/db"
	"time"
)

type StreakLeaderboardEntry struct {
	UserID     int    `json:"user_id"`
	Streak     int    `json:"streak"`
	LastPosted string `json:"last_posted"` // ISO 8601 format
}

// GET /streak-leaderboard
func GetStreakLeaderboard(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT user_id, streak, last_posted
		FROM user_streaks
		ORDER BY streak DESC
		LIMIT 10
	`)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var leaderboard []StreakLeaderboardEntry
	for rows.Next() {
		var entry StreakLeaderboardEntry
		var lastPosted time.Time
		if err := rows.Scan(&entry.UserID, &entry.Streak, &lastPosted); err == nil {
			entry.LastPosted = lastPosted.Format(time.RFC3339)
			leaderboard = append(leaderboard, entry)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaderboard)
}
