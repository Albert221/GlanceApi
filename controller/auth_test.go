package controller

import (
	"errors"
	"github.com/Albert221/ReddigramApi/repository/mock"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthController_AuthenticateHandler(t *testing.T) {
	redditRepoMock := new(mock.RedditRepository)
	redditRepoMock.
		On("FetchUsername", "correct-access-token").
		Return("johndoe", nil).
		On("FetchUsername", "incorrect-access-token").
		Return("", errors.New(""))

	authContr := NewAuthController(jwt.None(), redditRepoMock)

	tests := []struct {
		Name  string
		Body  string
		Check func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			"error when body is empty",
			"",
			func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, 400, rr.Code)
			},
		},
		{
			"error when access token is incorrect",
			"access_token=incorrect-access-token",
			func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, 400, rr.Code)
			},
		},
		{
			"token when access token is correct",
			"access_token=correct-access-token",
			func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, 200, rr.Code)
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			r := httptest.NewRequest("POST", "/authenticate", strings.NewReader(test.Body))
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()

			authContr.AuthenticateHandler(rr, r)

			test.Check(t, rr)
		})
	}
}
