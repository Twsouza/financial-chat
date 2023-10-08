package services

import (
	"server/core"
)

type Hub struct {
	Broadcast  chan *core.Message
	Register   chan *core.Client
	Rooms      map[string]*core.Room
	Unregister chan *core.Client
}

func NewHub() (*Hub, error) {
	return &Hub{
		Broadcast:  make(chan *core.Message, 10),
		Rooms:      make(map[string]*core.Room),
		Register:   make(chan *core.Client),
		Unregister: make(chan *core.Client),
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
			if _, ok := h.Rooms[m.RoomID]; ok {
				for _, cl := range h.Rooms[m.RoomID].Clients {
					cl.Message <- m
				}
			}
		}
	}
}
