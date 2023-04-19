package sroute

import (
	"net"
	"net/http"
)

var _ Server = &HTTPServer{}

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler
	Start(addr string) error
	addRoute(method string, path string, handleF HandleFunc)
}

type HTTPServer struct {
	*router
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Req:        req,
		RespWriter: writer,
	}
	h.serve(ctx)
}

func (h *HTTPServer) serve(ctx *Context) {
	n, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || n.handler == nil {
		// no route found
		ctx.RespWriter.WriteHeader(http.StatusNotFound)
		_, err := ctx.RespWriter.Write([]byte("Page Not Found"))
		if err != nil {
			panic("failed to write content")
		}
		return
	}
	n.handler(ctx)
}

func (h *HTTPServer) Get(path string, handleF HandleFunc) {
	h.addRoute(http.MethodGet, path, handleF)
}

func (h *HTTPServer) Post(path string, handleF HandleFunc) {
	h.addRoute(http.MethodPost, path, handleF)
}

func (h *HTTPServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return http.Serve(l, h)
}
