package sroute

import (
	"net"
	"net/http"
)

var _ Server = &HTTPServer{}

type HandleFunc func(ctx Context)

type Server interface {
	http.Handler
	Start(addr string) error
	AddRoute(method string, path string, handleF HandleFunc)
}

type HTTPServer struct {
	r router
}

func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Req:        req,
		RespWriter: writer,
	}
	h.serve(ctx)
}

func (h *HTTPServer) serve(ctx *Context) {

}

func (h *HTTPServer) Get(path string, handleF HandleFunc) {
	h.AddRoute(http.MethodGet, path, handleF)
}

func (h *HTTPServer) Post(path string, handleF HandleFunc) {
	h.AddRoute(http.MethodPost, path, handleF)
}

func (h *HTTPServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return http.Serve(l, h)
}

func (h *HTTPServer) AddRoute(method string, path string, handleF HandleFunc) {

}
