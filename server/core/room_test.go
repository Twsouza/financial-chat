package core_test

import (
	"server/core"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRoom(t *testing.T) {
	r, err := core.NewRoom("room name")
	require.Nil(t, err)
	require.NotNil(t, r)
	require.NotEmpty(t, r.ID)
	require.Equal(t, "room name", r.Name)
}

func TestNewRoomWithEmptyname(t *testing.T) {
	r, err := core.NewRoom("")
	require.Nil(t, r)
	require.Contains(t, err.Error(), "name: Missing required field")
}
