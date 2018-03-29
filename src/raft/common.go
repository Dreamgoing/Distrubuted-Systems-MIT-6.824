package raft

func Min(a int, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func MajorityOk(ok, total int) bool {
	return ok > total/2
}
