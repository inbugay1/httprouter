package httprouter

import (
	"context"
	"errors"
	"net/http"
	"regexp"
)

type Router interface {
	Match(request *http.Request) (Handler, error)

	Get(path string, handler Handler)
	Post(path string, handler Handler)
	Put(path string, handler Handler)
	Delete(path string, handler Handler)
	Patch(path string, handler Handler)
	Options(path string, handler Handler)
	Head(path string, handler Handler)
	Connect(path string, handler Handler)
	Trace(path string, handler Handler)
	Any(path string, methods []string, handler Handler)

	Group(callback func(r Router))
	Use(middlewares ...MiddlewareFunc)
	WithPrefix(prefix string)
}

type router struct {
	routes     []*route
	re         *regexp.Regexp
	middleware MiddlewareFunc
	prefix     string

	NotFoundHandler Handler
}

func New() *router { //nolint:golint,revive
	return &router{
		re: regexp.MustCompile(`{(?P<param>\w+):(?P<regex>.+)}`),
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func (r *router) route(path string, methods []string, handler Handler) {
	if r.middleware != nil {
		handler = r.middleware(handler)
	}

	newRoute := &route{methods: methods, handler: handler}

	if r.re.MatchString(path) {
		pathRegexStr := r.re.ReplaceAllString(path, "(?P<$1>$2)") // e.g modify /test/{id:\d+} to /test/(?P<id>\d+)

		if r.prefix != "" {
			pathRegexStr = "/" + r.prefix + pathRegexStr
		}

		newRoute.pathRegex = regexp.MustCompile("^" + pathRegexStr + "$")
	} else {
		if r.prefix != "" {
			path = "/" + r.prefix + path
		}

		newRoute.path = path
	}

	r.routes = append(r.routes, newRoute)
}

func (r *router) Get(path string, handler Handler) {
	r.route(path, []string{http.MethodGet}, handler)
}

func (r *router) Post(path string, handler Handler) {
	r.route(path, []string{http.MethodPost}, handler)
}

func (r *router) Put(path string, handler Handler) {
	r.route(path, []string{http.MethodPut}, handler)
}

func (r *router) Patch(path string, handler Handler) {
	r.route(path, []string{http.MethodPatch}, handler)
}

func (r *router) Delete(path string, handler Handler) {
	r.route(path, []string{http.MethodDelete}, handler)
}

func (r *router) Options(path string, handler Handler) {
	r.route(path, []string{http.MethodOptions}, handler)
}

func (r *router) Head(path string, handler Handler) {
	r.route(path, []string{http.MethodHead}, handler)
}

func (r *router) Connect(path string, handler Handler) {
	r.route(path, []string{http.MethodConnect}, handler)
}

func (r *router) Trace(path string, handler Handler) {
	r.route(path, []string{http.MethodTrace}, handler)
}

func (r *router) Any(path string, methods []string, handler Handler) {
	r.route(path, methods, handler)
}

func (r *router) Group(callback func(r Router)) {
	routerMiddleware := r.middleware
	routerPrefix := r.prefix

	callback(r)

	// remove group middleware and prefix
	r.middleware = routerMiddleware
	r.prefix = routerPrefix
}

// Use
// r.Use(middleware1)
// r.Use(middleware2)
// Or r.Use(middleware1, middleware2)
// -> middleware1(middleware2(next)).
func (r *router) Use(middlewares ...MiddlewareFunc) {
	for _, middleware := range middlewares {
		middleware := middleware // shadow for closure

		if r.middleware != nil {
			routerMiddleware := r.middleware
			r.middleware = func(next Handler) Handler {
				return routerMiddleware(middleware(next))
			}

			continue
		}

		r.middleware = middleware
	}
}

func (r *router) WithPrefix(prefix string) {
	if r.prefix != "" {
		r.prefix += "/" + prefix

		return
	}

	r.prefix = prefix
}

func (r *router) Match(request *http.Request) (Handler, error) { //nolint:ireturn
	methodNotAllowed := false

	for _, route := range r.routes {
		if route.path != "" {
			if route.path == request.URL.Path {
				if contains(route.methods, request.Method) {
					return route.handler, nil
				}

				methodNotAllowed = true

				continue
			}

			continue
		}

		if route.pathRegex == nil || !route.pathRegex.MatchString(request.URL.Path) {
			continue
		}

		if !contains(route.methods, request.Method) {
			methodNotAllowed = true

			continue
		}

		matches := route.pathRegex.FindAllStringSubmatch(request.URL.Path, -1)

		routeParamNames := route.pathRegex.SubexpNames()
		routeParams := make(map[string]string, len(routeParamNames))

		for paramKey, paramName := range routeParamNames {
			routeParams[paramName] = matches[0][paramKey]
		}

		ctx := context.WithValue(request.Context(), routeParamsKey, routeParams)
		*request = *request.WithContext(ctx)

		return route.handler, nil
	}

	if methodNotAllowed {
		return nil, ErrMethodNotAllowed
	}

	return nil, ErrRouteNotFound
}

func (r *router) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	handler, err := r.Match(request)
	if err != nil {
		switch {
		case errors.Is(err, ErrMethodNotAllowed):
			http.Error(responseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		case errors.Is(err, ErrRouteNotFound):
			if r.NotFoundHandler != nil {
				err = r.NotFoundHandler.Handle(responseWriter, request)
				if err != nil {
					http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				return
			}
			http.NotFound(responseWriter, request)
		}

		return
	}

	err = handler.Handle(responseWriter, request)
	if err != nil {
		http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
