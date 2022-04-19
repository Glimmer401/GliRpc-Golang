package main

import (
	"erpc/server"
	"erpc/client"
	"fmt"
	"log"
	"net"
	"time"
	"sync"
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
	// start a server
	// get addr back by channel
	addr := make(chan string)
	go startServer(addr)
	client, _ := client.Dial("tcp", <-addr)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)

	// start a client sending 5 request
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("erpc req %d", i)
			var reply string
			if err := client.Call("Gli.Add", args, &reply); err != nil {
				log.Fatal("call Gli.Add error:", err)
			}
			log.Println("reply:", reply)
		}(i)
	}
	wg.Wait()
}