package api

import (
	"errors"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"

	"gondos/app"
)

func NewHandler(app *app.App) http.Handler {
	h := &handler{app: app}
	router := httprouter.New()

	handle := func(fn Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			err := fn(w, r, p)
			if err != nil {
				h.errorHandler(w, r, p, err)
			}
		}
	}

	router.GET("/", handle(h.index))
	router.GET("/403", handle(h.thisShouldError))

	return router
}

type Handle func(http.ResponseWriter, *http.Request, httprouter.Params) error

type handler struct {
	app *app.App
}

func (h *handler) index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	_, err := w.Write([]byte{'H', 'e', 'l', 'l', 'o', ' ', 'W', 'o', 'r', 'l', 'd'})
	return err
}

func (h *handler) thisShouldError(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	return os.ErrPermission
}

func (h *handler) errorHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params, err error) {
	switch {
	case errors.Is(err, os.ErrPermission):
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("permission denied"))
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}
