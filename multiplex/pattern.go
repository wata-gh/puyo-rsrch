package main

import (
	"fmt"
	"os"
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

type Pattern interface {
	ValidEmpty(list []int) bool
	ValidPlace(list []int) bool
	Index2Field(idx int) [2]int
	Check(field <-chan []int, wg *sync.WaitGroup)
	FieldC() int
	ChainC() int
	AddCombi(c int)
	AddExecCombi()
	AddInvalidEmpty()
	AddCacheSkip()
	AddInvalidPlace()
	AddCheck()
	AddFound()
	ShowResult()
}

func Check(p *Pattern, field <-chan []int, wg *sync.WaitGroup) {
	for {
		puyos := <-field
		if len(puyos) == 0 {
			break
		}
		(*p).AddCheck()
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
			(*p).AddFound()
			fmt.Println(bf.MattulwanEditorParam())
			os.Mkdir("results", 0755)
			bf.ExportImage("results/" + bf.MattulwanEditorParam() + ".png")
		}
	}
	wg.Done()
}
