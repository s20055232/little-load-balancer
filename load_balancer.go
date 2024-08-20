package main

import (
	"log/slog"
	"net/http"
)

// 做一個自己的 Tuple，Python 使用者開心
//
//	type Pair[T, U any] struct {
//		Server   T
//		Behavior U
//	}
type Pair struct {
	Server   Server
	Behavior string
}

type LoadBalancer struct {
	host         string
	port         int
	count        int
	servers      []Server
	retryServers []Server
	// updateCh chan Pair[Server, string]
	updateCh chan Pair
	readCh   chan chan Server // 接受一個調用者提供 channel 用來將數值回傳給調用者
}

func newLoadBalancer(host string, port int, servers []Server) *LoadBalancer {
	if len(servers) < 1 {
		panic("no server to proxy")
	}
	lb := &LoadBalancer{
		host:     host,
		port:     port,
		count:    0,
		servers:  servers,
		updateCh: make(chan Pair),
		readCh:   make(chan chan Server),
		// addSerCh: make(chan Server),
	}
	go func() {
		for {
			select {
			case p := <-lb.updateCh:
				if p.Behavior == "del" {
					slog.Info("del addr", "addr", p.Server.Address())
					lb.removeAddr(p.Server)
					slog.Info("available servers now", "servers", lb.servers)
				} else if p.Behavior == "add" {
					lb.addAddrHelper(p.Server)
					slog.Info("retry servers", "retry", slog.Any("retry", lb.retryServers))
				}

			case respCh := <-lb.readCh:
				respCh <- lb.findNextService()

			}
		}
	}()
	return lb
}

func (lb *LoadBalancer) findNextService() Server {
	if len(lb.servers) == 0 {
		return nil
	}
	server := lb.servers[lb.count%len(lb.servers)]
	lb.count = (lb.count + 1) % len(lb.servers)
	return server
}

func (lb *LoadBalancer) Serve(rw http.ResponseWriter, req *http.Request) {
	slog.Info("get a request", "from", req.RemoteAddr)
	server := lb.findNextService()
	if server == nil {
		return
	}
	server.Serve(rw, req)
}

func (lb *LoadBalancer) delAddr(s Server) {
	lb.updateCh <- Pair{s, "del"}
}

func (lb *LoadBalancer) GetProxyAddresses() Server {
	respCh := make(chan Server)
	lb.readCh <- respCh
	return <-respCh
}

func (lb *LoadBalancer) removeAddr(server Server) {
	for i, s := range lb.servers {
		if s.Address() == server.Address() {
			// 將該 server 移動至候補清單，晚點再嘗試一次
			lb.retryServers = append(lb.retryServers, s)
			// 找到目標字符串，將其從 slice 中刪除
			lb.servers = append(lb.servers[:i], lb.servers[i+1:]...)
			return
		}
	}

}

func (lb *LoadBalancer) addAddr(s Server) {
	lb.updateCh <- Pair{s, "add"}
}

func (lb *LoadBalancer) addAddrHelper(s Server) {
	for i, rs := range lb.retryServers {
		if s.Address() == rs.Address() {
			lb.servers = append(lb.servers, s)
			lb.retryServers = append(lb.retryServers[:i], lb.retryServers[i+1:]...)
			return
		}
	}
}
