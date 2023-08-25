package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/inbugay1/httprouter"
)

func main() {
	router := httprouter.New(httprouter.NewPlaceholderRouteFactory())

	helloHandler := func(responseWriter http.ResponseWriter, request *http.Request) error {
		name := httprouter.RouteParam(request.Context(), "name")

		_, _ = responseWriter.Write([]byte("Hello " + name + "!"))

		return nil
	}

	router.Get(`/hello/:name`, httprouter.HandlerFunc(helloHandler))

	_ = http.ListenAndServe(":9015", router)
}
