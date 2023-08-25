package httprouter

import (
	"context"
	"net/http"
)

type RouteParams map[string]string

type RouteMatch struct {
	Handler Handler
	Params  RouteParams
}

type Route interface {
	Match(request *http.Request) (RouteMatch, error)
}

type ctxKey int

const routeParamsKey ctxKey = iota

func RouteParam(ctx context.Context, param string) string {
	routeParams, ok := ctx.Value(routeParamsKey).(RouteParams)
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
