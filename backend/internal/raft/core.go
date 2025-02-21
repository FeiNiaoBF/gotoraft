package raft

import (
	"sync"
	"time"
)

// SnapshotSink 是一个用于写入快照数据的接口
type SnapshotSink interface {
	Write(p []byte) (n int, err error)
	Close() error
	ID() string
	Cancel() error
}

// RaftState 表示节点状态
type RaftState int

const (
	Follower RaftState = iota
	Candidate
	Leader
)

// LogEntry 表示一条日志条目
type LogEntry struct {
	Term    uint64      // 任期号
	Index   uint64      // 索引号
	Command interface{} // 命令内容
}

// Storage 接口定义了存储层的方法
type Storage interface {
	// 持久化状态
	GetCurrentTerm() (uint64, error)
	SetCurrentTerm(term uint64) error
	GetVotedFor() (string, error)
	SetVotedFor(id string) error

	// 日志操作
	FirstIndex() (uint64, error)
	LastIndex() (uint64, error)
	GetLog(index uint64) (*LogEntry, error)
	StoreLogs(entries []*LogEntry) error
	DeleteRange(min, max uint64) error
}

// ServerID 是一个唯一标识服务器的字符串
type ServerID string

// ServerAddress 是服务器的网络地址
type ServerAddress string

// RPCType 表示 RPC 消息类型
type RPCType int

const (
	VoteRequest RPCType = iota
	VoteResponse
	AppendEntriesRequest
	AppendEntriesResponse
)

// RequestVoteArgs 请求投票的参数
type RequestVoteArgs struct {
	Term         uint64 // 候选人的任期号
	CandidateID  string // 请求选票的候选人的 ID
	LastLogIndex uint64 // 候选人的最后日志条目的索引值
	LastLogTerm  uint64 // 候选人最后日志条目的任期号
}

// RequestVoteReply 请求投票的响应
type RequestVoteReply struct {
	Term        uint64 // 当前任期号，以便于候选人去更新自己的任期号
	VoteGranted bool   // true 表示候选人赢得了此张选票
}

// AppendEntriesArgs 追加日志的参数
type AppendEntriesArgs struct {
	Term         uint64     // 领导人的任期号
	LeaderID     string     // 领导人的 ID
	PrevLogIndex uint64     // 新的日志条目紧随之前的索引值
	PrevLogTerm  uint64     // PrevLogIndex 条目的任期号
	Entries      []LogEntry // 准备存储的日志条目（表示心跳时为空）
	LeaderCommit uint64     // 领导人已经提交的日志的索引值
}

// AppendEntriesReply 追加日志的响应
type AppendEntriesReply struct {
	Term    uint64 // 当前的任期号，用于领导人去更新自己
	Success bool   // 如果跟随者包含了匹配上 PrevLogIndex 和 PrevLogTerm 的日志时为真
}

// RPC 是一个通用的 RPC 消息结构
type RPC struct {
	Type RPCType     // RPC 类型
	From string      // 发送者 ID
	To   string      // 接收者 ID
	Args interface{} // RPC 参数
}

// Transport 接口定义了节点间通信的方法
type Transport interface {
	// Send 发送 RPC 请求到指定节点
	Send(target string, rpc RPC) error
	// Consumer 返回接收 RPC 的通道
	Consumer() <-chan RPC
	// Close 关闭传输层
	Close() error
	// LocalAddr 返回本地地址
	LocalAddr() string
	// IsShutdown 检查传输层是否已关闭
	IsShutdown() bool
}

// CommandType 表示命令类型
type CommandType string

const (
	CommandTypeSet    CommandType = "set"
	CommandTypeGet    CommandType = "get"
	CommandTypeDelete CommandType = "delete"
)

// Command 表示状态机命令
type Command struct {
	Type  CommandType `json:"type"`
	Key   string      `json:"key"`
	Value string      `json:"value,omitempty"`
}

// Node 节点实现
// 以论文中的图2为基础
type Node struct {
	// 基本信息
	id    string          // 节点ID
	peers map[string]bool // 集群中的其他节点
	state RaftState       // 当前状态

	// 持久化状态
	currentTerm uint64     // 当前任期
	votedFor    string     // 投票给谁
	logs        []LogEntry // 日志条目

	// 易失性状态
	commitIndex uint64 // 已提交的最高日志索引
	lastApplied uint64 // 已应用到状态机的最高日志索引

	// 领导人易失性状态
	nextIndex  map[string]uint64 // 对于每个服务器，发送到该服务器的下一个日志条目的索引
	matchIndex map[string]uint64 // 对于每个服务器，已知的已复制到该服务器的最高日志条目的索引

	// 组件
	storage   Storage   // 存储层
	fsm       FSM       // 状态机
	transport Transport // 传输层

	// 定时器
	electionTimer  *time.Timer // 选举定时器
	heartbeatTimer *time.Timer // 心跳定时器

	// 通道
	shutdownCh chan struct{} // 关闭通道

	// 可视化
	visualizer *Visualizer // 可视化组件

	// 锁
	mu sync.Mutex
}
