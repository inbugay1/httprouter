package main

import (
	"net/http"

	"github.com/inbugay1/httprouter"
)

func main() {
	router := httprouter.New()

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
