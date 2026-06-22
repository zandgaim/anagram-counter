package domain

import (
    "sort"
)

// RuneSlice attaches methods to []rune to satisfy sort.Interface
type RuneSlice []rune

func (p RuneSlice) Len() int           { return len(p) }
func (p RuneSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p RuneSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Signature generates a sorted string representation of a word to group anagrams.
// e.g., "listen" -> "eilnst"
func Signature(word string) string {
    runes := []rune(word)
    sort.Sort(RuneSlice(runes))
    return string(runes)
}

// GroupStats holds the count and a representative word for the anagram group
type GroupStats struct {
    ExampleWord string
    Count       int
}