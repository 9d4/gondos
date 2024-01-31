package entity

import (
	"time"

	"gondos/jetgen/gondos/model"
	"gondos/jetgen/gondos/table"
)

var UsersTable = table.Users

type Users = model.Users

type UserFillable struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type User struct {
	UserFillable
	ID              int64     `json:"id"`
	CryptedPassword string    `json:"-"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
