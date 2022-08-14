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
	Eighth
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
		// column 1's eighth not permited. because col1 must connect to gtr.
		if i != 0 && os.Onebit(i, n+1) == 0 {
			return Eighth
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

func color(sbf *puyo2.ShapeBitField, targets []puyo2.Color, idx int) bool {
	if idx == sbf.ShapeCount() {
		sbfArray := sbf.ToChainShapesUInt64Array()

		bf := puyo2.NewBitField()
		for i, shape := range sbf.Shapes {
			for x := 0; x < 6; x++ {
				for y := 0; y < 14; y++ {
					if shape.Onebit(x, y) > 0 {
						bf.SetColor(targets[i], x, y)
					}
				}
			}
		}
		cbf := bf.Clone()
		result := cbf.Simulate()
		if result.Chains == sbf.ShapeCount() {
			sameChain := true
			for i, array := range bf.ToChainShapesUInt64Array() {
				if sbfArray[i][0] != array[0] || sbfArray[i][1] != array[1] {
					sameChain = false
					break
				}
			}
			if sameChain {
				// fmt.Println(bf.MattulwanEditorParam())
				return true
				// fmt.Println(csbf.ChainOrderedFieldString())
				// sbf.ShowDebug()
				// bf.ShowDebug()
			} else {
				fmt.Println("Whats!!!!!!!!!", bf.MattulwanEditorParam())
				return false
			}
		}
		return false
	}

	colors := []puyo2.Color{puyo2.Red, puyo2.Blue, puyo2.Green, puyo2.Yellow}
	rem := idx % len(colors)
	part := colors[0:rem]
	colors = colors[rem:]
	colors = append(colors, part...)
	for _, c := range colors {
		targets[idx] = c
		if color(sbf, targets, idx+1) {
			return true
		}
	}
	return false
}

func colorable(sbf *puyo2.ShapeBitField) bool {
	targets := make([]puyo2.Color, sbf.ShapeCount())
	return color(sbf, targets, 0)
}

func colorable2(sbf *puyo2.ShapeBitField) bool {
	csbf := sbf.Clone()
	csbf.Simulate()
	sbfArray := sbf.ToChainShapesUInt64Array()
	colors := []puyo2.Color{puyo2.Red, puyo2.Blue, puyo2.Green, puyo2.Yellow}
	targets := [6]puyo2.Color{puyo2.Empty, puyo2.Empty, puyo2.Empty, puyo2.Empty, puyo2.Empty, puyo2.Empty}
	for _, c1 := range colors {
		targets[0] = c1
		for _, c2 := range colors {
			targets[1] = c2
			for _, c3 := range colors {
				targets[2] = c3
				for _, c4 := range colors {
					targets[3] = c4
					for _, c5 := range colors {
						targets[4] = c5
						for _, c6 := range colors {
							targets[5] = c6
							bf := puyo2.NewBitField()
							for i, shape := range sbf.Shapes {
								for x := 0; x < 6; x++ {
									for y := 0; y < 13; y++ {
										if shape.Onebit(x, y) > 0 {
											bf.SetColor(targets[i], x, y)
										}
									}
								}
							}
							cbf := bf.Clone()
							result := cbf.Simulate()
							if result.Chains == sbf.ShapeCount() {
								sameChain := true
								for i, array := range bf.ToChainShapesUInt64Array() {
									if sbfArray[i][0] != array[0] || sbfArray[i][1] != array[1] {
										sameChain = false
										break
									}
								}
								if sameChain {
									// fmt.Println(bf.MattulwanEditorParam())
									return true
									// fmt.Println(csbf.ChainOrderedFieldString())
									// sbf.ShowDebug()
									// bf.ShowDebug()
								} else {
									fmt.Println("Whats!!!!!!!!!", bf.MattulwanEditorParam())
								}
							}
						}
					}
				}
			}
		}
	}
	// fmt.Println("not colorable")
	return false
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

func main2() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		param := strings.TrimRight(scanner.Text(), "\n")
		sbf := puyo2.NewShapeBitFieldWithFieldString(param)
		sbfc := sbf.Clone()

		r := adjacentColorCount(sbfc)
		s := fireable(sbfc)
		// fmt.Println("fireable", s)
		sbfc = sbf.Clone()
		if connectable(sbfc) && s != UnFireable {
			// fmt.Println("connectable && fireable")
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
}
