package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	"bytes"
	"sync"
	"sync/atomic"
	"time"
)

// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in part 3D you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh, but set CommandValid to false for these
// other uses.
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int

	// For 3D:
	SnapshotValid bool
	Snapshot      []byte
	SnapshotTerm  int
	SnapshotIndex int
}

type NodeState int

const (
	Follower NodeState = iota
	Candidate
	Leader
)

type LogEntry struct {
	Term    int
	Command interface{}
}

const (
	electionTimeoutMin    = 200 * time.Millisecond // 修改为200ms
	electionTimeoutMax    = 400 * time.Millisecond // 修改为400ms
	heartbeatInterval     = 100 * time.Millisecond // 固定心跳间隔 100ms
	electionTimeoutJitter = 50 * time.Millisecond  // 偏移量
)

// A Go object implementing a single Raft peer.
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]
	dead      int32               // set by Kill()

	applyCh        chan ApplyMsg
	applyCond      *sync.Cond
	replicatorCond []*sync.Cond
	state          NodeState

	// all peers
	currentTerm int
	voteFor     int
	log         []*LogEntry

	// all peers
	commitIndex int // 已提交的命令索引
	lastApplied int // 已应用的命令索引

	// leader
	nextIndex  []int // 对于每一台服务器，发送到该服务器的下一个日志条目的索引（初始值为**领导人最后的日志条目的索引+1**）
	matchIndex []int // 对于每一台服务器，已知的已经复制到该服务器的最高日志条目的索引（初始值为0，单调递增）

	//	timeouts
	electionTimer  *time.Timer // 选举计时器
	heartbeatTimer *time.Timer // 心跳计时器

	// 选举相关
	voteForCount int // 投票计数。-1 表示没有投票
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	return rf.currentTerm, rf.state == Leader
}

// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
// before you've implemented snapshots, you should pass nil as the
// second argument to persister.Save().
// after you've implemented snapshots, pass the current snapshot
// (or nil if there's not yet a snapshot).
func (rf *Raft) persist() {
	w := new(bytes.Buffer)
	e := labgob.NewEncoder(w)
	if err := e.Encode(rf.currentTerm); err != nil {
		DPrintf("persist Encode currentTerm error: %v", err)
	}
	if err := e.Encode(rf.voteFor); err != nil {
		DPrintf("persist Encode voteFor error: %v", err)
	}
	if err := e.Encode(rf.log); err != nil {
		DPrintf("persist Encode log error: %v", err)
	}
	if err := e.Encode(rf.commitIndex); err != nil {
		DPrintf("persist Encode commitIndex error: %v", err)
	}
	if err := e.Encode(rf.lastApplied); err != nil {
		DPrintf("persist Encode lastApplied error: %v", err)
	}
	data := w.Bytes()
	rf.persister.Save(data, nil)
}

// restore previously persisted state.
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	r := bytes.NewBuffer(data)
	d := labgob.NewDecoder(r)
	var term int
	var voteFor int
	var log []*LogEntry
	if d.Decode(&term) != nil ||
		d.Decode(&voteFor) != nil ||
		d.Decode(&log) != nil {
		DPrintf("readPersist: decode error")
	} else {
		rf.currentTerm = term
		rf.voteFor = voteFor
		rf.log = log
	}
}

// the service says it has created a snapshot that has
// all info up to and including index. this means the
// service no longer needs the log through (and including)
// that index. Raft should now trim its log as much as possible.
func (rf *Raft) Snapshot(index int, snapshot []byte) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// 如果已经有更新的快照，直接返回
	if index <= rf.lastApplied || index > rf.commitIndex {
		return
	}

	// 保留从 index 开始的日志条目
	newLog := make([]*LogEntry, 0)
	newLog = append(newLog, &LogEntry{}) // 添加一个空的条目在索引0
	newLog = append(newLog, rf.log[index+1:]...)
	rf.log = newLog

	// 更新相关索引
	rf.lastApplied = index
	rf.commitIndex = max(rf.commitIndex, index)

	// 持久化快照和 Raft 状态
	rf.persister.Save(rf.encodeRaftState(), snapshot)
}

func (rf *Raft) encodeRaftState() []byte {
	w := new(bytes.Buffer)
	e := labgob.NewEncoder(w)
	e.Encode(rf.currentTerm)
	e.Encode(rf.voteFor)
	e.Encode(rf.log)
	e.Encode(rf.lastApplied)
	e.Encode(rf.commitIndex)
	return w.Bytes()
}

