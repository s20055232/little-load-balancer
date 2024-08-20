package main

import (
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func healthCheck(lb *LoadBalancer) {
	for {
		toAdd := []Server{}
		for _, server := range lb.retryServers {
			addr := server.Address()
			_, err := net.Dial("tcp", addr)
			slog.Info("ping died server", "addr", addr)
			if err != nil {
				slog.Warn("server still dieing.", "addr", addr)
				continue
			}
			toAdd = append(toAdd, server)
		}
		for _, addr := range toAdd {
			lb.addAddr(addr)
		}

		toDels := []Server{}
		for _, server := range lb.servers {
			addr := server.Address()
			_, err := net.Dial("tcp", addr)
			slog.Info("ping addr", "addr", addr)
			if err != nil {
				slog.Warn("ping server failed.", "addr", addr)
				// WARN: 這邊不能直接呼叫刪除，動態變更 lb 的 servers 會導致髒讀
				toDels = append(toDels, server)
			}
		}
		for _, s := range toDels {
			lb.delAddr(s)
		}

		time.Sleep(time.Duration(5) * time.Second)
	}

}

func main() {
	cfg := readTOMLSetting()
	var s []Server
	for _, addr := range cfg.Servers {
		s = append(s, newSimpleServer(addr))
		slog.Info("add a server be proxy", "addr", addr)
	}
	lb := newLoadBalancer(cfg.Host, cfg.Port, s)
	time.Sleep(time.Duration(2) * time.Second)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		lb.Serve(rw, req)
		rw.Write([]byte("Hello world"))
	})
	addr := cfg.Host + ":" + strconv.Itoa(cfg.Port)
	slog.Info("Proxy Server Listening", "addr", addr)
	go healthCheck(lb)
	http.ListenAndServe(addr, mux)
}
