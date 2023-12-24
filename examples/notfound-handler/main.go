package main

import (
	"net/http"

	"github.com/inbugay1/httprouter"
)

func main() {
	router := httprouter.New()

	notFoundHandler := func(responseWriter http.ResponseWriter, request *http.Request) error {
		_, _ = responseWriter.Write([]byte("Oops... we cannot find what you want :'("))

		return nil
	}

	router.NotFoundHandler = httprouter.HandlerFunc(notFoundHandler)

	helloHandler := func(responseWriter http.ResponseWriter, request *http.Request) error {
		_, _ = responseWriter.Write([]byte("Hello World!"))

		return nil
	}

	router.Get("/hello", httprouter.HandlerFunc(helloHandler), "")

	_ = http.ListenAndServe(":9015", router)
}
