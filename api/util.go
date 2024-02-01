package api

import (
	"encoding/json"
	"net/http"
)

func parseJSON(r *http.Request, to any) error {
	return json.NewDecoder(r.Body).Decode(to)
}
