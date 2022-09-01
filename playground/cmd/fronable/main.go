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

func fronConnectable(sbf *puyo2.ShapeBitField) bool {
	for i := 1; i < len(sbf.Shapes); i++ {
		sbf.Simulate1()
	}
	last := sbf.Shapes[len(sbf.Shapes)-1]
	return last.Onebit(0, 1) != 0 || last.Onebit(1, 1) != 0
}

func fronFireable(sbf *puyo2.ShapeBitField) Status {
	os := sbf.OverallShape()
	sbfc := sbf.Clone()
	sbfc.Simulate()
	first := sbfc.ChainOrderedShapes[0][0]
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
			if i == 0 {
				return Eighth1
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

func fronableBase(sbf *puyo2.ShapeBitField) [2]*puyo2.FieldBits {
	s1 := puyo2.NewFieldBits()
	s1.SetOnebit(0, 1)
	s1.SetOnebit(0, 2)
	s1.SetOnebit(1, 1)
	sbf.InsertShape(s1)
	s2 := puyo2.NewFieldBits()
	s2.SetOnebit(0, 1)
	s2.SetOnebit(1, 1)
	s2.SetOnebit(2, 1)
	s2.SetOnebit(3, 1)
	sbf.InsertShape(s2)
	s3 := puyo2.NewFieldBits()
	s3.SetOnebit(1, 2)
	s3.SetOnebit(2, 1)
	s3.SetOnebit(2, 2)
	sbf.InsertShape(s3)
	return [2]*puyo2.FieldBits{s1, s3}
}

func checkFron(params chan string, wg *sync.WaitGroup) {
	for {
		param := <-params
		if param == "" {
			break
		}
		sbf := puyo2.NewShapeBitFieldWithFieldString(param)
		sbfb := sbf.Clone()

		r := adjacentColorCount(sbfb)
		sbfb = sbf.Clone()
		if fronConnectable(sbfb) {
			sbfb = sbf.Clone()
			fronableBase(sbfb)
			overall := sbfb.OverallShape()
			overall.SetOnebit(0, 0)
			overall.SetOnebit(1, 0)
			y1 := bits.Len64(overall.ColBits(0))
			y2 := bits.Len64(overall.ColBits(1) >> 16)

			keyPatterns := [][2][2]int{}
			if y1 <= 12 {
				keyPatterns = append(keyPatterns, [2][2]int{{0, y1}, {0, y1 + 1}})
			}
			keyPatterns = append(keyPatterns, [2][2]int{{0, y1}, {1, y2}})
			keyPatterns = append(keyPatterns, [2][2]int{{1, y2}, {0, y1}})
			if y2 <= 12 {
				keyPatterns = append(keyPatterns, [2][2]int{{1, y2}, {1, y2 + 1}})
			}
			found := false
			for _, keyPattern := range keyPatterns {
				sbfc := sbfb.Clone()
				s1 := sbfc.Shapes[len(sbfc.Shapes)-3]
				s3 := sbfc.Shapes[len(sbfc.Shapes)-1]
				s1.SetOnebit(keyPattern[0][0], keyPattern[0][1])
				s3.SetOnebit(keyPattern[1][0], keyPattern[1][1])
				s := fronFireable(sbfc)
				if s != UnFireable && colorable(sbfc) {
					fmt.Println(param, s, r)
					found = true
					break
				}
			}
			if found == false {
				fmt.Fprintf(os.Stderr, "%s\n", param)
			}
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", param)
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
		go checkFron(params, &wg)
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
