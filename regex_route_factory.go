package httprouter

import "regexp"

type regexRouteFactory struct {
	regexp *regexp.Regexp
}

func NewRegexRouteFactory() *regexRouteFactory { //nolint:golint,revive
	return &regexRouteFactory{
		regexp: regexp.MustCompile(`{(?P<param>\w+):(?P<regex>.+?)}`),
	}
}

func (f *regexRouteFactory) Name() string {
	return "regex"
}

func (f *regexRouteFactory) Handles(path string) bool {
	return f.regexp.MatchString(path)
}

func (f *regexRouteFactory) CreateRoute(path string, methods []string, handler Handler) Route {
	pathRegexStr := f.regexp.ReplaceAllString(path, "(?P<$1>$2)") // e.g modify /test/{id:\d+} to /test/(?P<id>\d+)

	return &RegexRoute{
		Methods: methods,
		Handler: handler,
		Regexp:  regexp.MustCompile("^" + pathRegexStr + "$"),
	}
}
