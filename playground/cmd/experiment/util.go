package main

import (
	"fmt"
	"sort"
	"strings"

	"gonum.org/v1/gonum/stat/combin"
)

var MULTIPLEX_20 = [...][2]int{
	{0, 7},
	{1, 7},
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

func noSortCacheKey(bits []int) string {
	var b strings.Builder
	for _, bit := range bits {
		fmt.Fprintf(&b, "_%d", bit)
	}
	return b.String()
}

func cacheKey(bits []int) string {
	sort.Ints(bits)
	var b strings.Builder
	for _, bit := range bits {
		fmt.Fprintf(&b, "_%d", bit)
	}
	return b.String()
}

func cloneArray(ints []int) []int {
	newInts := make([]int, len(ints))
	copy(newInts, ints)
	return newInts
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func exclude(b []int, e []int) ([]int, []int) {
	result := []int{}
	excluded := []int{}
	for i, v := range b {
		if contains(e, i) {
			excluded = append(excluded, v)
		} else {
			result = append(result, v)
		}
	}
	return result, excluded
}

func flip(emptyc []int) []int {
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

func index2Field(idx int) [2]int {
	return MULTIPLEX_27[idx]
}

func intArray2Bit(cv []int) int {
	bit := 0
	for _, c := range cv {
		bit |= 1 << c
	}
	return bit
}

func validEmpty(list []int) bool {
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

func genValidEmpties(fieldc int, emptyc int) [][]int {
	cache := make(map[string]struct{})
	emptycs := combin.Combinations(fieldc, emptyc)
	nemptycs := [][]int{}
	for _, emptyc := range emptycs {
		if validEmpty(emptyc) == false {
			continue
		}
		key := cacheKey(emptyc)
		if _, ok := cache[key]; ok {
			continue
		}
		femptyc := flip(emptyc)
		fkey := cacheKey(femptyc)
		cache[fkey] = struct{}{}

		nemptycs = append(nemptycs, emptyc)
	}
	return nemptycs
}

func validPlace(list []int) bool {
	poss := []int{}
	for _, v := range list {
		pos := index2Field(v)
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
