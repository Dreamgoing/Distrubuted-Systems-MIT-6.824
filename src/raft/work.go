package raft

import "sync"

func (rf *Raft) LeaderAppendEntries() {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	appendRes := make(chan bool)
	cnt := 1
	go countSuccess(appendRes, &cnt)
	total := len(rf.peers)
	var wg sync.WaitGroup
	wg.Add(total - 1)

	for it := range rf.peers {
		if it == rf.me {
			continue
		}
		go func(it int) {
			reply := &AppendEntriesReply{}
			lastLog := rf.GetLastLog()
			ok := rf.sendAppendEntries(it, &AppendEntriesArgs{rf.currentTerm,
				rf.me, lastLog.term, lastLog.index, rf.commitIndex},
				reply)
			if !ok || !reply.Success && reply.Term > rf.currentTerm {
				DPrintf("%v %v became follower", rf.state, rf.me)
				rf.currentTerm = reply.Term
				rf.state = FollowerState
				rf.votedFor = None
				rf.timer.Reset(rf.electionTimeout)
			}
			appendRes <- reply.Success
			wg.Done()
		}(it)
	}
	wg.Wait()
	close(appendRes)
	DPrintf("%v %v appendEntries %v/%v success", rf.state, rf.me, cnt, total)
}

func (rf *Raft) FollowerElectionTimeout() {

}

func (rf *Raft) CandidateRequestVotes() {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	rf.currentTerm++
	//rf.votedFor = None
	rf.timer.Reset(rf.electionTimeout)
	DPrintf("New %v %v", rf.state, rf.me)

	voteRes := make(chan bool)
	total := len(rf.peers)
	cnt := 1
	go countSuccess(voteRes, &cnt)
	var wg sync.WaitGroup
	wg.Add(total - 1)
	for it := range rf.peers {
		if it == rf.me {
			continue
		}
		go func(it int) {
			reply := &RequestVoteReply{}
			ok := rf.sendRequestVote(it, &RequestVoteArgs{rf.currentTerm,
				rf.me, None, None}, reply)
			if ok {

			}
			voteRes <- reply.VoteGrated
			if reply.VoteGrated {
				cnt++
			}
			wg.Done()
		}(it)
	}
	wg.Wait()
	close(voteRes)
	DPrintf("%v %v get %v/%v votes", rf.state, rf.me, cnt, total)
	if cnt > total/2 {
		DPrintf("%v %v became leader, term:%v", rf.state, rf.me, rf.currentTerm)
		rf.becameLeader()
	} else {
		rf.currentTerm--
		rf.state = FollowerState
		rf.votedFor = None
	}

}

func countSuccess(ch chan bool, res *int) {
	for it := range ch {
		if it {
			*res ++
		}
	}
}

func (rf *Raft) becameLeader() {
	rf.state = LeaderState
	rf.timer.Reset(rf.electionTimeout)
}

func (rf *Raft) GetLastLog() *Entry {
	if len(rf.logs) == 0 {
		return &Entry{None, None, None}
	} else {
		return rf.logs[len(rf.logs)-1]
	}

}
