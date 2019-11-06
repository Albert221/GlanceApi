package mock

import (
	"github.com/Albert221/ReddigramApi/domain"
	"github.com/stretchr/testify/mock"
)

type SubscriptionRepository struct {
	mock.Mock
}

func (s *SubscriptionRepository) GetUserSubscriptions(username string) ([]*domain.Subscription, error) {
	args := s.Called(username)
	return args.Get(0).([]*domain.Subscription), args.Error(1)
}

func (s *SubscriptionRepository) IsAlreadySubscribed(username, subredditID string) (bool, error) {
	args := s.Called(username, subredditID)
	return args.Bool(0), args.Error(1)
}

func (s *SubscriptionRepository) Subscribe(username, subredditID string) error {
	args := s.Called(username, subredditID)
	return args.Error(0)
}

func (s *SubscriptionRepository) Unsubscribe(username, subredditID string) error {
	args := s.Called(username, subredditID)
	return args.Error(0)
}

