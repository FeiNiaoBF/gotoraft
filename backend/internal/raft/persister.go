package raft

import (
	"sync"
)

// Persister 用于持久化 Raft 的状态
type Persister struct {
	mu        sync.Mutex
	raftstate []byte
	snapshot  []byte
}

// NewPersister 创建一个新的 Persister
func NewPersister() *Persister {
	return &Persister{}
}

// Save 保存 Raft 的状态和快照
func (ps *Persister) Save(raftstate []byte, snapshot []byte) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.raftstate = raftstate
	ps.snapshot = snapshot
}

// ReadRaftState 读取 Raft 的状态
func (ps *Persister) ReadRaftState() []byte {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return ps.raftstate
}

// ReadSnapshot 读取快照数据
func (ps *Persister) ReadSnapshot() []byte {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return ps.snapshot
}

// RaftStateSize 返回 Raft 状态的大小
func (ps *Persister) RaftStateSize() int {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return len(ps.raftstate)
}

// SnapshotSize 返回快照的大小
func (ps *Persister) SnapshotSize() int {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return len(ps.snapshot)
}

