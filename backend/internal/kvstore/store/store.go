package store

import (
	"gotoraft/config"
	"gotoraft/pkg/logger"
	"sync"
	"time"

	"github.com/hashicorp/raft"
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

// Store is a simple key-value store, where all changes are made via Raft consensus.
type Store struct {
	mu       sync.RWMutex
	data     map[string]string
	raftDir  string
	raftBind string
	inmem    bool       // true if the store is an in-memory store
	raft     *raft.Raft // raft 实体
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
