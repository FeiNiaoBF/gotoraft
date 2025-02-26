package raft

import (
	"testing"
	"time"
)

// mockTransport 是一个用于测试的传输层实现
type mockTransport struct {
	sendCh   chan RPC
	recvCh   chan RPC
	closeCh  chan struct{}
	isClosed bool
}

func newMockTransport() *mockTransport {
	return &mockTransport{
		sendCh:   make(chan RPC, 1),
		recvCh:   make(chan RPC, 1),
		closeCh:  make(chan struct{}),
		isClosed: false,
	}
}

func (t *mockTransport) Send(target string, rpc RPC) error {
	select {
	case <-t.closeCh:
		return ErrTransportShutdown
	case t.sendCh <- rpc:
		return nil
	}
}

func (t *mockTransport) Consumer() <-chan RPC {
	return t.recvCh
}

func (t *mockTransport) Close() error {
	if !t.isClosed {
		close(t.closeCh)
		t.isClosed = true
	}
	return nil
}

func (t *mockTransport) LocalAddr() string {
	return "mock-transport"
}

func (t *mockTransport) IsShutdown() bool {
	return t.isClosed
}

// TestNode_NewNode 测试节点创建
func TestNode_NewNode(t *testing.T) {
	config := &Config{
		NodeID:           "node1",
		PeerIDs:          []string{"node2", "node3"},
		HeartbeatTimeout: 50 * time.Millisecond,
		ElectionTimeout:  150 * time.Millisecond,
	}

	storage := NewMemoryStorage()
	fsm := NewMemoryFSM()
	transport := newMockTransport()

	node, err := NewNode(config, fsm, storage, transport)
	if err != nil {
		t.Fatalf("创建节点失败: %v", err)
	}

	if node.id != "node1" {
		t.Errorf("节点 ID 不匹配，期望 'node1'，实际为 %s", node.id)
	}

	if len(node.peers) != 2 {
		t.Errorf("对等节点数量不匹配，期望 2，实际为 %d", len(node.peers))
	}

	if node.state != Follower {
		t.Errorf("初始状态不匹配，期望 Follower，实际为 %v", node.state)
	}
}

// TestNode_StartStop 测试节点启动和停止
func TestNode_StartStop(t *testing.T) {
	config := &Config{
		NodeID:           "node1",
		PeerIDs:          []string{"node2", "node3"},
		HeartbeatTimeout: 50 * time.Millisecond,
		ElectionTimeout:  150 * time.Millisecond,
	}

	storage := NewMemoryStorage()
	fsm := NewMemoryFSM()
	transport := newMockTransport()

	node, _ := NewNode(config, fsm, storage, transport)

	// 测试启动
	if err := node.Start(); err != nil {
		t.Fatalf("启动节点失败: %v", err)
	}

	// 等待一段时间，确保节点正常运行
	time.Sleep(100 * time.Millisecond)

	// 测试停止
	if err := node.Stop(); err != nil {
		t.Fatalf("停止节点失败: %v", err)
	}
}

// TestNode_ElectionTimeout 测试选举超时
func TestNode_ElectionTimeout(t *testing.T) {
	config := &Config{
		NodeID:           "node1",
		PeerIDs:          []string{"node2", "node3"},
		HeartbeatTimeout: 50 * time.Millisecond,
		ElectionTimeout:  150 * time.Millisecond,
	}

	storage := NewMemoryStorage()
	fsm := NewMemoryFSM()
	transport := newMockTransport()

	node, _ := NewNode(config, fsm, storage, transport)
	node.Start()
	defer node.Stop()

	// 等待选举超时
	time.Sleep(200 * time.Millisecond)

	// 验证节点是否变为候选人
	if node.state != Candidate {
		t.Errorf("节点应该变为候选人，实际状态为 %v", node.state)
	}

	// 验证任期是否增加
	if node.currentTerm != 1 {
		t.Errorf("任期应该为 1，实际为 %d", node.currentTerm)
	}
}

