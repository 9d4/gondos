package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"

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

type ListStorage struct {
	db *sql.DB
}

func NewListStore(db *sql.DB) app.ListStore {
	return &ListStorage{db: db}
}

// CreateList implements app.ListStore.
func (s *ListStorage) CreateList(ctx context.Context, ownerID int64, list app.List) error {
	tx, _ := s.db.Begin()
	userStmt := table.Users.SELECT(table.Users.ID).WHERE(table.Users.ID.EQ(mysql.Int64(ownerID)))
	if err := userStmt.QueryContext(ctx, tx, &model.Users{}); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return app.ErrUserNotFound
		}
		return err
	}

	model := model.Lists{
		ID:          list.ID(),
		UserID:      &ownerID,
		Title:       list.Title(),
		Description: ptr(list.Description()),
		CreatedAt:   list.CreatedAt(),
		UpdatedAt:   list.UpdatedAt(),
	}

	insertStmt := table.Lists.INSERT(table.Lists.AllColumns).MODEL(model)
	if _, err := insertStmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}

var _ app.ListStore = &ListStorage{}

func ptr[T any](value T) *T {
	return &value
}
