package main

import (
	"net/http"

	"github.com/inbugay1/httprouter"
)

func main() {
	router := httprouter.New()

	helloHandler := func(responseWriter http.ResponseWriter, request *http.Request) error {
		_, _ = responseWriter.Write([]byte("Hello World!"))

		return nil
	}

	router.Get("/hello", httprouter.HandlerFunc(helloHandler))

	_ = http.ListenAndServe(":9015", router)
}
