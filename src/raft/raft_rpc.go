package raft

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
	Term       int
	VoteGrated bool
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	var term, voteGrated = None, false

	if args.Term < rf.currentTerm {
		voteGrated = false
	}
	// -1 indicate nil
	if rf.votedFor == -1 || args.CandidateID != -1 {
		voteGrated = true
	}

	rf.currentTerm = term
	rf.votedFor = args.CandidateID

	voteGrated = true
	reply.Term = term
	reply.VoteGrated = voteGrated

	// Your code here (2A, 2B).
}

type AppendEntriesArgs struct {
	Term         int
	LeaderID     int
	PrevLogIndex int
	PrevLogTerm  int
	LeaderCommit int
}

type AppendEntriesReply struct {
	Term    int
	Success bool
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	DPrintf("%v %v accept AppendEntries", rf.state, rf.me)
	var term, success = None, false
	if term < rf.currentTerm {
		success = false
	} else {
		rf.currentTerm = term
		rf.state = FollowerState
		rf.timer.Reset(rf.electionTimeout)
	}

	success = true

	term = rf.currentTerm
	reply.Term = term
	reply.Success = success
}
