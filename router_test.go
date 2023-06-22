package httprouter_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/inbugay1/httprouter"
	"github.com/stretchr/testify/assert"
)

func makeHandlerWithResponseContent(responseContent string) httprouter.Handler { //nolint:ireturn
	return httprouter.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) error {
		_, _ = responseWriter.Write([]byte(responseContent))

		return nil
	})
}

func makeHandlerWithError(text string) httprouter.Handler { //nolint:ireturn
	return httprouter.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) error {
		return errors.New(text)
	})
}

func makeDummyHandler() httprouter.Handler { //nolint:ireturn
	return httprouter.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) error {
		return nil
	})
}

//nolint:wrapcheck
func makeMiddleware(content string) httprouter.MiddlewareFunc {
	return func(next httprouter.Handler) httprouter.Handler {
		handler := func(responseWriter http.ResponseWriter, request *http.Request) error {
			_, _ = responseWriter.Write([]byte(content))

			return next.Handle(responseWriter, request)
		}

		return httprouter.HandlerFunc(handler)
	}
}

func TestRouterMatch(t *testing.T) { //nolint:funlen
	t.Parallel()

	t.Run("exact", func(t *testing.T) {
		t.Parallel()

		router := httprouter.New()
		router.Get("/test", makeHandlerWithResponseContent("test"))

		request := httptest.NewRequest(http.MethodGet, "/test", nil)
		responseWriter := httptest.NewRecorder()

		handler, err := router.Match(request)
		if assert.NoError(t, err) {
			_ = handler.Handle(responseWriter, request)
			assert.Equal(t, "test", responseWriter.Body.String())
		}
	})

	t.Run("regex", func(t *testing.T) {
		t.Parallel()

		router := httprouter.New()
		router.Get(`/test/{id:\d+}`, makeHandlerWithResponseContent("test"))

		request := httptest.NewRequest(http.MethodGet, "/test/1", nil)
		responseWriter := httptest.NewRecorder()

		handler, err := router.Match(request)
		if assert.NoError(t, err) {
			_ = handler.Handle(responseWriter, request)
			assert.Equal(t, "test", responseWriter.Body.String())
		}

		assert.Equal(t, "1", httprouter.RouteParam(request.Context(), "id"))
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		router := httprouter.New()
		router.Get("/test", makeDummyHandler())

		request := httptest.NewRequest(http.MethodGet, "/test2", nil)

		_, err := router.Match(request)
		assert.EqualError(t, err, httprouter.ErrRouteNotFound.Error())
	})

	t.Run("method not allowed", func(t *testing.T) {
		t.Parallel()

		router := httprouter.New()
		router.Get("/test", makeDummyHandler())
		router.Get(`/test/{id:\d+}`, makeDummyHandler())

		request := httptest.NewRequest(http.MethodPost, "/test", nil)

		_, err := router.Match(request)
		assert.EqualError(t, err, httprouter.ErrMethodNotAllowed.Error())

		request = httptest.NewRequest(http.MethodPost, "/test/1", nil)

		_, err = router.Match(request)
		assert.EqualError(t, err, httprouter.ErrMethodNotAllowed.Error())
	})

	t.Run("invalid route regex", func(t *testing.T) {
		t.Parallel()

		router := httprouter.New()
		router.Post(`/test/{slug:\s++}`, makeDummyHandler())

		request := httptest.NewRequest(http.MethodPost, "/test/foo", nil)

		_, err := router.Match(request)
		assert.ErrorContains(t, err, "invalid route regex")
	})
}

func TestRouterHTTPMethods(t *testing.T) {
	t.Parallel()

	router := httprouter.New()
	router.Get("/test/get", makeHandlerWithResponseContent(http.MethodGet))
	router.Post("/test/post", makeHandlerWithResponseContent(http.MethodPost))
	router.Put("/test/put", makeHandlerWithResponseContent(http.MethodPut))
	router.Patch("/test/patch", makeHandlerWithResponseContent(http.MethodPatch))
	router.Delete("/test/delete", makeHandlerWithResponseContent(http.MethodDelete))
	router.Options("/test/options", makeHandlerWithResponseContent(http.MethodOptions))
	router.Connect("/test/connect", makeHandlerWithResponseContent(http.MethodConnect))
	router.Trace("/test/trace", makeHandlerWithResponseContent(http.MethodTrace))
	router.Head("/test/head", makeHandlerWithResponseContent(http.MethodHead))

	testMethod := func(method string) {
		responseWriter := httptest.NewRecorder()
		request := httptest.NewRequest(method, "/test/"+strings.ToLower(method), nil)

		handler, err := router.Match(request)
		if assert.NoError(t, err) {
			_ = handler.Handle(responseWriter, request)
			assert.Equal(t, method, responseWriter.Body.String())
		}
	}

	testMethod(http.MethodGet)
	testMethod(http.MethodPost)
	testMethod(http.MethodPatch)
	testMethod(http.MethodPut)
	testMethod(http.MethodDelete)
	testMethod(http.MethodOptions)
	testMethod(http.MethodTrace)
	testMethod(http.MethodHead)
}

