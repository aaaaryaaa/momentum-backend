package handlers

import (
	"encoding/json"
	"net/http"
	"post-service/db"
)

type LeaderboardEntry struct {
	UserID    int    `json:"user_id"`
	Category  string `json:"category"`
	Nudges    int    `json:"nudges"`
	LastNudge string `json:"last_nudge"` // ISO 8601 format
}

// GET /leaderboard?category=Fitness
func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category == "" {
		http.Error(w, "Missing category", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query(`
		SELECT 
			p.user_id, 
			p.category, 
			COUNT(n.id) AS nudges,
			MAX(n.created_at) AS last_nudge
		FROM posts p
		JOIN nudges n ON n.post_id = p.id
		WHERE p.category = $1
		GROUP BY p.user_id, p.category
		ORDER BY nudges DESC
		LIMIT 10
	`, category)

	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var leaderboard []LeaderboardEntry
	for rows.Next() {
		var entry LeaderboardEntry
		var lastNudgeRaw string
		if err := rows.Scan(&entry.UserID, &entry.Category, &entry.Nudges, &lastNudgeRaw); err == nil {
			entry.LastNudge = lastNudgeRaw
			leaderboard = append(leaderboard, entry)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaderboard)
}
