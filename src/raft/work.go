package raft

func (rf *Raft) LeaderAppendEntries() {


}

func (rf *Raft) FollowerElectionTimeout() {


}

func (rf *Raft) CandidateRequestVotes() {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	rf.currentTerm++
	rf.votedFor = None
	DPrintf("New %v %v", rf.state, rf.me)
	voteRes := make(chan bool, 3)
	total := len(rf.peers)

	for it := range rf.peers {
		go func(it int) {
			reply := &RequestVoteReply{}
			ok := rf.sendRequestVote(it, &RequestVoteArgs{rf.currentTerm,
				rf.me, None, None}, reply)
			if ok {

			}
			voteRes <- reply.VoteGrated

		}(it)
	}
	num := countVotes(voteRes)
	if num > total/2 {
		DPrintf("%v %v became leader, term:%v", rf.state, rf.me, rf.currentTerm)
		rf.becameLeader()
	}

}

func countVotes(voteRes chan bool) int {
	res := 0
	for it := range voteRes {
		if it {
			res ++
		}
	}
	return res
}

func (rf *Raft) becameLeader() {
	rf.state = LeaderState
	rf.timer.Stop()
}
