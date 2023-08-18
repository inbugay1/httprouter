package main

import (
	"net/http"

	"github.com/inbugay1/httprouter"
)

func main() {
	router := httprouter.New(httprouter.NewRegexRouteFactory())

	listHandler := httprouter.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) error {
			_, _ = responseWriter.Write([]byte("users list"))

			return nil
		})

	showHandler := httprouter.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) error {
			_, _ = responseWriter.Write([]byte("user " + httprouter.RouteParam(request.Context(), "id")))

			return nil
		})

	createHandler := httprouter.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) error {
			_, _ = responseWriter.Write([]byte("create user"))

			return nil
		})

	router.WithPrefix("api")

	router.Group(func(router httprouter.Router) {
		router.WithPrefix("users")

		router.Get(`/{id:\d+}`, showHandler) // GET http://localhost:9015/api/users/1
		router.Post("", createHandler)       // POST http://localhost:9015/api/users
		router.Get("", listHandler)          // GET http://localhost:9015/api/users
	})

	_ = http.ListenAndServe(":9015", router)
}
