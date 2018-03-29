package raft

import "time"

const None = -1
const Zero = 0

type State string

const (
	FollowerState  = "Follower"
	CandidateState = "Candidate"
	LeaderState    = "Leader"
	End            = "End"
)

const (
	ElectionTimeout  = 300
	HeartbeatTimeout = 30
	CandidateTimeout = 150
)

const RandNum = 100

var TimeOutMapping = map[State]time.Duration{
	FollowerState:  ElectionTimeout,
	CandidateState: CandidateTimeout,
	LeaderState:    HeartbeatTimeout}
