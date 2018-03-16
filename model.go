package main

var counts map[int]int
var tokens map[string]int
var itemStrings map[int]string

type byCounts []int

func (s byCounts) Len() int {
	return len(s)
}
func (s byCounts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byCounts) Less(i, j int) bool {
	return counts[s[i]] > counts[s[j]]
}
