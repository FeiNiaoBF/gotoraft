package raft

import (
	"encoding/json"
	"testing"
)

// TestMemoryFSM_Apply 测试命令应用
func TestMemoryFSM_Apply(t *testing.T) {
	fsm := NewMemoryFSM()

	// 测试设置命令
	entry := &LogEntry{
		Term:  1,
		Index: 1,
		Command: &Command{
			Type:  CommandTypeSet,
			Key:   "name",
			Value: "gotoraft",
		},
	}

	result := fsm.Apply(entry)
	if result != nil {
		t.Errorf("应用命令失败: %v", result)
	}

	// 测试获取命令
	entry = &LogEntry{
		Term:  1,
		Index: 2,
		Command: &Command{
			Type: CommandTypeGet,
			Key:  "name",
		},
	}

	result = fsm.Apply(entry)
	if result != "gotoraft" {
		t.Errorf("获取值不匹配，期望 'gotoraft'，实际为 %v", result)
	}

	// 测试删除命令
	entry = &LogEntry{
		Term:  1,
		Index: 3,
		Command: &Command{
			Type: CommandTypeDelete,
			Key:  "name",
		},
	}

	result = fsm.Apply(entry)
	if result != nil {
		t.Errorf("删除命令失败: %v", result)
	}

	// 验证删除结果
	entry = &LogEntry{
		Term:  1,
		Index: 4,
		Command: &Command{
			Type: CommandTypeGet,
			Key:  "name",
		},
	}

	result = fsm.Apply(entry)
	if result != nil {
		t.Errorf("键应该已被删除，但获取到了值: %v", result)
	}
}

// TestMemoryFSM_Snapshot 测试快照功能
func TestMemoryFSM_Snapshot(t *testing.T) {
	fsm := NewMemoryFSM()

	// 添加一些数据
	entries := []*LogEntry{
		{
			Term:  1,
			Index: 1,
			Command: &Command{
				Type:  CommandTypeSet,
				Key:   "name",
				Value: "gotoraft",
			},
		},
		{
			Term:  1,
			Index: 2,
			Command: &Command{
				Type:  CommandTypeSet,
				Key:   "version",
				Value: "1.0",
			},
		},
	}

	for _, entry := range entries {
		fsm.Apply(entry)
	}

	// 创建快照
	snapshot, err := fsm.Snapshot()
	if err != nil {
		t.Fatalf("创建快照失败: %v", err)
	}

	// 获取快照数据
	snapshotData := snapshot.(*memorySnapshot).data

	// 序列化快照数据
	snapshotBytes, err := json.Marshal(snapshotData)
	if err != nil {
		t.Fatalf("序列化快照数据失败: %v", err)
	}

	// 创建新的状态机
	newFSM := NewMemoryFSM()

	// 从快照恢复
	if err := newFSM.Restore(snapshotBytes); err != nil {
		t.Fatalf("从快照恢复失败: %v", err)
	}

	// 验证恢复的数据
	entry := &LogEntry{
		Term:  1,
		Index: 1,
		Command: &Command{
			Type: CommandTypeGet,
			Key:  "name",
		},
	}
	result := newFSM.Apply(entry)
	if result != "gotoraft" {
		t.Errorf("恢复的值不匹配，期望 'gotoraft'，实际为 %v", result)
	}

	entry = &LogEntry{
		Term:  1,
		Index: 2,
		Command: &Command{
			Type: CommandTypeGet,
			Key:  "version",
		},
	}
	result = newFSM.Apply(entry)
	if result != "1.0" {
		t.Errorf("恢复的值不匹配，期望 '1.0'，实际为 %v", result)
	}
}

// MemorySnapshotSink 是一个用于测试的内存快照 sink
type MemorySnapshotSink struct {
	data []byte
}

func NewMemorySnapshotSink() *MemorySnapshotSink {
	return &MemorySnapshotSink{
		data: make([]byte, 0),
	}
}

func (m *MemorySnapshotSink) Write(p []byte) (n int, err error) {
	m.data = append(m.data, p...)
	return len(p), nil
}

func (m *MemorySnapshotSink) Close() error {
	return nil
}

func (m *MemorySnapshotSink) ID() string {
	return "memory"
}

func (m *MemorySnapshotSink) Cancel() error {
	return nil
}

func (m *MemorySnapshotSink) Data() []byte {
	return m.data
}
