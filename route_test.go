package httprouter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteParam(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		ctx            context.Context //nolint:containedctx
		param          string
		expectedResult string
	}{
		{
			name:           "ValidParam",
			ctx:            context.WithValue(context.Background(), routeParamsKey, RouteParams{"id": "123"}),
			param:          "id",
			expectedResult: "123",
		},
		{
			name:           "InvalidParam",
			ctx:            context.WithValue(context.Background(), routeParamsKey, RouteParams{"id": "123"}),
			param:          "name",
			expectedResult: "",
		},
		{
			name:           "NoContextValue",
			ctx:            context.Background(),
			param:          "id",
			expectedResult: "",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result := RouteParam(testCase.ctx, testCase.param)
			assert.Equal(t, testCase.expectedResult, result, "Result mismatch")
		})
	}
}
