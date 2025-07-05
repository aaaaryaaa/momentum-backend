// // models.go
// package main

// import (
// 	"time"
// )

// type User struct {
// 	ID        int       `json:"id" gorm:"primaryKey"`
// 	Name      string    `json:"name"`
// 	Email     string    `json:"email"`
// 	Username  string    `json:"username" gorm:"unique;not null"`
// 	Avatar    string    `json:"avatar"`
// 	Bio       string    `json:"bio"`
// 	IsPrivate bool      `json:"is_private" gorm:"default:false"`
// 	IsActive  bool      `json:"is_active" gorm:"default:true"`
// 	LastSeen  time.Time `json:"last_seen"`
// 	CreatedAt time.Time `json:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at"`

// 	// Relationships
// 	Following []Follow `json:"following,omitempty" gorm:"foreignKey:FollowerID"`
// 	Followers []Follow `json:"followers,omitempty" gorm:"foreignKey:FollowingID"`
// }

// type Follow struct {
// 	ID          int       `json:"id" gorm:"primaryKey"`
// 	FollowerID  int       `json:"follower_id" gorm:"not null"`
// 	FollowingID int       `json:"following_id" gorm:"not null"`
// 	Status      string    `json:"status" gorm:"default:'pending'"` // pending, accepted, blocked
// 	CreatedAt   time.Time `json:"created_at"`

// 	// Relationships
// 	Follower  User `json:"follower,omitempty" gorm:"foreignKey:FollowerID"`
// 	Following User `json:"following,omitempty" gorm:"foreignKey:FollowingID"`
// }

// type Conversation struct {
// 	ID        int       `json:"id" gorm:"primaryKey"`
// 	Type      string    `json:"type" gorm:"default:'direct'"` // direct, group
// 	Name      string    `json:"name"`
// 	Avatar    string    `json:"avatar"`
// 	CreatedBy int       `json:"created_by"`
// 	CreatedAt time.Time `json:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at"`

// 	// Relationships
// 	Participants  []User    `json:"participants,omitempty" gorm:"many2many:conversation_participants;"`
// 	Messages      []Message `json:"messages,omitempty"`
// 	LastMessage   *Message  `json:"last_message,omitempty" gorm:"foreignKey:ID;references:LastMessageID"`
// 	LastMessageID *int      `json:"last_message_id"`
// }

// type Message struct {
// 	ID             int       `json:"id" gorm:"primaryKey"`
// 	ConversationID int       `json:"conversation_id" gorm:"not null"`
// 	SenderID       int       `json:"sender_id" gorm:"not null"`
// 	Content        string    `json:"content"`
// 	Type           string    `json:"type" gorm:"default:'text'"` // text, image, video, audio, file
// 	MediaURL       string    `json:"media_url"`
// 	ReplyToID      *int      `json:"reply_to_id"`
// 	IsEdited       bool      `json:"is_edited" gorm:"default:false"`
// 	CreatedAt      time.Time `json:"created_at"`
// 	UpdatedAt      time.Time `json:"updated_at"`

// 	// Relationships
// 	Sender        User            `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
// 	Conversation  Conversation    `json:"conversation,omitempty" gorm:"foreignKey:ConversationID"`
// 	ReplyTo       *Message        `json:"reply_to,omitempty" gorm:"foreignKey:ReplyToID"`
// 	MessageStatus []MessageStatus `json:"message_status,omitempty"`
// }

// type MessageStatus struct {
// 	ID        int       `json:"id" gorm:"primaryKey"`
// 	MessageID int       `json:"message_id" gorm:"not null"`
// 	UserID    int       `json:"user_id" gorm:"not null"`
// 	Status    string    `json:"status" gorm:"default:'sent'"` // sent, delivered, read
// 	CreatedAt time.Time `json:"created_at"`

// 	// Relationships
// 	Message Message `json:"message,omitempty" gorm:"foreignKey:MessageID"`
// 	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
// }

// models.go
package main

import (
	"time"
)

type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	IsPrivate bool      `json:"is_private" gorm:"default:false"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	LastSeen  time.Time `json:"last_seen"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Following []Follow `json:"following,omitempty" gorm:"foreignKey:FollowerID"`
	Followers []Follow `json:"followers,omitempty" gorm:"foreignKey:FollowingID"`
}

type Follow struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	FollowerID  int       `json:"follower_id" gorm:"not null"`
	FollowingID int       `json:"following_id" gorm:"not null"`
	Status      string    `json:"status" gorm:"default:'pending'"` // pending, accepted, blocked
	CreatedAt   time.Time `json:"created_at"`

	// Relationships
	Follower  User `json:"follower,omitempty" gorm:"foreignKey:FollowerID"`
	Following User `json:"following,omitempty" gorm:"foreignKey:FollowingID"`
}

type Conversation struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	Type          string    `json:"type" gorm:"default:'direct'"` // direct, group
	Name          string    `json:"name"`
	Avatar        string    `json:"avatar"`
	CreatedBy     int       `json:"created_by"`
	LastMessageID *int      `json:"last_message_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relationships
	Participants []User    `json:"participants,omitempty" gorm:"many2many:conversation_participants;"`
	Messages     []Message `json:"messages,omitempty"`
	// Fixed: Use proper foreign key reference
	LastMessage *Message `json:"last_message,omitempty" gorm:"foreignKey:LastMessageID"`
}

type Message struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	ConversationID int       `json:"conversation_id" gorm:"not null"`
	SenderID       int       `json:"sender_id" gorm:"not null"`
	Content        string    `json:"content"`
	Type           string    `json:"type" gorm:"default:'text'"` // text, image, video, audio, file
	MediaURL       string    `json:"media_url"`
	ReplyToID      *int      `json:"reply_to_id"`
	IsEdited       bool      `json:"is_edited" gorm:"default:false"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relationships
	Sender        User            `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
	Conversation  Conversation    `json:"conversation,omitempty" gorm:"foreignKey:ConversationID"`
	ReplyTo       *Message        `json:"reply_to,omitempty" gorm:"foreignKey:ReplyToID"`
	MessageStatus []MessageStatus `json:"message_status,omitempty"`
}

type MessageStatus struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	MessageID int       `json:"message_id" gorm:"not null"`
	UserID    int       `json:"user_id" gorm:"not null"`
	Status    string    `json:"status" gorm:"default:'sent'"` // sent, delivered, read
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Message Message `json:"message,omitempty" gorm:"foreignKey:MessageID"`
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
