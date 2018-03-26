package raft

type Entry struct {
	term    int
	command interface{}
}
