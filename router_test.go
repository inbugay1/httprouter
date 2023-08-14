package httprouter_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inbugay1/httprouter" // Import the correct package path
	"github.com/stretchr/testify/assert"
)

type mockHandler struct {
	errToReturn error
}

func (h *mockHandler) Handle(_ http.ResponseWriter, _ *http.Request) error {
	return h.errToReturn
}

//nolint:funlen,cyclop
func TestRouter_AllHTTPMethods(t *testing.T) {
	t.Parallel()

	router := httprouter.New()

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodHead,
		http.MethodConnect,
		http.MethodTrace,
	}

	getRouteMethodFunction := func(router httprouter.Router, method string) func(string, httprouter.Handler) {
		switch method {
		case http.MethodGet:
			return router.Get
		case http.MethodPost:
			return router.Post
		case http.MethodPut:
			return router.Put
		case http.MethodPatch:
			return router.Patch
		case http.MethodDelete:
			return router.Delete
		case http.MethodOptions:
			return router.Options
		case http.MethodHead:
			return router.Head
		case http.MethodConnect:
			return router.Connect
		case http.MethodTrace:
			return router.Trace
		default:
			return nil
		}
	}

	// Register all the methods for the router
	handlers := make(map[string]*mockHandler, len(methods))

	// Register all the methods for the router
	for _, method := range methods {
		routeMethod := getRouteMethodFunction(router, method)
		route := "/test-" + method

		handler := &mockHandler{}
		handlers[method] = handler
		routeMethod(route, handler)
	}

	for _, method := range methods {
		method := method

		t.Run("HTTP Method "+method, func(t *testing.T) {
			t.Parallel()

			req, _ := http.NewRequestWithContext(context.Background(), method, "/test-"+method, nil)
			matchHandler, err := router.Match(req)

			if assert.NoError(t, err, "Expected no error") {
				assert.Equal(t, handlers[method], matchHandler, "Handler mismatch")
			}
		})
	}
}

func TestRouter_Any(t *testing.T) {
	t.Parallel()

	router := httprouter.New()

	handler := &mockHandler{}

	// Define a path to test
	path := "/test"

	// Add the route using the Any method
	router.Any(path, []string{http.MethodGet, http.MethodPost}, handler)

	// Test the path with GET and POST methods
	reqGet, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, path, nil)
	matchHandlerGet, errGet := router.Match(reqGet)
	if assert.NoError(t, errGet, "Expected no error for GET") {
		assert.Equal(t, handler, matchHandlerGet, "Handler mismatch for GET")
	}

	reqPost, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, path, nil)
	matchHandlerPost, errPost := router.Match(reqPost)
	if assert.NoError(t, errPost, "Expected no error for POST") {
		assert.Equal(t, handler, matchHandlerPost, "Handler mismatch for POST")
	}
}

func TestRouter_Match_NotFound(t *testing.T) {
	t.Parallel()

	router := httprouter.New()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/nonexistent", nil)
	_, err := router.Match(req)

	assert.ErrorIs(t, err, httprouter.ErrRouteNotFound, "Expected ErrRouteNotFound error")
}

func TestRouter_Match_MethodNotAllowed(t *testing.T) {
	t.Parallel()

	router := httprouter.New()
	router.Get("/test", &mockHandler{})
	router.Get("/test/{id:\\d+}", &mockHandler{}) // Regex route

	testCases := []struct {
		method      string
		path        string
		description string
	}{
		{method: http.MethodPost, path: "/test", description: "Unsupported method for /test"},
		{method: http.MethodPost, path: "/test/123", description: "Unsupported method for regex route /test/123"},
	}

	for _, testCase := range testCases {
		request, _ := http.NewRequestWithContext(context.Background(), testCase.method, testCase.path, nil)
		_, err := router.Match(request)
		assert.True(t, errors.Is(err, httprouter.ErrMethodNotAllowed), testCase.description)
	}
}

