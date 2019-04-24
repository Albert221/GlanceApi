package main

import (
	"encoding/json"
	"errors"
	"github.com/gbrlsnchs/jwt/v3"
	"net/http"
	"time"
)

func handleError(err error, w http.ResponseWriter, statusCode int) {
	body := map[string]string{"error": err.Error()}

	w.WriteHeader(statusCode)
	sendResponse(body, w)
}

func sendResponse(body interface{}, w http.ResponseWriter) {
	response, _ := json.Marshal(body)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		handleError(err, w, http.StatusBadRequest)
		return
	}

	accessToken := r.PostFormValue("access_token")
	if accessToken == "" {
		handleError(errors.New("access_token is required"), w, http.StatusBadRequest)
		return
	}

	username, err := FetchUsername(accessToken)
	if err != nil {
		handleError(err, w, http.StatusForbidden)
		return
	}

	now := time.Now()

	hs256 := jwt.NewHMAC(jwt.SHA256, []byte("some secret"))
	payload := jwt.Payload{
		Subject:        username,
		ExpirationTime: now.Add(time.Hour).Unix(),
		IssuedAt:       now.Unix(),
	}

	token, err := jwt.Sign(jwt.Header{}, payload, hs256)
	if err != nil {
		handleError(err, w, http.StatusInternalServerError)
		return
	}

	sendResponse(map[string]interface{}{"token": string(token)}, w)
}
