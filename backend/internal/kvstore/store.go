// 用于模拟键值存储
package kvstore

import (
	"gotoraft/internal/raft"
	"sync"
)

const ()

// cmd 是用于模拟键值存储的命令
type cmd struct {
	Op    string `json:"op,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// Store 是用于模拟键值存储的存储
type Store struct {
	storeMu sync.Mutex
	// 是否使用内存存储（登录后）
	inmem bool
	// Raft实例
	raft *raft.Raft
	// 数据
	data map[string]string
}
