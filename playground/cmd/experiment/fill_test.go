package main

import (
	"testing"

	"github.com/wata-gh/puyo2"
)

func TestWillNotDrop(t *testing.T) {
	shape := puyo2.NewFieldBits()
	shape.SetOnebit(0, 1)
	shape.SetOnebit(0, 2)
	shape.SetOnebit(0, 3)
	shape.SetOnebit(1, 2)

	overall := puyo2.NewFieldBits()
	for x := 0; x < 6; x++ {
		overall.SetOnebit(x, 0)
	}

	if willNotDrop(0, shape, overall) {
		t.Fatal("willNotDrop must be false")
	}
}
