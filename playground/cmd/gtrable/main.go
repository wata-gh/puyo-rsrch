package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
	"sync"

	"github.com/wata-gh/puyo2"
)

type Status int

const (
	UnFireable Status = iota
	Normal
	Eighth1
	Eighth2
	Eighth3
)

func connectable(sbf *puyo2.ShapeBitField) bool {
	for i := 1; i < len(sbf.Shapes); i++ {
		sbf.Simulate1()
	}
	last := sbf.Shapes[len(sbf.Shapes)-1]
	return last.Onebit(0, 1) != 0
}

func fireable(sbf *puyo2.ShapeBitField) Status {
	os := sbf.OriginalOverallShape()
	first := sbf.ChainOrderedShapes[0][0]
	for i := 0; i < 3; i++ {
		cb := first.ColBits(i)
		cb >>= i * 16
		n := bits.Len64(cb)
		if n == 0 {
			continue
		}
		// upper space is empty.
		if os.Onebit(i, n) == 0 {
			return Normal
		}
		if os.Onebit(i, n+1) == 0 {
			// 1st row's eighth not permitted.
			if i == 0 {
				continue
			}
			if i == 1 {
				return Eighth2
			}
			return Eighth3
		}
	}
	return UnFireable
}

func adjacentColorCount(sbf *puyo2.ShapeBitField) [3]map[int]struct{} {
	cnts := [3]map[int]struct{}{}

	for sbf.IsEmpty() == false {
		for c := 0; c < 3; c++ {
			n := sbf.ShapeNum(c, 1)
			if n == -1 {
				continue
			}
			if cnts[c] == nil {
				cnts[c] = map[int]struct{}{}
			}
			cnts[c][n] = struct{}{}
		}
		sbf.Simulate1()
	}
	return cnts
}

func colorable(sbf *puyo2.ShapeBitField) bool {
	bf := sbf.FillChainableColor()
	return bf != nil
}

func gtrable(sbf *puyo2.ShapeBitField) bool {
	overall := sbf.OverallShape()
	overall.SetOnebit(0, 0)
	y := bits.Len64(overall.ColBits(0))
	s1 := puyo2.NewFieldBits()
	s1.SetOnebit(2, 1)
	s1.SetOnebit(2, 2)
	s1.SetOnebit(3, 1)
	s1.SetOnebit(3, 2)
	sbf.InsertShape(s1)
	s2 := puyo2.NewFieldBits()
	s2.SetOnebit(0, 1)
	s2.SetOnebit(1, 1)
	s2.SetOnebit(1, 2)
	s2.SetOnebit(2, 2)
	sbf.InsertShape(s2)
	s3 := puyo2.NewFieldBits()
	s3.SetOnebit(0, 2)
	s3.SetOnebit(0, 3)
	s3.SetOnebit(1, 2)
	sbf.InsertShape(s3)
	s3.SetOnebit(0, y+3)
	return true
	// csbf := sbf.Clone()
	// result := csbf.Simulate()
	// return result.Chains == sbf.ShapeCount()
}

func check(params chan string, wg *sync.WaitGroup) {
	for {
		param := <-params
		if param == "" {
			break
		}
		sbf := puyo2.NewShapeBitFieldWithFieldString(param)
		sbfc := sbf.Clone()

		r := adjacentColorCount(sbfc)
		s := fireable(sbfc)
		sbfc = sbf.Clone()
		if connectable(sbfc) && s != UnFireable {
			sbfc = sbf.Clone()
			if gtrable(sbfc) {
				if colorable(sbfc) {
					fmt.Println(param, s, r)
				} else {
					fmt.Fprintf(os.Stderr, "%s\n", param)
				}
			}
		}
	}
	wg.Done()
}

func main() {
	params := make(chan string)
	var wg sync.WaitGroup
	grc := 8
	for i := 0; i < grc; i++ {
		wg.Add(1)
		go check(params, &wg)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		params <- strings.TrimRight(scanner.Text(), "\n")
	}
	for i := 0; i < grc; i++ {
		params <- ""
	}
	wg.Wait()
}
