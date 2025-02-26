package raft

import (
	"math/rand"
	"time"
)

// 选举相关的常量
const (
	minElectionTimeout = 150 * time.Millisecond
	maxElectionTimeout = 300 * time.Millisecond
)

// startElection 开始一轮选举
func (n *Node) startElection() {
	n.mu.Lock()
	if n.state != Candidate {
		n.mu.Unlock()
		return
	}
	
	// 增加任期
	n.currentTerm++
	savedCurrentTerm := n.currentTerm
	n.votedFor = n.id
	
	// 获取最后的日志信息
	lastLogIndex := uint64(len(n.logs))
	lastLogTerm := uint64(0)
	if lastLogIndex > 0 {
		lastLogTerm = n.logs[lastLogIndex-1].Term
	}
	n.mu.Unlock()
	
	// 准备请求投票的参数
	args := &RequestVoteArgs{
		Term:         savedCurrentTerm,
		CandidateID:  n.id,
		LastLogIndex: lastLogIndex,
		LastLogTerm:  lastLogTerm,
	}
	
	// 统计投票
	votesReceived := 1 // 给自己投票
	votesNeeded := len(n.peers)/2 + 1
	
	// 向其他节点请求投票
	for peerID := range n.peers {
		go func(peer string) {
			var reply RequestVoteReply
			rpc := RPC{
				Type: VoteRequest,
				To:   peer,
				Args: args,
			}
			
			if err := n.transport.Send(peer, rpc); err != nil {
				return
			}
			
			// 处理投票响应
			n.mu.Lock()
			defer n.mu.Unlock()
			
			// 如果任期已经改变，忽略响应
			if n.state != Candidate || n.currentTerm != savedCurrentTerm {
				return
			}
			
			// 如果收到更高的任期，转为追随者
			if reply.Term > savedCurrentTerm {
				n.becomeFollower(reply.Term)
				return
			}
			
			// 统计投票
			if reply.VoteGranted {
				votesReceived++
				if votesReceived >= votesNeeded {
					n.becomeLeader()
				}
			}
		}(peerID)
	}
}

// handleVoteRequest 处理投票请求
func (n *Node) handleVoteRequest(args *RequestVoteArgs) *RequestVoteReply {
	n.mu.Lock()
	defer n.mu.Unlock()
	
	reply := &RequestVoteReply{
		Term:        n.currentTerm,
		VoteGranted: false,
	}
	
	// 如果请求中的任期小于当前任期，拒绝投票
	if args.Term < n.currentTerm {
		return reply
	}
	
	// 如果请求中的任期大于当前任期，转为追随者
	if args.Term > n.currentTerm {
		n.becomeFollower(args.Term)
	}
	
	// 如果已经投票给其他人，拒绝投票
	if n.votedFor != "" && n.votedFor != args.CandidateID {
		return reply
	}
	
	// 检查日志是否至少和自己一样新
	lastLogIndex := uint64(len(n.logs))
	lastLogTerm := uint64(0)
	if lastLogIndex > 0 {
		lastLogTerm = n.logs[lastLogIndex-1].Term
	}
	
	if args.LastLogTerm < lastLogTerm || 
	   (args.LastLogTerm == lastLogTerm && args.LastLogIndex < lastLogIndex) {
		return reply
	}
	
	// 投票给候选人
	reply.VoteGranted = true
	n.votedFor = args.CandidateID
	
	// 重置选举定时器
	n.resetElectionTimer()
	
	return reply
}

// resetElectionTimer 重置选举定时器
func (n *Node) resetElectionTimer() {
	if n.electionTimer != nil {
		n.electionTimer.Stop()
	}
	
	// 随机化选举超时时间
	timeout := minElectionTimeout + 
		time.Duration(rand.Int63n(int64(maxElectionTimeout-minElectionTimeout)))
	
	n.electionTimer = time.AfterFunc(timeout, func() {
		n.becomeCandidate()
	})
}

// becomeCandidate 转变为候选人
func (n *Node) becomeCandidate() {
	n.mu.Lock()
	defer n.mu.Unlock()
	
	// 转变为候选人状态
	oldState := n.state
	n.state = Candidate
	n.currentTerm++
	n.votedFor = n.id
	if n.visualizer != nil {
		n.visualizer.OnStateChange(n.id, oldState, n.state, n.currentTerm)
	}
	
	// 开始选举
	go n.startElection()
}

// becomeFollower 转变为追随者
func (n *Node) becomeFollower(term uint64) {
	n.mu.Lock()
	defer n.mu.Unlock()
	
	// 转变为追随者状态
	oldState := n.state
	n.state = Follower
	if n.visualizer != nil {
		n.visualizer.OnStateChange(n.id, oldState, n.state, n.currentTerm)
	}
	
	// 重置选举定时器
	n.resetElectionTimer()
}
