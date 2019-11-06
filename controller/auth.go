package controller

import (
	"context"
	"github.com/Albert221/ReddigramApi/repository"
	"github.com/gbrlsnchs/jwt/v3"
	"log"
	"net/http"
	"strings"
	"time"
)

type AuthController struct {
	algorithm  jwt.Algorithm
	redditRepo repository.RedditRepository
}

func NewAuthController(algorithm jwt.Algorithm, redditRepo repository.RedditRepository) *AuthController {
	return &AuthController{
		algorithm:  algorithm,
		redditRepo: redditRepo,
	}
}

type jwtPayload struct {
	jwt.Payload
	RedditUsername string `json:"reddit_username"`
}

func (a *AuthController) AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	accessToken := r.PostFormValue("access_token")
	if accessToken == "" {
		writeJSON(w, map[string]string{"error": "access_token is required"}, http.StatusBadRequest)
		return
	}

	username, err := a.redditRepo.FetchUsername(accessToken)
	if err != nil {
		writeJSON(w, map[string]string{"error": "wrong access token"}, http.StatusBadRequest)
		return
	}

	now := time.Now()
	payload := jwtPayload{
		RedditUsername: username,
		Payload: jwt.Payload{
			ExpirationTime: jwt.NumericDate(now.Add(1 * time.Hour)),
			IssuedAt:       jwt.NumericDate(now),
		},
	}

	token, err := jwt.Sign(payload, a.algorithm)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, string(token), http.StatusOK)
}

func (a *AuthController) AuthenticationMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			h.ServeHTTP(w, r)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) < 2 || strings.ToLower(parts[0]) != "bearer" {
			writeJSON(w, map[string]string{"error": "authorization header is invalid"}, http.StatusUnauthorized)
			return
		}

		jwtToken := parts[1]

		var payload jwtPayload
		payloadValidator := jwt.ValidatePayload(&payload.Payload, jwt.ExpirationTimeValidator(time.Now()))
		_, err := jwt.Verify([]byte(jwtToken), a.algorithm, &payload, payloadValidator)
		if err != nil {
			writeJSON(w, map[string]string{"error": "given token is invalid or expired"}, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), usernameContextKey{}, payload.RedditUsername)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
