package raft

func (rf *Raft) FollowerLaunchElection() {
	rf.becomeCandidate()
}


