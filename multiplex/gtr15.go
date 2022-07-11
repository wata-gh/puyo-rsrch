package main

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/wata-gh/puyo2"
	"gonum.org/v1/gonum/stat/combin"
)

var GTR_R_15 = [...][2]int{
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

type Gtr15Pattern struct {
	CacheSkip         int
	CombiCount        int
	ExecCombiCount    int
	CheckCount        int
	FoundCount        int
	InvalidEmptyCount int
	InvalidPlaceCount int
	increment         chan *Increment
}

func (p *Gtr15Pattern) incrementer() {
outer:
	for {
		incr := <-p.increment
		switch incr.name {
		case "CacheSkip":
			p.CacheSkip += incr.value
		case "CombiCount":
			p.CombiCount += incr.value
		case "ExecCombiCount":
			p.ExecCombiCount += incr.value
		case "CheckCount":
			p.CheckCount += incr.value
		case "FoundCount":
			p.FoundCount += incr.value
		case "InvalidEmptyCount":
			p.InvalidEmptyCount += incr.value
		case "InvalidPlaceCount":
			p.InvalidPlaceCount += incr.value
		case "end":
			break outer
		default:
			panic(fmt.Sprintf("invalid increment name. %+v", incr))
		}
	}
}

func (p *Gtr15Pattern) flip(emptyc []int) []int {
	panic("Gtr15Pattern does not implements flip metho.")
}

func (p *Gtr15Pattern) Init() {
	p.increment = make(chan *Increment)
	go p.incrementer()
}

func (p *Gtr15Pattern) GenValidEmpties(pattern *Pattern, base []int, fieldc int, ctotal int) [][]int {
	return combin.Combinations(len(base), fieldc-ctotal)
}

func (p *Gtr15Pattern) ValidPlace(list []int) bool {
	poss := []int{}
	for _, v := range list {
		pos := p.Index2Field(v)
		poss = append(poss, pos[0])
	}
	sort.Ints(poss)
	for i, x := range poss[1:] {
		c := poss[i]
		if x != c && x != c+1 {
			return false
		}
	}
	return true
}

func (p *Gtr15Pattern) ValidEmpty(list []int) bool {
	for _, v := range list {
		if v < 4 {
			continue
		}
		if v < 8 {
			if contains(list, v-4) {
				continue
			}
			return false
		}
		if v < 11 {
			if contains(list, v-3) && contains(list, v-7) {
				continue
			}
			return false
		}
		return false
	}
	return true
}

func (p *Gtr15Pattern) Vanish(list []int) bool {
	fb := puyo2.NewFieldBits()
	for _, v := range list {
		pos := p.Index2Field(v)
		fb.SetOnebit(pos[0], pos[1])
	}
	return fb.FindVanishingBits().IsEmpty() == false
}

func (p *Gtr15Pattern) Index2Field(idx int) [2]int {
	return GTR_R_15[idx]
}

func (p *Gtr15Pattern) FieldC() int {
	return 15
}

func (p *Gtr15Pattern) ChainC() int {
	return 3
}

func (p *Gtr15Pattern) Incr(name string) {
	p.increment <- &Increment{
		name:  name,
		value: 1,
	}
}
func (p *Gtr15Pattern) Add(name string, c int) {
	p.increment <- &Increment{
		name:  name,
		value: c,
	}
}

func (p *Gtr15Pattern) ShowResult() {
	fmt.Fprintf(os.Stderr, "combi: %d\ninvalid empty: %d\nexec combi: %d(%0.2f%%)\ninvalid place: %d\ncache skip: %d\ncheck: %d(%0.2f%%)\nfound: %d\n",
		p.CombiCount,
		p.InvalidEmptyCount,
		p.ExecCombiCount,
		float64(p.ExecCombiCount*100)/float64(p.CombiCount),
		p.InvalidPlaceCount,
		p.CacheSkip,
		p.CheckCount,
		float64(p.CheckCount*100)/float64(p.ExecCombiCount),
		p.FoundCount,
	)
}

func (p *Gtr15Pattern) Check(field <-chan []int, opt options, wg *sync.WaitGroup) {
	for {
		puyos := <-field
		if len(puyos) == 0 {
			break
		}
		p.Incr("CheckCount")
		bf := puyo2.NewBitField()
		sort.Ints(puyos)
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
		result := bf.SimulateWithNewBitField()
		if result.BitField.Equals(bf) {
			nbf := bf.Clone()
			fb := puyo2.NewFieldBits()
			fb.SetOnebit(0, 2)
			nbf.Drop(fb)
			result = nbf.SimulateWithNewBitField()
			if result.Chains == p.ChainC() {
				p.Incr("FoundCount")
				fmt.Println(bf.MattulwanEditorParam())
				bf.ExportImage(opt.Dir + "/" + bf.MattulwanEditorParam() + ".png")
			}
		}
	}
	wg.Done()
}
