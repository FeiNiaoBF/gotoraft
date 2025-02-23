// Package store provides a simple distributed key-value store. The keys and
// associated values are changed via distributed consensus, meaning that the
// values are changed only when a majority of nodes in the cluster agree on
// the new value.
//
// Distributed consensus is provided via the Raft algorithm, specifically the
// Hashicorp implementation.
package store

import (
	"encoding/json"
	"fmt"
	"gotoraft/config"
	"gotoraft/pkg/logger"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"
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
	RaftDir  string
	RaftBind string
	nodeID   string
	inmem    bool

	mu sync.Mutex
	m  map[string]string // The key-value store for the system.

	raft         *raft.Raft // The consensus mechanism
	raftConfig   *config.RaftConfig
	logStore     raft.LogStore
	appliedIndex uint64
}

// New returns a new Store.
func New(inmem bool) *Store {
	return &Store{
		m:     make(map[string]string),
		inmem: inmem,
	}
}

// Open opens the store. If enableSingle is set, and there are no existing peers,
// then this node becomes the first node, and therefore leader, of the cluster.
// localID should be the server identifier for this node.
func (s *Store) Open(enableSingle bool, localID string) error {
	s.nodeID = localID

	// 创建或加载Raft配置
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(localID)

	// 创建Raft数据存储目录
	if err := os.MkdirAll(s.RaftDir, 0755); err != nil {
		return fmt.Errorf("failed to create raft directory: %v", err)
	}

	// 创建日志存储
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(s.RaftDir, "raft-log.db"))
	if err != nil {
		return fmt.Errorf("failed to create log store: %v", err)
	}
	s.logStore = logStore

	// 创建稳定存储
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(s.RaftDir, "raft-stable.db"))
	if err != nil {
		return fmt.Errorf("failed to create stable store: %v", err)
	}

	// 创建快照存储
	snapshotStore, err := raft.NewFileSnapshotStore(s.RaftDir, retainSnapshotCount, os.Stderr)
	if err != nil {
		return fmt.Errorf("failed to create snapshot store: %v", err)
	}

	// 设置传输层
	addr, err := net.ResolveTCPAddr("tcp", s.RaftBind)
	if err != nil {
		return fmt.Errorf("failed to resolve tcp address: %v", err)
	}

	transport, err := raft.NewTCPTransport(s.RaftBind, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return fmt.Errorf("failed to create tcp transport: %v", err)
	}

	// 创建Raft实例
	ra, err := raft.NewRaft(config, (*fsm)(s), logStore, stableStore, snapshotStore, transport)
	if err != nil {
		return fmt.Errorf("failed to create raft: %v", err)
	}
	s.raft = ra

	if enableSingle {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		ra.BootstrapCluster(configuration)
	}

	// 开始跟踪已应用的索引
	go s.monitorAppliedIndex()

	return nil
}

// monitorAppliedIndex 监控已应用的日志索引
func (s *Store) monitorAppliedIndex() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if s.raft != nil {
			s.appliedIndex = s.raft.AppliedIndex()
		}
	}
}

// Get returns the value for the given key.
func (s *Store) Get(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.m[key], nil
}

// Set sets the value for the given key.
func (s *Store) Set(key, value string) error {
	if s.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}

	c := &command{
		Op:    "set",
		Key:   key,
		Value: value,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := s.raft.Apply(b, raftTimeout)
	return f.Error()
}

// Delete deletes the given key.
func (s *Store) Delete(key string) error {
	if s.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}

	c := &command{
		Op:  "delete",
		Key: key,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := s.raft.Apply(b, raftTimeout)
	return f.Error()
}

