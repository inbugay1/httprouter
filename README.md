# HTTP Router

Just another simple http router with middleware and error handling support.

## Preview

HTTP router has a simple interface with several methods:

````
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
}
````

It also implements net/http Handler interface and can be used as net/http server handler.

## Examples

### Closure handler

````
func main() {
	router := httprouter.NewRouter()

	helloHandler := func(responseWriter http.ResponseWriter, request *http.Request) error {
		_, _ = responseWriter.Write([]byte("Hello World!"))

		return nil
	}

	router.Get("/hello", httprouter.HandlerFunc(helloHandler))

	_ = http.ListenAndServe(":9015", router)
}
````

### Non closure handler

Often your handlers will be located in their own package and will need some dependencies like logger, repositories and so on.
For sure, you can achieve that by using closure handler like this:

````
package handler

import (
	"net/http"

	"getting-started/internal/httprouter"
	"github.com/sirupsen/logrus"
)

func Test(logger *logrus.Logger) httprouter.Handler {
	handler := func(responseWriter http.ResponseWriter, request *http.Request) error {
		logger.Info("Some useful message could be here...")
		
		responseWriter.Header().Set(httphelper.HeaderContentType, httphelper.ContentTypeJSONUTF8)

		_, _ = responseWriter.Write([]byte("OK"))
		
		return nil
	}

	return httprouter.HandlerFunc(handler)
}
````

````
func main() {
	router := httprouter.NewRouter()
	
	logger := &logrus.Logger{}

	router.Get("/hello", handler.Test(logger))

	_ = http.ListenAndServe(":9015", router)
}
````

But with these functions, closures and HandlerFunc wrap it looks a bit cumbersome. 
A better approach probably would be to define a struct that implements httprouter.Handler interface:

````
package handler

import (
	"net/http"

	"getting-started/internal/httphelper"
	"github.com/sirupsen/logrus"
)

type Test struct {
	Logger *logrus.Logger
}

func (h *Test) Handle(responseWriter http.ResponseWriter, request *http.Request) error {
	h.Logger.Info("Some useful message could be here...")
	responseWriter.Header().Set(httphelper.HeaderContentType, httphelper.ContentTypeJSONUTF8)

	_, _ = responseWriter.Write([]byte("OK"))

	return nil
}
````

````
func main() {
	router := httprouter.NewRouter()
	
	logger := &logrus.Logger{}

	router.Get("/hello", &handler.Test{Logger: logger})

	_ = http.ListenAndServe(":9015", router)
}
````

### Middleware

````
func main() {
	router := httprouter.NewRouter()

	helloHandler := httprouter.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) error {
			_, _ = responseWriter.Write([]byte("Hello "))

			return nil
		})

	worldMiddleware := func(next httprouter.Handler) httprouter.Handler {
		handler := func(responseWriter http.ResponseWriter, request *http.Request) error {
			if err := next.Handle(responseWriter, request); err != nil {
				return err
			}

			_, _ = responseWriter.Write([]byte("World!"))

			return nil
		}

		return httprouter.HandlerFunc(handler)
	}

	router.Get("/hello", worldMiddleware(helloHandler))

	_ = http.ListenAndServe(":9015", router)
}
````

When you create a middleware you will probably want to place it in its own package and like for handlers 
you can define a struct and implement httprouter.Handler interface:

````
package middleware

import (
	"net/http"

	"getting-started/internal/httprouter"
	"github.com/sirupsen/logrus"
)

type test struct {
	logger *logrus.Logger
	next   httprouter.Handler
}

func (m *test) Handle(responseWriter http.ResponseWriter, request *http.Request) error {
	err := m.next.Handle(responseWriter, request)
	if err != nil {
		return err
	}

	m.logger.Info("This is log from test middleware")

	return nil
}

func Test(logger *logrus.Logger) httprouter.MiddlewareFunc {
	return func(next httprouter.Handler) httprouter.Handler {
		return &test{
			logger: logger,
			next:   next,
		}
	}
}
````

