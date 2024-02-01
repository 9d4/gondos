package app

import (
	"context"
)

func (app *App) CreateUser(ctx context.Context, user User) error {
	return app.UserStore.Add(ctx, user)
}
