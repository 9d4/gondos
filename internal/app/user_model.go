package app

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/godruoyi/go-snowflake"
)

type UserStore interface {
	Add(context.Context, User) error
	ByEmail(ctx context.Context, uc UserConstructor, email string) (User, error)
	ByID(ctx context.Context, uc UserConstructor, id int64) (User, error)
}

type User struct {
	id              int64
	name            string
	email           string
	cryptedPassword string
	createdAt       time.Time
	updatedAt       time.Time
}

// UserConstructor makes possible package outside app to create User
// under control of this package.
type UserConstructor func(id int64, name, email, cryptedPassword string, createdAt, updatedAt time.Time) User

var userConstructor UserConstructor = func(id int64, name, email, cryptedPassword string, createdAt, updatedAt time.Time) User {
	return User{id: id, name: name, email: email, cryptedPassword: cryptedPassword, createdAt: createdAt, updatedAt: updatedAt}
}

func NewUser(name string, email string, password string) (User, error) {
	var user User

	if err := validateFields(validate, map[string]interface{}{
		"name":     name,
		"email":    email,
		"password": password,
	}, map[string]string{
		"name":     "required,alpha,min=3",
		"email":    "required,email",
		"password": "required,min=8",
	}); err != nil {
		return user, err
	}

	user.id = int64(snowflake.ID())
	user.name = name
	user.email = email
	user.createdAt = time.Now()
	user.updatedAt = time.Now()

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
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	u.cryptedPassword = string(b)
	return nil
}
