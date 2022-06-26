package main

import (
	"testing"

	"github.com/wata-gh/puyo2"
)

func TestMain(t *testing.T) {
	puyo2.NewFieldBitsWithM([2]uint64{5066575350595584, 1048594}).ShowDebug()
}

func TestCreateGtr(t *testing.T) {
	bf := puyo2.NewBitFieldWithMattulwan("a61ba6ba3b2a4")
	createGtr(bf)
	c := bf.Color(0, 2)
	if c != puyo2.Green {
		t.Fatalf("color must be green. %v", c)
	}
	bf = puyo2.NewBitFieldWithMattulwan("a61ca6ca3c2a4")
	createGtr(bf)
	c = bf.Color(0, 2)
	if c != puyo2.Red {
		t.Fatalf("color must be red. %v", c)
	}
}

func TestCheckAvailColor(t *testing.T) {
	bf := puyo2.NewBitFieldWithMattulwan("a62da6da4da3")
	c := checkAvailColor(bf)
	if c != puyo2.Red {
		t.Fatalf("color must be red. %v", c)
	}
	bf = puyo2.NewBitFieldWithMattulwan("a62ba6ca4ea3")
	c = checkAvailColor(bf)
	if c != puyo2.Yellow {
		t.Fatalf("color must be yellow. %v", c)
	}

}
