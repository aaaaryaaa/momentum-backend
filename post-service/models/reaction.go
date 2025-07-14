package models

type Reaction struct {
	PostID int    `json:"post_id"`
	Emoji  string `json:"emoji"`
	UserID int    `json:"user_id"`
}

type Nudge struct {
	PostID int `json:"post_id"`
	UserID int `json:"user_id"`
}
