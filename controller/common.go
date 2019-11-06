package controller

import (
	"encoding/json"
	"net/http"
)

type usernameContextKey struct{}

func getUsername(r *http.Request) string {
	return r.Context().Value(usernameContextKey{}).(string)
}

func writeJSON(w http.ResponseWriter, body interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(body)
}