package core_test

import (
	"server/core"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	c, err := core.NewClient(&websocket.Conn{}, nil, "123", "user", "456")
	require.Nil(t, err)
	require.NotNil(t, c)
	require.Equal(t, "123", c.ID)
	require.Equal(t, "user", c.Username)
	require.Equal(t, "456", c.RoomID)
}

func TestNewClientWithEmptyname(t *testing.T) {
	c, err := core.NewClient(&websocket.Conn{}, nil, "", "", "")
	require.Nil(t, c)
	require.Contains(t, err.Error(), "id: Missing required field")
	require.Contains(t, err.Error(), "room_id: Missing required field")
	require.Contains(t, err.Error(), "user: Missing required field")
}
