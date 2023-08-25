package httprouter

import (
	"regexp"
)

type placeholderRouteFactory struct {
	regexp *regexp.Regexp
	tree   *tree
}

func NewPlaceholderRouteFactory() *placeholderRouteFactory { //nolint:golint,revive
	return &placeholderRouteFactory{
		regexp: regexp.MustCompile(`.*/:[^/]+.*`),
		tree:   NewTree(),
	}
}

func (f *placeholderRouteFactory) Name() string {
	return "placeholder"
}

func (f *placeholderRouteFactory) Handles(path string) bool {
	return f.regexp.MatchString(path)
}

func (f *placeholderRouteFactory) CreateRoute(path string, methods []string, handler Handler) Route {
	_ = f.tree.Insert(path, handler)

	return &PlaceholderRoute{
		Methods: methods,
		Tree:    f.tree,
	}
}
