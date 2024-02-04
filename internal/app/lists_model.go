package app

import (
	"context"
	"time"

	"github.com/godruoyi/go-snowflake"
)

type ListStore interface {
	CreateList(ctx context.Context, ownerID int64, list List) error
}

type List struct {
	id          int64
	owner       User
	title       string
	description string
	createdAt   time.Time
	updatedAt   time.Time
}

func (l List) ID() int64 {
	return l.id
}

func (l List) Owner() User {
	return l.owner
}

func (l List) Title() string {
	return l.title
}

func (l List) Description() string {
	return l.description
}

func (l List) CreatedAt() time.Time {
	return l.createdAt
}

func (l List) UpdatedAt() time.Time {
	return l.updatedAt
}

func NewList(title, description string) (List, error) {
	var list List

	if err := validateFields(validate, map[string]interface{}{
		"title":       title,
		"description": description,
	}, map[string]string{
		"title": "required,min=1,alphaunicode",
	}); err != nil {
		return list, err
	}

	list.id = int64(snowflake.ID())
	list.title = title
	list.description = description
	list.createdAt = time.Now()
	list.updatedAt = time.Now()

	return list, nil
}
