package store

import (
	"gotoraft/config"
	"gotoraft/internal/raft"
	"gotoraft/pkg/logger"
	"sync"
	"time"
)

const (
	retainSnapshotCount = 2
	raftTimeout         = 10 * time.Second
)

type command struct {
	Op    string `json:"op,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// Config 用于存储和管理配置
type Config struct {
	RaftDir  string // Raft 存储目录
	RaftBind string // Raft 绑定地址
	Inmem    bool   // 是否使用内存存储
}

// Store 是一个简单的键值存储，所有更改通过 Raft 共识进行。
type Store struct {
	mu       sync.RWMutex
	data     map[string]string
	raftDir  string
	raftBind string
	inmem    bool       // true 如果存储是内存存储
	raft     *raft.Raft // HashiCorp Raft 实体
}

// GetAppliedIndex 返回当前已应用的日志索引
func (s *Store) GetAppliedIndex() uint64 {
	if s.raft == nil {
		logger.Fatal("raft node is not initialized")
		return 0
	}
	return s.raft.AppliedIndex()
}

// InitStore 初始化kv服务
func InitStore() *Store {
	cfg := config.GetStoreConfig()
	if cfg == nil {
		logger.Fatal("store config is nil")
	}
	return &Store{
		raftDir:  cfg.RaftDir,
		raftBind: cfg.RaftBind,
		inmem:    cfg.Inmem,
	}
}

// 在Store中添加配置更新方法
func (s *Store) ReloadConfig(newConfig *config.StoreConfig) error {
	// 实现配置热更新逻辑
	// 例如更新Raft超时时间等
}

// GetRaft 返回 Raft 节点
func (s *Store) GetRaft() *raft.Raft {
	return s.raft
}

// NewStore 创建一个新的 Store 实例
func NewStore(peers []string, me string) *Store {
	return &Store{
		data: make(map[string]string),
		raft: raft.NewRaft(peers, me),
	}
}

func (s *Store) Set(key, value string) {
	// 使用 Raft 记录日志
	s.data[key] = value
}

func (s *Store) Get(key string) string {
	return s.data[key]
}

func (s *Store) Delete(key string) error {
	// 创建 Raft 日志条目
	return nil
}
