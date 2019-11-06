package mysql

import (
	"github.com/Albert221/ReddigramApi/domain"
	"github.com/jmoiron/sqlx"
)

type SubscriptionRepository struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) *SubscriptionRepository {
	repo := &SubscriptionRepository{
		db: db,
	}

	return repo
}

func (s *SubscriptionRepository) GetUserSubscriptions(username string) ([]*domain.Subscription, error) {
	var subscriptions []*domain.Subscription
	sql := "SELECT * FROM subscriptions WHERE username = ?"
	if err := s.db.Select(&subscriptions, sql, username); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (s *SubscriptionRepository) IsAlreadySubscribed(username, subredditID string) (bool, error) {
	var exists bool
	sql := "SELECT COUNT(id) FROM subscriptions WHERE username = ? AND subreddit = ?"
	if err := s.db.Get(&exists, sql, username, subredditID); err != nil {
		return false, err
	}

	return exists, nil
}

func (s *SubscriptionRepository) Subscribe(username, subredditID string) error {
	sql := "INSERT INTO subscriptions (username, subreddit) VALUES (?, ?)"
	_, err := s.db.Exec(sql, username, subredditID)

	return err
}

func (s *SubscriptionRepository) Unsubscribe(username, subredditID string) error {
	sql := "DELETE FROM subscriptions WHERE username = ? AND subreddit = ?"
	_, err := s.db.Exec(sql, username, subredditID)

	return err
}