// TestNode_HandleVoteRequest 测试处理投票请求
func TestNode_HandleVoteRequest(t *testing.T) {
	config := &Config{
		NodeID:           "node1",
		PeerIDs:          []string{"node2", "node3"},
		HeartbeatTimeout: 50 * time.Millisecond,
		ElectionTimeout:  150 * time.Millisecond,
	}

	storage := NewMemoryStorage()
	fsm := NewMemoryFSM()
	transport := newMockTransport()

	node, _ := NewNode(config, fsm, storage, transport)
	node.Start()
	defer node.Stop()

	// 创建投票请求
	args := &RequestVoteArgs{
		Term:         1,
		CandidateID:  "node2",
		LastLogIndex: 0,
		LastLogTerm:  0,
	}

	// 处理投票请求
	reply := node.handleVoteRequest(args)

	// 验证响应
	if !reply.VoteGranted {
		t.Error("应该授予投票")
	}
	if reply.Term != 1 {
		t.Errorf("任期不匹配，期望 1，实际为 %d", reply.Term)
	}

	// 验证节点状态
	if node.currentTerm != 1 {
		t.Errorf("节点任期不匹配，期望 1，实际为 %d", node.currentTerm)
	}
	if node.votedFor != "node2" {
		t.Errorf("投票记录不匹配，期望 'node2'，实际为 %s", node.votedFor)
	}
}

// TestNode_HandleAppendEntries 测试处理追加日志请求
func TestNode_HandleAppendEntries(t *testing.T) {
	config := &Config{
		NodeID:           "node1",
		PeerIDs:          []string{"node2", "node3"},
		HeartbeatTimeout: 50 * time.Millisecond,
		ElectionTimeout:  150 * time.Millisecond,
	}

	storage := NewMemoryStorage()
	fsm := NewMemoryFSM()
	transport := newMockTransport()

	node, _ := NewNode(config, fsm, storage, transport)
	node.Start()
	defer node.Stop()

	// 创建追加日志请求（心跳）
	args := &AppendEntriesArgs{
		Term:         1,
		LeaderID:     "node2",
		PrevLogIndex: 0,
		PrevLogTerm:  0,
		Entries:      nil,
		LeaderCommit: 0,
	}

	// 处理追加日志请求
	reply := node.handleAppendEntries(args)

	// 验证响应
	if !reply.Success {
		t.Error("应该接受心跳请求")
	}
	if reply.Term != 1 {
		t.Errorf("任期不匹配，期望 1，实际为 %d", reply.Term)
	}

	// 验证节点状态
	if node.state != Follower {
		t.Errorf("节点应该是追随者，实际为 %v", node.state)
	}
	if node.currentTerm != 1 {
		t.Errorf("节点任期不匹配，期望 1，实际为 %d", node.currentTerm)
	}
}

// TestNode_Submit 测试提交命令
func TestNode_Submit(t *testing.T) {
	config := &Config{
		NodeID:           "node1",
		PeerIDs:          []string{"node2", "node3"},
		HeartbeatTimeout: 50 * time.Millisecond,
		ElectionTimeout:  150 * time.Millisecond,
	}

	storage := NewMemoryStorage()
	fsm := NewMemoryFSM()
	transport := newMockTransport()

	node, _ := NewNode(config, fsm, storage, transport)
	node.Start()
	defer node.Stop()

	// 设置节点为领导者
	node.state = Leader
	node.currentTerm = 1

	// 提交命令
	command := &Command{
		Type:  CommandTypeSet,
		Key:   "name",
		Value: "gotoraft",
	}
	index, err := node.Submit(command)
	if err != nil {
		t.Fatalf("提交命令失败: %v", err)
	}

	// 验证日志
	if index != 1 {
		t.Errorf("日志索引不匹配，期望 1，实际为 %d", index)
	}

	if len(node.logs) != 1 {
		t.Errorf("日志数量不匹配，期望 1，实际为 %d", len(node.logs))
	}

	entry := node.logs[0]
	if entry.Term != 1 {
		t.Errorf("日志任期不匹配，期望 1，实际为 %d", entry.Term)
	}
	if entry.Index != 1 {
		t.Errorf("日志索引不匹配，期望 1，实际为 %d", entry.Index)
	}
}
