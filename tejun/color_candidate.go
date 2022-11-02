package main

import (
	"fmt"
	"math/bits"
	"strings"

	"github.com/wata-gh/puyo2"
)

type ColorCandidate struct {
	colors []puyo2.Color
	table  map[puyo2.Color]struct{}
}

type FieldColorCandidate struct {
	colors              []puyo2.Color
	BitField            *puyo2.BitField
	ColorCandidate      map[[2]int]*ColorCandidate
	ChigiriCount        int
	ShapeColorCandidate []*ColorCandidate
	ShapeBitField       *puyo2.ShapeBitField
	ShapeAdjacent       map[int][]int
}

func NewColorCandidate(colors []puyo2.Color) *ColorCandidate {
	cc := new(ColorCandidate)
	cc.colors = colors
	cc.resetColorTable()
	return cc
}

func (cc *ColorCandidate) resetColorTable() {
	cc.table = map[puyo2.Color]struct{}{}
	for _, c := range cc.colors {
		cc.table[c] = struct{}{}
	}
}

func (cc *ColorCandidate) Clone() *ColorCandidate {
	ncc := new(ColorCandidate)
	ncc.colors = make([]puyo2.Color, len(cc.colors))
	copy(ncc.colors, cc.colors)
	ncc.resetColorTable()
	return ncc
}

func (cc *ColorCandidate) Contains(c puyo2.Color) bool {
	_, e := cc.table[c]
	return e
}

func NewFieldColorCandidate(colors []puyo2.Color, sbf *puyo2.ShapeBitField) *FieldColorCandidate {
	fcc := new(FieldColorCandidate)
	fcc.ColorCandidate = map[[2]int]*ColorCandidate{}
	fcc.colors = colors
	fcc.ShapeBitField = sbf

	fcc.setShapeAdjacent()

	fcc.ShapeColorCandidate = make([]*ColorCandidate, sbf.ShapeCount())
	for i := 0; i < sbf.ShapeCount(); i++ {
		sc := make([]puyo2.Color, len(colors))
		copy(sc, colors)
		fcc.ShapeColorCandidate[i] = NewColorCandidate(sc)
	}
	for x := 0; x < 6; x++ {
		for y := 1; y < 15; y++ {
			n := sbf.ShapeNum(x, y)
			if n == -1 {
				sc := make([]puyo2.Color, len(colors))
				copy(sc, colors)
				fcc.ColorCandidate[[2]int{x, y}] = NewColorCandidate(sc)
			} else {
				fcc.ColorCandidate[[2]int{x, y}] = fcc.ShapeColorCandidate[n]
			}
		}
	}
	return fcc
}

func (fcc *FieldColorCandidate) setShapeAdjacent() {
	fcc.ShapeAdjacent = map[int][]int{}
	csbf := fcc.ShapeBitField.Clone()
	for i, shape := range csbf.Shapes {
		if shape.PopCount() == 3 {
			overall := csbf.OverallShape()
			overall.M[0] = ^overall.M[0]
			overall.M[1] = ^overall.M[1]
			overall = overall.MaskField13()
			expand := shape.Expand1(overall)
			csbf.Shapes[i] = shape.Or(expand)
		}
	}
	adjacentMap := map[int]map[int]struct{}{}
	for {
		for i, s1 := range csbf.Shapes {
			for j, s2 := range csbf.Shapes {
				if i == j {
					continue
				}
				if s1.Expand1(s2).IsEmpty() {
					continue
				}
				if _, ok := adjacentMap[i]; ok == false {
					adjacentMap[i] = map[int]struct{}{}
				}
				adjacentMap[i][j] = struct{}{}
			}
		}
		vfbn := csbf.Simulate1()
		if len(vfbn) == 0 {
			break
		}
	}

	for k, m := range adjacentMap {
		for v, _ := range m {
			fcc.ShapeAdjacent[k] = append(fcc.ShapeAdjacent[k], v)
		}
	}
}

