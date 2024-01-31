package entity

import (
	"time"

	"gondos/jetgen/gondos/model"
	"gondos/jetgen/gondos/table"
)

var ListsTable = table.Lists

type Lists = model.Lists

type ListFillable struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type List struct {
	ListFillable
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var ListItemsTable = table.ListItems

type ListItems = model.ListItems

type ListItemFillable struct {
	Body string `json:"body"`
}

type ListItem struct {
	ListItemFillable
	ID        int64     `json:"id"`
	ListID    int64     `json:"list_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
