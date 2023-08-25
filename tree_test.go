package httprouter_test

import (
	"encoding/json"
	"testing"

	"github.com/inbugay1/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	t.Parallel()

	tree := httprouter.NewTree()

	handler1 := &mockHandler{}
	handler2 := &mockHandler{}

	tree.Insert("/path/to/resource", handler1)
	tree.Insert("/path/to/resource2", handler2)
	tree.Insert("/path/to/:id", handler2)

	expected := `{"root":{"key":"","static_children":[{"key":"path","static_children":[{"key":"to","static_children":[{"key":"resource"},{"key":"resource2"}],"dynamic_child":{"key":":id"}}]}]}}`

	actual, err := json.Marshal(tree)
	if assert.NoError(t, err) {
		assert.JSONEq(t, expected, string(actual), "Tree JSON representation does not match")
	}
}

func TestSearch(t *testing.T) {
	t.Parallel()

	tree := httprouter.NewTree()

	handler1 := &mockHandler{}
	handler2 := &mockHandler{}

	tree.Insert("/path/to/resource", handler1)
	tree.Insert("/path/to/resource2", handler2)
	tree.Insert("/path/to/:id", handler2)

	tests := []struct {
		name    string
		path    string
		handler httprouter.Handler
		params  map[string]string
	}{
		{"StaticPath1", "/path/to/resource", handler1, map[string]string{}},
		{"StaticPath2", "/path/to/resource2", handler2, map[string]string{}},
		{"DynamicPath", "/path/to/123", handler2, map[string]string{"id": "123"}},
		{"PathWithoutHandler", "/path/to", nil, map[string]string{}},
		{"NonExistentPath", "/path/not/in/tree", nil, nil},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			handler, params := tree.Search(testCase.path)
			assert.Equal(t, testCase.handler, handler)
			assert.Equal(t, testCase.params, params)
		})
	}
}
