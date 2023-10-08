package dto

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateRoomReq struct {
	Name string `json:"name" form:"name" binding:"required"`
}
