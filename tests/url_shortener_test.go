package tests

import (
	"net/http"
	"net/url"
	"testing"

	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/pkg/random"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
)

const (
	host = "localhost:8082"
)

// blackbox testing get url
func TestURLShortener_HappyPath(t *testing.T) {
	testcases := []struct {
		name     string
		alias    string
		url      string
		login    string
		password string
		status   int
	}{
		{
			name:     "happy path - with alias",
			alias:    random.NewRandomString(10),
			url:      gofakeit.URL(),
			status:   http.StatusOK,
			login:    "myuser",
			password: "mypass",
		},
	}

	url := url.URL{Scheme: "http", Host: host}
	exp := httpexpect.Default(t, url.String())

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			expectedBody := httpexpect.NewObject(t, map[string]interface{}{
				"alias":  testcase.alias,
				"status": testcase.status,
			})

			exp.POST("/url").
				WithJSON(save.Request{
					URL:   testcase.url,
					Alias: testcase.alias,
				}).
				WithBasicAuth(testcase.login, testcase.password).
				Expect().
				Status(http.StatusOK).
				JSON().
				Object().
				ContainsSubset(expectedBody)
		})
	}
}
