package raft

type Entry struct {
	term    int
	index   int
	command interface{}
}
