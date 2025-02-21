// 用于Raft节点的快照和恢复

package raft

// Snapshot 是Raft节点的快照
type Snapshot struct {
	// SnapshotValid 是否有效
	snapshotValid bool
	// Snapshot 指向快照数据
	snapshot []byte
	// SnapshotTerm 快照中的Term
	snapshotTerm uint64
	// SnapshotIndex 快照中的索引
	snapshotIndex uint64
	// UseSnapshot 是否使用快照
	useSnapshot bool
}
