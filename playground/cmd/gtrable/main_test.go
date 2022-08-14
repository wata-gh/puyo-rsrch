package main

import (
	"fmt"
	"testing"

	"github.com/wata-gh/puyo2"
)

func TestShift(t *testing.T) {
	for i := 0; i < 6; i++ {
		array := []int{1, 2, 3, 4}
		remainder := i % 4
		part := array[0:remainder]
		array = array[remainder:]
		array = append(array, part...)
		fmt.Println(array)
	}
}

func TestGtrable(t *testing.T) {
	sbf := puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaa3aaaaa2aaaaa2aaaaa2aaaaa1aaaaa1aaaaa11aaaa23aaaa33aaaa")
	// sbf.ShowDebug()
	gtrable(sbf)
	len := len(sbf.Shapes)
	tmp := sbf.Shapes[len-3]
	sbf.Shapes[len-3] = sbf.Shapes[len-1]
	sbf.Shapes[len-1] = tmp
	result := sbf.Clone().Simulate()
	fmt.Println(result)
	r := colorable(sbf)
	fmt.Println(r)
	sbf.ShowDebug()
}

func TestFireable(t *testing.T) {
	sbf := puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa2aaaa22aaa333aaa211aaa311aaa")
	sbf.ShowDebug()
	sbf.Simulate()
	s := fireable(sbf)
	if s != UnFireable {
		t.Fatalf("must be unfireable. %d", s)
	}

	sbf = puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa3aaaaa2aaaaa2aaaaa2aaaaa1aaaaa1aaaaa12aaaa13aaaa33aaaa")
	sbf.ShowDebug()
	sbf.Simulate()
	s = fireable(sbf)
	if s != UnFireable {
		t.Fatalf("must be unfireable. %d", s)
	}

	sbf = puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa3aaaaa2aaaaa2aaaaa2aaaaa1aaaaa1aaaaa11aaaa23aaaa33aaaa")
	sbf.ShowDebug()
	sbf.Simulate()
	s = fireable(sbf)
	if s != Normal {
		t.Fatalf("must be fireable. %d", s)
	}

	sbf = puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa3aaaaa3aaaaa2aaaaa1aaaaa12aaaa11aaaa32aaaa32aaaa")
	sbf.ShowDebug()
	sbf.Simulate()
	s = fireable(sbf)
	if s != Eighth {
		t.Fatalf("must be fireable from eighth. %d", s)
	}
}

// func TestAdjacentColorCount(t *testing.T) {
// 	sbf := puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa3aaaaa2aaaaa2aaaaa2aaaaa1aaaaa1aaaaa11aaaa23aaaa33aaaa")
// 	sbf.ShowDebug()
// 	r := adjacentColorCount(sbf)
// 	if r[0] != 1 || r[1] != 1 || r[2] != 0 {
// 		t.Fatalf("result must be [1 1 0] but %v", r)
// 	}
// 	sbf = puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa2aaaaa11aaaa133aaa123aaa322aaa")
// 	sbf.ShowDebug()
// 	r = adjacentColorCount(sbf)
// 	if r[0] != 1 || r[1] != 2 || r[2] != 2 {
// 		t.Fatalf("result must be [1 2 2] but %v", r)
// 	}
// 	sbf = puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa3aaaa22aaaa123aaa112aaa313aaa")
// 	sbf.ShowDebug()
// 	r = adjacentColorCount(sbf)
// 	if r[0] != 1 || r[1] != 3 || r[2] != 1 {
// 		t.Fatalf("result must be [1 3 1] but %v", r)
// 	}
// }
