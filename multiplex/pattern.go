package main

import (
	"fmt"
	"sync"

	"github.com/wata-gh/puyo2"
)

var GTR_R_18 = [...][2]int{
	{1, 5},
	{2, 5},
	{3, 5},
	{0, 4},
	{1, 4},
	{2, 4},
	{3, 4},
	{0, 3},
	{1, 3},
	{2, 3},
	{3, 3},
	{1, 2},
	{2, 2},
	{3, 2},
	{0, 1},
	{1, 1},
	{2, 1},
	{3, 1},
}

var MULTIPLEX_27 = [...][2]int{
	{0, 9},
	{1, 9},
	{2, 9},
	{0, 8},
	{1, 8},
	{2, 8},
	{0, 7},
	{1, 7},
	{2, 7},
	{0, 6},
	{1, 6},
	{2, 6},
	{0, 5},
	{1, 5},
	{2, 5},
	{0, 4},
	{1, 4},
	{2, 4},
	{0, 3},
	{1, 3},
	{2, 3},
	{0, 2},
	{1, 2},
	{2, 2},
	{0, 1},
	{1, 1},
	{2, 1},
}

type Increment struct {
	name  string
	value int
}

type Pattern interface {
	Init()
	ValidEmpty(list []int) bool
	ValidPlace(list []int) bool
	Index2Field(idx int) [2]int
	Check(field <-chan []int, opt options, wg *sync.WaitGroup)
	FieldC() int
	ChainC() int
	Add(name string, c int)
	Incr(name string)
	GenValidEmpties(pattern *Pattern, base []int, fieldc int, ctotal int) [][]int
	ShowResult()
	Vanish(list []int) bool
	flip(emptyc []int) []int
}

func Check(p *Pattern, field <-chan []int, opt options, wg *sync.WaitGroup) {
	for {
		puyos := <-field
		if len(puyos) == 0 {
			break
		}
		(*p).Incr("CheckCount")
		bf := puyo2.NewBitField()
		for n, puyo := range puyos {
			for i := 0; puyo > 0; i++ {
				if puyo&1 == 1 {
					pos := (*p).Index2Field(i)
					color := []puyo2.Color{puyo2.Red, puyo2.Green, puyo2.Blue, puyo2.Yellow}[n%4]
					bf.SetColor(color, pos[0], pos[1])
				}
				puyo >>= 1
			}
		}
		result := bf.SimulateWithNewBitField()
		if result.Chains == (*p).ChainC() {
			(*p).Incr("FoundCount")
			fmt.Println(bf.MattulwanEditorParam())
			bf.ExportImage(opt.Dir + "/" + bf.MattulwanEditorParam() + ".png")
		}
	}
	wg.Done()
}

func createFlipShape(p *Pattern, puyos []int) *puyo2.ShapeBitField {
	sbf := puyo2.NewShapeBitField()
	for _, puyo := range puyos {
		shape := puyo2.NewFieldBits()
		for i := 0; puyo > 0; i++ {
			if puyo&1 == 1 {
				pos := (*p).Index2Field(i)
				switch pos[0] % 3 {
				case 0:
					pos[0] += 2
				case 1:
				case 2:
					pos[0] -= 2
				}

				shape.SetOnebit(pos[0], pos[1])
			}
			puyo >>= 1
		}
		sbf.AddShape(shape)
	}
	return sbf
}

func checkShape(p *Pattern, puyos []int, opt options) {
	(*p).Incr("CheckCount")
	hwm := [3]int{} // multi27 dependent
	sbf := puyo2.NewShapeBitField()
	for _, puyo := range puyos {
		shape := puyo2.NewFieldBits()
		for i := 0; puyo > 0; i++ {
			if puyo&1 == 1 {
				pos := (*p).Index2Field(i)
				shape.SetOnebit(pos[0], pos[1])
				if hwm[pos[0]] < pos[1] {
					hwm[pos[0]] = pos[1]
				}
			}
			puyo >>= 1
		}
		sbf.AddShape(shape)
	}

	flip := hwm[0] != hwm[2] // multi27 dependent

	result := sbf.Simulate()
	if result.Chains == (*p).ChainC() {
		(*p).Incr("FoundCount")
		fmt.Println(sbf.ChainOrderedFieldString())
		// sbf.ExportChainImage(fmt.Sprintf("%s/%s.png", opt.Dir, sbf.ChainOrderedFieldString()))
		if flip {
			fsbf := createFlipShape(p, puyos)
			fsbf.Simulate()
			if sbf.ChainOrderedFieldString() != fsbf.ChainOrderedFieldString() {
				(*p).Incr("FoundCount")
				fmt.Println(fsbf.ChainOrderedFieldString())
				// fsbf.ExportChainImage(fmt.Sprintf("%s/%s.png", opt.Dir, fsbf.ChainOrderedFieldString()))
			}
		}
	}
}

func CheckShape(p *Pattern, field <-chan []int, opt options, wg *sync.WaitGroup) {
	for {
		puyos := <-field
		if len(puyos) == 0 {
			break
		}
		checkShape(p, puyos, opt)
	}
	wg.Done()
}
