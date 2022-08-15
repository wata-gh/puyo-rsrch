package main

import (
	"fmt"
	"strings"
)

type Result struct {
	results [][]int
}

func (r Result) Equals(r2 Result) bool {
	for i, result := range r.results {
		for j, v := range result {
			if v != r2.results[i][j] {
				return false
			}
		}
	}
	return true
}

func (r Result) Len() int {
	return len(r.results)
}

func (r Result) Swap(i, j int) {
	r.results[i], r.results[j] = r.results[j], r.results[i]
}

func (r Result) Less(i, j int) bool {
	if r.results[i][0] < r.results[j][0] {
		return true
	}
	if r.results[i][0] == r.results[j][0] {
		return r.results[i][1] < r.results[j][1]
	}
	return false
}

func (r Result) ToString() string {
	var b strings.Builder
	for i, result := range r.results {
		if i != 0 {
			fmt.Fprintf(&b, ":")
		}
		fmt.Fprintf(&b, "%d_%d", result[0], result[1])
	}
	return b.String()
}

type Results []Result

func (r Results) Len() int {
	return len(r)
}

func (r Results) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r Results) Less(i, j int) bool {
	ri := r[i]
	rj := r[j]
	for i, r := range ri.results {
		if r[0] == rj.results[i][0] {
			if r[1] == rj.results[i][1] {
				continue
			} else {
				return r[1] < rj.results[i][1]
			}
		}
		return r[0] < rj.results[i][0]
	}
	return false
}
