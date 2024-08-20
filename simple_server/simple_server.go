package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
)

func handleFunc(rw http.ResponseWriter, req *http.Request) {
	fmt.Print("I got a request\n")
}

func handleTCP(con net.Conn) {
	defer con.Close()
	log.Println("Received TCP connection from:", con.RemoteAddr())
}

func main() {
	var portFlag = flag.Int("port", 13241, "port")
	var hostFlag = flag.String("host", "127.0.0.1", "host")
	flag.Parse()

	slog.Info("user's input", "port", *portFlag, "host", *hostFlag)

	addr := *hostFlag + ":" + strconv.Itoa(*portFlag)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", handleFunc)
		slog.Info("Listening", "address", addr)
		err := http.Serve(listener, mux)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}()

	for {
		con, err := listener.Accept()
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		go handleTCP(con)
	}
}
