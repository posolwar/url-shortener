package redirect

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"url-shortener/internal/http-server/handlers/url/redirect/mocks"
	"url-shortener/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name                 string
		alias                string
		expectedURL          string
		mockResult           *url.URL
		expectedMockError    error
		expectedHandlerError string
		expectedStatusCode   int
	}{
		{
			name:                 "Success",
			alias:                "test_alias",
			expectedURL:          "https://www.google.com/",
			mockResult:           &url.URL{Scheme: "https", Host: "www.google.com", Path: "/"},
			expectedMockError:    nil,
			expectedHandlerError: "",
			expectedStatusCode:   http.StatusOK,
		},
		{
			name:                 "Error alias not found",
			alias:                "not_found",
			expectedURL:          "",
			mockResult:           nil,
			expectedMockError:    storage.ErrURLNotFound,
			expectedHandlerError: "not found",
			expectedStatusCode:   http.StatusNotFound,
		},
		{
			name:                 "Unexpected get url error",
			alias:                "unexpected_error",
			expectedURL:          "",
			mockResult:           nil,
			expectedMockError:    errors.New("unexpected error"),
			expectedHandlerError: "failed to get url",
			expectedStatusCode:   http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock URLGetter
			mockURLGetter := mocks.NewURLGetter(t)
			mockURLGetter.On(
				"GetURL",
				mock.AnythingOfType("*context.valueCtx"),
				tc.alias,
			).Return(tc.mockResult, tc.expectedMockError)

			router := chi.NewRouter()
			router.Get("/{alias}", RedirectHandler(mockURLGetter))

			ts := httptest.NewServer(router)
			defer ts.Close()

			resp, err := http.Get(ts.URL + "/" + tc.alias)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.expectedStatusCode == http.StatusOK {
				assert.Equal(t, tc.expectedURL, resp.Request.URL.String())
			}
		})
	}
}
