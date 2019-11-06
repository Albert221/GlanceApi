package repository

import "github.com/Albert221/ReddigramApi/domain"

type SubscriptionRepository interface {
	GetUserSubscriptions(username string) ([]*domain.Subscription, error)
	IsAlreadySubscribed(username, subredditID string) (bool, error)
	Subscribe(username, subredditID string) error
	Unsubscribe(username, subredditID string) error
}
