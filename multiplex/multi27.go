package main

import (
	"fmt"
	"os"
	"sort"
	"sync"
)

type Multi27Pattern struct {
	CacheSkip         int
	CombiCount        int
	ExecCombiCount    int
	CheckCount        int
	FoundCount        int
	InvalidEmptyCount int
	InvalidPlaceCount int
	ChainCount        int
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

func (p *Multi27Pattern) Index2Field(idx int) [2]int {
	return MULTIPLEX_27[idx]
}

func (p *Multi27Pattern) FieldC() int {
	return 27
}

func (p *Multi27Pattern) ChainC() int {
	return p.ChainCount
}

func (p *Multi27Pattern) AddInvalidEmpty() {
	p.InvalidEmptyCount++
}
func (p *Multi27Pattern) AddCheck() {
	p.CheckCount++
}
func (p *Multi27Pattern) AddCacheSkip() {
	p.CacheSkip++
}
func (p *Multi27Pattern) AddCombi(c int) {
	p.CombiCount += c
}
func (p *Multi27Pattern) AddExecCombi() {
	p.ExecCombiCount++
}
func (p *Multi27Pattern) AddFound() {
	p.FoundCount++
}
func (p *Multi27Pattern) AddInvalidPlace() {
	p.InvalidPlaceCount++
}

func (p *Multi27Pattern) ShowResult() {
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

func (p *Multi27Pattern) Check(field <-chan []int, wg *sync.WaitGroup) {
	pattern := Pattern(p)
	Check(&pattern, field, wg)
}
