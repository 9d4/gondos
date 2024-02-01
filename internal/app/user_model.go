package app

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/godruoyi/go-snowflake"
)

type UserStore interface {
	Add(context.Context, User) error
}

type User struct {
	id              int64
	name            string
	email           string
	cryptedPassword string
	createdAt       time.Time
	updatedAt       time.Time
}

func NewUser(name string, email string, password string) (User, error) {
	user := User{
		id:        int64(snowflake.ID()),
		name:      name,
		email:     email,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
	if err := user.ChangePassword(password); err != nil {
		return user, err
	}

	return user, nil
}

func (u User) ID() int64 {
	return u.id
}

func (u User) Name() string {
	return u.name
}

func (u User) Email() string {
	return u.email
}

func (u User) CryptedPassword() string {
	return u.cryptedPassword
}

func (u User) CreatedAt() time.Time {
	return u.createdAt
}

func (u User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) ChangePassword(password string) error {
	if len(password) < 8 {
		return newUserError("validaton.error", "password length should be at least 8 characters")
	}

	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	u.cryptedPassword = string(b)
	return nil
}
