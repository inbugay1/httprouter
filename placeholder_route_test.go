package httprouter_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/inbugay1/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestPlaceholderRouteMatch(t *testing.T) {
	t.Parallel()

	tree := httprouter.NewTree()

	handler1 := &mockHandler{}
	handler2 := &mockHandler{}

	tree.Insert("/path/to/resource", handler1)
	tree.Insert("/path/to/resource2", handler2)
	tree.Insert("/path/to/:id", handler1)
	tree.Insert("/path/from/:from/to/:to", handler2)

	route := httprouter.PlaceholderRoute{
		Methods: []string{http.MethodGet, http.MethodPost},
		Tree:    tree,
	}

	var testCases = []struct {
		name                string
		method              string
		path                string
		shouldMatch         bool
		expectedRouteParams httprouter.RouteParams
		expectedErr         error
	}{
		{
			name:                "MatchStaticPath",
			method:              http.MethodGet,
			path:                "/path/to/resource",
			expectedRouteParams: httprouter.RouteParams{},
			shouldMatch:         true,
		},
		{
			name:   "MatchDynamicPath1",
			method: http.MethodPost,
			path:   "/path/to/123",
			expectedRouteParams: httprouter.RouteParams{
				"id": "123",
			},
			shouldMatch: true,
		},
		{
			name:   "MatchDynamicPath2",
			method: http.MethodPost,
			path:   "/path/from/123/to/456",
			expectedRouteParams: httprouter.RouteParams{
				"from": "123",
				"to":   "456",
			},
			shouldMatch: true,
		},
		{
			name:        "PathWithoutHandler",
			method:      http.MethodGet,
			path:        "/path/to",
			expectedErr: httprouter.ErrPathMismatch,
		},
		{
			name:        "NonExistentPath",
			method:      http.MethodPost,
			path:        "/path/not/in/tree",
			expectedErr: httprouter.ErrPathMismatch,
		},
		{
			name:        "MethodNotAllowed",
			method:      http.MethodPut,
			path:        "/path/to/resource",
			expectedErr: httprouter.ErrMethodNotAllowed,
		},
	}
	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			request, err := http.NewRequestWithContext(context.Background(), testCase.method, testCase.path, nil)
			assert.NoError(t, err)

			routeMatch, err := route.Match(request)
			assert.Equal(t, testCase.expectedErr, err)

			if testCase.shouldMatch {
				assert.NotNil(t, routeMatch.Handler)
				assert.Equal(t, testCase.expectedRouteParams, routeMatch.Params)
			} else {
				assert.Nil(t, routeMatch.Handler)
			}
		})
	}
}
