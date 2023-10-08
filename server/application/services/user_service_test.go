package services_test

import (
	"context"
	"os"
	"server/application/dto"
	mock_repositories "server/application/repositories/mock"
	"server/application/services"
	"server/application/utils"
	"server/core"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserService_CreateUser(t *testing.T) {
	ctx := context.Background()
	req := &dto.CreateUserReq{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "testpassword",
	}
	expected := &core.User{
		Username: req.Username,
		Email:    req.Email,
	}

	ctrl := gomock.NewController(t)
	mock := mock_repositories.NewMockUserRepository(ctrl)
	mock.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expected, nil)

	us := services.NewUserService(mock)
	res, err := us.CreateUser(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, res.ID)
	assert.Equal(t, expected.Username, res.Username)
	assert.Equal(t, expected.Email, res.Email)
}

func TestUserService_Login(t *testing.T) {
	ctx := context.Background()
	loginReq := &dto.LoginUserReq{
		Email:    "testuser@example.com",
		Password: "testpassword",
	}
	hashPwd, err := utils.HashPassword(loginReq.Password)
	if err != nil {
		t.Fatal(err)
	}
	expected := &core.User{
		Username: "testuser",
		Email:    loginReq.Email,
		Password: hashPwd,
	}

	ctrl := gomock.NewController(t)
	mock := mock_repositories.NewMockUserRepository(ctrl)
	mock.EXPECT().FindByEmail(gomock.Any(), loginReq.Email).Return(expected, nil)

	us := services.NewUserService(mock)
	res, err := us.Login(ctx, loginReq)
	assert.NoError(t, err)
	assert.NotNil(t, res.AccessToken)
	assert.NotNil(t, res.ID)
	assert.Equal(t, expected.Email, res.Email)
	assert.Equal(t, expected.Username, res.Username)
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	ctx := context.Background()
	loginReq := &dto.LoginUserReq{
		Email:    "testuser@example.com",
		Password: "testpassword",
	}
	hashPwd, err := utils.HashPassword("otherpassword")
	if err != nil {
		t.Fatal(err)
	}
	expected := &core.User{
		Username: "testuser",
		Email:    loginReq.Email,
		Password: hashPwd,
	}

	ctrl := gomock.NewController(t)
	mock := mock_repositories.NewMockUserRepository(ctrl)
	mock.EXPECT().FindByEmail(gomock.Any(), loginReq.Email).Return(expected, nil)

	us := services.NewUserService(mock)
	_, err = us.Login(ctx, loginReq)
	assert.Error(t, err)
}

// TestMain Setup the tests
func TestMain(m *testing.M) {
	// set JWT_SECRET env var, it will be used by Login
	os.Setenv("JWT_SECRET", "testjwtsecret")
	// run tests
	exitCode := m.Run()
	// teardown
	os.Unsetenv("JWT_SECRET")
	// exit
	os.Exit(exitCode)
}
