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
	"math"
	"sort"

	//	"bytes"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

)

// ApplyMsg as each Raft peer becomes aware that successive log entries are
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

// node state
type State int32

const (
	Follower State = iota
	Candidate
	Leader
)

// LogEntry is a struct to hold information about each log entry
type LogEntry struct {
	Metadata LogMetadata
	Command  interface{}
}

// LogMetadata 日志条目的元数据 {term, index}
type LogMetadata struct {
	Term  int // 日志条目被创建时的任期号
	Index int // 日志条目的索引 也是所有日志的index eg: in log[n] n = LogEntry.Metadata.Index
}

// Raft A Go object implementing a single Raft peer.
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*RPCEnd           // 修改为我们自己的 RPCEnd
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]
	dead      int32               // set by Kill()

	// Your data here (3A, 3B, 3C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

	// Server state
	// because of memory alignment that chose int32
	state State

	// Persistent state on all servers
	currentTerm int
	votedFor    int
	logs        []LogEntry //TODO: log[] struct
	//logLen      int
	// Volatile state on all servers
	commitIndex int
	lastApplied int

	// Volatile state on leaders
	nextIndex []int // for each server, index of the next log entry to send to that server 可能有更多的log需要发送
	// 因为 Raft 协议中的 Leader 需要对每个 Follower 都保持独立的日志复制进度，而不同的 Follower 在日志进度上可能不一致。
	matchIndex []int // index of highest log entry applied to state machine
	// channel
	applyCh chan ApplyMsg // CSP

	// election timeout & heartbeat timeout
	timer time.Time
}

// GetState return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isleader bool
	// Your code here (3A).
	rf.mu.Lock()
	defer rf.mu.Unlock()
	term = rf.currentTerm
	isleader = rf.state == Leader
	return term, isleader
}

// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
// before you've implemented snapshots, you should pass nil as the
// second argument to persister.Save().
// after you've implemented snapshots, pass the current snapshot
// (or nil if there's not yet a snapshot).
func (rf *Raft) persist() {

	// Your code here (3C).
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// raftstate := w.Bytes()
	// rf.persister.Save(raftstate, nil)
}

// restore previously persisted state.
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (3C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}

// Snapshot the service says it has created a snapshot that has
// all info up to and including index. this means the
// service no longer needs the log through (and including)
// that index. Raft should now trim its log as much as possible.
func (rf *Raft) Snapshot(index int, snapshot []byte) {
	// Your code here (3D).
}

// set rf.state
// in lock
func (rf *Raft) setState(s State) {
	// reset time
	// TODO: Do need to modify the reset time?
	rf.resetTimer()

	if s == Leader {
		// TODO: change become leader
		//if rf.state != Candidate {
		//	log.Fatalf("Server %v %p (Term %v) Invalid state change to leader from %v", rf.me, rf, rf.currentTerm, rf.state)
		//}
		rf.state = Leader
		rf.votedFor = -1
	} else if s == Candidate {
		// TODO: change become candidate
		//if rf.state != Follower {
		//	log.Fatalf("Server %v %p (Term %v) Invalid state change to candidate from %v", rf.me, rf, rf.currentTerm, rf.state)
		//}
		rf.state = Candidate
		rf.votedFor = rf.me
		// next term
		rf.currentTerm++
	} else {
		rf.state = Follower
		rf.votedFor = -1
	}
}

// log operate

func (rf *Raft) lastLogIndex() int {
	if len(rf.logs) == 0 {
		return 0
	}
	return rf.logs[len(rf.logs)-1].Metadata.Index
}

func (rf *Raft) lastLogTerm() int {
	if len(rf.logs) == 0 {
		return 0
	}
	return rf.logs[len(rf.logs)-1].Metadata.Term
}

func (rf *Raft) firstLogIndex() int {
	return rf.logs[0].Metadata.Index
}

func (rf *Raft) firstLogTerm() int {
	return rf.logs[0].Metadata.Term
}

func (rf *Raft) getLogTerm(index int) int {
	DPrintf("log entrise = %+v \n", rf.logs[index])
	return rf.logs[index].Metadata.Term
}

