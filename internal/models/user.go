package models

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID    `json:"id"`
	Password  string       `json:"password"`
	Cpf       string       `json:"cpf"`
	Email     string       `json:"email"`
	Phone     string       `json:"phone"`
	Name      string       `json:"name"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	NickName  string       `json:"nick_name"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

type UserRequest struct {
	Password  string `json:"password" validate:"required"`
	Cpf       string `json:"cpf" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Name      string `json:"name" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	NickName  string `json:"nick_name" validate:"required"`
}

type GetUserRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

type UserLoginRequest struct {
	Nickname string `json:"nickname" form:"nickname" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type UserResponse struct {
	ID        uuid.UUID    `json:"id"`
	Password  string       `json:"password,omitempty"`
	Cpf       string       `json:"cpf"`
	Email     string       `json:"email"`
	Phone     string       `json:"phone"`
	Name      string       `json:"name"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	NickName  string       `json:"nick_name"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at,omitempty"`
}
