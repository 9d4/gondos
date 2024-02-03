package store

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/mysql"

	"gondos/internal/app"
	"gondos/jetgen/gondos/model"
	"gondos/jetgen/gondos/table"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) app.UserStore {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) Add(ctx context.Context, user app.User) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	model := model.Users{
		ID:              user.ID(),
		Name:            user.Name(),
		Email:           user.Email(),
		CryptedPassword: user.CryptedPassword(),
		CreatedAt:       user.CreatedAt(),
		UpdatedAt:       user.UpdatedAt(),
	}

	stmt := table.Users.INSERT(table.Users.AllColumns).MODEL(model)
	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return s.handleErr(err)
	}
	tx.Commit()
	return nil
}

func (s *UserStorage) ByEmail(ctx context.Context, uc app.UserConstructor, email string) (app.User, error) {
	stmt := table.Users.
		SELECT(table.Users.AllColumns).
		WHERE(table.Users.Email.EQ(mysql.String(email)))
	var userModel model.Users
	if err := stmt.QueryContext(ctx, s.db, &userModel); err != nil {
		return app.User{}, s.handleErr(err)
	}

	return uc(
		userModel.ID,
		userModel.Name,
		userModel.Email,
		userModel.CryptedPassword,
		userModel.CreatedAt,
		userModel.UpdatedAt,
	), nil
}

// ByID implements app.UserStore.
func (s *UserStorage) ByID(ctx context.Context, uc app.UserConstructor, id int64) (app.User, error) {
	stmt := table.Users.
		SELECT(table.Users.AllColumns).
		WHERE(table.Users.ID.EQ(mysql.Int64(id)))
	var userModel model.Users
	if err := stmt.QueryContext(ctx, s.db, &userModel); err != nil {
		return app.User{}, s.handleErr(err)
	}

	return uc(
		userModel.ID,
		userModel.Name,
		userModel.Email,
		userModel.CryptedPassword,
		userModel.CreatedAt,
		userModel.UpdatedAt,
	), nil
}

// ensure implement app.UserStore
var _ app.UserStore = &UserStorage{}
