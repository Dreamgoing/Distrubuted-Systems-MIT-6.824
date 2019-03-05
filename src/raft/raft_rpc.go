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

type Message struct {
	Type       string
	To         int
	From       int
	Term       int
	LogTerm    int
	Index      int
	Entries    []Entry
	Commit     int
	Reject     bool
	RejectHint int
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	var term, voteGrated = None, true

	LevelDPrintf("%v %v Term: %v args.Term: %v", ShowVariable, rf.state, rf.me, rf.Term, args.Term)
	if args.Term < rf.Term {
		term = rf.Term
		voteGrated = false
	}
	// -1 indicate nil
	if voteGrated && rf.Vote == None || args.CandidateID != None {
		voteGrated = true
		LevelDPrintf("%v %v vote %v", ShowProcess, rf.state, rf.me, args.CandidateID)
		rf.Term = args.Term
		rf.Vote = args.CandidateID
		rf.becomeFollower()
	}

	reply.Term = term
	reply.VoteGrated = voteGrated

	// Your code here (2A, 2B).
}

type AppendEntriesArgs struct {
	Term         int
	LeaderID     int
	Logs         []Entry
	PrevLogIndex int
	PrevLogTerm  int
	Commit       int
}

type AppendEntriesReply struct {
	Term    int
	Success bool
	Index   int
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {

	// 1. term < Term
	if args.Term < rf.Term {
		reply.Term = rf.Term
		reply.Success = false
		return
	}
	rf.becomeFollower()

	LevelDPrintf("len(rf.logs): %v args.PrevLogIndex: %v", ShowVariable, len(rf.logs), args.PrevLogIndex)
	// 2. 当前Follower上面已提交日志的索引小于Leader发来的最后一个日志的索引
	// 这种情况需要Leader再补发之前未在本Follower提交的日志
	reply.Index = len(rf.logs)
	if len(rf.logs) < args.PrevLogIndex {
		reply.Success = false
		return
		//	3. 删除已存在冲突的日志条目以及之后所有的日志
	} else if len(rf.logs) > 0 && args.PrevLogIndex > 0 && rf.logs[args.PrevLogIndex-1].Term != args.PrevLogTerm {
		LevelDPrintf("delete conflict, PrevLogIndex: %v", ShowProcess, args.PrevLogIndex)
		rf.logs = rf.logs[:args.PrevLogIndex]

	}

	if len(args.Logs) == 0 {
		LevelDPrintf("%v %v receive HeartBeat from %v", ShowProcess, rf.state, rf.me, args.LeaderID)
	} else {
		// 4. 添加新增的日志
		LevelDPrintf("%v %v receive log from leader:%v, Term:%v, term:%v log:%v",
			ShowProcess, rf.state, rf.me, args.LeaderID, rf.Term, args.Term, args.Logs[0].Command)
		//LevelDPrintf("pre len: %v", len(rf.logs))
		rf.logs = append(rf.logs, args.Logs...)
		reply.Index = len(rf.logs)
	}

	// 5. 更新日志commitIndex
	if args.Commit > rf.commit {
		rf.commit = Min(args.Commit, len(rf.logs))
	}

	//LevelDPrintf("args.Commit:%v rf.commit:%v len(rf.logs)-1: %v", ShowProcess, args.Commit, rf.commit, len(rf.logs)-1)

	reply.Term = rf.Term
	reply.Success = true
}
