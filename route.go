package httprouter

import "context"

type ctxKey int

const routeParamsKey ctxKey = iota

func RouteParam(ctx context.Context, param string) string {
	routeParams, ok := ctx.Value(routeParamsKey).(map[string]string)
	if !ok {
		return ""
	}

	return routeParams[param]
}

type Route struct {
	path    string
	methods []string

	handler Handler
}
