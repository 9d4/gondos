package app

import (
	"context"
	"errors"
)

type ctxKey string

const (
	// userCtxKey is context key for User
	userCtxKey   ctxKey = "user"
	userIDCtxKey ctxKey = "userid"
)

var (
	ErrUserCtxUnset = errors.New("app: user not set")
)

func SetUserCtx(parent context.Context, user User) context.Context {
	return context.WithValue(parent, userCtxKey, user)
}

func UserFromCtx(ctx context.Context) (User, error) {
	user, ok := ctx.Value(userCtxKey).(User)
	if ok {
		return user, nil
	}
	return User{}, ErrUserCtxUnset
}

func SetUserIDCtx(parent context.Context, userID int64) context.Context {
	return context.WithValue(parent, userIDCtxKey, userID)
}

func UserIDFromCtx(ctx context.Context) (int64, error) {
	id, ok := ctx.Value(userIDCtxKey).(int64)
	if ok {
		return id, nil
	}
	return id, ErrUserCtxUnset
}
