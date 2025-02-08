package simulation

import (
	"sync"
	"time"

	"github.com/FeiNiaoBF/gotoraft/backend/internal/config"
	"github.com/FeiNiaoBF/gotoraft/backend/internal/kvstore"
	"github.com/FeiNiaoBF/gotoraft/backend/internal/raft"
	"github.com/FeiNiaoBF/gotoraft/backend/internal/websocket"
)

// SimulationManager controls the Raft simulation
type SimulationManager struct {
	mu sync.RWMutex

	// Components
	config     *config.Config
	raftNodes  []*raft.RaftNode
	kvStore    *kvstore.Store
	wsHub      *websocket.Hub

	// Simulation state
	running    bool
	speed     float64
	paused    bool
}

// NewSimulationManager creates a new simulation manager
func NewSimulationManager(cfg *config.Config, wsHub *websocket.Hub) *SimulationManager {
	return &SimulationManager{
		config:    cfg,
		wsHub:     wsHub,
		speed:     cfg.Config.Raft.AnimationSpeed,
		kvStore:   kvstore.NewStore(),
	}
}

// Start initializes and starts the simulation
func (sm *SimulationManager) Start() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.running {
		return nil
	}

	// Initialize Raft nodes
	sm.initRaftNodes()

	// Start simulation loop
	go sm.runSimulation()

	sm.running = true
	return nil
}

// Stop stops the simulation
func (sm *SimulationManager) Stop() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if !sm.running {
		return
	}

	// Stop all Raft nodes
	for _, node := range sm.raftNodes {
		node.Stop()
	}

	sm.running = false
}

// SetSpeed sets the simulation speed
func (sm *SimulationManager) SetSpeed(speed float64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.speed = speed
}

// Pause pauses the simulation
func (sm *SimulationManager) Pause() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.paused = true
}

// Resume resumes the simulation
func (sm *SimulationManager) Resume() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.paused = false
}

// initRaftNodes initializes the Raft nodes
func (sm *SimulationManager) initRaftNodes() {
	nodeCount := sm.config.Config.Raft.NodeCount
	sm.raftNodes = make([]*raft.RaftNode, nodeCount)

	// Create peer list for each node
	for i := 0; i < nodeCount; i++ {
		peers := make([]string, 0, nodeCount-1)
		for j := 0; j < nodeCount; j++ {
			if j != i {
				peers = append(peers, sm.getNodeID(j))
			}
		}
		sm.raftNodes[i] = raft.NewRaftNode(sm.getNodeID(i), peers)
	}
}

// runSimulation runs the main simulation loop
func (sm *SimulationManager) runSimulation() {
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	for range ticker.C {
		sm.mu.RLock()
		if !sm.running || sm.paused {
			sm.mu.RUnlock()
			continue
		}

		// Collect current state
		state := sm.collectState()
		sm.mu.RUnlock()

		// Broadcast state to WebSocket clients
		sm.wsHub.BroadcastState("raft_state", state)
	}
}

// collectState collects the current state of all nodes
func (sm *SimulationManager) collectState() interface{} {
	nodes := make([]interface{}, len(sm.raftNodes))
	for i, node := range sm.raftNodes {
		nodes[i] = map[string]interface{}{
			"id":    node.ID,
			"state": node.State,
			"term":  node.CurrentTerm,
		}
	}

	return map[string]interface{}{
		"nodes":     nodes,
		"kvstore":   sm.kvStore.GetLogs(),
		"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
	}
}

// getNodeID generates a node ID
func (sm *SimulationManager) getNodeID(index int) string {
	return fmt.Sprintf("node%d", index+1)
}

// SimulateNetworkDelay simulates network delay
func (sm *SimulationManager) SimulateNetworkDelay() time.Duration {
	if !sm.config.Config.Raft.EnableNetworkDelay {
		return 0
	}

	min := sm.config.Config.Raft.NetworkDelayMin
	max := sm.config.Config.Raft.NetworkDelayMax
	delay := min + rand.Intn(max-min+1)
	return time.Duration(delay) * time.Millisecond
}
