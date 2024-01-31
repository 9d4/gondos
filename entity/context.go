package entity

import (
	"context"
	"errors"
)

type ctxKey string

const (
	userCtxKey ctxKey = "user"
)

var ErrNoUser = errors.New("app: no user set")

// NewUserContext returns new context containing user.
func NewUserContext(parent context.Context, user any) context.Context {
	return context.WithValue(parent, userCtxKey, user)
}

// UserFromContext returns user from context.
// If user unset, NoUserErr will be returned.
func UserFromContext(ctx context.Context) (any, error) {
	return nil, ErrNoUser
}