func TestRouterAny(t *testing.T) {
	t.Parallel()

	router := httprouter.New()

	router.Any("/test", []string{http.MethodGet, http.MethodPost}, makeHandlerWithResponseContent("test"))

	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	responseWriter := httptest.NewRecorder()

	handler, err := router.Match(request)
	if assert.NoError(t, err) {
		_ = handler.Handle(responseWriter, request)
		assert.Equal(t, "test", responseWriter.Body.String())
	}

	request = httptest.NewRequest(http.MethodPost, "/test", nil)
	responseWriter = httptest.NewRecorder()

	handler, err = router.Match(request)
	if assert.NoError(t, err) {
		_ = handler.Handle(responseWriter, request)
		assert.Equal(t, "test", responseWriter.Body.String())
	}
}

func TestRouterUse(t *testing.T) {
	t.Parallel()

	router := httprouter.New()

	middleware1 := makeMiddleware("hello ")
	middleware2 := makeMiddleware("world")

	router.Use(middleware1, middleware2)

	router.Get("/test", makeDummyHandler())

	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	responseWriter := httptest.NewRecorder()

	handler, err := router.Match(request)
	if assert.NoError(t, err) {
		_ = handler.Handle(responseWriter, request)

		assert.Equal(t, "hello world", responseWriter.Body.String())
	}
}

func TestRouterGroup(t *testing.T) {
	t.Parallel()

	router := httprouter.New()

	middleware1 := makeMiddleware("hello ")
	middleware2 := makeMiddleware("world")

	router.Use(middleware1)

	router.Group(func(r httprouter.Router) {
		r.Use(middleware2)

		router.Get("/helloworld", makeDummyHandler())
	})

	router.Get("/hello", makeDummyHandler())

	request := httptest.NewRequest(http.MethodGet, "/helloworld", nil)
	responseWriter := httptest.NewRecorder()

	handler, err := router.Match(request)
	if assert.NoError(t, err) {
		_ = handler.Handle(responseWriter, request)

		assert.Equal(t, "hello world", responseWriter.Body.String())
	}

	request = httptest.NewRequest(http.MethodGet, "/hello", nil)
	responseWriter = httptest.NewRecorder()

	handler, err = router.Match(request)
	if assert.NoError(t, err) {
		_ = handler.Handle(responseWriter, request)

		assert.Equal(t, "hello ", responseWriter.Body.String())
	}
}

func TestRouterServeHTTP(t *testing.T) { //nolint: funlen
	t.Parallel()

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		request := httptest.NewRequest(http.MethodGet, "/test", nil)
		responseWriter := httptest.NewRecorder()

		router := httprouter.New()

		router.ServeHTTP(responseWriter, request)

		assert.Equal(t, responseWriter.Code, http.StatusNotFound)
	})

	t.Run("method not allowed", func(t *testing.T) {
		t.Parallel()

		request := httptest.NewRequest(http.MethodPost, "/test", nil)
		responseWriter := httptest.NewRecorder()

		router := httprouter.New()
		router.Get("/test", makeDummyHandler())

		router.ServeHTTP(responseWriter, request)

		assert.Equal(t, responseWriter.Code, http.StatusMethodNotAllowed)
	})

	t.Run("internal server error caused by router match", func(t *testing.T) {
		t.Parallel()

		request := httptest.NewRequest(http.MethodGet, "/test", nil)
		responseWriter := httptest.NewRecorder()

		router := httprouter.New()
		router.Get(`/test/{slug:\s++}`, makeDummyHandler())

		router.ServeHTTP(responseWriter, request)

		assert.Equal(t, responseWriter.Code, http.StatusInternalServerError)
	})

	t.Run("internal server error caused by handler", func(t *testing.T) {
		t.Parallel()

		request := httptest.NewRequest(http.MethodGet, "/test", nil)
		responseWriter := httptest.NewRecorder()

		router := httprouter.New()
		router.Get(`/test`, makeHandlerWithError("some error"))

		router.ServeHTTP(responseWriter, request)

		assert.Equal(t, responseWriter.Code, http.StatusInternalServerError)
	})

	t.Run("internal server error caused by notfound handler", func(t *testing.T) {
		t.Parallel()

		request := httptest.NewRequest(http.MethodGet, "/test", nil)
		responseWriter := httptest.NewRecorder()

		router := httprouter.New()

		router.NotFoundHandler = makeHandlerWithError("some error")

		router.ServeHTTP(responseWriter, request)

		assert.Equal(t, responseWriter.Code, http.StatusInternalServerError)
	})
}
