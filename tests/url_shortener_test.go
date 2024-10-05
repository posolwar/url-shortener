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

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	e.
		POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.NewRandomString(10),
		}).
		WithBasicAuth("myuser", "mypass").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("alias")
}