func (fcc *FieldColorCandidate) Clone() *FieldColorCandidate {
	nfcc := new(FieldColorCandidate)

	nfcc.colors = fcc.colors

	if fcc.BitField != nil {
		nfcc.BitField = fcc.BitField.Clone()
	}

	nfcc.ColorCandidate = map[[2]int]*ColorCandidate{}
	for k, v := range fcc.ColorCandidate {
		nfcc.ColorCandidate[k] = v.Clone()
	}

	nfcc.ChigiriCount = fcc.ChigiriCount

	nfcc.ShapeAdjacent = fcc.ShapeAdjacent

	nfcc.ShapeColorCandidate = make([]*ColorCandidate, fcc.ShapeBitField.ShapeCount())
	for i, cc := range fcc.ShapeColorCandidate {
		ncc := make([]puyo2.Color, len(cc.colors))
		copy(ncc, cc.colors)
		nfcc.ShapeColorCandidate[i] = NewColorCandidate(ncc)
	}
	nfcc.ShapeBitField = fcc.ShapeBitField.Clone()

	return nfcc
}

func (fcc *FieldColorCandidate) CountPlaced() int {
	fb := fcc.BitField.Bits(puyo2.Empty)
	fb.M[0] = ^fb.M[0]
	fb.M[1] = ^fb.M[1]
	return fcc.ShapeBitField.OverallShape().And(fb).PopCount()
}

func (fcc *FieldColorCandidate) GetColorCandidate(x int, y int) *ColorCandidate {
	if x < 0 || x >= 6 || y < 1 || y > 14 {
		panic(fmt.Sprintf("no such position. x => %d y => %d.", x, y))
	}
	n := fcc.ShapeBitField.ShapeNum(x, y)
	if n != -1 {
		return fcc.ShapeColorCandidate[n]
	}
	return fcc.ColorCandidate[[2]int{x, y}]
}

func (fcc *FieldColorCandidate) RemoveColorCandidate(x int, y int, colors []puyo2.Color) {
	n := fcc.ShapeBitField.ShapeNum(x, y)
	if n != -1 { // In-Shape
		scc := fcc.ShapeColorCandidate[n]
		ncc := []puyo2.Color{}
		for _, c := range colors {
			if scc.Contains(c) {
				continue
			}
			ncc = append(ncc, c)
		}
		scc.colors = append(scc.colors[:0], ncc...)
		scc.resetColorTable()
	} else { // Outer-Shape
		rcc := NewColorCandidate(colors)
		cc := fcc.ColorCandidate[[2]int{x, y}]
		ncc := []puyo2.Color{}
		for _, c := range cc.colors {
			if rcc.Contains(c) {
				continue
			}
			ncc = append(ncc, c)
		}
		cc.colors = append(cc.colors[:0], ncc...)
		cc.resetColorTable()
	}
}

func (fcc *FieldColorCandidate) SetColorCandidate(x int, y int, colors []puyo2.Color) {
	cc := NewColorCandidate(colors)
	n := fcc.ShapeBitField.ShapeNum(x, y)
	if n != -1 { // In-Shape
		fcc.ShapeColorCandidate[n].colors = append(fcc.ShapeColorCandidate[n].colors[:0], colors...)
		fcc.ShapeColorCandidate[n].resetColorTable()

		for _, i := range fcc.ShapeAdjacent[n] {
			scc := fcc.ShapeColorCandidate[i]
			ncc := []puyo2.Color{}
			for _, c := range scc.colors {
				if cc.Contains(c) {
					continue
				}
				ncc = append(ncc, c)
			}
			scc.colors = append(scc.colors[:0], ncc...)
			scc.resetColorTable()
		}

		nshape := fcc.ShapeBitField.Shapes[n]
		overall := fcc.ShapeBitField.OverallShape()
		overall.M[0] = ^overall.M[0]
		overall.M[1] = ^overall.M[1]
		nshape.Expand1(overall.MaskField13()).IterateBitWithMasking(func(fb *puyo2.FieldBits) *puyo2.FieldBits {
			x := 0
			y := 0
			for x = 0; x < 6; x++ {
				col := fb.ColBits(x)
				if col > 0 {
					sb := x
					if x > 3 {
						sb = x - 4
					}
					y = bits.Len64(col>>(16*sb)) - 1
					break
				}
			}
			fcc.RemoveColorCandidate(x, y, colors)
			return fb
		})
	} else { // Outer-Shape
		cc := fcc.ColorCandidate[[2]int{x, y}]
		cc.colors = append(cc.colors[:0], colors...)
		cc.resetColorTable()
	}
}

