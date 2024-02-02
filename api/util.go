package api

import (
	"encoding/json"
	"net/http"
)

func parseJSON(r *http.Request, to any) error {
	return json.NewDecoder(r.Body).Decode(to)
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
