package app

import (
	"context"
	"errors"
)

type ctxKey string

const (
	// userCtx is context key for User
	userCtx ctxKey = "user"
)

var (
	ErrUserCtxUnset = errors.New("app: user not set")
)

func SetUserCtx(parent context.Context, user User) context.Context {
	return context.WithValue(parent, userCtx, user)
}

func UserFromCtx(ctx context.Context) (User, error) {
	user, ok := ctx.Value(userCtx).(User)
	if ok {
		return user, nil
	}
	return User{}, ErrUserCtxUnset
}
