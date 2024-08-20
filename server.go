package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Server interface {
	Serve(rw http.ResponseWriter, req *http.Request)
	Address() string
}

type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

func (s *simpleServer) Address() string {
	return s.addr
}

func newSimpleServer(addr string) *simpleServer {
	serverUrl, err := url.Parse("http://" + addr)
	check(err)
	return &simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}