// Join joins a node, identified by nodeID and located at addr, to this store.
// The node must be ready to respond to Raft communications at that address.
func (s *Store) Join(nodeID, addr string) error {
	logger.Infof("received join request for remote node %s at %s", nodeID, addr)

	configFuture := s.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		logger.Errorf("failed to get raft configuration: %v", err)
		return err
	}

	for _, srv := range configFuture.Configuration().Servers {
		// If a node already exists with either the joining node's ID or address,
		// that node may need to be removed from the config first.
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			// However if *both* the ID and the address are the same, then nothing -- not even
			// a join operation -- is needed.
			if srv.Address == raft.ServerAddress(addr) && srv.ID == raft.ServerID(nodeID) {
				logger.Infof("node %s at %s already member of cluster, ignoring join request", nodeID, addr)
				return nil
			}

			future := s.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return fmt.Errorf("error removing existing node %s at %s: %s", nodeID, addr, err)
			}
		}
	}

	f := s.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}
	logger.Infof("node %s at %s joined successfully", nodeID, addr)
	return nil
}

// GetRaft 返回底层的Raft实例
func (s *Store) GetRaft() *raft.Raft {
	return s.raft
}

// GetAppliedIndex 返回已经应用到状态机的最新日志索引
func (s *Store) GetAppliedIndex() uint64 {
	if s.raft == nil {
		return 0
	}
	return s.raft.AppliedIndex()
}

// GetLastLogIndex 返回最新的日志索引
func (s *Store) GetLastLogIndex() uint64 {
	if s.raft == nil {
		return 0
	}
	return s.raft.LastIndex()
}

// GetCurrentTerm 返回当前的任期号
func (s *Store) GetCurrentTerm() uint64 {
	if s.raft == nil {
		return 0
	}
	return s.raft.CurrentTerm()
}

// GetLastLogTerm 返回最新日志的任期号
func (s *Store) GetLastLogTerm() uint64 {
	if s.raft == nil {
		return 0
	}
	lastIndex := s.raft.LastIndex()
	if lastIndex == 0 {
		return 0
	}

	// 获取最新日志任期
	if s.logStore == nil {
		logger.Error("日志存储未初始化")
		return 0
	}

	log := new(raft.Log)
	if err := s.logStore.GetLog(lastIndex, log); err != nil {
		logger.Error("获取最新日志任期失败", "error", err)
		return 0
	}

	return log.Term
}

// GetCommitIndex 返回已提交的最新日志索引
func (s *Store) GetCommitIndex() uint64 {
	if s.raft == nil {
		return 0
	}
	return s.raft.CommitIndex()
}

type fsm Store

// Apply applies a Raft log entry to the key-value store.
func (f *fsm) Apply(l *raft.Log) interface{} {
	var c command
	if err := json.Unmarshal(l.Data, &c); err != nil {
		panic(fmt.Sprintf("failed to unmarshal command: %s", err.Error()))
	}

	switch c.Op {
	case "set":
		return f.applySet(c.Key, c.Value)
	case "delete":
		return f.applyDelete(c.Key)
	default:
		panic(fmt.Sprintf("unrecognized command op: %s", c.Op))
	}
}

// Snapshot returns a snapshot of the key-value store.
func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Clone the map.
	o := make(map[string]string)
	for k, v := range f.m {
		o[k] = v
	}
	return &fsmSnapshot{store: o}, nil
}

// Restore stores the key-value store to a previous state.
func (f *fsm) Restore(rc io.ReadCloser) error {
	o := make(map[string]string)
	if err := json.NewDecoder(rc).Decode(&o); err != nil {
		return err
	}

	// Set the state from the snapshot, no lock required according to
	// Hashicorp docs.
	f.m = o
	return nil
}

func (f *fsm) applySet(key, value string) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.m[key] = value
	return nil
}

func (f *fsm) applyDelete(key string) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.m, key)
	return nil
}

type fsmSnapshot struct {
	store map[string]string
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		// Encode data.
		b, err := json.Marshal(f.store)
		if err != nil {
			return err
		}

		// Write data to sink.
		if _, err := sink.Write(b); err != nil {
			return err
		}

		// Close the sink.
		return sink.Close()
	}()

	if err != nil {
		sink.Cancel()
	}

	return err
}

func (f *fsmSnapshot) Release() {}
