package core

import (
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type Room struct {
	ID      string             `json:"id" valid:"-"`
	Name    string             `json:"name" valid:"notnull"`
	Clients map[string]*Client `json:"clients" valid:"-"`
}

func NewRoom(name string) (*Room, error) {
	r := &Room{
		Name: name,
	}

	if err := r.prepare(); err != nil {
		return nil, err
	}

	if err := r.Validate(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Room) prepare() error {
	r.ID = uuid.NewV4().String()

	return nil
}

func (r *Room) Validate() error {
	_, err := govalidator.ValidateStruct(r)
	if err != nil {
		return err
	}

	return nil
}
