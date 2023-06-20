package httprouter

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testHandler struct {
}

func (h *testHandler) Handle(http.ResponseWriter, *http.Request) error {
	return nil
}

type testMiddleware struct {
	Next Handler
}

func (m *testMiddleware) Handle(responseWriter http.ResponseWriter, request *http.Request) error {
	return m.Next.Handle(responseWriter, request) // nolint:wrapcheck
}

type test2Middleware struct {
	Next Handler
}

func (m *test2Middleware) Handle(responseWriter http.ResponseWriter, request *http.Request) error {
	return m.Next.Handle(responseWriter, request) // nolint:wrapcheck
}

type test3Middleware struct {
	Next Handler
}

func (m *test3Middleware) Handle(responseWriter http.ResponseWriter, request *http.Request) error {
	return m.Next.Handle(responseWriter, request) // nolint:wrapcheck
}

func TestRouter_MatchExact(t *testing.T) {
	t.Parallel()

	testHandler := &testHandler{}

	router := NewRouter()
	router.Get("/test", testHandler)

	request := httptest.NewRequest(http.MethodGet, "/test", nil)

	handler, err := router.Match(request)
	if assert.NoError(t, err) {
		assert.Equal(t, testHandler, handler)
	}
}

func TestRouter_MatchRegex(t *testing.T) {
	t.Parallel()

	testHandler := &testHandler{}

	router := NewRouter()
	router.Get(`/test/{id:\d+}`, testHandler)

	request := httptest.NewRequest(http.MethodGet, "/test/1", nil)

	handler, err := router.Match(request)
	if assert.NoError(t, err) {
		assert.Equal(t, testHandler, handler)
	}

	assert.Equal(t, "1", RouteParam(request.Context(), "id"))
}

func TestRouter_MatchErrorNotFound(t *testing.T) {
	t.Parallel()

	router := NewRouter()
	router.Get("/test", &testHandler{})

	request := httptest.NewRequest(http.MethodGet, "/test2", nil)

	_, err := router.Match(request)
	assert.EqualError(t, err, ErrRouteNotFound.Error())
}

func TestRouter_MatchMethodNotAllowed(t *testing.T) {
	t.Parallel()

	testHandler := &testHandler{}

	router := NewRouter()
	router.Get("/test", testHandler)
	router.Get(`/test/{id:\d+}`, testHandler)

	request := httptest.NewRequest(http.MethodPost, "/test", nil)

	_, err := router.Match(request)
	assert.EqualError(t, err, ErrMethodNotAllowed.Error())

	request = httptest.NewRequest(http.MethodPost, "/test/1", nil)

	_, err = router.Match(request)
	assert.EqualError(t, err, ErrMethodNotAllowed.Error())
}

func TestRouter_MatchInvalidRegex(t *testing.T) {
	t.Parallel()

	testHandler := &testHandler{}

	router := NewRouter()
	router.Post(`/test/{slug:\s++}`, testHandler)

	request := httptest.NewRequest(http.MethodPost, "/test/foo", nil)

	_, err := router.Match(request)
	assert.ErrorContains(t, err, "httprouter, httprouter.Match, regexp.Compile, err:")
}

func TestRouter_Get(t *testing.T) {
	t.Parallel()

	router := NewRouter()

	testHandler := &testHandler{}
	router.Get("/test", testHandler)

	assert.Equal(t, 1, len(router.routes))
	assert.Equal(t, 1, len(router.routes[0].methods))

	assert.Contains(t, router.routes[0].methods, http.MethodGet)
	assert.Equal(t, "/test", router.routes[0].path)
	assert.Equal(t, testHandler, router.routes[0].handler)
}

func TestRouter_Post(t *testing.T) {
	t.Parallel()

	router := NewRouter()

	testHandler := &testHandler{}
	router.Post("/test", testHandler)

	assert.Equal(t, 1, len(router.routes))
	assert.Equal(t, 1, len(router.routes[0].methods))

	assert.Contains(t, router.routes[0].methods, http.MethodPost)
	assert.Equal(t, "/test", router.routes[0].path)
	assert.Equal(t, testHandler, router.routes[0].handler)
}

