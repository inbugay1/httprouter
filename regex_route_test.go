package httprouter_test

import (
	"fmt"
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
			id := httprouter.RouteParam(r.Context(), "id")
			fmt.Fprintf(w, "ID: %s!", id)

			return nil
		}),
		Regexp: regexp.MustCompile(`^/test/(?P<id>\d+)$`),
	}

	testCases := []struct {
		name         string
		request      *http.Request
		expectedResp string
		expectedErr  error
	}{
		{
			name:         "Matching Path",
			request:      httptest.NewRequest(http.MethodGet, "/test/123", nil),
			expectedResp: "ID: 123!",
			expectedErr:  nil,
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
			matchHandler, err := route.Match(testCase.request)
			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)
				assert.Nil(t, matchHandler)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, matchHandler)
				err = matchHandler.Handle(recorder, testCase.request)
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedResp, recorder.Body.String())
			}
		})
	}
}
