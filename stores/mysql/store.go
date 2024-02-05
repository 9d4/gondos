package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

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

var getOnelistStmt = func(userID, listID int64) mysql.SelectStatement {
	return table.Lists.
		SELECT(table.Lists.UserID).
		WHERE(mysql.AND(
			table.Lists.UserID.EQ(mysql.Int64(userID)),
			table.Lists.ID.EQ(mysql.Int64(listID)),
		))
}

// CreateList implements app.ListStore.
func (s *ListStorage) CreateList(ctx context.Context, list app.List) error {
	tx, _ := s.db.Begin()
	userStmt := table.Users.SELECT(table.Users.ID).WHERE(table.Users.ID.EQ(mysql.Int64(list.OwnerID())))
	if err := userStmt.QueryContext(ctx, tx, &model.Users{}); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return app.ErrUserNotFound
		}
		return err
	}

	model := model.Lists{
		ID:          list.ID(),
		UserID:      ptr(list.OwnerID()),
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

// Lists implements app.ListStore.
func (s *ListStorage) Lists(ctx context.Context, listConstructor app.ListConstructor, userID int64) ([]app.List, error) {
	var listModels []model.Lists

	stmt := table.Lists.SELECT(table.Lists.AllColumns).WHERE(table.Lists.UserID.EQ(mysql.Int64(userID)))
	if err := stmt.QueryContext(ctx, s.db, &listModels); err != nil {
		return nil, err
	}

	var lists []app.List
	for _, v := range listModels {
		lists = append(lists, listConstructor(v.ID, *v.UserID, v.Title, *v.Description, v.CreatedAt, v.UpdatedAt))
	}

	return lists, nil
}

// UpdateList implements app.ListStore.
func (si *ListStorage) UpdateList(ctx context.Context, userID int64, listID int64, title, description string) error {
	tx, _ := si.db.Begin()
	if err := getOnelistStmt(userID, listID).QueryContext(ctx, tx, &model.Lists{}); err != nil {
		return err
	}

	updateStmt := table.Lists.
		UPDATE(table.Lists.Title, table.Lists.Description, table.Lists.UpdatedAt).
		SET(title, description, time.Now()).
		WHERE(table.Lists.ID.EQ(mysql.Int64(listID)))
	if _, err := updateStmt.ExecContext(ctx, tx); err != nil {
		return err
	}
	return tx.Commit()
}

// AddItemToList implements app.ListStore.
func (s *ListStorage) AddItemToList(ctx context.Context, userID int64, item app.ListItem) error {
	tx, _ := s.db.Begin()

	if err := getOnelistStmt(userID, item.ListID()).QueryContext(ctx, tx, &model.Lists{}); err != nil {
		return err
	}

	model := model.ListItems{
		ID:        item.ID(),
		ListID:    item.ListID(),
		Body:      item.Body(),
		CreatedAt: item.CreatedAt(),
		UpdatedAt: item.UpdatedAt(),
	}
	insertStmt := table.ListItems.INSERT(table.ListItems.AllColumns).MODEL(model)
	if _, err := insertStmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}

var _ app.ListStore = &ListStorage{}

func ptr[T any](value T) *T {
	return &value
}
