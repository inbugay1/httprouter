package httprouter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inbugay1/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestLiteralRouteMatch(t *testing.T) {
	t.Parallel()

	handlerFunc := func(responseWriter http.ResponseWriter, request *http.Request) error {
		_, _ = responseWriter.Write([]byte("handler called"))

		return nil
	}

	route := &httprouter.LiteralRoute{
		Methods: []string{http.MethodGet},
		Handler: httprouter.HandlerFunc(handlerFunc),
		Path:    "/test",
	}

	testCases := []struct {
		name         string
		request      *http.Request
		expectedErr  error
		expectedResp string
	}{
		{"Matching GET request", httptest.NewRequest(http.MethodGet, "/test", nil), nil, "handler called"},
		{"Non-matching path", httptest.NewRequest(http.MethodGet, "/wrongpath", nil), httprouter.ErrPathMismatch, ""},
		{"Non-matching method", httptest.NewRequest(http.MethodPost, "/test", nil), httprouter.ErrMethodNotAllowed, ""},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			recorder := httptest.NewRecorder()
			routeMatch, err := route.Match(testCase.request)
			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)
				assert.Empty(t, routeMatch, "RouteMatch should be empty")
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, routeMatch, "RouteMatch should not be empty")
				err = routeMatch.Handler.Handle(recorder, testCase.request)
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedResp, recorder.Body.String())
			}
		})
	}
}
