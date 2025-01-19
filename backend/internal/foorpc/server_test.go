package foorpc

import (
	"context"
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

func TestServerAccept(t *testing.T) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	server := NewServer()
	go server.Accept(lis)
}

func startServer(addr chan string) {
	var foo Foo
	if err := Register(&foo); err != nil {
		log.Fatal("register error:", err)
	}
	// pick a free port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	Accept(l)
}

func TestServer(t *testing.T) {
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)

	client, err := Dial("tcp", <-addr)
	if err != nil {
		t.Fatal("dialing error:", err)
	}
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &Args{Num1: i, Num2: i * i}
			var reply int
			if err := client.Call(context.Background(), "Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}
