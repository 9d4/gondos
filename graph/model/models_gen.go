// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type List struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ListInput struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
}

type ListItem struct {
	ID        string `json:"id"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ListItemInput struct {
	Body string `json:"body"`
}

type Mutation struct {
}

type Query struct {
}