// return true that this raft peer is up-to-date
// If votedFor is null or candidateId,
// and candidate's log is at least as up-to-date as receiver's log, grant vote
func (rf *Raft) isLogUpdate(term, index int) bool {
	lastLogIndex, lastLogTerm := rf.lastLogIndex(), rf.lastLogTerm()
	return term > lastLogTerm || (term == lastLogTerm && index >= lastLogIndex)
}

// RequestVoteArgs example RequestVote RPC arguments structure.
// field names must start with capital letters!
type RequestVoteArgs struct {
	// Your data here (3A, 3B).
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

// RequestVoteReply example RequestVote RPC reply structure.
// field names must start with capital letters!
type RequestVoteReply struct {
	// Your data here (3A).
	Term        int
	VoteGranted bool // true means candidate received vote
}

// example RequestVote RPC handler.
// candidates to gather votes
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (3A, 3B).

	rf.mu.Lock()
	defer rf.mu.Unlock()
	// Reply false if term < currentTerm
	// candidate to leader
	DPrintf("Server %v %p (in term: %v) received RequestVote from %v, "+
		"rf.votedFor: %v\n",
		rf.me, rf, rf.currentTerm, args, rf.votedFor)
	// 自身的term大于请求的
	// Reply false if term < currentTerm
	if rf.currentTerm > args.Term {
		reply.Term, reply.VoteGranted = rf.currentTerm, false
		return
	}
	// 自身的term小于请求的
	// 说明需要更新log，转化为follow, 重置投票记录
	if rf.currentTerm < args.Term {
		rf.setState(Follower)
		rf.currentTerm = args.Term
	}
	// 检查是否可以投票
	// 如果未投票或已投给该 Candidate 且日志是最新的，给予投票
	// If votedFor is `null` or `candidateId`,
	// and candidate's log is at least as up-to-date as receiver's log, grant vote
	if (rf.votedFor == -1 || rf.votedFor == args.CandidateId) && rf.isLogUpdate(args.Term, args.LastLogIndex) {
		rf.votedFor = args.CandidateId
		rf.resetTimer() // 重置选举计时器
		reply.Term, reply.VoteGranted = rf.currentTerm, true
	} else {
		reply.Term, reply.VoteGranted = rf.currentTerm, false
	}
	DPrintf("Server %v %p (in term: %v) received RequestVote from %v, "+
		"rf.votedFor: %v\n",
		rf.me, rf, rf.currentTerm, args, rf.votedFor)
}

// 广播请求投票的消息
// in lock
func (rf *Raft) broadcastRequestVote() {
	DPrintf("Server %v (Term: %v) broadcast RequestVote", rf.me, rf.currentTerm)
	if rf.state != Candidate {
		return
	}
	args := &RequestVoteArgs{
		Term:         rf.currentTerm,
		CandidateId:  rf.me,
		LastLogIndex: rf.lastLogIndex(),
		LastLogTerm:  rf.lastLogTerm(),
	}
	receivedVotes := 1 // me
	for peer := range rf.peers {
		if peer != rf.me {
			// use goroutine collect vote for all peers
			go func(peer int) {
				reply := &RequestVoteReply{}
				if rf.sendRequestVote(peer, args, reply) {
					// 这里的lock不会dead lock
					// 这是用来保证 votes
					rf.mu.Lock()
					defer rf.mu.Unlock()

					if rf.state != Candidate {
						return
					}
					if rf.currentTerm == args.Term {
						// get vote
						if reply.VoteGranted {
							receivedVotes++
							if receivedVotes >= len(rf.peers)/2 {
								rf.setState(Leader)
								// 立即发送一个heat time
								// TODO: send heat time
							}
						} else {
							// 说明自己的term低于其他节点
							if reply.Term > rf.currentTerm {
								rf.currentTerm = reply.Term
								// rf.votedFor = -1
								rf.setState(Follower)
							}
						}
					}
				}
			}(peer)
		}
	}

}

// AppendEntriesArgs
type AppendEntriesArgs struct {
	Term         int        // leader's term
	LeaderId     int        // so follower can redirect clients
	PrevLogIndex int        // index of log entry immediately preceding new ones
	PrevLogTerm  int        // term of prevLogIndex entry
	Entries      []LogEntry // log entries to store (empty for heartbeat; may send more than one for efficiency)
	LeaderCommit int        // leader's commitIndex
}

