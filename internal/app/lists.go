package app

import "context"

func (app *App) UserCreateList(ctx context.Context, list List) error {
	userID, err := UserIDFromCtx(ctx)
	if err != nil {
		return err
	}
	list.ownerID = userID

	if err := app.d.ListStore.CreateList(ctx, list); err != nil {
		return err
	}
	return nil
}

func (app *App) UserLists(ctx context.Context) ([]List, error) {
	userID, err := UserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	return app.d.ListStore.Lists(ctx, listConstructor, userID)
}

func (app *App) UserUpdateList(ctx context.Context, listID int64, title, description string) error {
	userID, err := UserIDFromCtx(ctx)
	if err != nil {
		return err
	}

	return app.d.ListStore.UpdateList(ctx, userID, listID, title, description)
}

func (app *App) UserAddItemToList(ctx context.Context, item ListItem) error {
	userID, err := UserIDFromCtx(ctx)
	if err != nil {
		return err
	}

	return app.d.ListStore.AddItemToList(ctx, userID, item)
}
