package httprouter

import (
	"net/http"
	"regexp"
)

type RegexRoute struct {
	Methods []string
	Handler Handler
	Regexp  *regexp.Regexp
}

func (regexRoute *RegexRoute) Match(request *http.Request) (RouteMatch, error) {
	var routeMatch RouteMatch
	if !regexRoute.Regexp.MatchString(request.URL.Path) {
		return routeMatch, ErrPathMismatch
	}

	if !contains(regexRoute.Methods, request.Method) {
		return routeMatch, ErrMethodNotAllowed
	}

	matches := regexRoute.Regexp.FindAllStringSubmatch(request.URL.Path, -1)

	routeParamNames := regexRoute.Regexp.SubexpNames()
	routeParams := make(map[string]string, len(routeParamNames))

	for paramKey, paramName := range routeParamNames {
		if paramName == "" {
			continue
		}
		routeParams[paramName] = matches[0][paramKey]
	}

	routeMatch.Handler = regexRoute.Handler
	routeMatch.Params = routeParams

	return routeMatch, nil
}