// Log Operations
// lastLogIndex 返回最后一个日志条目的索引
func (rf *Raft) lastLogIndex() int {
	return len(rf.log) - 1
}

// lastLogTerm 返回最后一个日志条目的任期
func (rf *Raft) lastLogTerm() int {
	return int(rf.log[rf.lastLogIndex()].Term)
}

// RequestVote RPC arguments structure.
// field names must start with capital letters!
type RequestVoteArgs struct {
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

// RequestVote RPC reply structure.
// field names must start with capital letters!
type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

// RequestVote RPC handler.
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// 如果请求的任期小于当前任期，则不投票
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
		return
	}

	// 如果请求的任期大于当前任期，则更新当前任期，并改变状态为Follower
	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.voteFor = -1 // 重置投票状态
		rf.persist()    // 持久化状态
		rf.changeState(Follower)
	}

	// 如果在当前任期已经投票给其他候选人，则拒绝投票
	if rf.voteFor != -1 && rf.voteFor != args.CandidateId {
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
		return
	}

	// 安全性检查：拒绝掉那些日志没有自己新的投票请求
	if !rf.isLogUpToDate(args.LastLogIndex, args.LastLogTerm) {
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
		return
	}

	// 投票给候选人
	reply.Term = rf.currentTerm
	reply.VoteGranted = true
	rf.voteFor = args.CandidateId
	rf.persist() // 持久化状态
	rf.resetElectionTimer()
}

func (rf *Raft) isLogUpToDate(lastLogIndex, lastLogTerm int) bool {
	// Implement logic to check if the candidate's log is up-to-date
	return lastLogTerm > rf.lastLogTerm() ||
		(lastLogTerm == rf.lastLogTerm() && lastLogIndex >= rf.lastLogIndex())

}

// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

// startElection 开始选举
// 1. 增加任期
// 2. 改变状态为Candidate
// 3. 给自己投票
// $. 向其他节点发送RequestVote RPC请求
// !Lock 保护的函数
func (rf *Raft) startElection() {
	// DPrintf("startElection>>>")
	// DPrintf("{Peer %d Term %d}", rf.me, rf.currentTerm)
	// rf.resetElectionTimer()
	rf.changeState(Candidate)

	for peer := range rf.peers {
		if peer != rf.me {
			args := RequestVoteArgs{
				Term:         rf.currentTerm,
				CandidateId:  rf.me,
				LastLogIndex: rf.lastLogIndex(),
				LastLogTerm:  rf.lastLogTerm(),
			}
			// DPrintf("{Args: %+v to peer %d}", args, peer)
			go rf.sendRequestVoteToPeer(peer, args)
		}
	}
}

// handlerRequestAll 处理所有节点的投票请求
func (rf *Raft) sendRequestVoteToPeer(peer int, args RequestVoteArgs) {
	var reply RequestVoteReply
	if ok := rf.sendRequestVote(peer, &args, &reply); ok {
		rf.handleResponseVote(reply)
	}
}

// handleResponseVote 处理投票后的结果
func (rf *Raft) handleResponseVote(reply RequestVoteReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	// DPrintf("handleResponseVote>>>")
	// 如果回复的任期大于当前任期，则更新当前任期，并改变状态为Follower
	if reply.Term > rf.currentTerm {
		rf.currentTerm = reply.Term
		rf.voteFor = -1
		rf.changeState(Follower)
		rf.persist()
		return
	}

	if reply.VoteGranted && rf.state == Candidate {
		rf.voteForCount++
		if rf.voteForCount > len(rf.peers)/2 {
			// DPrintf("{VoteForCount: %d} become leader", rf.voteForCount)
			rf.changeState(Leader)
		}
	}
}

// AppendEntries RPC arguments structure.
type AppendEntriesArgs struct {
	Term         int         //  领导人的任期
	LeaderId     int         //  领导人的ID
	PrevLogIndex int         //  上一个日志条目的索引
	PrevLogTerm  int         //  上一个日志条目的任期
	Entries      []*LogEntry //  需要追加的日志条目, 为空时代表心跳
	LeaderCommit int         //  领导人的提交索引
}

// AppendEntries RPC reply structure.
type AppendEntriesReply struct {
	Term          int
	Success       bool
	LastLogIndex  int // follower 的最后一条日志索引
	ConflictIndex int
	ConflictTerm  int
}

