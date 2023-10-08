package core

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type Message struct {
	Content   string    `json:"content" valid:"notnull"`
	UserID    string    `json:"user_id" valid:"notnull"`
	Username  string    `json:"user" valid:"notnull"`
	RoomID    string    `json:"room_id" valid:"notnull"`
	CreatedAt time.Time `json:"created_at" valid:"-"`
}

func NewMessage(content, userId, username, roomId string) (*Message, error) {
	m := &Message{
		Content:   content,
		UserID:    userId,
		Username:  username,
		RoomID:    roomId,
		CreatedAt: time.Now(),
	}

	if err := m.Validate(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Message) Validate() error {
	_, err := govalidator.ValidateStruct(m)
	if err != nil {
		return err
	}

	return nil
}
