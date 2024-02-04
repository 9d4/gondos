package app

import "context"

func (app *App) UserCreateList(ctx context.Context, list List) error {
	userID, err := UserIDFromCtx(ctx)
	if err != nil {
		return err
	}

	if err := app.d.ListStore.CreateList(ctx, userID, list); err != nil {
		return err
	}
	return nil
}
