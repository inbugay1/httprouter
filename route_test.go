package httprouter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteParam(t *testing.T) {
	t.Parallel()

	t.Run("get empty string because of invalid map type assertion", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(context.Background(), routeParamsKey, nil)

		assert.Equal(t, "", RouteParam(ctx, "test"))
	})

	t.Run("get empty string by non-existing key", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(context.Background(), routeParamsKey, map[string]string{})

		assert.Equal(t, "", RouteParam(ctx, "test"))
	})

	t.Run("get value by existing key", func(t *testing.T) {
		t.Parallel()

		params := map[string]string{
			"key": "value",
		}

		ctx := context.WithValue(context.Background(), routeParamsKey, params)

		assert.Equal(t, "value", RouteParam(ctx, "key"))
	})
}
