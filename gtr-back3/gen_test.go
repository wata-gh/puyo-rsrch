package main

import (
	"testing"
)

func TestValidEmpty(t *testing.T) {
	if validEmpty([]int{0, 1, 2}) == false {
		t.Error("validEmpty() All 4th row must be valid")
	}
	if validEmpty([]int{1, 2, 3}) == false {
		t.Error("validEmpty() All 4th row must be valid")
	}
	if validEmpty([]int{0, 2, 3}) == false {
		t.Error("validEmpty() All 4th row must be valid")
	}
	if validEmpty([]int{0, 2, 4}) == false {
		t.Error("validEmpty() 3rd row with 4th row must be valid")
	}
	if validEmpty([]int{2, 3, 7}) == false {
		t.Error("validEmpty() 3rd row with 4th row must be valid")
	}
	if validEmpty([]int{0, 3, 7}) == false {
		t.Error("validEmpty() 3rd row with 4th row must be valid")
	}
	if validEmpty([]int{1, 5, 8}) == false {
		t.Error("validEmpty() 2nd row with 3rd and 4th row must be valid")
	}
	if validEmpty([]int{3, 7, 10}) == false {
		t.Error("validEmpty() 2nd row with 3rd and 4th row must be valid")
	}
	if validEmpty([]int{0, 4, 8}) {
		t.Errorf("validEmpty() 2nd row without 3rd and 4th row must be invalid")
	}
	if validEmpty([]int{0, 5, 8}) {
		t.Errorf("validEmpty() 2nd row with 3rd row but without 4th row must be invalid")
	}
}
