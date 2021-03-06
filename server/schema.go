package server

import (
	"github.com/VikasSherawat/raft/labrpc"
	"github.com/VikasSherawat/raft/storage"
	"sync"
	"time"
)

type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *storage.Persister  // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]
	dead      int32               // set by Kill()

	//Persistent state on all servers:
	CurrentTerm int
	VotedFor    int
	Log         LogType

	//Volatile state on all servers:
	CommitIndex int
	LastApplied int

	//Volatile state on leaders:
	NextIndex  []int
	MatchIndex []int

	//Volatile state on all servers:
	//(added by me
	State serverState //0:follower; 1:candidate; 2:leader
	//tickerResetChannel chan bool
	ApplyCh   chan ApplyMsg
	timer     *time.Timer
	timerLock sync.Mutex
	// Look at the paper's Figure 2 for a description of what
	// state a Raft storage must maintain.

	//added in lab2B:
	newLogCome *sync.Cond //only valid in leader state

	VoteCount                      int
	heartbeatTimer                 time.Timer
	heartbeatTimerTerminateChannel chan bool

	//added in lab2D
	snapshot  []byte
	applyCond *sync.Cond
}

type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int

	// For 2D:
	SnapshotValid bool
	Snapshot      []byte
	SnapshotTerm  int
	SnapshotIndex int
}

type LogEntry struct {
	Term    int
	Command interface{}
}

type LogType struct {
	Entries           []LogEntry
	LastIncludedIndex int
	LastIncludedTerm  int
}

type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	Term         int
	CandidateID  int
	LastLogIndex int
	LastLogTerm  int
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	Term        int
	VoteGranted bool
}

type AppendEntriesArgs struct {
	Term         int
	LeaderID     int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int
}

type AppendEntriesReply struct {
	Term                   int
	Success                bool
	ConflictEntryTerm      int
	ConflictTermFirstIndex int
}
type InstallSnapshotArgs struct {
	Term              int
	LeaderID          int
	LastIncludedIndex int
	LastIncludedTerm  int
	Data              []byte
}

type InstallSnapshotReply struct {
	Term int
}

type serverState int
