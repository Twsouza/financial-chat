package dto

import "time"

type Message struct {
	Content   string    `json:"content"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"user"`
	RoomID    string    `json:"room_id"`
	CreatedAt time.Time `json:"created_at"`
}
