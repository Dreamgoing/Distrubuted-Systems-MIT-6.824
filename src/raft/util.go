package raft

import "log"

// Debugging
const Debug = 1

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}

func Min(a int, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}
