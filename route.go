package httprouter

import "context"

type ctxKey int

const routeParamsKey ctxKey = iota

func RouteParam(ctx context.Context, param string) string {
	routeParams := ctx.Value(routeParamsKey).(map[string]string)

	return routeParams[param]
}

type Route struct {
	path    string
	methods []string

	handler Handler
}