func (fcc *FieldColorCandidate) ShowDebug() {
	bitFieldStrings := []string{"14:", "13:", "12:", "11:", "10:", "09:", "08:", "07:", "06:", "05:", "04:", "03:", "02:", "01:", "00:"}
	if fcc.BitField != nil {
		bitFieldStrings = strings.Split(fcc.BitField.ToString(), "\n")
	}
	var b strings.Builder
	var ccb strings.Builder
	for y := 13; y > 0; y-- {
		var l strings.Builder
		fmt.Fprintf(&b, "%s ", bitFieldStrings[14-y])
		for x := 0; x < 6; x++ {
			n := fcc.ShapeBitField.ShapeNum(x, y)
			if n != -1 {
				fmt.Fprintf(&b, "%d", n)
				fmt.Fprint(&l, "_")
			} else {
				cc := fcc.ColorCandidate[[2]int{x, y}]
				fmt.Fprintf(&b, ".")
				if x == 0 && y == 13 {
					fmt.Print()
				}
				if len(cc.colors) != len(fcc.colors) {
					fmt.Fprintf(&l, "%d", len(cc.colors))
					fmt.Fprintf(&ccb, "%d,%d %v\n", x, y, cc.colors)
				} else {
					fmt.Fprint(&l, ".")
				}
			}
			// cc := fcc.ColorCandidate[[2]int{x, y}]
			// found := false
			// for i := 0; i < len(fcc.ShapeColorCandidate); i++ {
			// 	scc := fcc.ShapeColorCandidate[i]
			// 	// fmt.Printf("%d,%d %p %d:%p\n", x, y, cc, i, scc)
			// 	if cc == scc {
			// 		found = true
			// 		fmt.Fprintf(&b, "%d", i)
			// 		fmt.Fprint(&l, ".")
			// 		break
			// 	}
			// }
			// if !found {
			// 	fmt.Fprintf(&b, ".")
			// 	if x == 0 && y == 13 {
			// 		fmt.Print()
			// 	}
			// 	if len(cc.colors) != len(fcc.colors) {
			// 		fmt.Fprintf(&l, "%d", len(cc.colors))
			// 		fmt.Fprintf(&ccb, "%d,%d %v\n", x, y, cc.colors)
			// 	} else {
			// 		fmt.Fprint(&l, ".")
			// 	}
			// }
		}
		fmt.Fprintf(&b, " %s\n", l.String())
	}
	fmt.Print(b.String())
	fmt.Printf("<colors>\n%v\n", fcc.colors)
	fmt.Println("<shape adjacent>")
	for i, adj := range fcc.ShapeAdjacent {
		fmt.Printf("%d %v\n", i, adj)
	}
	fmt.Println("<shape color tables>")
	for i, scc := range fcc.ShapeColorCandidate {
		fmt.Printf("#%d %p%v\n", i, scc, scc)
	}
	fmt.Println("<position color tables>")
	fmt.Print(ccb.String())

	if fcc.BitField != nil {
		fmt.Println("<coverage>")
		pc := fcc.CountPlaced()
		sc := fcc.ShapeBitField.OverallShape().PopCount()
		fmt.Printf("%d/%d(%.02f%%)\n", pc, fcc.ShapeBitField.OverallShape().PopCount(), float64(pc*100)/float64(sc))
	}
}
