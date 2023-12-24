package httprouter

import (
	"context"
	"errors"
	"net/http"
)

type Router interface {
	Match(request *http.Request) (RouteMatch, error)

	Get(path string, handler Handler, routeName string)
	Post(path string, handler Handler, routeName string)
	Put(path string, handler Handler, routeName string)
	Delete(path string, handler Handler, routeName string)
	Patch(path string, handler Handler, routeName string)
	Options(path string, handler Handler, routeName string)
	Head(path string, handler Handler, routeName string)
	Connect(path string, handler Handler, routeName string)
	Trace(path string, handler Handler, routeName string)
	Any(path string, methods []string, handler Handler, routeName string)

	Group(callback func(r Router))
	Use(middlewares ...MiddlewareFunc)
	WithPrefix(prefix string)
}

type router struct {
	routes            []Route
	routeFactoriesSet map[string]struct{}
	routeFactories    []RouteFactory
	middleware        MiddlewareFunc
	prefix            string

	NotFoundHandler Handler
}

func New(routeFactories ...RouteFactory) *router { //nolint:golint,revive
	router := &router{
		routeFactoriesSet: make(map[string]struct{}),
	}

	for _, routeFactory := range routeFactories {
		router.registerRouteFactory(routeFactory)
	}

	router.registerRouteFactory(NewLiteralRouteFactory())

	return router
}

func (r *router) registerRouteFactory(routeFactory RouteFactory) {
	if _, ok := r.routeFactoriesSet[routeFactory.Name()]; ok {
		return
	}

	r.routeFactoriesSet[routeFactory.Name()] = struct{}{}

	r.routeFactories = append(r.routeFactories, routeFactory)
}

func (r *router) route(path string, methods []string, handler Handler, routeName string) {
	if r.middleware != nil {
		handler = r.middleware(handler)
	}

	if r.prefix != "" {
		path = "/" + r.prefix + path
	}

	for _, routeFactory := range r.routeFactories {
		if routeFactory.Handles(path) {
			route := routeFactory.CreateRoute(path, methods, handler, routeName)
			r.routes = append(r.routes, route)

			return
		}
	}
}

func (r *router) Get(path string, handler Handler, routeName string) {
	r.route(path, []string{http.MethodGet}, handler, routeName)
}

func (r *router) Post(path string, handler Handler, routeName string) {
	r.route(path, []string{http.MethodPost}, handler, routeName)
}

func (r *router) Put(path string, handler Handler, routeName string) {
	r.route(path, []string{http.MethodPut}, handler, routeName)
}

func (r *router) Patch(path string, handler Handler, routeName string) {
	r.route(path, []string{http.MethodPatch}, handler, routeName)
}

func (r *router) Delete(path string, handler Handler, routeName string) {
	r.route(path, []string{http.MethodDelete}, handler, routeName)
}

func (r *router) Options(path string, handler Handler, routeName string) {
	r.route(path, []string{http.MethodOptions}, handler, routeName)
}

func (r *router) Head(path string, handler Handler, routeName string) {
	r.route(path, []string{http.MethodHead}, handler, routeName)
}

func (r *router) Connect(path string, handler Handler, routeName string) {
	r.route(path, []string{http.MethodConnect}, handler, routeName)
}

func (r *router) Trace(path string, handler Handler, routeName string) {
	r.route(path, []string{http.MethodTrace}, handler, routeName)
}

func (r *router) Any(path string, methods []string, handler Handler, routeName string) {
	r.route(path, methods, handler, routeName)
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

func (r *router) Match(request *http.Request) (RouteMatch, error) { //nolint:ireturn
	var routeMatch RouteMatch
	var methodNotAllowed bool

	for _, route := range r.routes {
		routeMatch, err := route.Match(request)
		if err != nil {
			switch {
			case errors.Is(err, ErrPathMismatch):
				continue
			case errors.Is(err, ErrMethodNotAllowed):
				methodNotAllowed = true

				continue
			}

			return routeMatch, err //nolint:wrapcheck
		}

		return routeMatch, nil
	}

	if methodNotAllowed {
		return routeMatch, ErrMethodNotAllowed
	}

	return routeMatch, ErrRouteNotFound
}

func (r *router) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	routeMatch, err := r.Match(request)
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

	ctx := context.WithValue(request.Context(), routeParamsKey, routeMatch.Params)
	ctx = context.WithValue(ctx, routeNameKey, routeMatch.RouteName)

	err = routeMatch.Handler.Handle(responseWriter, request.WithContext(ctx))
	if err != nil {
		http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
