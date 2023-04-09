package sched

func min(a, b int) (bool, int) {
	if a < b {
		return true, a
	}
	return false, b
}

func remove(slice []*Job, s int) []*Job {
	if len(slice) == 1 {
		return []*Job{}
	}
	return append(slice[:s], slice[s+1:]...)
}
