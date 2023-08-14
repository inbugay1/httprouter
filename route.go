package httprouter

import (
	"context"
	"regexp"
)

type ctxKey int

const routeParamsKey ctxKey = iota

func RouteParam(ctx context.Context, param string) string {
	routeParams, ok := ctx.Value(routeParamsKey).(map[string]string)
	if !ok {
		return ""
	}

	return routeParams[param]
}

type route struct {
	path      string
	pathRegex *regexp.Regexp
	methods   []string

	handler Handler
}
