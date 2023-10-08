package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	u, err := NewUser("username", "email@mail.com", "password")
	require.Nil(t, err)
	require.NotNil(t, u)
}

func TestNewUser_WithInvalidEmail(t *testing.T) {
	u, err := NewUser("username", "invalid-email", "password")
	require.Nil(t, u)
	require.Contains(t, err.Error(), "email: invalid-email does not validate as email")
}

func TestNewUser_WithoutEmailPassword(t *testing.T) {
	u, err := NewUser("username", "", "")
	require.Nil(t, u)
	require.Contains(t, err.Error(), "email: Missing required field")
	require.Contains(t, err.Error(), "password: Missing required field")
}

func TestPrepare(t *testing.T) {
	u := &User{}
	err := u.prepare()
	require.Nil(t, err)
	require.NotEmpty(t, u.ID)
}

func TestUserValidation(t *testing.T) {
	u := &User{
		Username: "testuser",
		Email:    "invalidemail",
		Password: "testpassword",
	}

	err := u.validate()
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "email")
}
