package raft

import (
	"gotoraft/internal/foorpc"
	"log"
	"net"
	"sync"
)

type Raft struct {
	mu          sync.Mutex
	state       string // "leader", "follower", "candidate"
	currentTerm int
	votedFor    string
	log         []LogEntry
	commitIndex int
	lastApplied int
	peers       []string
	me          string
	rpcClient   *foorpc.Client // RPC 客户端
}

type LogEntry struct {
	Term    int
	Command interface{}
}

// NewRaft 创建一个新的 Raft 实例
func NewRaft(peers []string, me string, rpcClient *foorpc.Client) *Raft {
	return &Raft{
		state:       "follower",
		currentTerm: 0,
		votedFor:    "",
		log:         []LogEntry{},
		commitIndex: 0,
		lastApplied: 0,
		peers:       peers,
		me:          me,
		rpcClient:   rpcClient,
	}
}

// StartElection 启动选举
func (r *Raft) StartElection() {
	r.mu.Lock()
	defer r.mu.Unlock()
	// 选举逻辑
}

// AppendEntries 追加日志条目
func (r *Raft) AppendEntries(term int, leaderId string, entries []LogEntry) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// 处理追加日志的逻辑
}

// RequestVote 请求投票
func (r *Raft) RequestVote(term int, candidateId string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// 处理投票请求的逻辑
}

// StartRPCServer 启动 RPC 服务器
func (r *Raft) StartRPCServer() {
	listener, err := net.Listen("tcp", ":port") // 替换为实际端口
	if err != nil {
		log.Fatalf("failed to start RPC server: %v", err)
	}
	go foorpc.Accept(listener)
}

// 其他 Raft 方法...
