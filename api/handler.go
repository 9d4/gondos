package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"

	"gondos/internal/app"
)

func NewHandler(app *app.App) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(middleware.StripSlashes)

	router.Use(jwtauth.Verifier(tokenAuth))
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rcx := chi.RouteContext(r.Context())
			fmt.Printf("%#v", rcx)
			h.ServeHTTP(w, r)
		})
	})

	handler := newServer(app)

	return HandlerFromMux(handler, router)
}

func newServer(app *app.App) ServerInterface {
	return &serverImpl{
		app: app,
	}
}

type serverImpl struct {
	app *app.App
}

// GetUser implements ServerInterface.
func (si *serverImpl) GetUser(w http.ResponseWriter, r *http.Request) {
	if err := authenticate(r); err != nil {
		si.deliverErr(w, r, err)
		return
	}

	user, err := si.app.AuthenticatedUser(r.Context())
	if err != nil {
		si.deliverErr(w, r, err)
		return
	}

	sendJSON(w, http.StatusOK, User{
		Id:        user.ID(),
		Email:     user.Email(),
		Name:      user.Name(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	})
}

// PostUserLists implements ServerInterface.
func (si *serverImpl) UserCreateList(w http.ResponseWriter, r *http.Request) {
	if err := authenticate(r); err != nil {
		si.deliverErr(w, r, err)
		return
	}

	var request ListCreateRequest
	if err := parseJSON(r, &request); err != nil {
		si.deliverErr(w, r, err)
		return
	}
	list, err := app.NewList(request.Title, request.Description)
	if err != nil {
		si.deliverErr(w, r, err)
		return
	}

	if err := si.app.UserCreateList(r.Context(), list); err != nil {
		si.deliverErr(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// UserLists implements ServerInterface.
func (si *serverImpl) UserLists(w http.ResponseWriter, r *http.Request) {
	if err := authenticate(r); err != nil {
		si.deliverErr(w, r, err)
		return
	}

	lists, err := si.app.UserLists(r.Context())
	if err != nil {
		si.deliverErr(w, r, err)
		return
	}

	res := make([]List, 0)
	for _, v := range lists {
		res = append(res, List{
			Id:          v.ID(),
			Title:       v.Title(),
			Description: v.Description(),
			CreatedAt:   v.CreatedAt(),
			UpdatedAt:   v.UpdatedAt(),
		})
	}

	sendJSON(w, http.StatusOK, res)
}

// UserListUpdate implements ServerInterface.
func (si *serverImpl) UserListUpdate(w http.ResponseWriter, r *http.Request, listId int) {
	if err := authenticate(r); err != nil {
		si.deliverErr(w, r, err)
		return
	}

	var request ListUpdateRequest
	if err := parseJSON(r, &request); err != nil {
		si.deliverErr(w, r, err)
		return
	}
	list, err := app.NewList(request.Title, request.Description)
	if err != nil {
		si.deliverErr(w, r, err)
		return
	}

	// request transformed to app.List even can just use request.*
	// because need to validate, and for now the list for creating new
	// and updating is the same, so just use NewList.

	if err := si.app.UserUpdateList(
		r.Context(),
		int64(listId),
		list.Title(),
		list.Description(),
	); err != nil {
		si.deliverErr(w, r, err)
		return
	}
}

// UserDeleteList implements ServerInterface.
func (si *serverImpl) UserDeleteList(w http.ResponseWriter, r *http.Request, listId int) {
	if err := authenticate(r); err != nil {
		si.deliverErr(w, r, err)
		return
	}

	if err := si.app.UserDeleteList(r.Context(), int64(listId)); err != nil {
		si.deliverErr(w, r, err)
		return
	}
}

// ListAddItem implements ServerInterface.
func (si *serverImpl) ListAddItem(w http.ResponseWriter, r *http.Request, listId int) {
	if err := authenticate(r); err != nil {
		si.deliverErr(w, r, err)
		return
	}

	var request AddItemToListRequest
	if err := parseJSON(r, &request); err != nil {
		si.deliverErr(w, r, err)
		return
	}

	item, err := app.NewListItem(int64(listId), request.Body)
	if err != nil {
		si.deliverErr(w, r, err)
		return
	}

	if err := si.app.UserAddItemToList(r.Context(), item); err != nil {
		si.deliverErr(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// ListItems implements ServerInterface.
func (si *serverImpl) ListItems(w http.ResponseWriter, r *http.Request, listId int) {
	if err := authenticate(r); err != nil {
		si.deliverErr(w, r, err)
		return
	}

	items, err := si.app.UserListItems(r.Context(), int64(listId))
	if err != nil {
		si.deliverErr(w, r, err)
		return
	}

	res := make([]ListItem, 0)
	for _, v := range items {
		res = append(res, ListItem{
			Id:        v.ID(),
			Body:      v.Body(),
			CreatedAt: v.CreatedAt(),
			UpdatedAt: v.UpdatedAt(),
		})
	}

	sendJSON(w, http.StatusOK, res)
}

// UpdateListItem implements ServerInterface.
func (si *serverImpl) UpdateListItem(w http.ResponseWriter, r *http.Request, listId int, itemId int) {
	if err := authenticate(r); err != nil {
		si.deliverErr(w, r, err)
		return
	}

	var request ListItemUpdateRequest
	if err := parseJSON(r, &request); err != nil {
		si.deliverErr(w, r, err)
		return
	}
	item, err := app.NewListItem(int64(listId), request.Body)
	if err != nil {
		si.deliverErr(w, r, err)
		return
	}

	if err := si.app.UserUpdateListItem(r.Context(), int64(itemId), item.Body()); err != nil {
		si.deliverErr(w, r, err)
		return
	}
}

// DeleteListItem implements ServerInterface.
func (si *serverImpl) DeleteListItem(w http.ResponseWriter, r *http.Request, listId int, itemId int) {
	if err := authenticate(r); err != nil {
		si.deliverErr(w, r, err)
		return
	}

	if err := si.app.UserDeleteListItem(r.Context(), int64(itemId)); err != nil {
		si.deliverErr(w, r, err)
		return
	}
}
