package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/wata-gh/puyo2"
)

var MULTI_4 = [...][2]int{
	{0, 1},
	{0, 2},
	{1, 1},
	{1, 2},
}

type Multi4Pattern struct {
	CacheSkip         int
	CombiCount        int
	ExecCombiCount    int
	CheckCount        int
	FoundCount        int
	InvalidEmptyCount int
}

func (p *Multi4Pattern) ValidEmpty(list []int) bool {
	for _, v := range list {
		for i := v - 2; i >= 0; i -= 2 {
			if contains(list, i) {
				continue
			}
			return false
		}
	}
	return true
}

func (p *Multi4Pattern) Index2Field(idx int) [2]int {
	return MULTI_4[idx]
}

func (p *Multi4Pattern) FieldC() int {
	return 4
}

func (p *Multi4Pattern) ChainC() int {
	return 2
}

func (p *Multi4Pattern) AddInvalidEmpty() {
	p.InvalidEmptyCount++
}
func (p *Multi4Pattern) AddCheck() {
	p.CheckCount++
}
func (p *Multi4Pattern) AddCacheSkip() {
	p.CacheSkip++
}
func (p *Multi4Pattern) AddCombi(c int) {
	p.CombiCount += c
}
func (p *Multi4Pattern) AddExecCombi() {
	p.ExecCombiCount++
}
func (p *Multi4Pattern) AddFound() {
	p.FoundCount++
}
func (p *Multi4Pattern) ShowResult() {
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

func (p *Multi4Pattern) Check(field <-chan []int, wg *sync.WaitGroup) {
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
		bf.ExportImage("results/" + bf.MattulwanEditorParam() + ".png")
	}
	wg.Done()

}
