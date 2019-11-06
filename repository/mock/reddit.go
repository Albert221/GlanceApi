package mock

import "github.com/stretchr/testify/mock"

type RedditRepository struct {
	mock.Mock
}

func (r *RedditRepository) FetchUsername(accessToken string) (string, error) {
	args := r.Called(accessToken)
	return args.String(0), args.Error(1)
}
