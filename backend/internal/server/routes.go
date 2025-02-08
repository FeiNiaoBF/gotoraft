package server

import (
	"github.com/gin-gonic/gin"
)

// setupRoutes configures all the routes for the server
func (s *Server) setupRoutes(r *gin.Engine) {
	// API group
	api := r.Group("/api")
	{
		// Raft status and control
		api.GET("/status", s.getStatus)
		api.POST("/command", s.submitCommand)
		api.GET("/logs", s.getLogs)

		// Node management
		nodes := api.Group("/nodes")
		{
			nodes.GET("", s.getNodes)
			nodes.POST("", s.addNode)
			nodes.DELETE("/:id", s.removeNode)
		}

		// Configuration
		config := api.Group("/config")
		{
			config.GET("", s.getConfig)
			config.PUT("", s.updateConfig)
		}
	}

	// WebSocket endpoint
	r.GET("/ws", func(c *gin.Context) {
		s.handleWebSocket(c.Writer, c.Request)
	})
}

// getNodes returns all nodes in the cluster
func (s *Server) getNodes(c *gin.Context) {
	// TODO: Implement getting all nodes from Raft cluster
	c.JSON(200, gin.H{
		"nodes": []interface{}{},
	})
}

// addNode adds a new node to the cluster
func (s *Server) addNode(c *gin.Context) {
	var node struct {
		ID   string `json:"id"`
		Addr string `json:"addr"`
	}

	if err := c.BindJSON(&node); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement adding node to Raft cluster
	c.JSON(200, gin.H{
		"status": "added",
		"node":   node,
	})
}

// removeNode removes a node from the cluster
func (s *Server) removeNode(c *gin.Context) {
	id := c.Param("id")

	// TODO: Implement removing node from Raft cluster
	c.JSON(200, gin.H{
		"status": "removed",
		"id":     id,
	})
}

// getConfig returns the current configuration
func (s *Server) getConfig(c *gin.Context) {
	// TODO: Implement getting configuration
	c.JSON(200, gin.H{
		"config": map[string]interface{}{},
	})
}

// updateConfig updates the configuration
func (s *Server) updateConfig(c *gin.Context) {
	var config struct {
		HeartbeatTimeout  int `json:"heartbeatTimeout"`
		ElectionTimeout   int `json:"electionTimeout"`
		ReplicationFactor int `json:"replicationFactor"`
	}

	if err := c.BindJSON(&config); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement updating configuration
	c.JSON(200, gin.H{
		"status": "updated",
		"config": config,
	})
}
