package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"post-service/db"
	"post-service/middleware"
)

func GetUserStreak(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var streak int
	err := db.DB.QueryRow("SELECT streak FROM user_streaks WHERE user_id = $1", user.ID).Scan(&streak)
	if err != nil {
		log.Println("Streak fetch error:", err)
		streak = 0
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"streak": streak})
}
