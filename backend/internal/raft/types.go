package raft


// Node represents a node in the Raft cluster
type Node struct {
	ID           string    `json:"id"`
	State        NodeState `json:"state"`
	CurrentTerm  uint64    `json:"currentTerm"`
	VotedFor     string    `json:"votedFor"`
	LastLogIndex uint64    `json:"lastLogIndex"`
	LastLogTerm  uint64    `json:"lastLogTerm"`
	CommitIndex  uint64    `json:"commitIndex"`
	LastApplied  uint64    `json:"lastApplied"`
}

// // LogEntry represents a single entry in the Raft log
// type LogEntry struct {
// 	Term    uint64      `json:"term"`
// 	Index   uint64      `json:"index"`
// 	Command interface{} `json:"command"`
// }

// Cluster represents the state of the entire Raft cluster
type Cluster struct {
	Nodes       []Node     `json:"nodes"`
	CurrentTerm uint64     `json:"currentTerm"`
	Log         []LogEntry `json:"log"`
	LeaderID    string     `json:"leaderId"`
}
