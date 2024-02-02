package api

import (
	"net/http"
	"strconv"
	"time"

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

// Login
// (POST /auth/login)
func (si serverImpl) AuthLogin(w http.ResponseWriter, r *http.Request) {
	var request AuthLoginRequest
	if err := parseJSON(r, &request); err != nil {
		si.deliverErr(w, r, err)
		return
	}

	user, err := si.app.AuthEmail(r.Context(), request.Email, request.Password)
	if err != nil {
		si.deliverErr(w, r, err)
		return
	}

	_, tk, err := tokenAuth.Encode(map[string]interface{}{
		"sub": strconv.Itoa(int(user.ID())),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(10 * time.Minute).Unix(),
	})
	if err != nil {
		si.deliverErr(w, r, err)
		return
	}

	sendJSON(w, http.StatusOK, AuthLoginResponse{
		AccessToken: tk,
	})
}
