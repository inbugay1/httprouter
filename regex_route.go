package httprouter

import (
	"context"
	"net/http"
	"regexp"
)

type RegexRoute struct {
	Methods []string
	Handler Handler
	Regexp  *regexp.Regexp
}

func (regexRoute *RegexRoute) Match(request *http.Request) (Handler, error) {
	if !regexRoute.Regexp.MatchString(request.URL.Path) {
		return nil, ErrPathMismatch
	}

	if !contains(regexRoute.Methods, request.Method) {
		return nil, ErrMethodNotAllowed
	}

	matches := regexRoute.Regexp.FindAllStringSubmatch(request.URL.Path, -1)

	routeParamNames := regexRoute.Regexp.SubexpNames()
	routeParams := make(map[string]string, len(routeParamNames))

	for paramKey, paramName := range routeParamNames {
		routeParams[paramName] = matches[0][paramKey]
	}

	ctx := context.WithValue(request.Context(), routeParamsKey, routeParams)
	*request = *request.WithContext(ctx)

	return regexRoute.Handler, nil
}