// AppendEntries RPC handler
func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// 如果请求的任期小于当前任期，则不成功
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.Success = false
		return
	}

	// 如果请求的任期大于当前任期，则更新当前任期，并改变状态为Follower
	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.changeState(Follower)
		rf.persist()
	}
	// 重置选举计时器
	rf.resetElectionTimer()

	// 检查日志一致性
	// 1. 如果 prevLogIndex 超出了日志范围
	if args.PrevLogIndex >= len(rf.log) {
		reply.Success = false
		reply.Term = rf.currentTerm
		reply.ConflictIndex = len(rf.log)
		return
	}

	// 2. 如果在 prevLogIndex 位置的任期不匹配
	if args.PrevLogIndex > 0 && rf.log[args.PrevLogIndex].Term != args.PrevLogTerm {
		reply.Success = false
		reply.Term = rf.currentTerm
		reply.ConflictTerm = rf.log[args.PrevLogIndex].Term

		// 找到冲突任期的第一个日志条目
		for i := args.PrevLogIndex - 1; i >= 0; i-- {
			if rf.log[i].Term != reply.ConflictTerm {
				reply.ConflictIndex = i + 1
				break
			}
			if i == 0 {
				reply.ConflictIndex = 0
			}
		}
		return
	}

	// 从 PrevLogIndex+1 开始检查并追加新的日志条目
	i := 0
	// 检查已存在的日志与新日志的一致性
	for ; i < len(args.Entries); i++ {
		logIndex := args.PrevLogIndex + 1 + i
		if logIndex < len(rf.log) {
			// 遇到冲突则截断日志，从此位置开始由新日志覆盖
			if rf.log[logIndex].Term != args.Entries[i].Term {
				rf.log = rf.log[:logIndex]
				break
			}
		} else {
			break
		}
	}
	// 将剩余的新日志条目全部追加
	for ; i < len(args.Entries); i++ {
		rf.log = append(rf.log, args.Entries[i])
	}
	rf.persist()

	// 如果 leaderCommit 大于当前 commitIndex，则更新 commitIndex
	if args.LeaderCommit > rf.commitIndex {
		lastNewIndex := len(rf.log) - 1
		if args.LeaderCommit < lastNewIndex {
			rf.commitIndex = args.LeaderCommit
		} else {
			rf.commitIndex = lastNewIndex
		}
		// 异步应用日志：把新提交的日志应用到状态机
		go rf.applyLogs()
	}

	reply.Success = true
	reply.Term = rf.currentTerm
}

// applyLogs 将已提交的日志应用到状态机
func (rf *Raft) applyLogs() {
	for !rf.killed() {
		rf.mu.Lock()

		// 等待新的日志提交
		for rf.lastApplied >= rf.commitIndex {
			rf.applyCond.Wait()
			if rf.killed() {
				rf.mu.Unlock()
				return
			}
		}

		// 严格按顺序应用日志
		if rf.lastApplied < rf.commitIndex && rf.lastApplied+1 < len(rf.log) {
			// 记录当前要应用的日志
			nextIndex := rf.lastApplied + 1
			msg := ApplyMsg{
				CommandValid: true,
				Command:      rf.log[nextIndex].Command,
				CommandIndex: nextIndex,
			}

			// 先解锁再发送，避免死锁
			rf.mu.Unlock()

			// 发送日志到应用通道
			rf.applyCh <- msg

			// 重新获取锁并更新 lastApplied
			rf.mu.Lock()
			// 确保没有其他 goroutine 修改了 lastApplied
			if rf.lastApplied == nextIndex-1 {
				rf.lastApplied = nextIndex
			}
			rf.mu.Unlock()
		} else {
			rf.mu.Unlock()
		}

		// 如果还有更多日志需要应用，继续发送信号
		rf.mu.Lock()
		if rf.lastApplied < rf.commitIndex {
			rf.applyCond.Signal()
		}
		rf.mu.Unlock()
	}
}

// 发送AppendEntries RPC请求到指定节点
func (rf *Raft) sendAppendEntries(peer int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
	return rf.peers[peer].Call("Raft.AppendEntries", args, reply)
}

