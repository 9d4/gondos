package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.43

import (
	"context"
	"gondos/graph/model"
	"gondos/internal/app"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

// CreateList is the resolver for the createList field.
func (r *mutationResolver) CreateList(ctx context.Context, list model.ListInput) (*string, error) {
	desc := ""
	if list.Description != nil {
		desc = *list.Description
	}
	input, err := app.NewList(list.Title, desc)
	if err != nil {
		return nil, err
	}

	if err := r.App.UserCreateList(ctx, input); err != nil {
		return nil, err
	}
	return nil, nil
}

// UpdateList is the resolver for the updateList field.
func (r *mutationResolver) UpdateList(ctx context.Context, listID string, list model.ListInput) (*string, error) {
	listIDInt, err := strconv.Atoi(listID)
	if err != nil {
		return nil, err
	}
	desc := ""
	if list.Description != nil {
		desc = *list.Description
	}
	input, err := app.NewList(list.Title, desc)
	if err != nil {
		return nil, err
	}

	if err := r.App.UserUpdateList(ctx, int64(listIDInt), input.Title(), input.Description()); err != nil {
		return nil, err
	}
	return nil, nil
}

// DeleteList is the resolver for the deleteList field.
func (r *mutationResolver) DeleteList(ctx context.Context, listID string) (*string, error) {
	listIDInt, err := strconv.Atoi(listID)
	if err != nil {
		return nil, err
	}

	if err := r.App.UserDeleteList(ctx, int64(listIDInt)); err != nil {
		return nil, err
	}
	return nil, nil
}

// AddItemToList is the resolver for the addItemToList field.
func (r *mutationResolver) AddItemToList(ctx context.Context, listID string, item model.ListItemInput) (*string, error) {
	listIDInt, err := strconv.Atoi(listID)
	if err != nil {
		return nil, err
	}
	input, err := app.NewListItem(int64(listIDInt), item.Body)
	if err != nil {
		return nil, err
	}

	if err := r.App.UserAddItemToList(ctx, input); err != nil {
		return nil, err
	}
	return nil, nil
}

// UpdateItem is the resolver for the updateItem field.
func (r *mutationResolver) UpdateItem(ctx context.Context, itemID string, item model.ListItemInput) (*string, error) {
	itemIDInt, err := strconv.Atoi(itemID)
	if err != nil {
		return nil, err
	}
	input, err := app.NewListItem(1, item.Body)
	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	if err := r.App.UserUpdateListItem(ctx, int64(itemIDInt), input.Body()); err != nil {
		return nil, err
	}
	return nil, nil
}

// DeleteItem is the resolver for the deleteItem field.
func (r *mutationResolver) DeleteItem(ctx context.Context, itemID string) (*string, error) {
	itemIDInt, err := strconv.Atoi(itemID)
	if err != nil {
		return nil, err
	}

	if err := r.App.UserDeleteListItem(ctx, int64(itemIDInt)); err != nil {
		return nil, err
	}
	return nil, nil
}

// Lists is the resolver for the lists field.
func (r *queryResolver) Lists(ctx context.Context) ([]*model.List, error) {
	lists, err := r.App.UserLists(ctx)
	if err != nil {
		return nil, err
	}

	modelLists := make([]*model.List, 0)
	for _, v := range lists {
		modelLists = append(modelLists, &model.List{
			ID:          strconv.Itoa(int(v.ID())),
			Title:       v.Title(),
			Description: ptr(v.Description()),
			CreatedAt:   v.CreatedAt().Format(time.RFC3339Nano),
			UpdatedAt:   v.UpdatedAt().Format(time.RFC3339Nano),
		})
	}

	return modelLists, nil
}

// ListItems is the resolver for the listItems field.
func (r *queryResolver) ListItems(ctx context.Context, listID string) ([]*model.ListItem, error) {
	listIDInt, err := strconv.Atoi(listID)
	if err != nil {
		return nil, err
	}

	items, err := r.App.UserListItems(ctx, int64(listIDInt))
	if err != nil {
		return nil, err
	}

	modelListItems := make([]*model.ListItem, 0)
	for _, v := range items {
		modelListItems = append(modelListItems, &model.ListItem{
			ID:        strconv.Itoa(int(v.ID())),
			Body:      v.Body(),
			CreatedAt: v.CreatedAt().Format(time.RFC3339Nano),
			UpdatedAt: v.CreatedAt().Format(time.RFC3339Nano),
		})
	}

	return modelListItems, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
