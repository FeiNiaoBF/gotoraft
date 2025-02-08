package main

import (
	"flag"
	"log"

	"github.com/FeiNiaoBF/gotoraft/backend/internal/config"
	"github.com/FeiNiaoBF/gotoraft/backend/internal/raft"
	"github.com/FeiNiaoBF/gotoraft/backend/internal/server"
)

func main() {
	// Command line flags
	nodeID := flag.String("id", "node1", "Node ID")
	addr := flag.String("addr", ":8080", "HTTP server address")
	peers := flag.String("peers", "", "Comma-separated list of peer node IDs")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create Raft node
	peerList := []string{} // TODO: Parse peers from command line
	raftNode := raft.NewRaftNode(*nodeID, peerList)

	// Create and start HTTP server
	srv := server.NewServer(raftNode)
	log.Printf("Starting server on %s\n", *addr)
	if err := srv.Start(*addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
