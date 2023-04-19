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
