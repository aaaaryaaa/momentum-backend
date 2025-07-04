package models

type Video struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	VideoName string `json:"video_name"`
	S3URL     string `json:"s3_url"`
	Timestamp string `json:"timestamp"`
}
