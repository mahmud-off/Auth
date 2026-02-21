package domain

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type User struct {
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RegisteredAt time.Time `json:"registered_at"`
}

type SignUpInput struct {
	Name     string `json:"name" db:"name" validate:"required,gte=2"`
	Email    string `json:"email" db:"name" validate:"required,email"`
	Password string `json:"password" db:"name" validate:"required,gte=8"`
}

func (i SignUpInput) Validate() error {
	return validate.Struct(i)
}
