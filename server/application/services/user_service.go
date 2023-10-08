package services

import (
	"context"
	"fmt"
	"server/application/dto"
	"server/application/repositories"
	"server/application/utils"
	"server/core"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserService struct {
	UserRepository repositories.UserRepository
}

const (
	// TODO: replace for an env var
	Timeout = 5
	// TODO: replace for an env var
	SECRET = "secret"
)

func NewUserService(ur repositories.UserRepository) *UserService {
	return &UserService{
		UserRepository: ur,
	}
}

func (us *UserService) CreateUser(c context.Context, req *dto.CreateUserReq) (*dto.CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, time.Duration(Timeout)*time.Second)
	defer cancel()

	u, err := core.NewUser(req.Username, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	_, err = us.UserRepository.Insert(ctx, u)
	if err != nil {
		return nil, err
	}

	return &dto.CreateUserRes{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}, nil
}

func (us *UserService) Login(c context.Context, req *dto.LoginUserReq) (*dto.LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, time.Duration(Timeout)*time.Second)
	defer cancel()

	u, err := us.UserRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if utils.VerifyPassword(u.Password, req.Password) != nil {
		return nil, fmt.Errorf("password is not correct: %w", err)
	}

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Issuer:    u.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(SECRET))
	if err != nil {
		return nil, err
	}

	return &dto.LoginUserRes{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		AccessToken: ss,
	}, nil
}
