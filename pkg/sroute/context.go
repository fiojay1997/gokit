package sroute

import "net/http"

type Context struct {
	Req        *http.Request
	RespWriter http.ResponseWriter
}
