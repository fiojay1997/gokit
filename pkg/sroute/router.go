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

func (r *router) AddRoute(method string, path string, handleF HandleFunc) {
	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	path = path[1:]
	segs := strings.Split(path, "/")
	for _, seg := range segs {
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
