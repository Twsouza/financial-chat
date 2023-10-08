package core

import (
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID       string `json:"id" valid:"notnull" gorm:"primaryKey;type:uuid"`
	Username string `json:"username" valid:"notnull"`
	Email    string `json:"email" valid:"notnull,email"`
	Password string `json:"password" valid:"notnull"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

// NewUser creates a new user, the ID is generated automatically
func NewUser(username, email, password string) (*User, error) {
	u := &User{
		Username: username,
		Email:    email,
		Password: password,
	}

	if err := u.prepare(); err != nil {
		return nil, err
	}

	if err := u.validate(); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) prepare() error {
	u.ID = uuid.NewV4().String()

	return nil
}

func (u *User) validate() error {
	_, err := govalidator.ValidateStruct(u)
	if err != nil {
		return err
	}

	return nil
}
