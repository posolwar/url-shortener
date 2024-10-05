package remove

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"url-shortener/internal/http-server/handlers/url/remove/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRemover(t *testing.T) {
	testcases := []struct {
		name         string
		alias        string
		mockError    error
		expectedErr  error
		expectedCode int
	}{
		{
			name:         "Success",
			alias:        "test",
			mockError:    nil,
			expectedErr:  nil,
			expectedCode: http.StatusNoContent,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			urlRemoverMock := mocks.NewURLRemover(t)
			urlRemoverMock.On("DeleteURL", mock.AnythingOfType("*context.valueCtx"), tc.alias).Return(tc.mockError)

			router := chi.NewRouter()
			router.Delete("/{alias}", RemoveURLHandler(urlRemoverMock))

			testServer := httptest.NewServer(router)
			defer testServer.Close()

			req, err := http.NewRequest(http.MethodDelete, testServer.URL+"/"+tc.alias, nil)
			assert.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedCode, resp.StatusCode)
		})
	}
}
