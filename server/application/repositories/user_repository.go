package repositories

import (
	"context"
	"server/application/utils"
	"server/core"

	"gorm.io/gorm"
)

type UserRepository interface {
	Insert(ctx context.Context, u *core.User) (*core.User, error)
	FindByEmail(ctx context.Context, email string) (*core.User, error)
}

type UserRepositoryImpl struct {
	Db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{Db: db}
}

func (ur UserRepositoryImpl) Insert(ctx context.Context, u *core.User) (*core.User, error) {
	hashedPwd, err := utils.HashPassword(u.Password)
	if err != nil {
		return nil, err
	}

	u.Password = hashedPwd
	if err := ur.Db.WithContext(ctx).Create(u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (ur UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*core.User, error) {
	var u core.User
	err := ur.Db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}

	return &u, nil
}
