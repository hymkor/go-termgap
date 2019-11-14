package lunewhydos

//go:generate go run script.go

func Query(n rune) int {
	if w, ok := table[n]; ok {
		return w
	} else {
		return -1
	}
}