// AppendEntriesReply
type AppendEntriesReply struct {
	Term    int
	Success bool // true if follower contained entry matching prevLogIndex and prevLogTerm
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// 如果 leader 的 term 比当前节点小，拒绝请求
	// 1. Reply false if term < currentTerm
	if args.Term < rf.currentTerm {
		reply.Term, reply.Success = rf.currentTerm, false
		return
	}
	// 重置心跳时间
	rf.resetTimer()

	// 如果收到更高的任期，更新自己为 follower 并更新 term
	// 这个可以将还是在Candidate的peer变为follower
	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.setState(Follower)
	}

	// TODO: Log Persistence
	if len(args.Entries) == 0 {
		// 心跳函数
		DPrintf("[heart]Server %v receive heartbeat to leader &%v. \n", rf.me, args.LeaderId)
	} else {
		DPrintf("[log repliction]Server %v receive heartbeat to leader &%v. \n", rf.me, args.LeaderId)
	}
	// 检查follow中的log是否为最新的log
	// 2.
	// Reply false if log doesn't contain an entry at prevLogIndex whose term matches prevLogTerm
	// 如果 PrevLogIndex 超出 Follower 当前日志的索引范围，
	// 这意味着 Follower 的日志比 Leader 预期的更短

	if args.PrevLogIndex > rf.lastLogIndex() || rf.logs[args.PrevLogIndex].Metadata.Term != args.PrevLogTerm {
		reply.Term, reply.Success = rf.currentTerm, false
		return
	}
	// 3.
	// If an existing entry conflicts with a new one (same index but different terms),
	// delete the existing entry and all that follow it
	if len(rf.logs) > 0 && args.PrevLogIndex < rf.lastLogIndex() && rf.logs[args.PrevLogIndex].Metadata.Term != args.PrevLogTerm {
		rf.logs = rf.logs[:args.PrevLogIndex]
	}
	// Append any new entries not already in the log
	// TODO: Log Replication
	rf.logs = append(rf.logs, args.Entries...)
	DPrintf("Server %v successfully received AppendEntries from %v. \n", rf.me, args.LeaderId)
	reply.Term, reply.Success = rf.currentTerm, true

	// commitIndex
	// If leaderCommit > commitIndex,
	// set commitIndex = min(leaderCommit, index of last new entry)
	if args.LeaderCommit > rf.commitIndex {
		rf.commitIndex = int(math.Min(float64(args.LeaderCommit), float64(rf.lastLogIndex())))
	}
}

func (rf *Raft) broadcastAppendEntries() {
	DPrintf("Server %v (Term: %v) broadcast HeartBeat", rf.me, rf.currentTerm)
	for !rf.killed() {
		if rf.state != Leader {
			return
		}
		for peer := range rf.peers {
			if peer != rf.me {
				args := &AppendEntriesArgs{
					Term:         rf.currentTerm,
					LeaderId:     rf.me,
					PrevLogIndex: rf.nextIndex[peer] - 1,
					PrevLogTerm:  rf.getLogTerm(rf.nextIndex[peer] - 1),
					// Entries:      rf.logs,
					LeaderCommit: rf.commitIndex,
				}
				// 确定发送的是同步还是心跳
				// If last log index ≥ nextIndex for a follower
				if rf.lastLogIndex() >= rf.nextIndex[peer] {
					DPrintf("Server %v received AppendEntries from %v for log replication. \n", rf.me, peer)
					args.Entries = rf.logs[rf.nextIndex[peer]:]
				} else {
					DPrintf("Server %v received AppendEntries from %v for heartBeat. \n", rf.me, peer)
					args.Entries = make([]LogEntry, 0)
				}
				go rf.handleHeartBeat(peer, args)
			}
		}
	}
}

