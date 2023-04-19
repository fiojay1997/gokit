package sroute

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestAddRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
	}

	mockHandler := func(ctx Context) {}

	r := newRouter()
	for _, route := range testRoutes {
		r.AddRoute(route.method, route.path, mockHandler)
	}

	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path: "/",
				children: map[string]*node{
					"user": {
						path: "user",
						children: map[string]*node{
							"home": {
								path:     "home",
								children: map[string]*node{},
								handler:  mockHandler,
							},
						},
					},
				},
			},
		},
	}

	if equal, error := wantRouter.equals(r); !equal && error != nil {
		t.Errorf("router content not equal: %+v", error)
	}
}

func (r *router) equals(other *router) (bool, error) {
	for k, v := range r.trees {
		dst, ok := other.trees[k]
		if !ok {
			return false, errors.New("method not equal")
		}
		equal := v.equals(dst)
		if !equal {
			return false, errors.New("node value not equal")
		}
	}
	return true, nil
}

func (n *node) equals(other *node) bool {
	if n.path != other.path {
		return false
	}
	if len(n.children) != len(other.children) {
		return false
	}

	// compare handlers
	nHandler := reflect.ValueOf(n.handler)
	oHandler := reflect.ValueOf(other.handler)
	if nHandler != oHandler {
		return false
	}

	// compare node values
	for path, c := range n.children {
		dst, ok := other.children[path]
		if !ok {
			return false
		}
		ok = c.equals(dst)
		if !ok {
			return false
		}
	}
	return true
}
