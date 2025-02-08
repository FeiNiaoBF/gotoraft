package handler

import (
	"gotoraft/internal/raft"
	"net/http"

	"github.com/gin-gonic/gin"
)

var cluster = &raft.Cluster{
	Nodes: []raft.Node{
		{ID: "node1", State: raft.Follower, CurrentTerm: 1},
		{ID: "node2", State: raft.Follower, CurrentTerm: 1},
		{ID: "node3", State: raft.Leader, CurrentTerm: 1},
		{ID: "node4", State: raft.Follower, CurrentTerm: 1},
		{ID: "node5", State: raft.Follower, CurrentTerm: 1},
	},
	CurrentTerm: 1,
	LeaderID:    "node3",
}

// GetClusterState returns the current state of the Raft cluster
func GetClusterState(c *gin.Context) {
	c.JSON(http.StatusOK, cluster)
}

// UpdateNodeState updates the state of a specific node
func UpdateNodeState(c *gin.Context) {
	nodeID := c.Param("id")
	var node raft.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := range cluster.Nodes {
		if cluster.Nodes[i].ID == nodeID {
			cluster.Nodes[i].State = node.State
			c.JSON(http.StatusOK, cluster.Nodes[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Node not found"})
}

// AppendLogEntry adds a new log entry to the Raft log
func AppendLogEntry(c *gin.Context) {
	var entry raft.LogEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry.Index = uint64(len(cluster.Log))
	entry.Term = cluster.CurrentTerm
	cluster.Log = append(cluster.Log, entry)

	c.JSON(http.StatusOK, entry)
}

// GetLog returns the current Raft log
func GetLog(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"log": cluster.Log})
}

// StartElection initiates a leader election
func StartElection(c *gin.Context) {
	nodeID := c.Param("id")
	
	// Simple election logic for demonstration
	for i := range cluster.Nodes {
		if cluster.Nodes[i].ID == nodeID {
			cluster.Nodes[i].State = raft.Candidate
			cluster.CurrentTerm++
			cluster.Nodes[i].CurrentTerm = cluster.CurrentTerm
			
			// Simulate winning the election
			cluster.Nodes[i].State = raft.Leader
			cluster.LeaderID = nodeID
			
			// Update other nodes
			for j := range cluster.Nodes {
				if j != i {
					cluster.Nodes[j].State = raft.Follower
					cluster.Nodes[j].CurrentTerm = cluster.CurrentTerm
				}
			}
			
			c.JSON(http.StatusOK, gin.H{
				"message": "Election completed",
				"leader":  nodeID,
				"term":    cluster.CurrentTerm,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Node not found"})
}
