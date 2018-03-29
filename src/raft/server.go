package raft

func (rf *Raft) ApplyStateMachine() {
	if rf.commitIndex > rf.lastApplied {
		rf.lastApplied++
	}
	//log[rf.lastApplied] => StateMachine
}

func (rf *Raft) ToFollower() {
	rf.timer.Reset(rf.electionTimeout)
	rf.votedFor = None
	rf.state = FollowerState
}

func (rf *Raft) PrepareElection() {
	rf.currentTerm++
	rf.timer.Reset(rf.electionTimeout)
}

func (rf *Raft) ToLeader() {
	rf.state = LeaderState
	rf.votedFor = None
	rf.timer.Reset(rf.electionTimeout)
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
