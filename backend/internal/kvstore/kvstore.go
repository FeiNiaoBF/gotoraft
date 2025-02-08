package kvstore

import (
	"sync"
)

// Store represents a simple key-value store that simulates the state machine in Raft
type Store struct {
	mu    sync.RWMutex
	data  map[string]string
	logs  []LogEntry
}

// LogEntry represents a single operation in the key-value store
type LogEntry struct {
	Index   int    // Log index
	Term    int    // Term when entry was received by leader
	Command string // Command (e.g., "set key value" or "get key")
	Key     string
	Value   string
}

// NewStore creates a new key-value store
func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
		logs: make([]LogEntry, 0),
	}
}

// Set stores a key-value pair and records the operation in logs
func (s *Store) Set(term int, index int, key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Record the operation in logs
	s.logs = append(s.logs, LogEntry{
		Index:   index,
		Term:    term,
		Command: "set",
		Key:     key,
		Value:   value,
	})

	// Update the key-value store
	s.data[key] = value
	return nil
}

// Get retrieves a value by key
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.data[key]
	return value, exists
}

// GetLogs returns all recorded log entries
func (s *Store) GetLogs() []LogEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy of logs to prevent external modifications
	logsCopy := make([]LogEntry, len(s.logs))
	copy(logsCopy, s.logs)
	return logsCopy
}

// GetLastLogIndex returns the index of the last log entry
func (s *Store) GetLastLogIndex() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.logs) == 0 {
		return 0
	}
	return s.logs[len(s.logs)-1].Index
}

// GetLastLogTerm returns the term of the last log entry
func (s *Store) GetLastLogTerm() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.logs) == 0 {
		return 0
	}
	return s.logs[len(s.logs)-1].Term
}
