package api

import (
	"errors"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"gondos/internal/app"
)

func NewHandler(app *app.App) http.Handler {
	r := chi.NewRouter()
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

func (si serverImpl) deliverErr(w http.ResponseWriter, r *http.Request, err error) {
	var (
		appUserErr app.UserError
	)

	switch {
	case errors.As(err, &appUserErr):
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(appUserErr.Message()))
		return
	case errors.Is(err, os.ErrPermission):
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("forbidden"))
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
	log.Debug().Caller().Err(err).Send()
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
