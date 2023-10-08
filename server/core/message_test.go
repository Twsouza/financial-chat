package core_test

import (
	"server/core"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMessage(t *testing.T) {
	m, err := core.NewMessage("message content", "ABC", "user", "123")
	require.Nil(t, err)
	require.NotNil(t, m)
	require.Equal(t, "message content", m.Content)
	require.Equal(t, "user", m.Username)
	require.Equal(t, "123", m.RoomID)
}

func TestNewMessageWithEmptyname(t *testing.T) {
	m, err := core.NewMessage("", "", "", "")
	require.Nil(t, m)
	require.Contains(t, err.Error(), "content: Missing required field")
	require.Contains(t, err.Error(), "user_id: Missing required field")
	require.Contains(t, err.Error(), "user: Missing required field")
	require.Contains(t, err.Error(), "room_id: Missing required field")
}
