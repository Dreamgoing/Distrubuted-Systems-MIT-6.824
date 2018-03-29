package raft

import (
	"sync"
	"sync/atomic"
)

func (rf *Raft) LeaderAppendEntries() {
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
			DPrintf("prevLogTerm: %v prevLogIndex: %v", lastLog.Term, lastLog.Index)
			ok := rf.sendAppendEntries(it, &AppendEntriesArgs{rf.currentTerm,
				rf.me, rf.logs[rf.nextIndex[it]:], lastLog.Index,
				lastLog.Term, rf.commitIndex},
				reply)

			if !ok || !reply.Success && reply.Term > rf.currentTerm {
				DPrintf("%v %v became follower", rf.state, rf.me)
				rf.currentTerm = reply.Term
				rf.state = FollowerState
				rf.votedFor = None
				rf.timer[FollowerState].Reset(rf.electionTimeout)
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

	if !MajorityOk(int(cnt), total) {
		rf.ToFollower()
	}

}

func (rf *Raft) GetLastLog() Entry {
	if len(rf.logs) <= 1 {
		return Entry{None, None, None}
	} else {
		return rf.logs[rf.commitIndex-1]
	}

}