// broadcastAppendEntries 向所有节点发送 AppendEntries（心跳或日志）
// !Lock 保护的函数
func (rf *Raft) broadcastAppendEntries(heartbeat bool) {
	for peer := range rf.peers {
		if peer == rf.me {
			continue
		}
		// If it's a heartbeat, or if nextIndex is behind, send AppendEntries
		if heartbeat || rf.nextIndex[peer] <= rf.lastLogIndex() {
			go rf.sendAppendEntry(peer)
		}
	}
}

// sendAppendEntry 向指定 peer 发送 AppendEntries RPC
func (rf *Raft) sendAppendEntry(peer int) {
	rf.mu.Lock()
	if rf.state != Leader {
		rf.mu.Unlock()
		return
	}

	prevLogIndex := rf.nextIndex[peer] - 1
	if prevLogIndex < 0 {
		prevLogIndex = 0
	}

	// 构造 AppendEntriesArgs
	args := AppendEntriesArgs{
		Term:         rf.currentTerm,
		LeaderId:     rf.me,
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  rf.log[prevLogIndex].Term,
		Entries:      make([]*LogEntry, 0),
		LeaderCommit: rf.commitIndex,
	}

	// 如果 nextIndex 小于日志长度，发送从 nextIndex 开始的所有条目
	if rf.nextIndex[peer] <= rf.lastLogIndex() {
		args.Entries = append(args.Entries, rf.log[rf.nextIndex[peer]:]...)
	}

	reply := AppendEntriesReply{}
	rf.mu.Unlock()

	ok := rf.sendAppendEntries(peer, &args, &reply)

	rf.mu.Lock()
	defer rf.mu.Unlock()
	if !ok || rf.state != Leader || rf.currentTerm != args.Term {
		return
	}

	if reply.Success {
		// 更新 matchIndex 和 nextIndex
		rf.matchIndex[peer] = args.PrevLogIndex + len(args.Entries)
		rf.nextIndex[peer] = rf.matchIndex[peer] + 1

		// 检查是否可以提交新的日志
		for N := rf.lastLogIndex(); N > rf.commitIndex; N-- {
			count := 1 // Leader 自身已包含该日志
			for p := range rf.peers {
				if p != rf.me && rf.matchIndex[p] >= N {
					count++
				}
			}
			if count > len(rf.peers)/2 && rf.log[N].Term == rf.currentTerm {
				rf.commitIndex = N
				go rf.applyLogs()
				break
			}
		}
	} else {
		// 处理冲突：根据 ConflictTerm 和 ConflictIndex 回退 nextIndex
		if reply.Term > rf.currentTerm {
			rf.currentTerm = reply.Term
			rf.voteFor = -1
			rf.changeState(Follower)
			rf.persist()
			return
		}
		if reply.ConflictTerm != 0 {
			// Search backwards in log for the last index with ConflictTerm
			idx := -1
			for i := len(rf.log) - 1; i >= 0; i-- {
				if rf.log[i].Term == reply.ConflictTerm {
					idx = i
					break
				}
			}
			if idx != -1 {
				rf.nextIndex[peer] = idx + 1
			} else {
				rf.nextIndex[peer] = reply.ConflictIndex
			}
			// 立即重试发送（避免等待心跳间隔）
			go rf.sendAppendEntry(peer)
		} else {
			rf.nextIndex[peer] = reply.ConflictIndex
		}
		// 立即重试发送（避免等待心跳间隔）
		go rf.sendAppendEntry(peer)
	}
}

// changeState raft 改变状态
// 改变状态时，需要更新选举计时器
// !Lock 保护的函数
func (rf *Raft) changeState(state NodeState) {
	// DPrintf("changeState>>>")
	// DPrintf("{Peer %d Term %d} change state to %v", rf.me, rf.currentTerm, state)
	switch state {
	case Follower:
		// DPrintf("{Peer %d to Follower in term %d}", rf.me, rf.currentTerm)
		rf.resetElectionTimer()
		// rf.resetHeartbeatTimer()
		rf.state = state
		rf.voteFor = -1
		rf.voteForCount = 0
	case Candidate:
		// DPrintf("{Peer %d to Candidate in term %d}", rf.me, rf.currentTerm)
		rf.resetElectionTimer()
		// rf.resetHeartbeatTimer()
		rf.state = state
		rf.currentTerm++
		rf.voteFor = rf.me
		rf.voteForCount = 1
	case Leader:
		// DPrintf("{Peer %d become leader in term %d}", rf.me, rf.currentTerm)
		rf.resetElectionTimer()
		rf.resetHeartbeatTimer()
		rf.state = state
		rf.voteFor = -1
		rf.voteForCount = 0
		// 初始化 nextIndex 和 matchIndex
		lastLogIndex := rf.lastLogIndex()
		for i := range rf.peers {
			rf.nextIndex[i] = lastLogIndex + 1
			rf.matchIndex[i] = 0
		}
		// 立即发送心跳
		go func() {
			rf.mu.Lock()
			rf.broadcastAppendEntries(true)
			rf.mu.Unlock()
		}()
	}
	// 持久化状态
	rf.persist()

	// DPrintf("{Peer %d} change state to %v", rf.me, state)
}

// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
// Start函数根据入参构造LogEntry，添加到自己的log数组里，
// 并发送AppendEntry给Follower，进行一次agreement。
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	index, term, isLeader := -1, (rf.currentTerm), (rf.state == Leader)

	if rf.state != Leader {
		return index, term, isLeader
	}

	// 将新条目追加到日志末尾
	index = rf.lastLogIndex() + 1
	newEntry := &LogEntry{
		Term:    term,
		Command: command,
	}
	rf.log = append(rf.log, newEntry)
	rf.persist()

	// 异步通知其他节点尽快复制日志
	go func() {
		rf.mu.Lock()
		defer rf.mu.Unlock()
		rf.broadcastAppendEntries(false)
	}()
	return index, term, isLeader
}

// the tester doesn't halt goroutines created by Raft after each test,
// but it does call the Kill() method. your code can use killed() to
// check whether Kill() has been called. the use of atomic avoids the
// need for a lock.
//
// the issue is that long-running goroutines use memory and may chew
// up CPU time, perhaps causing later tests to fail and generating
// confusing debug output. any goroutine with a long-running loop
// should call killed() to check whether it should stop.
func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
	// Your code here, if desired.
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}

// ticker 负责选举和心跳的逻辑
func (rf *Raft) ticker() {
	for !rf.killed() {
		rf.mu.Lock()
		state := rf.state
		rf.mu.Unlock()
		select {
		case <-rf.heartbeatTimer.C:
			if state == Leader {
				rf.mu.Lock()
				rf.broadcastAppendEntries(true)
				rf.resetHeartbeatTimer()
				rf.mu.Unlock()
			}
		case <-rf.electionTimer.C:
			if state != Leader {
				rf.mu.Lock()
				rf.startElection()
				rf.resetElectionTimer()
				rf.mu.Unlock()
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg,
) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me
	rf.dead = 0
	rf.state = Follower           // 初始状态为Follower
	rf.currentTerm = 0            // 初始化当前任期为0
	rf.voteFor = -1               // 初始化投票给-1，表示没有投票
	rf.log = make([]*LogEntry, 1) // 初始化日志为1
	// 初始化日志, 从0开始，但是log索引为1
	for i := range rf.log {
		rf.log[i] = &LogEntry{Term: 0, Command: nil}
	}
	rf.nextIndex = make([]int, len(peers))  // 初始化nextIndex
	rf.matchIndex = make([]int, len(peers)) // 初始化matchIndex

	rf.applyCh = applyCh
	// 初始化条件变量, applyCond用于通知ApplyMsg的提交
	rf.applyCond = sync.NewCond(&rf.mu)
	// 初始化replicatorCond, 用于通知AppendEntries的提交
	rf.replicatorCond = make([]*sync.Cond, len(peers))
	for i := range rf.replicatorCond {
		rf.replicatorCond[i] = sync.NewCond(&rf.mu)
	}

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())
	// 修改：使用固定的心跳定时器，而非随机超时
	rf.heartbeatTimer = time.NewTimer(heartbeatInterval)
	// 选举计时器保持随机化（加上偏移量）
	rf.electionTimer = time.NewTimer(randomElectionTimeout(electionTimeoutMin, electionTimeoutMax, electionTimeoutJitter))
	// 启动 ticker goroutine 进行选举监控及心跳发送
	go rf.ticker()

	return rf
}

// 重置选举计时器
func (rf *Raft) resetElectionTimer() {
	rf.electionTimer.Reset(randomElectionTimeout(electionTimeoutMin, electionTimeoutMax, electionTimeoutJitter))
}

// 重置心跳计时器
func (rf *Raft) resetHeartbeatTimer() {
	rf.heartbeatTimer.Reset(heartbeatInterval)
}
