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
	rf.electionTimeout = time.Duration(rand.Int63n(RandNum)+ElectionTimeout) * time.Millisecond
	rf.timer = make(map[State]*time.Timer, 3)

	rf.timer[FollowerState] = time.NewTimer(rf.electionTimeout)
	rf.timer[CandidateState] = time.NewTimer(rf.electionTimeout)
	rf.timer[LeaderState] = time.NewTimer(rf.electionTimeout)

	rf.state = FollowerState
	rf.currentTerm = None
	rf.votedFor = None
	rf.commitIndex = None
	rf.nextIndex = make([]int, len(rf.peers))
}

func (rf *Raft) ToFollower() {
	rf.votedFor = None
	rf.state = FollowerState
	rf.ResetTimer(rf.state)
}

func (rf *Raft) ToCandidate() {
	rf.state = CandidateState
}

func (rf *Raft) PrepareElection() {
	rf.currentTerm++
	rf.votedFor = rf.me
	rf.ResetTimer(rf.state)
}

func (rf *Raft) ToLeader() {
	rf.state = LeaderState
	rf.votedFor = None
	rf.ResetTimer(rf.state)
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

func (rf *Raft) ResetTimer(state State) {
	timeout, _ := TimeOutMapping[state]
	rf.timer[state].Stop()
	rf.electionTimeout = (timeout + time.Duration(rand.Int63n(RandNum))) * time.Millisecond
	rf.timer[state] = time.NewTimer(rf.electionTimeout)
	rf.timer[state].Reset(rf.electionTimeout)
}

func (rf *Raft) AcquireLock() {
	rf.mu.Lock()
	rf.isLock = true
}

func (rf *Raft) ReleaseLock() {
	rf.isLock = false
	rf.mu.Unlock()
}