func TestRouter_Put(t *testing.T) {
	t.Parallel()

	router := NewRouter()

	testHandler := &testHandler{}
	router.Put("/test", testHandler)

	assert.Equal(t, 1, len(router.routes))
	assert.Equal(t, 1, len(router.routes[0].methods))

	assert.Contains(t, router.routes[0].methods, http.MethodPut)
	assert.Equal(t, "/test", router.routes[0].path)
	assert.Equal(t, testHandler, router.routes[0].handler)
}

func TestRouter_Delete(t *testing.T) {
	t.Parallel()

	router := NewRouter()

	testHandler := &testHandler{}
	router.Delete("/test", testHandler)

	assert.Equal(t, 1, len(router.routes))
	assert.Equal(t, 1, len(router.routes[0].methods))

	assert.Contains(t, router.routes[0].methods, http.MethodDelete)
	assert.Equal(t, "/test", router.routes[0].path)
	assert.Equal(t, testHandler, router.routes[0].handler)
}

func TestRouter_Patch(t *testing.T) {
	t.Parallel()

	router := NewRouter()

	testHandler := &testHandler{}
	router.Patch("/test", testHandler)

	assert.Equal(t, 1, len(router.routes))
	assert.Equal(t, 1, len(router.routes[0].methods))

	assert.Contains(t, router.routes[0].methods, http.MethodPatch)
	assert.Equal(t, "/test", router.routes[0].path)
	assert.Equal(t, testHandler, router.routes[0].handler)
}

func TestRouter_Options(t *testing.T) {
	t.Parallel()

	router := NewRouter()

	testHandler := &testHandler{}
	router.Options("/test", testHandler)

	assert.Equal(t, 1, len(router.routes))
	assert.Equal(t, 1, len(router.routes[0].methods))

	assert.Contains(t, router.routes[0].methods, http.MethodOptions)
	assert.Equal(t, "/test", router.routes[0].path)
	assert.Equal(t, testHandler, router.routes[0].handler)
}

func TestRouter_Head(t *testing.T) {
	t.Parallel()

	router := NewRouter()

	testHandler := &testHandler{}
	router.Head("/test", testHandler)

	assert.Equal(t, 1, len(router.routes))
	assert.Equal(t, 1, len(router.routes[0].methods))

	assert.Contains(t, router.routes[0].methods, http.MethodHead)
	assert.Equal(t, "/test", router.routes[0].path)
	assert.Equal(t, testHandler, router.routes[0].handler)
}

func TestRouter_Connect(t *testing.T) {
	t.Parallel()

	router := NewRouter()

	testHandler := &testHandler{}
	router.Connect("/test", testHandler)

	assert.Equal(t, 1, len(router.routes))
	assert.Equal(t, 1, len(router.routes[0].methods))

	assert.Contains(t, router.routes[0].methods, http.MethodConnect)
	assert.Equal(t, "/test", router.routes[0].path)
	assert.Equal(t, testHandler, router.routes[0].handler)
}

func TestRouter_Trace(t *testing.T) {
	t.Parallel()

	router := NewRouter()

	testHandler := &testHandler{}
	router.Trace("/test", testHandler)

	assert.Equal(t, 1, len(router.routes))
	assert.Equal(t, 1, len(router.routes[0].methods))

	assert.Contains(t, router.routes[0].methods, http.MethodTrace)
	assert.Equal(t, "/test", router.routes[0].path)
	assert.Equal(t, testHandler, router.routes[0].handler)
}

func TestRouter_Any(t *testing.T) {
	t.Parallel()

	router := NewRouter()

	testHandler := &testHandler{}
	router.Any("/test", []string{http.MethodGet, http.MethodPost}, testHandler)

	assert.Equal(t, 1, len(router.routes))
	assert.Equal(t, 2, len(router.routes[0].methods))

	assert.Contains(t, router.routes[0].methods, http.MethodGet)
	assert.Contains(t, router.routes[0].methods, http.MethodPost)

	assert.Equal(t, "/test", router.routes[0].path)
	assert.Equal(t, testHandler, router.routes[0].handler)
}

