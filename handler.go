package httprouter

import (
	"net/http"
)

type Handler interface {
	Handle(responseWriter http.ResponseWriter, request *http.Request) error
}

type HandlerFunc func(responseWriter http.ResponseWriter, request *http.Request) error

func (f HandlerFunc) Handle(responseWriter http.ResponseWriter, request *http.Request) error {
	return f(responseWriter, request)
}

type MiddlewareFunc func(next Handler) Handler
