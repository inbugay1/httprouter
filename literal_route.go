package httprouter

import "net/http"

type LiteralRoute struct {
	Methods []string
	Handler Handler
	Path    string
}

func (literalRoute *LiteralRoute) Match(request *http.Request) (RouteMatch, error) {
	var routeMatch RouteMatch

	if literalRoute.Path != request.URL.Path {
		return routeMatch, ErrPathMismatch
	}

	if !contains(literalRoute.Methods, request.Method) {
		return routeMatch, ErrMethodNotAllowed
	}

	routeMatch.Handler = literalRoute.Handler

	return routeMatch, nil
}
