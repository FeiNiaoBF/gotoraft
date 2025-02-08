package server

import (
	"encoding/json"
	"net/http"

	"github.com/FeiNiaoBF/gotoraft/backend/internal/raft"
	"github.com/FeiNiaoBF/gotoraft/backend/internal/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// Server represents the HTTP server that handles Raft visualization requests
type Server struct {
	hub  *websocket.Hub
	raft *raft.RaftNode
}

// NewServer creates a new HTTP server instance
func NewServer(raftNode *raft.RaftNode) *Server {
	return &Server{
		hub:  websocket.NewHub(),
		raft: raftNode,
	}
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	go s.hub.Run()

	r := gin.Default()
	
	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// WebSocket endpoint
	r.GET("/ws", func(c *gin.Context) {
		s.handleWebSocket(c.Writer, c.Request)
	})

	// API endpoints
	api := r.Group("/api")
	{
		api.GET("/status", s.getStatus)
		api.POST("/command", s.submitCommand)
		api.GET("/logs", s.getLogs)
	}

	return r.Run(addr)
}

// handleWebSocket handles WebSocket connections
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &websocket.Client{
		conn: conn,
		send: make(chan []byte, 256),
	}

	s.hub.Register <- client

	go client.WritePump()
	go client.ReadPump(s.hub)
}

// getStatus returns the current status of the Raft cluster
func (s *Server) getStatus(c *gin.Context) {
	// TODO: Implement status retrieval from Raft node
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// submitCommand submits a new command to the Raft cluster
func (s *Server) submitCommand(c *gin.Context) {
	var cmd struct {
		Command string `json:"command"`
	}

	if err := c.BindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Submit command to Raft node
	c.JSON(http.StatusOK, gin.H{
		"status": "accepted",
	})
}

// getLogs returns the Raft log entries
func (s *Server) getLogs(c *gin.Context) {
	// TODO: Implement log retrieval from Raft node
	c.JSON(http.StatusOK, gin.H{
		"logs": []interface{}{},
	})
}
