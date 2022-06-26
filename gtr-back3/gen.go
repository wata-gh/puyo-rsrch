package main

import (
	"sort"

	"github.com/wata-gh/puyo2"
	"gonum.org/v1/gonum/stat/combin"
)

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

func remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func validEmpty(list []int) bool {
	for _, v := range list {
		if v <= 3 {
			continue
		}
		if v <= 6 {
			if contains(list, v-4) {
				continue
			}
			return false
		}
		if v <= 10 {
			if contains(list, v-4) && (v == 7 || contains(list, v-8)) {
				continue
			}
			return false
		}
		if v <= 13 {
			if contains(list, v-3) && contains(list, v-7) && contains(list, v-11) {
				continue
			}
			return false
		}
		if v == 14 {
			return false
		}
		if v <= 17 {
			if contains(list, v-4) && contains(list, v-7) && contains(list, v-11) && contains(list, v-15) {
				continue
			}
			return false
		}
		panic(list)
	}
	return true
}

func genCombinations(fc int, base []int, c1c int, c2c int, c3c int, c4c int, field chan<- []int) {
	ctotal := c1c + c2c + c3c + c4c
	if fc > ctotal {
		ccs := combin.Combinations(len(base), fc-ctotal)
		for _, cs := range ccs {
			if validEmpty(cs) {
				base, _ := exclude(base, cs)
				genColorCombinations(base, c1c, c2c, c3c, c4c, field)
			}
		}
	} else {
		genColorCombinations(base, c1c, c2c, c3c, c4c, field)
	}
}

func genColorCombinations(base []int, c1c int, c2c int, c3c int, c4c int, field chan<- []int) {
	cache := map[[4]int]struct{}{}
	c1s := combin.Combinations(len(base), c1c)
	for _, c1 := range c1s {
		bitsC1 := [4]int{0, 0, 0, 0}
		board := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		base, c1V := exclude(base, c1)
		for _, c := range c1V {
			bitsC1[0] |= 1 << c
			board[c] = int(puyo2.Red)
		}
		c2s := combin.Combinations(len(base), c2c)
		for _, c2 := range c2s {
			bitsC2 := [4]int{bitsC1[0], 0, 0, 0}
			boardC2 := make([]int, len(board))
			copy(boardC2, board)
			base, c2V := exclude(base, c2)
			for _, c := range c2V {
				bitsC2[1] |= 1 << c
				boardC2[c] = int(puyo2.Green)
			}
			c3s := combin.Combinations(len(base), c3c)
			for _, c3 := range c3s {
				bitsC3 := [4]int{bitsC2[0], bitsC2[1], 0, 0}
				boardC3 := make([]int, len(board))
				copy(boardC3, boardC2)
				base, c3V := exclude(base, c3)
				for _, c := range c3V {
					bitsC3[2] |= 1 << c
					boardC3[c] = int(puyo2.Yellow)
				}
				c4s := combin.Combinations(len(base), c4c)
				for _, c4 := range c4s {
					bitsC4 := [4]int{bitsC3[0], bitsC3[1], bitsC3[2], 0}
					boardC4 := make([]int, len(board))
					copy(boardC4, boardC3)
					_, c4V := exclude(base, c4)
					for _, c := range c4V {
						bitsC4[3] |= 1 << c
						boardC4[c] = int(puyo2.Blue)
					}

					sort.Ints(bitsC4[:])
					if _, ok := cache[bitsC4]; ok {
					} else {
						cache[bitsC4] = struct{}{}
						field <- boardC4
					}
				}
			}
		}
	}
}

func Gen(field chan<- []int, grc int) {
	fieldc := 18
	puyocnt := []int{4}
	for _, c1c := range puyocnt {
		for _, c2c := range puyocnt {
			for _, c3c := range puyocnt {
				for _, c4c := range puyocnt {
					puyoc := c1c + c2c + c3c + c4c
					if puyoc <= fieldc {
						base := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17}
						genCombinations(fieldc, base, c1c, c2c, c3c, c4c, field)
					}
				}
			}
		}
	}
	for i := 0; i < grc; i++ {
		field <- []int{}
	}
}
