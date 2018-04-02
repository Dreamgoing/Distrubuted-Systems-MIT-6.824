package raft

import (
	"sync"
	"sync/atomic"
)

func (rf *Raft) CandidateRequestVotes() {
	rf.PrepareElection()
	LevelDPrintf("New %v %v", ShowProcess, rf.state, rf.me)

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
	LevelDPrintf("%v %v get %v/%v votes", ShowProcess, rf.state, rf.me, cnt, total)
	if int(cnt) > total/2 {
		LevelDPrintf("%v %v became leader, term:%v", ShowProcess, rf.state, rf.me, rf.currentTerm)
		rf.ToLeader()
	} else {
		rf.ToFollower()
		rf.currentTerm--
	}

}
