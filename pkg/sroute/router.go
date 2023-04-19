package sroute

type router struct {
	trees map[string]*node
}

type node struct {
	path     string
	children map[string]*node
}

func (r *router) AddRoute(method string, path string, handleF HandleFunc) {

}
