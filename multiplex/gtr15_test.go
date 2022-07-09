package main

import (
	"testing"
)

func TestValidPlace(t *testing.T) {
	list := []int{4, 8, 11, 12}
	p := Gtr15Pattern{}
	if p.ValidPlace(list) == false {
		t.Fatal("ValidPlace must be true.")
	}

	list = []int{5, 9, 10, 13}
	p = Gtr15Pattern{}
	if p.ValidPlace(list) == false {
		t.Fatal("ValidPlace must be true.")
	}

	list = []int{3, 6, 7, 14}
	p = Gtr15Pattern{}
	if p.ValidPlace(list) == false {
		t.Fatal("ValidPlace must be true.")
	}
}
