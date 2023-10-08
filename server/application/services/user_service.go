package services

import (
	"context"
	"fmt"
	"os"
	"server/application/dto"
	"server/application/repositories"
	"server/application/utils"
	"server/core"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserService struct {
	UserRepository repositories.UserRepository
}

var (
	// timeout_repository is the timeout for the repository layer, default to 5 seconds
	// can be set using TIMEOUT_REPOSITORY env var
	timeout_repository = 5
)

func init() {
	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err == nil {
		timeout_repository = timeout
	}
}

func NewUserService(ur repositories.UserRepository) *UserService {
	return &UserService{
		UserRepository: ur,
	}
}

func (us *UserService) CreateUser(c context.Context, req *dto.CreateUserReq) (*dto.CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, time.Duration(timeout_repository)*time.Second)
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
	ctx, cancel := context.WithTimeout(c, time.Duration(timeout_repository)*time.Second)
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
	ss, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
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
