package sched

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func remove(slice []interface{}, s int) []interface{} {
	return append(slice[:s], slice[s+1:]...)
}
