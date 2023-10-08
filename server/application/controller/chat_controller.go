package controller

import (
	"net/http"
	"server/application/dto"
	"server/application/services"
	"server/core"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatController struct {
	hub *services.Hub
}

func NewChatController(hub *services.Hub) ChatController {
	return ChatController{
		hub: hub,
	}
}

func (ct *ChatController) CreateRoom(gc *gin.Context) {
	req := dto.CreateRoomReq{}
	if err := gc.Bind(&req); err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, err := core.NewRoom(req.Name)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ct.hub.Rooms[room.ID] = room
	ct.hub.Rooms[room.ID].Clients = make(map[string]*core.Client)
	res := &dto.CreateRoomRes{
		ID:   room.ID,
		Name: room.Name,
	}

	gc.JSON(http.StatusCreated, res)
}

func (ct *ChatController) GetRooms(gc *gin.Context) {
	rooms := make([]dto.CreateRoomRes, 0)

	for _, room := range ct.hub.Rooms {
		rooms = append(rooms, dto.CreateRoomRes{
			ID:   room.ID,
			Name: room.Name,
		})
	}

	gc.JSON(http.StatusOK, rooms)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ct *ChatController) JoinRoom(gc *gin.Context) {
	roomId := gc.Param("id")
	clientId := gc.Query("userId")
	clientName := gc.Query("username")

	conn, err := upgrader.Upgrade(gc.Writer, gc.Request, nil)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "error setting up connection"})
		return
	}

	client, err := core.NewClient(conn, ct.hub.Unregister, clientId, clientName, roomId)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "error joining the client in the room"})
		return
	}

	ct.hub.Register <- client

	go client.SendMessage()
	client.ReceiveMessage(ct.hub.Broadcast)
}

func (ct *ChatController) GetClients(c *gin.Context) {
	var clients []core.Client
	roomId := c.Param("id")

	if _, ok := ct.hub.Rooms[roomId]; !ok {
		clients = make([]core.Client, 0)
		c.JSON(http.StatusOK, clients)
	}

	for _, c := range ct.hub.Rooms[roomId].Clients {
		clients = append(clients, core.Client{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}
