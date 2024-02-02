package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"gondos/internal/app"
)

func NewHandler(app *app.App) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	handler := newServer(app)

	return HandlerFromMux(handler, r)
}

func newServer(app *app.App) ServerInterface {
	return &serverImpl{
		app: app,
	}
}

type serverImpl struct {
	app *app.App
}

// Register a new account
// (POST /auth/register)
func (si serverImpl) AuthRegister(w http.ResponseWriter, r *http.Request) {
	var request AuthRegisterRequest
	if err := parseJSON(r, &request); err != nil {
		si.deliverErr(w, r, err)
		return
	}

	user, err := app.NewUser(request.Name, request.Email, request.Password)
	if err != nil {
		si.deliverErr(w, r, err)
		return
	}

	if err := si.app.CreateUser(r.Context(), user); err != nil {
		si.deliverErr(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