````
func main() {
	router := httprouter.NewRouter()

	logger := &logrus.Logger{}

	middleware := middleware.Test(logger)

	router.Get("/hello", middleware(&handler.Test{Logger: logger}))

	_ = http.ListenAndServe(":9015", router)
}
````

### Router Use method

It takes a list of middlewares and wrap a handler:

````
Router.Use(middleware1, middleware2)
Router.Get("/hello", helloHandler)
````
It will be the same as:
````
Router.Get("/hello", middleware1(middleware2(helloHandler)))
````

You can try the Use method to apply a middleware to several handlers:

````
func main() {
	router := httprouter.NewRouter()

	helloHandler := httprouter.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) error {
			_, _ = responseWriter.Write([]byte("Hello "))

			return nil
		})

	goodbyeHandler := httprouter.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) error {
			_, _ = responseWriter.Write([]byte("Goodbye "))

			return nil
		})

	worldMiddleware := func(next httprouter.Handler) httprouter.Handler {
		handler := func(responseWriter http.ResponseWriter, request *http.Request) error {
			if err := next.Handle(responseWriter, request); err != nil {
				return err
			}

			_, _ = responseWriter.Write([]byte("World!"))

			return nil
		}

		return httprouter.HandlerFunc(handler)
	}

	router.Use(worldMiddleware)

	router.Get("/hello", helloHandler)
	router.Get("/goodbye", goodbyeHandler)

	_ = http.ListenAndServe(":9015", router)
}
````

### Router Group method
Sometimes you may find yourself wanting to group routes in order to apply a middleware to them.
With the Group method you can easily do this:

````
func main() {
	router := httprouter.NewRouter()

	helloHandler := httprouter.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) error {
			_, _ = responseWriter.Write([]byte("Hello "))

			return nil
		})

	goodbyeHandler := httprouter.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) error {
			_, _ = responseWriter.Write([]byte("Goodbye "))

			return nil
		})

	worldMiddleware := func(next httprouter.Handler) httprouter.Handler {
		handler := func(responseWriter http.ResponseWriter, request *http.Request) error {
			if err := next.Handle(responseWriter, request); err != nil {
				return err
			}

			_, _ = responseWriter.Write([]byte("World!"))

			return nil
		}

		return httprouter.HandlerFunc(handler)
	}

	router.Group(func(router httprouter.Router) {
		router.Use(worldMiddleware)

		router.Get("/hello", helloHandler)
		router.Get("/goodbye", goodbyeHandler)
	})

	router.Get("/test", httprouter.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) error {
			_, _ = responseWriter.Write([]byte("Test"))

			return nil
		}))

	_ = http.ListenAndServe(":9015", router)
}
````

### Custom NotFound handler

Easy peasy

````
func main() {
	router := httprouter.NewRouter()

	notFoundHandler := func(responseWriter http.ResponseWriter, request *http.Request) error {
		_, _ = responseWriter.Write([]byte("Oops... we cannot find what you want :'("))

		return nil
	}

	router.NotFoundHandler = httprouter.HandlerFunc(notFoundHandler)

	helloHandler := func(responseWriter http.ResponseWriter, request *http.Request) error {
		_, _ = responseWriter.Write([]byte("Hello World!"))

		return nil
	}

	router.Get("/hello", httprouter.HandlerFunc(helloHandler))

	_ = http.ListenAndServe(":9015", router)
}
````

### RouteParam

HTTP router supports parameterized routes

````
func main() {
	router := httprouter.NewRouter()

	helloHandler := func(responseWriter http.ResponseWriter, request *http.Request) error {
		name := httprouter.RouteParam(request.Context(), "name")

		_, _ = responseWriter.Write([]byte("Hello " + name + "!"))

		return nil
	}

	router.Get(`/hello/{name:[a-z]+}`, httprouter.HandlerFunc(helloHandler))

	_ = http.ListenAndServe(":9015", router)
}
````