package raft

import "time"

const None = -1
const Zero = 0

type State string

const (
	StateFollower  = "Follower"
	StateCandidate = "Candidate"
	StateLeader    = "Leader"
	End            = "End"
)

const (
	ElectionTimeout  = 350
	HeartbeatTimeout = 10
	CandidateTimeout = 100
)

const RandNum = 100

var TimeOutMapping = map[State]time.Duration{
	StateFollower:  ElectionTimeout,
	StateCandidate: CandidateTimeout,
	StateLeader:    HeartbeatTimeout,}
