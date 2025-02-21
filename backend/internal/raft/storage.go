package raft

import (
	"sync"
)

// MemoryStorage 是一个基于内存的存储实现
type MemoryStorage struct {
	mu          sync.RWMutex
	currentTerm uint64
	votedFor    string
	logs        []*LogEntry
	snapshot    []byte
}

// NewMemoryStorage 创建一个新的内存存储
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		logs: make([]*LogEntry, 0),
	}
}

// GetCurrentTerm 获取当前任期
func (m *MemoryStorage) GetCurrentTerm() (uint64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentTerm, nil
}

// SetCurrentTerm 设置当前任期
func (m *MemoryStorage) SetCurrentTerm(term uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentTerm = term
	return nil
}

// GetVotedFor 获取投票给谁
func (m *MemoryStorage) GetVotedFor() (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.votedFor, nil
}

// SetVotedFor 设置投票给谁
func (m *MemoryStorage) SetVotedFor(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.votedFor = id
	return nil
}

// FirstIndex 获取第一个日志条目的索引
func (m *MemoryStorage) FirstIndex() (uint64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.logs) == 0 {
		return 0, nil
	}
	return m.logs[0].Index, nil
}

// LastIndex 获取最后一个日志条目的索引
func (m *MemoryStorage) LastIndex() (uint64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.logs) == 0 {
		return 0, nil
	}
	return m.logs[len(m.logs)-1].Index, nil
}

// GetLog 获取指定索引的日志条目
func (m *MemoryStorage) GetLog(index uint64) (*LogEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, log := range m.logs {
		if log.Index == index {
			return log, nil
		}
	}
	return nil, ErrLogNotFound
}

// StoreLogs 存储多个日志条目
func (m *MemoryStorage) StoreLogs(entries []*LogEntry) error {
	if len(entries) == 0 {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.logs = append(m.logs, entries...)
	return nil
}

// DeleteRange 删除指定范围的日志条目
func (m *MemoryStorage) DeleteRange(min, max uint64) error {
	if min > max {
		return ErrInvalidLogRange
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	var newLogs []*LogEntry
	for _, log := range m.logs {
		if log.Index < min || log.Index > max {
			newLogs = append(newLogs, log)
		}
	}
	m.logs = newLogs
	return nil
}

// StoreSnapshot 存储快照
func (m *MemoryStorage) StoreSnapshot(snapshot []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.snapshot = make([]byte, len(snapshot))
	copy(m.snapshot, snapshot)
	return nil
}

// GetSnapshot 获取快照
func (m *MemoryStorage) GetSnapshot() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.snapshot) == 0 {
		return nil, ErrNoSnapshot
	}

	snapshot := make([]byte, len(m.snapshot))
	copy(snapshot, m.snapshot)
	return snapshot, nil
}

// 错误定义
var (
	ErrLogNotFound     = NewError("log entry not found")
	ErrInvalidLogRange = NewError("invalid log range")
	ErrNoSnapshot      = NewError("no snapshot available")
)

// Error 是一个自定义错误类型
type Error struct {
	message string
}

// NewError 创建一个新的错误
func NewError(message string) *Error {
	return &Error{message: message}
}

// Error 实现了 error 接口
func (e *Error) Error() string {
	return e.message
}
