package rpc

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)

	// 使用client
	client, err := Dial("tcp", <-addr)

	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer client.Close()

	time.Sleep(time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("foorpc req %d", i)
			var reply string
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatalf("failed to call Foo.Sum: %v", err)
			}
			log.Printf("reply: %s", reply)
		}(i)
	}
	wg.Wait()
}
