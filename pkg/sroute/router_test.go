package sroute

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodPost,
			path:   "/",
		},
		{
			method: http.MethodPost,
			path:   "/user",
		},
		{
			method: http.MethodPost,
			path:   "/user/home",
		},
		{
			method: http.MethodDelete,
			path:   "user",
		},
		{
			method: http.MethodDelete,
			path:   "order/",
		},
	}

	mockHandler := func(ctx Context) {}

	r := newRouter()
	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}

	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path:    "/",
				handler: mockHandler,
				children: map[string]*node{
					"order": {
						path: "order",
						children: map[string]*node{
							"detail": {
								path:     "detail",
								children: map[string]*node{},
								handler:  mockHandler,
							},
						},
					},
					"user": {
						path:    "user",
						handler: mockHandler,
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
			http.MethodPost: {
				path:    "/",
				handler: mockHandler,
				children: map[string]*node{
					"user": {
						path:    "user",
						handler: mockHandler,
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
			http.MethodDelete: {
				path: "/",
				children: map[string]*node{
					"user": {
						path:    "user",
						handler: mockHandler,
					},
					"order": {
						path:    "order",
						handler: mockHandler,
					},
				},
			},
		},
	}

	if equal, error := wantRouter.equals(r); !equal && error != nil {
		t.Errorf("router content not equal: %+v", error)
	}

	r = newRouter()
	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "", mockHandler)
	})
	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "/a//", mockHandler)
	})
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
