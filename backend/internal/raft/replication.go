package raft

import (
	"time"
)

// 日志复制相关的常量
const (
	heartbeatInterval = 50 * time.Millisecond
)

// startReplication 开始日志复制
func (n *Node) startReplication() {
	n.mu.Lock()
	if n.state != Leader {
		n.mu.Unlock()
		return
	}

	// 初始化 nextIndex 和 matchIndex
	for peerID := range n.peers {
		n.nextIndex[peerID] = uint64(len(n.logs))
		n.matchIndex[peerID] = 0
	}
	n.mu.Unlock()

	// 启动心跳定时器
	n.resetHeartbeatTimer()
}

// sendHeartbeat 发送心跳
func (n *Node) sendHeartbeat() {
	n.mu.Lock()
	if n.state != Leader {
		n.mu.Unlock()
		return
	}

	// 获取当前状态
	term := n.currentTerm
	commitIndex := n.commitIndex
	n.mu.Unlock()

	// 向每个节点发送心跳
	for peerID := range n.peers {
		go func(peer string) {
			n.mu.Lock()
			prevLogIndex := n.nextIndex[peer] - 1
			prevLogTerm := uint64(0)
			if prevLogIndex > 0 && int(prevLogIndex) <= len(n.logs) {
				prevLogTerm = n.logs[prevLogIndex-1].Term
			}

			// 准备要发送的日志
			entries := make([]LogEntry, 0)
			if n.nextIndex[peer] < uint64(len(n.logs)) {
				entries = n.logs[n.nextIndex[peer]:]
			}
			n.mu.Unlock()

			// 准备追加日志的参数
			args := &AppendEntriesArgs{
				Term:         term,
				LeaderID:     n.id,
				PrevLogIndex: prevLogIndex,
				PrevLogTerm:  prevLogTerm,
				Entries:      entries,
				LeaderCommit: commitIndex,
			}

			// 发送 RPC
			rpc := RPC{
				Type: AppendEntriesRequest,
				To:   peer,
				Args: args,
			}

			if err := n.transport.Send(peer, rpc); err != nil {
				return
			}
		}(peerID)
	}

	// 重置心跳定时器
	n.resetHeartbeatTimer()
}

// handleAppendEntries 处理追加日志请求
func (n *Node) handleAppendEntries(args *AppendEntriesArgs) *AppendEntriesReply {
	n.mu.Lock()
	defer n.mu.Unlock()

	reply := &AppendEntriesReply{
		Term:    n.currentTerm,
		Success: false,
	}

	// 如果请求中的任期小于当前任期，拒绝请求
	if args.Term < n.currentTerm {
		return reply
	}

	// 如果收到更高的任期，转为追随者
	if args.Term > n.currentTerm {
		n.becomeFollower(args.Term)
	}

	// 重置选举定时器
	n.resetElectionTimer()

	// 日志一致性检查
	if args.PrevLogIndex > 0 {
		if args.PrevLogIndex > uint64(len(n.logs)) ||
			n.logs[args.PrevLogIndex-1].Term != args.PrevLogTerm {
			return reply
		}
	}

	// 追加新日志
	if len(args.Entries) > 0 {
		index := args.PrevLogIndex
		for i, entry := range args.Entries {
			if index+uint64(i) >= uint64(len(n.logs)) {
				// 追加新日志
				n.logs = append(n.logs, entry)
			} else if n.logs[index+uint64(i)].Term != entry.Term {
				// 删除冲突的日志
				n.logs = n.logs[:index+uint64(i)]
				n.logs = append(n.logs, entry)
			}
		}
	}

	// 更新提交索引
	if args.LeaderCommit > n.commitIndex {
		n.commitIndex = min(args.LeaderCommit, uint64(len(n.logs)))
	}

	reply.Success = true
	return reply
}

// becomeLeader 转变为领导者
func (n *Node) becomeLeader() {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.state != Candidate {
		return
	}

	oldState := n.state
	n.state = Leader

	// 初始化领导者状态
	n.nextIndex = make(map[string]uint64)
	n.matchIndex = make(map[string]uint64)

	// 转变为领导者状态
	n.state = Leader
	if n.visualizer != nil {
		n.visualizer.OnStateChange(n.id, oldState, n.state, n.currentTerm)
	}

	// 开始日志复制
	go n.startReplication()
}

// resetHeartbeatTimer 重置心跳定时器
func (n *Node) resetHeartbeatTimer() {
	if n.heartbeatTimer != nil {
		n.heartbeatTimer.Stop()
	}

	n.heartbeatTimer = time.AfterFunc(heartbeatInterval, func() {
		n.sendHeartbeat()
	})
}

// min returns the smaller of x or y.
func min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}
