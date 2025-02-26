// Package raft 提供Raft节点实现
// Raft节点实现了Raft共识算法，包括选举，投票，日志复制等操作
// 来自hashicorp/raft的实现

package raft

import (
	"fmt"
	"sync"
	"time"
)

// NewNode 创建一个新的 Raft 节点
func NewNode(config *Config, fsm FSM, storage Storage, transport Transport) (*Node, error) {
	if err := ValidateConfig(config); err != nil {
		return nil, err
	}

	node := &Node{
		id:          config.NodeID,
		peers:       make(map[string]bool),
		state:       Follower,
		currentTerm: 0,
		votedFor:    "",
		logs:        make([]LogEntry, 0),
		nextIndex:   make(map[string]uint64),
		matchIndex:  make(map[string]uint64),
		storage:     storage,
		fsm:         fsm,
		transport:   transport,
		shutdownCh:  make(chan struct{}),
	}

	// 初始化 peers
	for _, peerID := range config.PeerIDs {
		if peerID != config.NodeID {
			node.peers[peerID] = true
		}
	}

	return node, nil
}

// Start 启动节点
func (n *Node) Start() error {
	// 从存储中恢复状态
	if term, err := n.storage.GetCurrentTerm(); err == nil {
		n.currentTerm = term
	}
	if votedFor, err := n.storage.GetVotedFor(); err == nil {
		n.votedFor = votedFor
	}
	// 启动选举定时器
	n.resetElectionTimer()

	// 启动主循环
	go n.run()

	return nil
}

// Stop 停止节点
func (n *Node) Stop() error {
	close(n.shutdownCh)

	if n.electionTimer != nil {
		n.electionTimer.Stop()
	}

	if n.heartbeatTimer != nil {
		n.heartbeatTimer.Stop()
	}

	return nil
}

// run 是节点的主循环
func (n *Node) run() {
	for {
		select {
		case <-n.shutdownCh:
			return

		case rpc := <-n.transport.Consumer():
			n.handleRPC(rpc)
		}
	}
}

// handleRPC 处理接收到的 RPC 消息
func (n *Node) handleRPC(rpc RPC) {
	switch rpc.Type {
	case VoteRequest:
		args := rpc.Args.(*RequestVoteArgs)
		reply := n.handleVoteRequest(args)

		// 发送响应
		response := RPC{
			Type: VoteResponse,
			To:   rpc.From,
			Args: reply,
		}
		n.transport.Send(rpc.From, response)

	case AppendEntriesRequest:
		args := rpc.Args.(*AppendEntriesArgs)
		reply := n.handleAppendEntries(args)

		// 发送响应
		response := RPC{
			Type: AppendEntriesResponse,
			To:   rpc.From,
			Args: reply,
		}
		n.transport.Send(rpc.From, response)
	}
}

// GetState 返回节点的当前状态
func (n *Node) GetState() (term uint64, isLeader bool) {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.currentTerm, n.state == Leader
}

// Submit 提交一个新的命令到日志中
func (n *Node) Submit(command interface{}) (uint64, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.state != Leader {
		return 0, fmt.Errorf("not leader")
	}

	// 创建新的日志条目
	entry := LogEntry{
		Term:    n.currentTerm,
		Index:   uint64(len(n.logs) + 1),
		Command: command,
	}

	// 追加到本地日志
	n.logs = append(n.logs, entry)

	return entry.Index, nil
}

// RaftPeer 表示一个Raft节点的peer信息
type RaftPeer struct {
	ID        string
	Address   string
	Transport Transport

	// 连接状态
	IsActive bool
	LastSeen time.Time

	// 用于模拟网络延迟和丢包
	LatencyMin time.Duration
	LatencyMax time.Duration
	PacketLoss float64
}

// leaderState 根据RAFT论文中的图2来定义
type leaderState struct {
	leaderMu sync.RWMutex
	// 关于Leader
	leaderID   string
	leaderAddr string
	// leader
	nextIndex  map[string]uint64 // 对于每一台服务器，发送到该服务器的下一个日志条目的索引（初始值为**领导人最后的日志条目的索引+1**）
	matchIndex map[string]uint64 // 对于每一台服务器，已知的已经复制到该服务器的最高日志条目的索引（初始值为0，单调递增）
	// leaderCh 用于向其他节点通知领导者变更
	leaderCh chan bool
}

// voteState 根据RAFT论文中的图2来定义
type voteState struct {
	// 选举相关字段
	voteMu sync.Mutex
	// voteFor 记录当前节点的投票结果
	voteFor string
	// voteCount 记录当前节点的投票次数
	voteCount int // -1表示没有投票
	voteCh    chan bool
}

// ApplyMsg 用于将从日志复制或者从网络传输中获取的应用消息传递给应用层
type ApplyMsg struct {
	// CommandValid 是否有命令有效
	CommandValid bool
	// Command 命令
	Command interface{}
	// CommandIndex 命令索引
	CommandIndex uint64
	// 可视化用
	ProposeTime time.Duration // 命令提出时间
	CommitTime  time.Duration // 提交时间
	AppliedTime time.Duration // 应用时间
	ClientID    string        // 关联客户端
	// 对于快照
	Snapshot Snapshot
}
