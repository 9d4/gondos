package app

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (app *App) CreateUser(ctx context.Context, user User) error {
	return app.d.UserStore.Add(ctx, user)
}

func (app *App) AuthEmail(ctx context.Context, email string, password string) (User, error) {
	user, err := app.d.UserStore.ByEmail(ctx, userConstructor, email)
	if err != nil {
		return user, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.cryptedPassword), []byte(password)); err != nil {
		return user, fmt.Errorf("%w: %w", ErrCredentialsIncorrect, err)
	}

	return user, nil
}
