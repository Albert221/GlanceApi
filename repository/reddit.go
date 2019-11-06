package repository

type RedditRepository interface {
	FetchUsername(accessToken string) (string, error)
}