func TestRouter_Match_Regex(t *testing.T) {
	t.Parallel()

	router := httprouter.New()

	handler := &mockHandler{}
	router.Get("/test/{id:\\d+}", handler)

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test/123", nil)
	matchHandler, err := router.Match(req)

	if assert.NoError(t, err, "Expected no error") {
		assert.Equal(t, handler, matchHandler, "Handler mismatch")
	}
}

func TestRouter_Group(t *testing.T) {
	t.Parallel()

	router := httprouter.New()

	handler := &mockHandler{}
	router.Group(func(group httprouter.Router) {
		group.WithPrefix("v1")
		group.Get("/test", handler)
	})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/test", nil)
	matchHandler, err := router.Match(req)

	if assert.NoError(t, err) {
		assert.Equal(t, handler, matchHandler, "Handler mismatch")
	}
}

func TestRouter_Use(t *testing.T) {
	t.Parallel()

	router := httprouter.New()

	handler := &mockHandler{}

	var middlewareOrder []string

	// Custom middleware that records the order of execution
	recordMiddleware := func(name string) httprouter.MiddlewareFunc {
		return func(next httprouter.Handler) httprouter.Handler {
			return httprouter.HandlerFunc(func(w http.ResponseWriter, router *http.Request) error {
				middlewareOrder = append(middlewareOrder, name)

				return next.Handle(w, router) //nolint:wrapcheck
			})
		}
	}

	router.Use(recordMiddleware("middleware1"), recordMiddleware("middleware2"))

	// Add the route using the Get method
	router.Get("/test", handler)

	// Test the route with the applied middleware
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)
	matchHandler, err := router.Match(req)

	if assert.NoError(t, err, "Expected no error") {
		assert.NotNil(t, matchHandler, "Matched handler should not be nil")

		// Call the Handle function to populate middlewareOrder slice
		err := matchHandler.Handle(nil, req)
		assert.NoError(t, err, "Expected no error")

		// Verify the order of middleware execution
		expectedOrder := []string{"middleware1", "middleware2"}
		assert.Equal(t, expectedOrder, middlewareOrder, "Middleware order mismatch")
	}
}

func TestRouter_GroupUse(t *testing.T) {
	t.Parallel()

	router := httprouter.New()

	handler := &mockHandler{}
	var middlewareOrder []string

	// Custom middleware that records the order of execution
	recordMiddleware := func(name string) httprouter.MiddlewareFunc {
		return func(next httprouter.Handler) httprouter.Handler {
			return httprouter.HandlerFunc(func(w http.ResponseWriter, router *http.Request) error {
				middlewareOrder = append(middlewareOrder, name)

				return next.Handle(w, router) //nolint:wrapcheck
			})
		}
	}

	// Add a route outside the group with its own middleware
	router.Use(recordMiddleware("OutsideMiddleware"))
	router.Get("/outside-route", handler)

	// Create a route group with middleware
	router.Group(func(group httprouter.Router) {
		group.Use(recordMiddleware("GroupMiddleware1"), recordMiddleware("GroupMiddleware2"))
		group.Get("/group-route", handler)
	})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/outside-route", nil)
	matchHandler, err := router.Match(req)

	if assert.NoError(t, err, "Expected no error") {
		middlewareOrder = nil             // Reset middleware order
		_ = matchHandler.Handle(nil, nil) // Trigger middleware execution

		assert.Equal(t, []string{"OutsideMiddleware"}, middlewareOrder, "Middleware order mismatch")
	}

	req, _ = http.NewRequestWithContext(context.Background(), http.MethodGet, "/group-route", nil)
	matchHandler, err = router.Match(req)

	if assert.NoError(t, err, "Expected no error") {
		middlewareOrder = nil             // Reset middleware order
		_ = matchHandler.Handle(nil, nil) // Trigger middleware execution

		assert.Equal(t, []string{"OutsideMiddleware", "GroupMiddleware1", "GroupMiddleware2"}, middlewareOrder, "Middleware order mismatch")
	}
}

