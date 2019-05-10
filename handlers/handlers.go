package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Albert221/ReddigramApi/reddit"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

const usernameKey = "username"

type Controller struct {
	db          *sqlx.DB
	jwtSigner   jwt.Signer
	jwtVerifier jwt.Verifier
}

func NewController(db *sqlx.DB, secret string) *Controller {
	hs256 := jwt.NewHMAC(jwt.SHA256, []byte(secret))

	return &Controller{
		db:          db,
		jwtSigner:   hs256,
		jwtVerifier: hs256,
	}
}

func (c *Controller) handleError(err error, w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)

	if err != nil {
		body := map[string]string{"error": err.Error()}
		c.sendResponse(body, w)
	}
}

func (c *Controller) sendResponse(body interface{}, w http.ResponseWriter) {
	response, _ := json.Marshal(body)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (c *Controller) getUsername(r *http.Request) string {
	return r.Context().Value(usernameKey).(string)
}

func (c *Controller) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fragments := strings.Split(authHeader, " ")
		if len(fragments) != 2 || strings.ToLower(fragments[0]) != "bearer" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		token := fragments[1]
		rawToken, err := jwt.Parse([]byte(token))
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if rawToken.Verify(c.jwtVerifier) != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		var payload jwt.Payload
		_, err = rawToken.Decode(&payload)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if payload.Validate(jwt.ExpirationTimeValidator(time.Now(), false)) != nil {
			c.handleError(errors.New("token is expired"), w, http.StatusForbidden)
			return
		}

		username := payload.Subject
		ctx := context.WithValue(r.Context(), usernameKey, username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *Controller) createToken(username string) ([]byte, error) {
	now := time.Now()

	payload := jwt.Payload{
		Subject:        username,
		ExpirationTime: now.Add(1 * time.Hour).Unix(),
		IssuedAt:       now.Unix(),
	}

	return jwt.Sign(jwt.Header{}, payload, c.jwtSigner)
}

func (c *Controller) AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		c.handleError(err, w, http.StatusBadRequest)
		return
	}

	accessToken := r.PostFormValue("access_token")
	if accessToken == "" {
		c.handleError(errors.New("access_token is required"), w, http.StatusBadRequest)
		return
	}

	username, err := reddit.FetchUsername(accessToken)
	if err != nil {
		c.handleError(errors.New("wrong access token"), w, http.StatusForbidden)
		return
	}

	token, err := c.createToken(username)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.sendResponse(string(token), w)
}

type Subscription struct {
	Id           int       `db:"id"`
	Username     string    `db:"username"`
	Subreddit    string    `db:"subreddit"`
	SubscribedAt time.Time `db:"subscribed_at"`
}

//create table subscriptions
//(
//	id int auto_increment
//	primary key,
//	username varchar(255) not null,
//	subreddit varchar(255) not null,
//	subscribed_at datetime default CURRENT_TIMESTAMP not null
//);
//
//create index subscriptions_username_index
//	on subscriptions (username);

func (c *Controller) ListSubsHandler(w http.ResponseWriter, r *http.Request) {
	username := c.getUsername(r)

	var subscriptions []*Subscription
	err := c.db.Select(&subscriptions, "SELECT * FROM subscriptions WHERE username = ?", username)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	subs := []string{}
	for _, subscription := range subscriptions {
		subs = append(subs, subscription.Subreddit)
	}

	c.sendResponse(subs, w)
}

func (c *Controller) subExists(username, subreddit string) bool {
	var exists bool
	query := "SELECT COUNT(id) FROM subscriptions WHERE username = ? AND subreddit = ?"
	err := c.db.Get(&exists, query, username, subreddit)
	if err != nil {
		log.Print(err)
		return false
	}

	return exists
}

func (c *Controller) AddSubHandler(w http.ResponseWriter, r *http.Request) {
	username := c.getUsername(r)
	subreddit := mux.Vars(r)["name"]

	if c.subExists(username, subreddit) {
		c.handleError(errors.New("subreddit already subscribed"), w, http.StatusBadRequest)
		return
	}

	_, err := c.db.Exec("INSERT INTO subscriptions (username, subreddit) VALUES (?, ?)", username, subreddit)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *Controller) RemoveSubHandler(w http.ResponseWriter, r *http.Request) {
	username := c.getUsername(r)
	subreddit := mux.Vars(r)["name"]

	if !c.subExists(username, subreddit) {
		c.handleError(errors.New("subreddit isn't subscribed"), w, http.StatusBadRequest)
		return
	}

	_, err := c.db.Exec("DELETE FROM subscriptions WHERE username = ? AND subreddit = ?", username, subreddit)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
