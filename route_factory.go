package httprouter

type RouteFactory interface {
	Name() string
	Handles(path string) bool
	CreateRoute(path string, methods []string, handler Handler, routeName string) Route
}
