package main

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/wata-gh/puyo2"
)

var MULTI_9 = [...][2]int{
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

type Multi9Pattern struct {
	CacheSkip         int
	CombiCount        int
	ExecCombiCount    int
	CheckCount        int
	FoundCount        int
	InvalidEmptyCount int
	InvalidPlaceCount int
}

func (p *Multi9Pattern) ValidEmpty(list []int) bool {
	width := 3
	for _, v := range list {
		for i := v - width; i >= 0; i -= width {
			if contains(list, i) {
				continue
			}
			return false
		}
	}
	return true
}

func (p *Multi9Pattern) ValidPlace(list []int) bool {
	poss := []int{}
	for _, v := range list {
		pos := p.Index2Field(v)
		poss = append(poss, pos[0])
	}
	sort.Ints(poss)
	c := poss[0]
	for _, x := range poss {
		if c == x || x == c+1 {
			c = x
			continue
		} else {
			return false
		}
	}
	return true
}

func (p *Multi9Pattern) Index2Field(idx int) [2]int {
	return MULTI_9[idx]
}

func (p *Multi9Pattern) FieldC() int {
	return 9
}

func (p *Multi9Pattern) ChainC() int {
	return 2
}

func (p *Multi9Pattern) AddInvalidEmpty() {
	p.InvalidEmptyCount++
}
func (p *Multi9Pattern) AddCheck() {
	p.CheckCount++
}
func (p *Multi9Pattern) AddCacheSkip() {
	p.CacheSkip++
}
func (p *Multi9Pattern) AddCombi(c int) {
	p.CombiCount += c
}
func (p *Multi9Pattern) AddExecCombi() {
	p.ExecCombiCount++
}
func (p *Multi9Pattern) AddFound() {
	p.FoundCount++
}
func (p *Multi9Pattern) AddInvalidPlace() {
	p.InvalidPlaceCount++
}

func (p *Multi9Pattern) ShowResult() {
	fmt.Fprintf(os.Stderr, "combi: %d\ninvalid empty: %d\nexec combi: %d(%0.2f%%)\ncache skip: %d(%0.2f%%)\ncheck: %d(%0.2f%%)\nfound: %d\n",
		p.CombiCount,
		p.InvalidEmptyCount,
		p.ExecCombiCount,
		float64(p.ExecCombiCount*100)/float64(p.CombiCount),
		p.CacheSkip,
		float64(p.CacheSkip*100)/float64(p.ExecCombiCount),
		p.CheckCount,
		float64(p.CheckCount*100)/float64(p.ExecCombiCount),
		p.FoundCount,
	)
}

func (p *Multi9Pattern) Check(field <-chan []int, wg *sync.WaitGroup) {
	for {
		puyos := <-field
		if len(puyos) == 0 {
			break
		}
		p.AddCheck()
		bf := puyo2.NewBitField()
		for n, puyo := range puyos {
			for i := 0; puyo > 0; i++ {
				if puyo&1 == 1 {
					pos := p.Index2Field(i)
					color := []puyo2.Color{puyo2.Red, puyo2.Green, puyo2.Blue, puyo2.Yellow}[n%4]
					bf.SetColor(color, pos[0], pos[1])
				}
				puyo >>= 1
			}
		}
		p.AddFound()
		fmt.Println(bf.MattulwanEditorParam())
		os.Mkdir("results", 0755)
		bf.ExportImage("results/" + bf.MattulwanEditorParam() + ".png")
	}
	wg.Done()
}
