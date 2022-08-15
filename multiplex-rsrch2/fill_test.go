package main

import (
	"fmt"
	"testing"
)

func TestFillUp(t *testing.T) {
	results := fillUp([]int{3, 3, 2, 0, 0, 0}, 2)
	fmt.Println(results)
	removeDuplication(results)
}

func TestFill(t *testing.T) {
	var pattern Result
	pattern.results = [][]int{{1, 0, 2, 1, 1}, {3, 0, 1, 2, 1}}
	Fill(pattern, 2)
}
