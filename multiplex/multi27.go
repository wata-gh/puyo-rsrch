package main

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/wata-gh/puyo2"
	"gonum.org/v1/gonum/stat/combin"
)

type Multi27Pattern struct {
	CacheSkipCount    int
	CombiCount        int
	ExecCombiCount    int
	CheckCount        int
	FoundCount        int
	InvalidEmptyCount int
	InvalidPlaceCount int
	ChainCount        int
	increment         chan *Increment
}

func (p *Multi27Pattern) incrementer() {
outer:
	for {
		incr := <-p.increment
		switch incr.name {
		case "CacheSkipCount":
			p.CacheSkipCount += incr.value
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

func (p *Multi27Pattern) Init() {
	p.increment = make(chan *Increment)
	go p.incrementer()
}

func (p *Multi27Pattern) Close() {
	p.increment <- &Increment{
		name:  "end",
		value: 0,
	}
}

func (p *Multi27Pattern) flip(emptyc []int) []int {
	nemptyc := make([]int, len(emptyc))
	for i, v := range emptyc {
		switch v % 3 {
		case 0:
			nemptyc[i] = v + 2
		case 1:
			nemptyc[i] = v
		case 2:
			nemptyc[i] = v - 2
		}
	}
	return nemptyc
}

func (p *Multi27Pattern) GenValidEmpties(pattern *Pattern, base []int, fieldc int, ctotal int) [][]int {
	cache := make(map[string]struct{})
	emptycs := combin.Combinations(len(base), fieldc-ctotal)
	nemptycs := [][]int{}
	for _, emptyc := range emptycs {
		if (*pattern).ValidEmpty(emptyc) == false {
			(*pattern).Incr("InvalidEmptyCount")
			continue
		}
		if opt.ShapeOnly {
			key := cacheKey(emptyc)
			if _, ok := cache[key]; ok {
				continue
			}
			femptyc := p.flip(emptyc)
			fkey := cacheKey(femptyc)
			cache[fkey] = struct{}{}
		}

		nemptycs = append(nemptycs, emptyc)
	}
	return nemptycs
}

func (p *Multi27Pattern) ValidPlace(list []int) bool {
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

func (p *Multi27Pattern) ValidEmpty(list []int) bool {
	for _, v := range list {
		for i := v - 3; i >= 0; i -= 3 {
			if contains(list, i) {
				continue
			}
			return false
		}
	}
	return true
}

func (p *Multi27Pattern) Vanish(list []int) bool {
	fb := puyo2.NewFieldBits()
	for _, v := range list {
		pos := p.Index2Field(v)
		fb.SetOnebit(pos[0], pos[1])
	}
	return fb.FindVanishingBits().IsEmpty() == false
}

func (p *Multi27Pattern) Index2Field(idx int) [2]int {
	return MULTIPLEX_27[idx]
}

func (p *Multi27Pattern) FieldC() int {
	return 27
}

func (p *Multi27Pattern) ChainC() int {
	return p.ChainCount
}

func (p *Multi27Pattern) Incr(name string) {
	p.increment <- &Increment{
		name:  name,
		value: 1,
	}
}
func (p *Multi27Pattern) Add(name string, c int) {
	p.increment <- &Increment{
		name:  name,
		value: c,
	}
}

func (p *Multi27Pattern) ShowResult() {
	fmt.Fprintf(os.Stderr, "combi: %d\ninvalid empty: %d\nexec combi: %d(%0.2f%%)\ninvalid place: %d\ncache skip: %d\ncheck: %d(%0.2f%%)\nfound: %d\n",
		p.CombiCount,
		p.InvalidEmptyCount,
		p.ExecCombiCount,
		float64(p.ExecCombiCount*100)/float64(p.CombiCount),
		p.InvalidPlaceCount,
		p.CacheSkipCount,
		p.CheckCount,
		float64(p.CheckCount*100)/float64(p.ExecCombiCount),
		p.FoundCount,
	)
}

func (p *Multi27Pattern) Check(field <-chan []int, opt options, wg *sync.WaitGroup) {
	pattern := Pattern(p)
	if opt.ShapeOnly {
		CheckShape(&pattern, field, opt, wg)
	} else {
		Check(&pattern, field, opt, wg)
	}

}
