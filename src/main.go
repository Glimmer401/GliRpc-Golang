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
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", listener.Addr())
	addr <- listener.Addr().String()
	server.Accept(listener)
}

func main() {
	// 通过管道回传 addr
	addr := make(chan string)
	go startServer(addr)

	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	time.Sleep(time.Second)
	// send options
	_ = json.NewEncoder(conn).Encode(server.DefaultOption)
	cc := codec.NewGobCodec(conn)
	// send request & receive response
	for i := 0; i < 5; i++ {
		h := &codec.Header{
			MethodName: "Gli.Add",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("erpc req %d", h.Seq))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("reply:", reply)
	}
}