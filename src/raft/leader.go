package raft

import (
	"sync"
	"sync/atomic"
)

func (rf *Raft) LeaderAppendEntries() bool {
	// TODO 将heart beat拆分出来
	cnt := int32(1)
	total := len(rf.peers)

	var wg sync.WaitGroup
	wg.Add(total - 1)

	reply := &AppendEntriesReply{}

	for it := range rf.peers {
		if it == rf.me {
			continue
		}
		lastLog := rf.GetPrevLog(it)
		var log []Entry
		if rf.nextIndex[it] < len(rf.logs) {
			log = rf.logs[rf.nextIndex[it]:]
		}

		go func(it int, log []Entry) {
			if len(log) == 0 {
				//LevelDPrintf("%v %v send HeartBeat to %v", ShowProcess, rf.state, rf.me, it)
				rf.sendHeartbeat(it)
			} else {
				//LevelDPrintf("%v %v send log %v index: %v to %v", ShowProcess,
				//	rf.state, rf.me, log[0].Command, log[0].Index, it)
				LevelDPrintf("prevLogTerm: %v prevLogIndex: %v rf.nextIndex[%v]: %v",
					ShowVariable, lastLog.Term, lastLog.Index, it, rf.nextIndex[it])

			}
			ok := rf.sendAppendEntries(it, &AppendEntriesArgs{rf.Term,
				rf.me, log, lastLog.Index,
				lastLog.Term, rf.commit},
				reply)

			if ok {
				// FIXME Lock
				if !reply.Success && reply.Term > rf.Term {
					LevelDPrintf("%v %v became follower", ShowProcess, rf.state, rf.me)
					rf.Term = reply.Term
					rf.becomeFollower()
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

	//LevelDPrintf("%v %v appendEntries %v/%v success", ShowProcess, rf.state, rf.me, cnt, total)

	if !MajorityOk(int(cnt), total) {
		return false
		rf.becomeFollower()
	}
	rf.ResetTimer()
	return true

}

func (rf *Raft) GetPrevLog(it int) Entry {
	if len(rf.logs) <= 1 {
		return Entry{None, None, None}
	} else {
		idx := rf.nextIndex[it] - 1
		return rf.logs[idx]
	}

}

func (rf *Raft) tickHeartbeat() {
	// leader 心跳函数
	rf.heartbeatElapsed++
	rf.electionElapsed++

}

func (rf *Raft) sendHeartbeat(to int) {
	reply := &AppendEntriesReply{}
	ok := rf.sendAppendEntries(to, &AppendEntriesArgs{Term: rf.Term, LeaderID: rf.me, Commit: rf.commit}, reply)
	if ok {
		LevelDPrintf("sendHeartbeat %+v", ShowProcess, reply)
	}

}
