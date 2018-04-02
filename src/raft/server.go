package raft

import (
	"time"
	"math/rand"
)

func (rf *Raft) ApplyStateMachine() {
	if rf.commitIndex > rf.lastApplied {
		rf.lastApplied++
	}
	//log[rf.lastApplied] => StateMachine
}

func (rf *Raft) Init() {
	rf.timeout = time.Duration(rand.Int63n(RandNum)+ElectionTimeout) * time.Millisecond
	rf.timer = make(map[State]*time.Timer, 3)

	rf.timer[FollowerState] = time.NewTimer(rf.timeout)
	rf.timer[CandidateState] = time.NewTimer(rf.timeout)
	rf.timer[LeaderState] = time.NewTimer(rf.timeout)

	for _, timer := range rf.timer {
		timer.Stop()
	}

	rf.trans = make(chan State)
	rf.state = FollowerState
	rf.ToFollower()

	rf.currentTerm = None
	rf.votedFor = None
	rf.commitIndex = None
	rf.nextIndex = make([]int, len(rf.peers))
	rf.matchIndex = make([]int, len(rf.peers))
}

func (rf *Raft) ToFollower() {
	rf.timer[rf.state].Stop()
	rf.votedFor = None
	rf.state = FollowerState
	rf.ResetTimer()
}

func (rf *Raft) ToCandidate() {
	rf.timer[rf.state].Stop()
	rf.state = CandidateState
}

func (rf *Raft) PrepareElection() {
	rf.currentTerm++
	rf.votedFor = rf.me
	rf.ResetTimer()
}

func (rf *Raft) ToLeader() {
	rf.timer[rf.state].Stop()
	rf.state = LeaderState
	rf.votedFor = None
	rf.ResetTimer()
	rf.nextIndex = make([]int, len(rf.peers))
	for i := range rf.peers {
		rf.nextIndex[i] = rf.commitIndex + 1
	}

}

func (rf *Raft) ApplyCommit() {
	if rf.commitIndex > rf.lastApplied {
		rf.lastApplied++
		rf.applyChan <- ApplyMsg{
			true,
			rf.logs[rf.lastApplied].Command,
			rf.lastApplied}
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
