package core

import (
	"log"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/websocket"
	"github.com/microcosm-cc/bluemonday"
)

var (
	p = bluemonday.StripTagsPolicy()
)

type Client struct {
	Conn       *websocket.Conn `json:"-" valid:"-"`
	Message    chan *Message   `json:"-" valid:"-"`
	Unregister chan *Client    `json:"-" valid:"-"`
	ID         string          `json:"id" valid:"notnull"`
	RoomID     string          `json:"room_id,omitempty" valid:"notnull"`
	Username   string          `json:"user" valid:"notnull"`
}

func NewClient(conn *websocket.Conn, unregisterChan chan *Client, userId, username, roomId string) (*Client, error) {
	c := &Client{
		Conn:       conn,
		Message:    make(chan *Message, 10),
		Unregister: unregisterChan,
		ID:         userId,
		RoomID:     roomId,
		Username:   username,
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) Validate() error {
	_, err := govalidator.ValidateStruct(c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SendMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		c.Conn.WriteJSON(message)
	}
}

func (c *Client) ReceiveMessage(broadcast chan *Message) {
	defer func() {
		c.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg, _ := NewMessage(p.Sanitize(string(m)), c.ID, c.Username, c.RoomID)

		broadcast <- msg
	}
}
