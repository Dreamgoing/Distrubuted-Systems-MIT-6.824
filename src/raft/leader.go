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

	reply := &AppendEntriesReply{}
	lastLog := rf.GetLastLog()

	for it := range rf.peers {
		if it == rf.me {
			continue
		}

		var log []Entry
		if rf.nextIndex[it] < len(rf.logs) {
			log = rf.logs[rf.nextIndex[it]:]
		}

		go func(it int, log []Entry) {
			if len(log) == 0 {
				LevelDPrintf("%v %v send HeartBeat to %v", ShowProcess, rf.state, rf.me, it)
			} else {
				LevelDPrintf("%v %v send log %v %v to %v", ShowProcess,
					rf.state, rf.me, log[0].Command, log[0].Index, it)
			}

			LevelDPrintf("prevLogTerm: %v prevLogIndex: %v rf.nextIndex[%v]: %v",
				ShowVariable, lastLog.Term, lastLog.Index, it, rf.nextIndex[it])

			ok := rf.sendAppendEntries(it, &AppendEntriesArgs{rf.currentTerm,
				rf.me, log, lastLog.Index,
				lastLog.Term, rf.commitIndex},
				reply)

			if ok {
				// FIXME Lock
				if !reply.Success && reply.Term > rf.currentTerm {
					LevelDPrintf("%v %v became follower", ShowProcess, rf.state, rf.me)
					rf.currentTerm = reply.Term
					rf.ToFollower()
				}
				if reply.Success {
					atomic.AddInt32(&cnt, 1)
					rf.matchIndex[it] = reply.Index
				}
				rf.nextIndex[it] = reply.Index
			}
			//LevelDPrintf("%v %v send AppendEntries to %v %v", ShowProcess, rf.state, rf.me, it, ok)

			wg.Done()
		}(it, log)
	}
	wg.Wait()

	LevelDPrintf("%v %v appendEntries %v/%v success", ShowProcess, rf.state, rf.me, cnt, total)

	if !MajorityOk(int(cnt), total) {
		rf.ToFollower()
	}
	rf.ResetTimer()

}

func (rf *Raft) GetLastLog() Entry {
	if len(rf.logs) <= 1 {
		return Entry{None, None, None}
	} else {
		idx := len(rf.logs) - 2
		return rf.logs[idx]
	}

}
