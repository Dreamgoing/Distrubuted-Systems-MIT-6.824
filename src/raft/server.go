package raft

import (
	"time"
	"math/rand"
)

func (rf *Raft) Init() {
	rf.timeout = time.Duration(rand.Int63n(RandNum)+ElectionTimeout) * time.Millisecond
	rf.timer = make(map[State]*time.Timer, 3)

	rf.timer[StateFollower] = time.NewTimer(rf.timeout)
	rf.timer[StateCandidate] = time.NewTimer(rf.timeout)
	rf.timer[StateLeader] = time.NewTimer(rf.timeout)

	for _, timer := range rf.timer {
		timer.Stop()
	}

	rf.trans = make(chan State)
	rf.state = StateFollower
	rf.becomeFollower()

	rf.Term = None
	rf.Vote = None
	rf.commit = Zero
	rf.lastApplied = Zero
	rf.nextIndex = make([]int, len(rf.peers))
	rf.matchIndex = make([]int, len(rf.peers))
}

func (rf *Raft) becomeFollower() {
	rf.timer[rf.state].Stop()
	rf.Vote = None
	rf.state = StateFollower
	rf.ResetTimer()
}

func (rf *Raft) becomeCandidate() {
	rf.timer[rf.state].Stop()
	rf.state = StateCandidate
}

func (rf *Raft) PrepareElection() {
	rf.Term++
	rf.Vote = rf.me
	rf.ResetTimer()
}

func (rf *Raft) becomeLeader() {
	rf.timer[rf.state].Stop()
	rf.state = StateLeader
	rf.Vote = None
	rf.ResetTimer()
	rf.nextIndex = make([]int, len(rf.peers))
	for i := range rf.peers {
		rf.nextIndex[i] = rf.commit + 1
	}

}

func (rf *Raft) ApplyCommit() {
	//DPrintf("%v %v %v", rf.me, rf.commit, rf.lastApplied)
	if rf.commit > rf.lastApplied {
		LevelDPrintf("%v %v apply commit %v", ShowProcess, rf.state, rf.me, rf.logs[rf.lastApplied].Command)
		rf.applyChan <- ApplyMsg{
			true,
			rf.logs[rf.lastApplied].Command,
			rf.logs[rf.lastApplied].Index}
		rf.lastApplied++
	}

}

func (rf *Raft) ResetTimer() {
	timeout, _ := TimeOutMapping[rf.state]
	rf.timeout = (timeout + time.Duration(rand.Int63n(RandNum))) * time.Millisecond
	rf.timer[rf.state].Reset(rf.timeout)
}

func (rf *Raft) AcquireLock() {
	rf.mu.Lock()
	rf.isLock = true
}

func (rf *Raft) ReleaseLock() {
	rf.isLock = false
	rf.mu.Unlock()
}
