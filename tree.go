package httprouter

import (
	"strings"
)

type node struct {
	Key            string  `json:"key"`
	Value          Handler `json:"-"`
	StaticChildren []*node `json:"static_children,omitempty"`
	DynamicChild   *node   `json:"dynamic_child,omitempty"`
}

func (n *node) findStaticChildByKey(key string) *node {
	for _, child := range n.StaticChildren {
		if child.Key == key {
			return child
		}
	}

	return nil
}

type tree struct {
	Root *node `json:"root"`
}

func NewTree() *tree { //nolint:golint,revive
	return &tree{
		Root: &node{
			StaticChildren: make([]*node, 0),
		},
	}
}

func (tree *tree) Insert(path string, handler Handler) *node {
	currentNode := tree.Root

	start := 0

	path = strings.Trim(path, "/")

	//nolint
	for idx := 0; idx < len(path); idx++ {
		if path[idx] == '/' || idx == len(path)-1 {
			end := idx
			if idx == len(path)-1 {
				end = idx + 1
			}
			segment := path[start:end]

			if strings.HasPrefix(segment, ":") {
				if currentNode.DynamicChild == nil {
					currentNode.DynamicChild = &node{
						Key: segment,
					}
				}
				currentNode = currentNode.DynamicChild
			} else if child := currentNode.findStaticChildByKey(segment); child != nil {
				currentNode = child
			} else {
				newNode := &node{
					Key: segment,
				}

				currentNode.StaticChildren = append(currentNode.StaticChildren, newNode)
				currentNode = newNode
			}

			start = idx + 1
		}
	}

	currentNode.Value = handler

	return currentNode
}

func (tree *tree) Search(path string) (Handler, map[string]string) {
	currentNode := tree.Root
	start := 0
	params := make(map[string]string)

	path = strings.Trim(path, "/")

	for idx := 0; idx < len(path); idx++ {
		if path[idx] == '/' || idx == len(path)-1 {
			end := idx
			if idx == len(path)-1 {
				end = idx + 1
			}
			segment := path[start:end]

			if child := currentNode.findStaticChildByKey(segment); child != nil {
				currentNode = child
			} else if currentNode.DynamicChild != nil {
				currentNode = currentNode.DynamicChild
				params[currentNode.Key[1:]] = segment
			} else {
				return nil, nil
			}

			start = idx + 1
		}
	}

	return currentNode.Value, params
}
