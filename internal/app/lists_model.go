package app

import (
	"context"
	"time"

	"github.com/godruoyi/go-snowflake"
)

type ListStore interface {
	CreateList(ctx context.Context, list List) error
	Lists(ctx context.Context, listConstructor ListConstructor, userID int64) ([]List, error)
	UpdateList(ctx context.Context, userID, listID int64, title, description string) error
	DeleteList(ctx context.Context, userID, listID int64) error
	AddItemToList(ctx context.Context, userID int64, item ListItem) error
	ListItems(ctx context.Context, listItemConstructor ListItemConstructor, userID, listID int64) ([]ListItem, error)
	UpdateListItem(ctx context.Context, userID, itemID int64, body string) error
	DeleteListItem(ctx context.Context, userID, itemID int64) error
}

type ListConstructor func(id, ownerID int64, title, description string, createdAt, updatedAt time.Time) List

var listConstructor ListConstructor = func(id, ownerID int64, title, description string, createdAt, updatedAt time.Time) List {
	return List{
		id:          id,
		ownerID:     ownerID,
		title:       title,
		description: description,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

type List struct {
	id          int64
	ownerID     int64
	title       string
	description string
	createdAt   time.Time
	updatedAt   time.Time
}

func (l List) ID() int64 {
	return l.id
}

func (l List) OwnerID() int64 {
	return l.ownerID
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

type ListItemConstructor func(id, listID int64, body string, createdAt, updatedAt time.Time) ListItem

var listItemConstructor ListItemConstructor = func(id, listID int64, body string, createdAt, updatedAt time.Time) ListItem {
	return ListItem{
		id:        id,
		listID:    listID,
		body:      body,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

type ListItem struct {
	id        int64
	listID    int64
	body      string
	createdAt time.Time
	updatedAt time.Time
}

func (i ListItem) ID() int64 {
	return i.id
}

func (i ListItem) ListID() int64 {
	return i.listID
}

func (i ListItem) Body() string {
	return i.body
}

func (i ListItem) CreatedAt() time.Time {
	return i.createdAt
}

func (i ListItem) UpdatedAt() time.Time {
	return i.updatedAt
}

func NewListItem(listID int64, body string) (ListItem, error) {
	var item ListItem

	if err := validateFields(validate, map[string]interface{}{
		"list_id": listID,
		"body":    body,
	}, map[string]string{
		"list_id": "required,gt=0",
		"body":    "required,min=1",
	}); err != nil {
		return item, err
	}
	item.id = int64(snowflake.ID())
	item.listID = listID
	item.body = body
	item.createdAt = time.Now()
	item.updatedAt = time.Now()

	return item, nil
}
