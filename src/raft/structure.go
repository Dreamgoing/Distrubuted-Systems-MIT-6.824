package raft

type Entry struct {
	Term    int
	Index   int
	Command interface{}
}
