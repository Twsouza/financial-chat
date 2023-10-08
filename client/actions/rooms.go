package actions

import (
	"client/dto"
	"client/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gobuffalo/buffalo"
)

func RoomsHandler(c buffalo.Context) error {
	res, err := utils.CallApi(c, "/ws/rooms", http.MethodGet, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	rooms := []dto.RoomRes{}
	if err := json.NewDecoder(res.Body).Decode(&rooms); err != nil {
		return err
	}

	c.Set("rooms", rooms)

	return c.Render(http.StatusOK, r.HTML("room/index.plush.html"))
}

func CreateRoomHandler(c buffalo.Context) error {
	roomReq := dto.CreateRoomReq{}
	if err := c.Bind(&roomReq); err != nil {
		return err
	}

	body, err := json.Marshal(roomReq)
	if err != nil {
		return err
	}

	res, err := utils.CallApi(c, "/ws/rooms", http.MethodPost, body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		// convert response body to string
		var bodyBytes []byte
		if res.Body != nil {
			bodyBytes, _ = io.ReadAll(res.Body)
		}

		log.Printf("response: %v", string(bodyBytes))
		c.Flash().Add("danger", "Error creating room")
	}

	return c.Redirect(http.StatusSeeOther, "/rooms")
}

func JoinRoomHandler(c buffalo.Context) error {
	roomID := c.Param("room_id")
	userId, _ := c.Cookies().Get("userId")
	username, _ := c.Cookies().Get("username")

	c.Set("roomID", roomID)
	c.Set("userID", userId)
	c.Set("username", username)
	c.Set("wsURL", fmt.Sprintf("%s/ws/rooms/%s/join", os.Getenv("WS_URL"), roomID))

	return c.Render(http.StatusOK, r.HTML("room/room.plush.html"))
}
