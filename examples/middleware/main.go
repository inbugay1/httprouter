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

	router.Get("/hello", worldMiddleware(helloHandler), "")

	_ = http.ListenAndServe(":9015", router)
}
