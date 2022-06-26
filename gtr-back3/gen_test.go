package main

import (
	"testing"
)

func TestValidEmpty(t *testing.T) {
	if validEmpty([]int{0, 1, 2}) == false {
		t.Error("validEmpty() All 4th row must be valid")
	}
	if validEmpty([]int{3}) == false {
		t.Error("validEmpty() 3 must be valid")
	}
	if validEmpty([]int{3, 7}) == false {
		t.Error("validEmpty() 3,7 must be valid")
	}
	if validEmpty([]int{0, 4, 8}) == false {
		t.Error("validEmpty() 0,4,8 row must be valid")
	}
	if validEmpty([]int{14}) {
		t.Error("validEmpty() 14 must be invalid")
	}
	if validEmpty([]int{0, 4, 8, 11, 15}) == false {
		t.Error("validEmpty() 0,4,8,11,15 must be valid")
	}
	if validEmpty([]int{0, 5}) {
		t.Error("validEmpty() 0,5 must be invalid")
	}
	if validEmpty([]int{2, 6, 10, 13, 17}) == false {
		t.Error("validEmpty() 2,6,10,13,17 must be valid")
	}
}
