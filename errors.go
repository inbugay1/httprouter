package httprouter

import "errors"

var ErrMethodNotAllowed = errors.New("httprouter: method not allowed")
var ErrRouteNotFound = errors.New("httprouter: route not found")
