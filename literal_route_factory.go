package httprouter

type literalRouteFactory struct {
}

func NewLiteralRouteFactory() *literalRouteFactory { //nolint:golint,revive
	return &literalRouteFactory{}
}

func (f *literalRouteFactory) Name() string {
	return "literal"
}

func (f *literalRouteFactory) Handles(_ string) bool {
	return true
}

func (f *literalRouteFactory) CreateRoute(path string, methods []string, handler Handler) Route {
	return &LiteralRoute{
		Methods: methods,
		Handler: handler,
		Path:    path,
	}
}