func (rf *Raft) handleHeartBeat(serverTo int, args *AppendEntriesArgs) {
	reply := &AppendEntriesReply{}
	// copy
	sendArgs := *args
	if rf.sendAppendEntries(serverTo, args, reply) {
		rf.mu.Lock()
		defer rf.mu.Unlock()
		// 在同Term的修改是有用的
		if sendArgs.Term != rf.currentTerm {
			return
		}
		//
		if reply.Success {
			// 更新当前此follow的index和下一个的index
			rf.matchIndex[serverTo] = args.PrevLogIndex + len(args.Entries)
			rf.nextIndex[serverTo] = rf.matchIndex[serverTo] + 1
			// 检查是否可以更新 commitIndex
			// If there exists an N such that N > commitIndex,
			// a majority of matchIndex[i] ≥ N,
			// and log[N].term == currentTerm: set commitIndex = N(§5.3, §5.4).
			rf.updateCommit()
		} else if rf.state == Leader {
			// 超Term的
			if reply.Term > rf.currentTerm {
				rf.currentTerm = reply.Term
				rf.setState(Follower)
				return
			}
			if reply.Term == rf.currentTerm {
				// TODO: BACKROCK
				rf.nextIndex[serverTo]--
			}
		} else {
			return
		}
	}
}

func (rf *Raft) updateCommit() {
	// Leader 会在 matchIndex[] 中查找大多数节点已经复制的最高日志条目索引
	// 需要查看的大多数
	sortMatch := make([]int, len(rf.peers))
	copy(sortMatch, rf.matchIndex)
	sort.Ints(sortMatch)
	// 取中位数作为大多数节点的匹配索引
	// N is major peers
	N := sortMatch[len(rf.peers)/2] // matchIndex[i] ≥ N
	// 如果存在某个 N，使得 N > commitIndex，
	// 且多数 `matchIndex[i] ≥ N`，
	// 且 `log[N].term == currentTerm`
	if N > rf.commitIndex && rf.logs[N].Metadata.Term == rf.currentTerm {
		rf.commitIndex = N
		// TODO: apply logs
	}
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

func (rf *Raft) sendAppendEntries(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
	ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	return ok
}

// Start the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
term. the third return value is true if this server believes it is
// the leader.
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	// Your code here (3B).
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if rf.state != Leader {
		isLeader = false
	} else {
		term = rf.currentTerm
		// log index ++
		index = rf.lastLogIndex() + 1
		rf.logs = append(rf.logs, LogEntry{
			Metadata: LogMetadata{
				Term:  term,
				Index: index,
			},
			Command: command,
		})
	}

	return index, term, isLeader
}

// Kill the tester doesn't halt goroutines created by Raft after each test,
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

func (rf *Raft) ticker() {
	for rf.killed() == false {

		// Your code here (3A)
		// Check if a leader election should be started.
		rf.mu.Lock()
		// election time ~150-300ms
		if rf.state == Candidate {
			// 上一次选举超时
			if time.Since(rf.timer) > time.Duration(200+rand.Int63()%300)*time.Millisecond {
				rf.setState(Candidate)
				// rf.resetTimer()
				// rf.currentTerm++
				rf.broadcastRequestVote()
			}
		} else if rf.state == Leader {
			// heartbeat timeout ~100-150ms
			// Lab03A-Hint:
			// The tester requires that the leader send heartbeat RPCs no
			// more than ten times per second.
			if time.Since(rf.timer) > time.Duration(100+rand.Int63()%50)*time.Millisecond {
				rf.resetTimer()
				rf.broadcastAppendEntries()
			}
		} else {
			// follow to candidate
			if time.Since(rf.timer) > time.Duration(150+rand.Int63()%300)*time.Millisecond {
				rf.setState(Candidate)
				// rf.resetTimer()
				rf.broadcastRequestVote()
			}
		}
		// pause for a random amount of time between 50 and 350
		// milliseconds.
		rf.mu.Unlock()
		ms := 50 + (rand.Int63() % 300) // ~50-300
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}

// 重置时间
func (rf *Raft) resetTimer() {
	rf.timer = time.Now()
}

// Make the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
func Make(peers []*RPCEnd, me int, persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me
	rf.applyCh = applyCh
	// Your initialization code here (3A, 3B, 3C).
	// create new raft server peer
	rf.state = Follower

	rf.currentTerm = 0
	rf.votedFor = -1 // this is candidateId that received vote in current term
	rf.logs = make([]LogEntry, 0)

	rf.commitIndex = 0
	rf.lastApplied = 0

	n := len(peers)
	rf.nextIndex = make([]int, n)
	//initialized to leader last log index + 1
	for i, _ := range rf.nextIndex {
		rf.nextIndex[i] = 1
	}
	rf.matchIndex = make([]int, n)

	rf.resetTimer()
	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	// start ticker goroutine to start elections
	go rf.ticker()

	return rf
}
