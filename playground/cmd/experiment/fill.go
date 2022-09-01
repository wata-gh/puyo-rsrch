package main

import (
	"fmt"
	"math/bits"
	"os"
	"sync"

	"github.com/wata-gh/puyo2"
	"gonum.org/v1/gonum/stat/combin"
)

func FillSearch(fields chan []int, wg *sync.WaitGroup) {
	for {
		field := <-fields
		if len(field) == 0 {
			break
		}
		fmt.Fprintln(os.Stderr, field)
		patterns := fill2(field) // <- CHANGE HERE
		for _, pattern := range patterns {
			Fill(pattern)
		}
	}
	wg.Done()
}

func shapes(c int, x int) []*puyo2.FieldBits {
	var shapes []*puyo2.FieldBits
	switch c {
	case 0:
		s := puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+2, 1)
		s.SetOnebit(x+3, 1)
	case 1:
		s := puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x, 2)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+2, 1)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x, 2)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+2, 2)
		shapes = append(shapes, s)
	case 2:
		s := puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+2, 1)
		s.SetOnebit(x+2, 2)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 2)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+2, 1)
		s.SetOnebit(x+2, 2)
		shapes = append(shapes, s)
	case 3:
		s := puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+2, 1)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 2)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+2, 2)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 2)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+2, 1)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+2, 2)
		shapes = append(shapes, s)
	case 4:
		s := puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x, 2)
		s.SetOnebit(x, 3)
		s.SetOnebit(x+1, 1)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x, 2)
		s.SetOnebit(x, 3)
		s.SetOnebit(x+1, 2)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x, 2)
		s.SetOnebit(x, 3)
		s.SetOnebit(x+1, 3)
		shapes = append(shapes, s)
	case 5:
		s := puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+1, 3)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 2)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+1, 3)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 3)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+1, 3)
		shapes = append(shapes, s)
	case 6:
		s := puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x, 2)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+1, 2)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x, 2)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+1, 3)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 2)
		s.SetOnebit(x, 3)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+1, 2)
		shapes = append(shapes, s)
	case 7:
		s := puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x, 2)
		s.SetOnebit(x, 3)
		s.SetOnebit(x, 4)
		shapes = append(shapes, s)
	}
	return shapes
}

func willNotDrop(y int, shape *puyo2.FieldBits, overall *puyo2.FieldBits) bool {
	for x := 0; x < 6; x++ {
		if shape.Onebit(x, y+1) > 0 { // shape あり
			if overall.Onebit(x, y) == 0 { // 下なし
				return false
			}
		} else {
			col := shape.ColBits(x)
			if col > 0 {
				for z := y; z < bits.TrailingZeros64(col>>(x*16)); z++ {
					if overall.Onebit(x, z) == 0 {
						return false
					}
				}
			}
		}
	}
	return true
}

func place(osbf *puyo2.ShapeBitField, perm []int, clusters [][]int, last *puyo2.FieldBits) {
	if len(perm) == 0 {
		nsbf := osbf.Clone()
		result := nsbf.Simulate()
		if result.Chains == 2 { // CHANGE HERE
			fmt.Println(nsbf.ChainOrderedFieldString())
			// osbf.ShowDebug()
			// } else {
			// 	fmt.Println("--> NG: ", osbf.FieldString())
			// 	fmt.Printf("%+v\n", result)
			// 	osbf.ShowDebug()
			// 	for _, shape := range osbf.Shapes {
			// 		shape.ShowDebug()
			// 	}
		}
		return
	}
	c1 := clusters[perm[0]]
	overall := osbf.OverallShape()
	for x := 0; x < 6; x++ {
		overall.SetOnebit(x, 0)
	}
	for _, shape := range shapes(c1[0], c1[1]) {
		for yOffset := 0; yOffset < 13; yOffset++ {
			s := shape.FastLift(yOffset)
			if s.Equals(s.MaskField13()) == false {
				continue
			}
			if willNotDrop(yOffset, s, overall) {
				sfb := osbf.Clone()
				sfb.InsertShape(s)
				place(sfb, perm[1:], clusters, shape)
			}
		}
	}
}

func Fill(clusters [][]int) {
	perm := combin.Permutations(len(clusters), len(clusters))
	for _, p := range perm {
		sbf := puyo2.NewShapeBitField()
		place(sbf, p, clusters, puyo2.NewFieldBits())
	}
}
