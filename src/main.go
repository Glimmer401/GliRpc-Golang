package main

import (
	"encoding/json"
	"fmt"
	"erpc/server"
	"erpc/codec"
	"log"
	"net"
	"time"
)

func startServer(addr chan string) {
	server := server.Server{}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("can not create tcp listener:", err)
	}
	log.Println("start erpc server on", listener.Addr())
	addr <- listener.Addr().String()

	server.Accept(listener)
}

func main() {
	// get addr back by channel
	addr := make(chan string)
	go startServer(addr)

	// behaviour similar to clients
	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()
	time.Sleep(time.Second)
	// send 1 options with 5 request
	_ = json.NewEncoder(conn).Encode(server.DefaultOption)
	cc := codec.NewGobCodec(conn)
	// send request & receive response
	for i := 0; i < 5; i++ {
		h := &codec.Header{
			MethodName: "Gli.Add",
			Seq:         uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("erpc req %d", h.Seq))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("reply:", reply)
	}
}