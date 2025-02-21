package raft

import (
	"encoding/json"
	"sync"
)

// FSMSnapshot 接口定义了状态机快照的方法
type FSMSnapshot interface {
	// Persist 持久化快照到一个 sink
	Persist(sink SnapshotSink) error

	// Release 释放快照资源
	Release()
}

// FSM 接口定义了状态机的方法
type FSM interface {
	// Apply 将日志条目应用到状态机
	Apply(log *LogEntry) interface{}

	// Snapshot 创建状态机的快照
	Snapshot() (FSMSnapshot, error)

	// Restore 从快照恢复状态机
	Restore(snapshot []byte) error
}

// MemoryFSM 是一个基于内存的状态机实现
type MemoryFSM struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewMemoryFSM 创建一个新的内存状态机
func NewMemoryFSM() *MemoryFSM {
	return &MemoryFSM{
		data: make(map[string]string),
	}
}

// Apply 将日志条目应用到状态机
func (m *MemoryFSM) Apply(log *LogEntry) interface{} {
	if log == nil {
		return nil
	}

	cmd, ok := log.Command.(*Command)
	if !ok {
		return nil
	}

	switch cmd.Type {
	case CommandTypeSet:
		m.mu.Lock()
		m.data[cmd.Key] = cmd.Value
		m.mu.Unlock()
		return nil

	case CommandTypeGet:
		m.mu.RLock()
		value, exists := m.data[cmd.Key]
		m.mu.RUnlock()
		if !exists {
			return nil
		}
		return value

	case CommandTypeDelete:
		m.mu.Lock()
		delete(m.data, cmd.Key)
		m.mu.Unlock()
		return nil

	default:
		return nil
	}
}

// Snapshot 创建状态机的快照
func (m *MemoryFSM) Snapshot() (FSMSnapshot, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 创建数据的副本
	data := make(map[string]string)
	for k, v := range m.data {
		data[k] = v
	}

	return &memorySnapshot{data: data}, nil
}

// Restore 从快照恢复状态机
func (m *MemoryFSM) Restore(snapshot []byte) error {
	var data map[string]string
	if err := json.Unmarshal(snapshot, &data); err != nil {
		return err
	}

	m.mu.Lock()
	m.data = data
	m.mu.Unlock()

	return nil
}

// Get 获取键值
func (m *MemoryFSM) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, exists := m.data[key]
	return value, exists
}

// Set 设置键值
func (m *MemoryFSM) Set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

// Delete 删除键值
func (m *MemoryFSM) Delete(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.data[key]; exists {
		delete(m.data, key)
		return true
	}
	return false
}

// Size 返回状态机中的键值对数量
func (m *MemoryFSM) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// MemorySnapshot 是一个基于内存的快照实现
type MemorySnapshot struct {
	data []byte
}

// Persist 持久化快照
func (s *MemorySnapshot) Persist(sink SnapshotSink) error {
	_, err := sink.Write(s.data)
	return err
}

// Release 释放快照资源
func (s *MemorySnapshot) Release() {
	s.data = nil
}

// memorySnapshot 实现了 FSMSnapshot 接口
type memorySnapshot struct {
	data map[string]string
}

// Persist 实现了 FSMSnapshot 接口
func (s *memorySnapshot) Persist(sink SnapshotSink) error {
	data, err := json.Marshal(s.data)
	if err != nil {
		return err
	}

	if _, err := sink.Write(data); err != nil {
		return err
	}

	return sink.Close()
}

// Release 实现了 FSMSnapshot 接口
func (s *memorySnapshot) Release() {
	// 在内存实现中不需要做任何事情
}
