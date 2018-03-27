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
		term = rf.currentTerm
		voteGrated = false
	}
	// -1 indicate nil
	if rf.votedFor == None || args.CandidateID != None {
		voteGrated = true
		DPrintf("%v %v vote %v", rf.state, rf.me, args.CandidateID)
		rf.currentTerm = args.Term
		rf.votedFor = args.CandidateID
	}

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

	var term, success = None, false
	if args.Term < rf.currentTerm {
		term = rf.currentTerm
		success = false
	} else {
		DPrintf("%v %v accept AppendEntries, leader:%v,currentTerm:%v, term:%v",
			rf.state, rf.me, args.LeaderID, rf.currentTerm, args.Term)
		rf.currentTerm = args.Term
		success = true
		rf.state = FollowerState
		rf.timer.Reset(rf.electionTimeout)
		term = rf.currentTerm
	}

	reply.Term = term
	reply.Success = success
}
