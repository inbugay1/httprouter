package httprouter_test

import (
	"testing"

	"github.com/inbugay1/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestRegexRouteFactory(t *testing.T) {
	t.Parallel()

	factory := httprouter.NewRegexRouteFactory()
	assert.Equal(t, "regex", factory.Name(), "The factory name should be 'regex'")

	testCases := []struct {
		path            string
		shouldHandle    bool
		expectedPattern string
	}{
		{
			path:            `/test/{id:\d+}`,
			shouldHandle:    true,
			expectedPattern: `^/test/(?P<id>\d+)$`,
		},
		{
			path:            `/test/{name:\w+}/id/{id:\d+}`,
			shouldHandle:    true,
			expectedPattern: `^/test/(?P<name>\w+)/id/(?P<id>\d+)$`,
		},
		{
			path:            `/some_path/{id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/some_path2/{name:\w+}`,
			shouldHandle:    true,
			expectedPattern: `^/some_path/(?P<id>[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})/some_path2/(?P<name>\w+)$`,
		},
		{
			path:         "/test/no/regex",
			shouldHandle: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.path, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.shouldHandle, factory.Handles(testCase.path))

			if testCase.shouldHandle {
				route := factory.CreateRoute(testCase.path, nil, nil)
				regexRoute, ok := route.(*httprouter.RegexRoute)
				assert.True(t, ok, "Expected route to be of type *RegexRoute")
				assert.Equal(t, testCase.expectedPattern, regexRoute.Regexp.String())
			}
		})
	}
}
