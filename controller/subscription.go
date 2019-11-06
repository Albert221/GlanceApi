package controller

import (
	"github.com/Albert221/ReddigramApi/repository"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type SubscriptionController struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionController(repo repository.SubscriptionRepository) *SubscriptionController {
	return &SubscriptionController{
		repo: repo,
	}
}

func (s *SubscriptionController) ListHandler(w http.ResponseWriter, r *http.Request) {
	username := getUsername(r)
	if username == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	subscriptions, err := s.repo.GetUserSubscriptions(username)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var subIDs []string
	for _, subscription := range subscriptions {
		subIDs = append(subIDs, subscription.Subreddit)
	}

	writeJSON(w, subIDs, http.StatusOK)
}

func (s *SubscriptionController) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	subreddit := mux.Vars(r)["id"]
	username := getUsername(r)
	if username == "" {
		writeJSON(w, nil, http.StatusForbidden)
		return
	}

	subscribed, err := s.repo.IsAlreadySubscribed(username, subreddit)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if subscribed {
		writeJSON(w, map[string]string{"error": "subreddit is already subscribed"}, http.StatusBadRequest)
	}

	if err := s.repo.Subscribe(username, subreddit); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *SubscriptionController) UnsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	subreddit := mux.Vars(r)["id"]
	username := getUsername(r)
	if username == "" {
		writeJSON(w, nil, http.StatusForbidden)
		return
	}

	subscribed, err := s.repo.IsAlreadySubscribed(username, subreddit)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !subscribed {
		writeJSON(w, map[string]string{"error": "subreddit is not subscribed"}, http.StatusBadRequest)
	}

	if err := s.repo.Unsubscribe(username, subreddit); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}