func TestRouter_WithPrefix(t *testing.T) {
	t.Parallel()

	router := httprouter.New()

	handler := &mockHandler{}

	// Add a route with a prefix using the WithPrefix method
	router.WithPrefix("api")
	router.Get("/test", handler)
	router.Get("/test/{id:\\d+}", handler)

	// Test the route with the applied prefix
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/test", nil)
	matchHandler, err := router.Match(req)

	if assert.NoError(t, err, "Expected no error") {
		assert.NotNil(t, matchHandler, "Matched handler should not be nil")
	}

	// Test the route with the applied prefix
	req, _ = http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/test/123", nil)
	matchHandler, err = router.Match(req)

	if assert.NoError(t, err, "Expected no error") {
		assert.NotNil(t, matchHandler, "Matched handler should not be nil")
	}
}

func TestRouter_GroupWithPrefix(t *testing.T) {
	t.Parallel()

	router := httprouter.New()

	handler := &mockHandler{}

	// Add a route outside the group with its own prefix
	router.WithPrefix("outside")
	router.Get("/route1", handler)

	// Create a route group with a specific prefix
	router.Group(func(group httprouter.Router) {
		group.WithPrefix("group1")
		group.Get("/route2", handler)
	})

	// Test routes with prefixes applied only to the respective groups and outside route
	testCases := []struct {
		path     string
		expected bool
	}{
		{"/outside/route1", true},
		{"/outside/group1/route2", true},
		{"/outside/group1/route1", false}, // Should not match group1 prefix
		{"/group1/route2", false},         // Should not match group1 prefix without outside prefix
		{"/group1/route1", false},         // Should not match group1 prefix without outside prefix
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.path, func(t *testing.T) {
			t.Parallel()

			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, testCase.path, nil)
			matchHandler, err := router.Match(req)

			if testCase.expected {
				assert.NoError(t, err, "Expected no error")
				assert.NotNil(t, matchHandler, "Matched handler should not be nil")
			} else {
				assert.Error(t, err, "Expected error")
				assert.Nil(t, matchHandler, "Matched handler should be nil")
			}
		})
	}
}

func TestRouter_ServeHTTP_MethodNotAllowed(t *testing.T) {
	t.Parallel()

	router := httprouter.New()

	handler := &mockHandler{}
	router.Get("/test", handler)

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/test", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code, "Expected StatusMethodNotAllowed")
}

func TestRouter_ServeHTTP_RouteFound(t *testing.T) {
	t.Parallel()

	router := httprouter.New()

	handler := &mockHandler{}
	router.Get("/test", handler)

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Status code mismatch")
}

func TestRouter_ServeHTTP_RouteNotFound(t *testing.T) {
	t.Parallel()

	router := httprouter.New()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/nonexistent", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, "Status code mismatch")
}

func TestRouter_ServeHTTP_InternalServerError_HandlerError(t *testing.T) {
	t.Parallel()

	router := httprouter.New()

	handler := &mockHandler{errToReturn: errors.New("internal handler error")}
	router.Get("/test", handler)

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code mismatch")
}

func TestRouter_ServeHTTP_InternalServerError_NotFoundHandlerError(t *testing.T) {
	t.Parallel()

	router := httprouter.New()

	router.NotFoundHandler = &mockHandler{errToReturn: errors.New("internal not found handler error")}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/nonexistent", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code mismatch")
}

func BenchmarkServeHTTPWithoutParams(b *testing.B) {
	router := httprouter.New()

	// Define routes here
	router.Get("/sample", &mockHandler{})
	router.Get("/example", &mockHandler{})
	// ...

	// Create a sample HTTP request for benchmarking
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/sample", nil)
	recorder := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(recorder, req)
	}
}

func BenchmarkServeHTTPWithParams(b *testing.B) {
	router := httprouter.New()

	// Define routes here
	router.Get("/sample/{id:\\d+}", &mockHandler{})
	router.Get("/example/{id:\\d+}", &mockHandler{})
	// ...

	// Create a sample HTTP request for benchmarking
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/sample/123", nil)
	recorder := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(recorder, req)
	}
}
