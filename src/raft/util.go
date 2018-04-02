package raft

import "log"

// Debugging
const Debug = 1
const ShowVerbose = ShowProcess
const ShowVariable = 1
const ShowProcess = 1 << 1

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}

func LevelDPrintf(format string, flag int, a ...interface{}) (n int, err error) {
	if ShowVerbose&flag > 0 {
		return DPrintf(format, a...)
	}
	return
}
