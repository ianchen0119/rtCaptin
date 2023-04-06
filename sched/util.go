package sched

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func remove(slice []*Job, s int) []*Job {
	if len(slice) == 1 {
		return []*Job{}
	}
	return append(slice[:s], slice[s+1:]...)
}
