// chat-service/models/follow.go
package models

type Follow struct {
	FollowerID  int `json:"follower_id"`
	FollowingID int `json:"following_id"`
}
