package main

import (
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

func genCombinations(base []int, c1c int, c2c int, c3c int, c4c int, field chan<- []int) {
	c1s := combin.Combinations(len(base), c1c)
	for _, c1 := range c1s {
		board := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		base, c1v := exclude(base, c1)
		for _, c := range c1v {
			board[c] = int(puyo2.Red)
		}
		c2s := combin.Combinations(len(base), c2c)
		for _, c2 := range c2s {
			boardC2 := make([]int, len(board))
			copy(boardC2, board)
			base, c2v := exclude(base, c2)
			for _, c := range c2v {
				boardC2[c] = int(puyo2.Green)
			}
			c3s := combin.Combinations(len(base), c3c)
			for _, c3 := range c3s {
				boardC3 := make([]int, len(board))
				copy(boardC3, boardC2)
				base, c3v := exclude(base, c3)
				for _, c := range c3v {
					boardC3[c] = int(puyo2.Yellow)
				}
				c4s := combin.Combinations(len(base), c4c)
				for _, c4 := range c4s {
					boardC4 := make([]int, len(board))
					copy(boardC4, boardC3)
					base, c4v := exclude(base, c4)
					for _, c := range c4v {
						boardC4[c] = int(puyo2.Blue)
					}
					// 残りは空白マスとして扱う
					if validEmpty(base) == false {
						continue
					}
					// for _, c := range base {
					// 	board_c3[c] = int(puyo2.Empty)
					// }
					field <- boardC4
				}
			}
		}
	}
}

func Gen(field chan<- []int, routines int) {
	fieldc := 15
	puyocnt := []int{4, 5, 6, 7}
	for _, c1c := range puyocnt {
		for _, c2c := range puyocnt {
			for _, c3c := range puyocnt {
				for _, c4c := range []int{0, 1, 2, 3} {
					puyoc := c1c + c2c + c3c + c4c
					if puyoc <= fieldc {
						base := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
						genCombinations(base, c1c, c2c, c3c, c4c, field)
					}
				}
			}
		}
	}
	for i := 0; i < routines; i++ {
		field <- []int{}
	}
}
