package models

import "time"

type Post struct {
	ID        int            `json:"id"`
	UserID    int            `json:"user_id"`
	VideoID   int            `json:"video_id"`
	Caption   string         `json:"caption"`
	Category  string         `json:"category"`
	ThreadID  *int           `json:"thread_id,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	Nudges    int            `json:"nudges"`
	Reactions map[string]int `json:"reactions"` // emoji: count
}

type CreatePostRequest struct {
	VideoID  int    `json:"video_id"`
	Caption  string `json:"caption"`
	Category string `json:"category"`
	ThreadID *int   `json:"thread_id,omitempty"`
}

type UserStreak struct {
	UserID     int       `json:"user_id"`
	LastPosted time.Time `json:"last_posted"`
	Streak     int       `json:"streak"`
}