func TestRouter_Use(t *testing.T) {
	t.Parallel()

	router := NewRouter()

	middleware := func(next Handler) Handler {
		return &testMiddleware{
			Next: next,
		}
	}

	router.Use(middleware)

	testHandler := &testHandler{}
	router.Get("/test", testHandler)

	request := httptest.NewRequest(http.MethodGet, "/test", nil)

	handler, err := router.Match(request)
	if assert.NoError(t, err) {
		assert.Equal(t, middleware(testHandler), handler)
	}
}

func TestRouter_Group(t *testing.T) {
	t.Parallel()

	router := NewRouter()

	middleware1 := func(next Handler) Handler {
		return &testMiddleware{Next: next}
	}

	router.Use(middleware1)

	testHandler := &testHandler{}
	router.Get("/test", testHandler)

	middleware2 := func(next Handler) Handler {
		return &test2Middleware{Next: next}
	}

	middleware3 := func(next Handler) Handler {
		return &test3Middleware{Next: next}
	}

	router.Group(func(router Router) {
		router.Use(middleware2, middleware3)

		router.Get("/test2", testHandler)
	})

	t.Run("Group middleware", func(t *testing.T) {
		t.Parallel()

		request := httptest.NewRequest(http.MethodGet, "/test2", nil)

		handler, err := router.Match(request)
		if assert.NoError(t, err) {
			assert.Equal(t, middleware1(middleware2(middleware3(testHandler))), handler)
		}
	})

	t.Run("Common middleware", func(t *testing.T) {
		t.Parallel()

		request := httptest.NewRequest(http.MethodGet, "/test", nil)

		handler, err := router.Match(request)
		if assert.NoError(t, err) {
			assert.Equal(t, middleware1(testHandler), handler)
		}
	})
}

func TestRouter_ServeHTTPNotFound(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	responseWriter := httptest.NewRecorder()

	router := NewRouter()

	router.ServeHTTP(responseWriter, request)

	assert.Equal(t, responseWriter.Code, http.StatusNotFound)
}

func TestRouter_ServeHTTPMethodNotAllowed(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodPost, "/test", nil)
	responseWriter := httptest.NewRecorder()

	router := NewRouter()
	router.Get("/test", &testHandler{})

	router.ServeHTTP(responseWriter, request)

	assert.Equal(t, responseWriter.Code, http.StatusMethodNotAllowed)
}

func TestRouter_ServeHTTPInternalServerError(t *testing.T) {
	t.Parallel()

	t.Run("Router match error", func(t *testing.T) {
		t.Parallel()

		request := httptest.NewRequest(http.MethodGet, "/test", nil)
		responseWriter := httptest.NewRecorder()

		router := NewRouter()
		router.Get(`/test/{slug:\s++}`, &testHandler{})

		router.ServeHTTP(responseWriter, request)

		assert.Equal(t, responseWriter.Code, http.StatusInternalServerError)
	})

	t.Run("Handler error", func(t *testing.T) {
		t.Parallel()

		request := httptest.NewRequest(http.MethodGet, "/test", nil)
		responseWriter := httptest.NewRecorder()

		router := NewRouter()
		router.Get(`/test`, HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) error {
			return errors.New("some error")
		}))

		router.ServeHTTP(responseWriter, request)

		assert.Equal(t, responseWriter.Code, http.StatusInternalServerError)
	})

	t.Run("NotFound handler error", func(t *testing.T) {
		t.Parallel()

		request := httptest.NewRequest(http.MethodGet, "/test", nil)
		responseWriter := httptest.NewRecorder()

		router := NewRouter()

		router.NotFoundHandler = HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) error {
			return errors.New("some error")
		})

		router.ServeHTTP(responseWriter, request)

		assert.Equal(t, responseWriter.Code, http.StatusInternalServerError)
	})
}
