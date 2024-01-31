package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func NewHandler() http.Handler {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte{'H', 'e', 'l', 'l', 'o', ' ', 'W', 'o', 'r', 'l', 'd'})
	})

	return router
}
