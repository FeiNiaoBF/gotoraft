// Package observer 提供Raft状态的观察和监控功能
package observer

import (
	"encoding/json"
	"fmt"
	"gotoraft/internal/kvstore/store"
	"gotoraft/internal/websocket"
	"gotoraft/pkg/logger"
	"strconv"
	"sync"
	"time"
)

// RaftMetrics 表示Raft的度量指标
type RaftMetrics struct {
	// 基本状态
	State        string `json:"state"`
	Term         uint64 `json:"term"`
	LastLogIndex uint64 `json:"lastLogIndex"`
	LastLogTerm  uint64 `json:"lastLogTerm"`
	CommitIndex  uint64 `json:"commitIndex"`
	AppliedIndex uint64 `json:"appliedIndex"`

	// 性能指标
	Progress float64 `json:"progress"`
	Speed    float64 `json:"speed"`

	// 集群信息
	Leader   string   `json:"leader"`
	VotedFor string   `json:"votedFor"`
	Peers    []string `json:"peers"`

	// 统计信息
	NumLogs     uint64    `json:"numLogs"`
	PendingLogs uint64    `json:"pendingLogs"`
	LastContact time.Time `json:"lastContact"`
}

// RaftStateMessage 表示Raft状态消息
type RaftStateMessage struct {
	Type      string      `json:"type"` // 消息类型
	NodeID    string      `json:"nodeId"`
	Timestamp time.Time   `json:"timestamp"`
	Metrics   RaftMetrics `json:"metrics"`
}

// RaftStateObserver 观察Raft状态的观察器
type RaftStateObserver struct {
	store     *store.Store
	wsManager *websocket.Manager

	// 用于计算速率
	mu               sync.RWMutex
	lastAppliedIndex uint64
	lastUpdateTime   time.Time
}

// NewRaftStateObserver 创建一个新的Raft状态观察器
func NewRaftStateObserver(store *store.Store, wsManager *websocket.Manager) *RaftStateObserver {
	return &RaftStateObserver{
		store:          store,
		wsManager:      wsManager,
		lastUpdateTime: time.Now(),
	}
}

// Start 开始观察Raft状态
func (o *RaftStateObserver) Start() {
	go o.observeState()
}

// observeState 持续观察Raft状态
func (o *RaftStateObserver) observeState() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := o.collectAndBroadcastState(); err != nil {
			logger.Error("收集并广播Raft状态失败", "error", err)
		}
	}
}

// collectMetrics 收集Raft度量指标
func (o *RaftStateObserver) collectMetrics() (*RaftMetrics, error) {
	raftNode := o.store.GetRaft()
	if raftNode == nil {
		return nil, fmt.Errorf("raft node not initialized")
	}

	stats := raftNode.Stats()
	state := raftNode.State()

	// 计算进度和速率
	o.mu.Lock()
	currentAppliedIndex := o.store.GetAppliedIndex()
	currentTime := time.Now()
	timeDiff := currentTime.Sub(o.lastUpdateTime).Seconds()
	speed := float64(currentAppliedIndex-o.lastAppliedIndex) / timeDiff

	// 更新上次的值
	o.lastAppliedIndex = currentAppliedIndex
	o.lastUpdateTime = currentTime
	o.mu.Unlock()

	// 计算总体进度
	lastLogIndex := o.store.GetLastLogIndex()
	var progress float64
	if lastLogIndex > 0 {
		progress = float64(currentAppliedIndex) / float64(lastLogIndex)
	}

	// 解析集群信息
	var peers []string
	if configFuture := raftNode.GetConfiguration(); configFuture.Error() == nil {
		for _, server := range configFuture.Configuration().Servers {
			peers = append(peers, string(server.ID))
		}
	}

	// 构建度量指标
	metrics := &RaftMetrics{
		State:        state.String(),
		Term:         o.store.GetCurrentTerm(),
		LastLogIndex: lastLogIndex,
		LastLogTerm:  o.store.GetLastLogTerm(),
		CommitIndex:  o.store.GetCommitIndex(),
		AppliedIndex: currentAppliedIndex,
		Progress:     progress,
		Speed:        speed,
		Leader:       stats["leader_id"],
		VotedFor:     stats["voted_for"],
		Peers:        peers,
		NumLogs:      lastLogIndex,
		PendingLogs:  lastLogIndex - currentAppliedIndex,
	}

	// 解析最后联系时间
	if lastContact, err := strconv.ParseInt(stats["last_contact"], 10, 64); err == nil {
		metrics.LastContact = time.Unix(0, lastContact)
	}

	return metrics, nil
}

// collectAndBroadcastState 收集并广播Raft状态
func (o *RaftStateObserver) collectAndBroadcastState() error {
	metrics, err := o.collectMetrics()
	if err != nil {
		return err
	}

	message := RaftStateMessage{
		Type:      "raft_state",
		NodeID:    o.nodeID,
		Timestamp: time.Now(),
		Metrics:   *metrics,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	o.wsManager.Broadcast(data)
	return nil
}
