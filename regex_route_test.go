package httprouter_test

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/inbugay1/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestRegexRoute(t *testing.T) {
	t.Parallel()

	route := &httprouter.RegexRoute{
		Methods: []string{http.MethodGet},
		Handler: httprouter.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			return nil
		}),
		Regexp: regexp.MustCompile(`^/test/(?P<id>\d+)$`),
	}

	testCases := []struct {
		name                string
		request             *http.Request
		expectedRouteParams httprouter.RouteParams
		expectedErr         error
	}{
		{
			name:    "Matching Path",
			request: httptest.NewRequest(http.MethodGet, "/test/123", nil),
			expectedRouteParams: httprouter.RouteParams{
				"id": "123",
			},
			expectedErr: nil,
		},
		{
			name:        "Mismatching Path",
			request:     httptest.NewRequest(http.MethodGet, "/test/abc", nil),
			expectedErr: httprouter.ErrPathMismatch,
		},
		{
			name:        "Mismatching Method",
			request:     httptest.NewRequest(http.MethodPost, "/test/123", nil),
			expectedErr: httprouter.ErrMethodNotAllowed,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // Capture variable
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
				assert.Equal(t, testCase.expectedRouteParams, routeMatch.Params)
			}
		})
	}
}
