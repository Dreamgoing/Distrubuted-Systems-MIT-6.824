package raft

import (
	"sync"
	"sync/atomic"
)

func (rf *Raft) LeaderAppendEntries() {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	cnt := int32(1)
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
				rf.me, rf.logs[rf.nextIndex[it]:], lastLog.Term,
				lastLog.Index, rf.commitIndex},
				reply)

			if !ok || !reply.Success && reply.Term > rf.currentTerm {
				DPrintf("%v %v became follower", rf.state, rf.me)
				rf.currentTerm = reply.Term
				rf.state = FollowerState
				rf.votedFor = None
				rf.timer.Reset(rf.electionTimeout)
			}
			if ok && reply.Success {
				atomic.AddInt32(&cnt, 1)
				rf.nextIndex[it] = rf.commitIndex + 1
			}
			wg.Done()
		}(it)
	}
	wg.Wait()
	DPrintf("%v %v appendEntries %v/%v success", rf.state, rf.me, cnt, total)
}

func (rf *Raft) FollowerElectionTimeout() {

}

func (rf *Raft) CandidateRequestVotes() {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	rf.currentTerm++
	rf.votedFor = None
	rf.timer.Reset(rf.electionTimeout)
	DPrintf("New %v %v", rf.state, rf.me)

	total := len(rf.peers)
	cnt := int32(1)

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
			if ok && reply.VoteGrated {
				atomic.AddInt32(&cnt, 1)
			}
			wg.Done()
		}(it)
	}
	wg.Wait()

	DPrintf("%v %v get %v/%v votes", rf.state, rf.me, cnt, total)
	if int(cnt) > total/2 {
		DPrintf("%v %v became leader, term:%v", rf.state, rf.me, rf.currentTerm)
		rf.ToLeader()
	} else {
		rf.ToFollower()
		rf.currentTerm--
	}

}

func countSuccess(ch chan bool, res *int) {
	for it := range ch {
		if it {
			*res ++
		}
	}
}

func (rf *Raft) GetLastLog() Entry {
	if len(rf.logs) == 0 {
		return Entry{None, None, None}
	} else {
		return rf.logs[len(rf.logs)-1]
	}

}
