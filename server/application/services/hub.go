package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"queue"
	"regexp"
	"server/core"
	"time"

	"github.com/streadway/amqp"
)

type Hub struct {
	Broadcast     chan *core.Message
	MessageBroker *queue.RabbitMQ
	Queue         chan amqp.Delivery
	Register      chan *core.Client
	Rooms         map[string]*core.Room
	Unregister    chan *core.Client
}

func NewHub(messageBroker *queue.RabbitMQ) (*Hub, error) {
	return &Hub{
		Broadcast:     make(chan *core.Message, 10),
		MessageBroker: messageBroker,
		Queue:         make(chan amqp.Delivery),
		Rooms:         make(map[string]*core.Room),
		Register:      make(chan *core.Client),
		Unregister:    make(chan *core.Client),
	}, nil
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			if _, ok := h.Rooms[cl.RoomID]; ok {
				r := h.Rooms[cl.RoomID]

				if _, ok := r.Clients[cl.ID]; !ok {
					r.Clients[cl.ID] = cl
					break
				}
			}

		case cl := <-h.Unregister:
			if _, ok := h.Rooms[cl.RoomID]; ok {
				if _, ok := h.Rooms[cl.RoomID].Clients[cl.ID]; ok {
					delete(h.Rooms[cl.RoomID].Clients, cl.ID)
					close(cl.Message)
				}
			}

		case m := <-h.Broadcast:
			if m == nil {
				continue
			}

			if isCommand(m.Content) {
				err := h.SendCommand(m)
				if err != nil {
					h.notifyUserAboutCommandError(m, err)
				}
				continue
			}

			if _, ok := h.Rooms[m.RoomID]; ok {
				for _, cl := range h.Rooms[m.RoomID].Clients {
					cl.Message <- m
				}
			}

		case q := <-h.Queue:
			m := &core.Message{}
			err := json.Unmarshal(q.Body, m)
			if err != nil {
				log.Printf("error unmarshalling message: %v", err)
				continue
			}

			if _, ok := h.Rooms[m.RoomID]; ok {
				for _, cl := range h.Rooms[m.RoomID].Clients {
					cl.Message <- m
				}
			}

			q.Ack(false)
		}
	}
}

func (h *Hub) SendCommand(m *core.Message) error {
	body, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return h.MessageBroker.Notify(
		string(body),
		"application/json",
		os.Getenv("RABBITMQ_NOTIFICATION_EX"),
		os.Getenv("RABBITMQ_NOTIFICATION_ROUTING_KEY"),
	)
}

func (h *Hub) notifyUserAboutCommandError(m *core.Message, err error) {
	log.Printf("error sending command to the message broker: %v\n", err)
	errorMessage := &core.Message{
		Content:   fmt.Sprintf(`error sending command "%s" to the message broker`, m.Content),
		Username:  "bot",
		RoomID:    m.RoomID,
		CreatedAt: time.Now(),
	}

	if _, ok := h.Rooms[m.RoomID]; ok {
		for _, cl := range h.Rooms[m.RoomID].Clients {
			if cl.ID == m.UserID {
				cl.Message <- errorMessage
			}
		}
	}
}

// isCommand returns true if the content follows the command pattern:
// /{command}={args}
func isCommand(content string) bool {
	pattern := `^/[^=]+=.+$`
	return regexp.MustCompile(pattern).MatchString(content)
}
