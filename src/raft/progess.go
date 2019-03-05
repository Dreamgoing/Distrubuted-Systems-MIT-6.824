package raft

type Progress struct {
	Match, Next uint64
}
