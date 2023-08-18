package httprouter

import (
	"context"
	"net/http"
)

type Route interface {
	Match(request *http.Request) (Handler, error)
}

type ctxKey int

const routeParamsKey ctxKey = iota

func RouteParam(ctx context.Context, param string) string {
	routeParams, ok := ctx.Value(routeParamsKey).(map[string]string)
	if !ok {
		return ""
	}

	return routeParams[param]
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
