package reddit

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (Repository) FetchUsername(accessToken string) (string, error) {
	req, _ := http.NewRequest("GET", "https://oauth.reddit.com/api/v1/me", nil)
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("User-Agent", "Glance API Server (by /u/Albert221)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)

	if username, ok := data["name"].(string); ok {
		return username, nil
	}

	return "", errors.New("there is no name in response")
}
