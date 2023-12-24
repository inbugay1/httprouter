package httprouter

import (
	"context"
	"net/http"
)

type RouteParams map[string]string

type RouteMatch struct {
	Handler   Handler
	Params    RouteParams
	RouteName string
}

type Route interface {
	Match(request *http.Request) (RouteMatch, error)
}

type ctxKey int

const (
	routeParamsKey ctxKey = iota
	routeNameKey
)

func RouteParam(ctx context.Context, param string) string {
	routeParams, ok := ctx.Value(routeParamsKey).(RouteParams)
	if !ok {
		return ""
	}

	return routeParams[param]
}

func RouteName(ctx context.Context) string {
	routeName, ok := ctx.Value(routeNameKey).(string)
	if !ok {
		return ""
	}

	return routeName
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
