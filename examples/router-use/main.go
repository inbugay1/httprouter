package main

import (
	"net/http"

	"github.com/inbugay1/httprouter"
)

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
