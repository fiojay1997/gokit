package sroute

import (
	"strings"
)

type router struct {
	trees map[string]*node
}

type node struct {
	path     string
	children map[string]*node
	handler  HandleFunc
}

func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

// AddRoute adds a path variables with a handler to the router
func (r *router) addRoute(method string, path string, handleF HandleFunc) {
	if strings.Contains(path, " ") || path == "" || len(path) == 0 {
		panic("path variables cannot be empty")
	}

	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	if path[0] != '/' {
		path = "/" + path
	}

	if path != "/" && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	if path == "/" {
		root.handler = handleF
		return
	}

	path = path[1:]
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		if seg == "" {
			panic("path variables cannot be empty")
		}
		child := root.childOrCreate(seg)
		root = child
	}
	root.handler = handleF
}

func (r *router) findRoute(method string, path string) (*node, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	path = strings.Trim(path, "/")
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		child, found := root.childOf(seg)
		if !found {
			return nil, false
		}
		root = child
	}
	return root, root.handler != nil
}

func (n *node) childOf(path string) (*node, bool) {
	if n.children == nil {
		return nil, false
	}
	child, ok := n.children[path]
	return child, ok
}

func (n *node) childOrCreate(seg string) *node {
	if n.children == nil {
		n.children = map[string]*node{}
	}
	res, ok := n.children[seg]
	if !ok {
		res = &node{
			path: seg,
		}
		n.children[seg] = res
	}
	return res
}
