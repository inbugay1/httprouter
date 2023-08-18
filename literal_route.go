package httprouter

import "net/http"

type LiteralRoute struct {
	Methods []string
	Handler Handler
	Path    string
}

func (literalRoute *LiteralRoute) Match(request *http.Request) (Handler, error) {
	if literalRoute.Path != request.URL.Path {
		return nil, ErrPathMismatch
	}

	if !contains(literalRoute.Methods, request.Method) {
		return nil, ErrMethodNotAllowed
	}

	return literalRoute.Handler, nil
}
