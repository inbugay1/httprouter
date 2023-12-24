package httprouter

import (
	"net/http"
)

type PlaceholderRoute struct {
	Methods []string
	Tree    *tree
	Name    string
}

func (route *PlaceholderRoute) Match(request *http.Request) (RouteMatch, error) {
	var routeMatch RouteMatch

	handler, params := route.Tree.Search(request.URL.Path)
	if handler == nil {
		return routeMatch, ErrPathMismatch
	}

	if !contains(route.Methods, request.Method) {
		return routeMatch, ErrMethodNotAllowed
	}

	routeMatch.Handler = handler
	routeMatch.Params = params
	routeMatch.RouteName = route.Name

	return routeMatch, nil
}
