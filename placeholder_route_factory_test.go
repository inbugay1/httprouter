package httprouter_test

import (
	"testing"

	"github.com/inbugay1/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestPlaceholderRouteFactory(t *testing.T) {
	t.Parallel()

	factory := httprouter.NewPlaceholderRouteFactory()

	assert.Equal(t, "placeholder", factory.Name())

	testCases := []struct {
		name         string
		path         string
		shouldHandle bool
	}{
		{
			name:         "PathWithPlaceholder",
			path:         "/path/to/:id",
			shouldHandle: true,
		},
		{
			name: "PathWithoutPlaceholder",
			path: "/path/to/resource",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.shouldHandle, factory.Handles(testCase.path))

			if testCase.shouldHandle {
				handler := &mockHandler{}
				methods := []string{"GET", "POST"}

				route := factory.CreateRoute("/path/to/:id", methods, handler, "")

				assert.IsType(t, &httprouter.PlaceholderRoute{}, route)
				placeholderRoute := route.(*httprouter.PlaceholderRoute)
				assert.Equal(t, methods, placeholderRoute.Methods)
				assert.NotNil(t, placeholderRoute.Tree)
			}
		})
	}
}
