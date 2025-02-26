package store

import (
	"io"

	"github.com/hashicorp/raft"
)

// FSM 实现 Raft 的状态机
type FSM struct{}

// Apply 应用状态变化
func (f *FSM) Apply(log *raft.Log) interface{} {
	// 处理日志条目
	// 例如，更新数据存储
	return nil
}

// Snapshot 创建快照
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	// 返回当前状态快照
	return nil, nil
}

// Restore 从快照恢复状态
func (f *FSM) Restore(rc io.ReadCloser) error {
	// 从快照恢复状态
	return nil
}
