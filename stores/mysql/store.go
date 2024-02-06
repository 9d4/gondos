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
	t  tables
}

func NewUserStore(db *sql.DB, schema string) app.UserStore {
	return &UserStorage{
		db: db,
		t:  newTables(schema),
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

	stmt := s.t.users.INSERT(s.t.users.AllColumns).MODEL(model)
	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return s.handleErr(err)
	}
	tx.Commit()
	return nil
}

func (s *UserStorage) ByEmail(ctx context.Context, uc app.UserConstructor, email string) (app.User, error) {
	stmt := s.t.users.
		SELECT(s.t.users.AllColumns).
		WHERE(s.t.users.Email.EQ(mysql.String(email)))
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
	stmt := s.t.users.
		SELECT(s.t.users.AllColumns).
		WHERE(s.t.users.ID.EQ(mysql.Int64(id)))
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
	t  tables
}

func NewListStore(db *sql.DB, schema string) app.ListStore {
	return &ListStorage{
		db: db,
		t:  newTables(schema),
	}
}

func (s *ListStorage) getOnelistStmt(userID, listID int64) mysql.SelectStatement {
	return s.t.lists.
		SELECT(s.t.lists.UserID).
		WHERE(mysql.AND(
			s.t.lists.UserID.EQ(mysql.Int64(userID)),
			s.t.lists.ID.EQ(mysql.Int64(listID)),
		))
}

// CreateList implements app.ListStore.
func (s *ListStorage) CreateList(ctx context.Context, list app.List) error {
	tx, _ := s.db.Begin()
	userStmt := s.t.users.SELECT(s.t.users.ID).WHERE(s.t.users.ID.EQ(mysql.Int64(list.OwnerID())))
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

	insertStmt := s.t.lists.INSERT(s.t.lists.AllColumns).MODEL(model)
	if _, err := insertStmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}

// Lists implements app.ListStore.
func (s *ListStorage) Lists(ctx context.Context, listConstructor app.ListConstructor, userID int64) ([]app.List, error) {
	var listModels []model.Lists

	stmt := s.t.lists.SELECT(s.t.lists.AllColumns).WHERE(s.t.lists.UserID.EQ(mysql.Int64(userID)))
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
func (s *ListStorage) UpdateList(ctx context.Context, userID int64, listID int64, title, description string) error {
	tx, _ := s.db.Begin()
	if err := s.getOnelistStmt(userID, listID).QueryContext(ctx, tx, &model.Lists{}); err != nil {
		return err
	}

	updateStmt := s.t.lists.
		UPDATE(s.t.lists.Title, s.t.lists.Description, s.t.lists.UpdatedAt).
		SET(title, description, time.Now()).
		WHERE(s.t.lists.ID.EQ(mysql.Int64(listID)))
	if _, err := updateStmt.ExecContext(ctx, tx); err != nil {
		return err
	}
	return tx.Commit()
}

// DeleteList implements app.ListStore.
func (s *ListStorage) DeleteList(ctx context.Context, userID int64, listID int64) error {
	tx, _ := s.db.Begin()
	if err := s.getOnelistStmt(userID, listID).QueryContext(ctx, tx, &model.Lists{}); err != nil {
		return err
	}

	stmt := s.t.lists.DELETE().WHERE(s.t.lists.ID.EQ(mysql.Int64(listID)))
	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}
	stmt = table.ListItems.DELETE().WHERE(table.ListItems.ListID.EQ(mysql.Int64(listID)))
	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}

// AddItemToList implements app.ListStore.
func (s *ListStorage) AddItemToList(ctx context.Context, userID int64, item app.ListItem) error {
	tx, _ := s.db.Begin()

	if err := s.getOnelistStmt(userID, item.ListID()).QueryContext(ctx, tx, &model.Lists{}); err != nil {
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

// ListItems implements app.ListStore.
func (s *ListStorage) ListItems(ctx context.Context, listItemConstructor app.ListItemConstructor, userID, listID int64) ([]app.ListItem, error) {
	stmt := table.ListItems.SELECT(table.ListItems.AllColumns).
		FROM(table.ListItems.INNER_JOIN(s.t.lists, s.t.lists.ID.EQ(table.ListItems.ListID))).
		WHERE(mysql.AND(
			s.t.lists.ID.EQ(mysql.Int64(listID)),
			s.t.lists.UserID.EQ(mysql.Int64(userID)),
		))
	var models []model.ListItems
	if err := stmt.QueryContext(ctx, s.db, &models); err != nil {
		return nil, err
	}

	items := make([]app.ListItem, 0)
	for _, v := range models {
		items = append(items, listItemConstructor(v.ID, v.ListID, v.Body, v.CreatedAt, v.UpdatedAt))
	}

	return items, nil
}

func (s *ListStorage) getListItemStmt(itemID int64, userID int64, sel mysql.Projection, selects ...mysql.Projection) mysql.SelectStatement {
	return table.ListItems.
		SELECT(sel, selects...).
		FROM(
			table.ListItems.INNER_JOIN(s.t.lists, s.t.lists.ID.EQ(table.ListItems.ListID)),
		).
		WHERE(mysql.AND(
			table.ListItems.ID.EQ(mysql.Int64(itemID)),
			s.t.lists.UserID.EQ(mysql.Int64(userID)),
		))
}

// UpdateListItem implements app.ListStore.
func (s *ListStorage) UpdateListItem(ctx context.Context, userID int64, itemID int64, body string) error {
	tx, _ := s.db.Begin()
	stmt := s.getListItemStmt(itemID, userID, table.ListItems.ID)
	if err := stmt.QueryContext(ctx, tx, &model.ListItems{}); err != nil {
		return err
	}

	model := model.ListItems{
		Body:      body,
		UpdatedAt: time.Now(),
	}
	updateStmt := table.ListItems.
		UPDATE(table.ListItems.Body, table.ListItems.UpdatedAt).
		MODEL(model).
		WHERE(table.ListItems.ID.EQ(mysql.Int64(itemID)))
	if _, err := updateStmt.ExecContext(ctx, tx); err != nil {
		return err
	}
	return tx.Commit()
}

// DeleteListItem implements app.ListStore.
func (s *ListStorage) DeleteListItem(ctx context.Context, userID int64, itemID int64) error {
	tx, _ := s.db.Begin()
	stmt := s.getListItemStmt(itemID, userID, table.ListItems.ID)
	if err := stmt.QueryContext(ctx, tx, &model.ListItems{}); err != nil {
		return err
	}
	delStmt := table.ListItems.DELETE().WHERE(table.ListItems.ID.EQ(mysql.Int64(itemID)))
	if _, err := delStmt.ExecContext(ctx, tx); err != nil {
		return err
	}
	return tx.Commit()
}

var _ app.ListStore = &ListStorage{}

func ptr[T any](value T) *T {
	return &value
}
