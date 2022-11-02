package main

import (
	"fmt"
	"math/bits"
	"reflect"
	"sort"
	"testing"

	"github.com/wata-gh/puyo2"
)

// func test(m map[int][]string, i int) []string {
// 	return m[i]
// }

func TestSetShapeAdjacent(t *testing.T) {
	sbf := puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa5123445112334223554")
	fcc := NewFieldColorCandidate([]puyo2.Color{puyo2.Red, puyo2.Blue, puyo2.Yellow, puyo2.Green}, sbf)

	if reflect.DeepEqual([]int{1}, fcc.ShapeAdjacent[0]) == false {
		t.Fatalf("ShapeAdjacent[0] must be [1] but %v", fcc.ShapeAdjacent[0])
	}
	expects := [][]int{
		{1},
		{0, 2},
		{1, 3, 4},
		{2, 4},
		{2, 3},
	}
	for i, expect := range expects {
		sort.Ints(fcc.ShapeAdjacent[i])
		if reflect.DeepEqual(expect, fcc.ShapeAdjacent[i]) == false {
			t.Fatalf("ShapeAdjacent[%d] must be %v but %v", i, expect, fcc.ShapeAdjacent[i])
		}
	}
}

func TestSetColorCandidate(t *testing.T) {
	sbf := puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa6a6aaa565123554112333226444")
	fcc := NewFieldColorCandidate([]puyo2.Color{puyo2.Red, puyo2.Blue, puyo2.Yellow, puyo2.Green}, sbf)

	expects := [][]puyo2.Color{
		{puyo2.Red, puyo2.Blue, puyo2.Yellow, puyo2.Green},
		{puyo2.Red, puyo2.Blue, puyo2.Yellow, puyo2.Green},
		{puyo2.Red, puyo2.Blue, puyo2.Yellow, puyo2.Green},
		{puyo2.Red, puyo2.Blue, puyo2.Yellow, puyo2.Green},
		{puyo2.Red, puyo2.Blue, puyo2.Yellow, puyo2.Green},
		{puyo2.Red, puyo2.Blue, puyo2.Yellow, puyo2.Green},
	}
	for i, expect := range expects {
		scc := fcc.ShapeColorCandidate[i]
		if reflect.DeepEqual(expect, scc.colors) == false {
			t.Fatalf("fcc.ShapeColorCandidate[%d] must be %v but %v", i, expect, scc.colors)
		}
	}

	fcc.SetColorCandidate(0, 13, []puyo2.Color{puyo2.Red})
	for i, expect := range expects {
		scc := fcc.ShapeColorCandidate[i]
		if reflect.DeepEqual(expect, scc.colors) == false {
			t.Fatalf("fcc.ShapeColorCandidate[%d] must be %v but %v", i, expect, scc.colors)
		}
	}
	cc := fcc.GetColorCandidate(0, 13)
	if reflect.DeepEqual([]puyo2.Color{puyo2.Red}, cc.colors) == false {
		t.Fatalf("fcc.GetColorCandidate(0, 13) must be %v but %v", []puyo2.Color{puyo2.Red}, cc.colors)
	}

	fcc.SetColorCandidate(0, 3, []puyo2.Color{puyo2.Red})
	expects = [][]puyo2.Color{
		{puyo2.Red},
		{puyo2.Blue, puyo2.Yellow, puyo2.Green},
		{puyo2.Red, puyo2.Blue, puyo2.Yellow, puyo2.Green},
		{puyo2.Red, puyo2.Blue, puyo2.Yellow, puyo2.Green},
		{puyo2.Red, puyo2.Blue, puyo2.Yellow, puyo2.Green},
		{puyo2.Red, puyo2.Blue, puyo2.Yellow, puyo2.Green},
	}
	for i, expect := range expects {
		scc := fcc.ShapeColorCandidate[i]
		if reflect.DeepEqual(expect, scc.colors) == false {
			fcc.ShowDebug()
			t.Fatalf("fcc.ShapeColorCandidate[%d] must be %v but %v", i, expect, scc.colors)
		}
	}

	fcc.SetColorCandidate(2, 1, []puyo2.Color{puyo2.Blue})
	expects = [][]puyo2.Color{
		{puyo2.Red},
		{puyo2.Yellow, puyo2.Green},
		{puyo2.Red, puyo2.Yellow, puyo2.Green},
		{puyo2.Red, puyo2.Yellow, puyo2.Green},
		{puyo2.Red, puyo2.Yellow, puyo2.Green},
		{puyo2.Blue},
	}
	for i, expect := range expects {
		scc := fcc.ShapeColorCandidate[i]
		if reflect.DeepEqual(expect, scc.colors) == false {
			fcc.ShowDebug()
			t.Fatalf("fcc.ShapeColorCandidate[%d] must be %v but %v", i, expect, scc.colors)
		}
	}
}

func TestFieldColorCandidate(t *testing.T) {
	sbf := puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa5123445112334223554")
	fcc := NewFieldColorCandidate([]puyo2.Color{puyo2.Red, puyo2.Blue, puyo2.Yellow, puyo2.Green}, sbf)
	fcc.ShowDebug()
	fcc.SetColorCandidate(0, 13, []puyo2.Color{puyo2.Red})
	fcc.ShowDebug()
	fcc.SetColorCandidate(0, 3, []puyo2.Color{puyo2.Red})
	fcc.ShowDebug()
	fmt.Printf("%v\n", fcc.GetColorCandidate(0, 2))
}

func TestSetColorCanidiateRemoveCandidate(t *testing.T) {
	target := puyo2.NewFieldBits()
	target.SetOnebit(4, 10)
	target.IterateBitWithMasking(func(fb *puyo2.FieldBits) *puyo2.FieldBits {
		x := 0
		y := 0
		for x = 0; x < 6; x++ {
			col := fb.ColBits(x)
			if col > 0 {
				sb := x
				if x > 3 {
					sb = 4 - x
				}
				y = bits.Len64(col>>(16*sb)) - 1
				break
			}
		}
		fmt.Printf("x: %d, y %d\n", x, y)
		return fb
	})
}
