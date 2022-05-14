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
			} else {
				return false
			}
		}
		if v < 10 {
			if contains(list, v-3) && contains(list, v-7) {
				continue
			} else {
				return false
			}
		}
	}
	return true
}

func Gen(field chan<- [15]int) {
	list := combin.Combinations(15, 3)
	for _, empties := range list {
		board := [15]int{}
		base := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
		if validEmpty(empties) == false {
			continue
		}
		base, _ = exclude(base, empties)
		// 空マスは 0 なので埋める必要なし
		// for _, e := range empties_v {
		// 	board[e] = int(puyo2.Empty)
		// }
		c1s := combin.Combinations(len(base), 4)
		for _, c1 := range c1s {
			base, c1_v := exclude(base, c1)
			for _, c := range c1_v {
				board[c] = int(puyo2.Red)
			}
			c2s := combin.Combinations(len(base), 4)
			for _, c2 := range c2s {
				base, c2_v := exclude(base, c2)
				for _, c := range c2_v {
					board[c] = int(puyo2.Green)
				}
				c3s := combin.Combinations(len(base), 4)
				for _, c3 := range c3s {
					_, c3_v := exclude(base, c3)
					for _, c := range c3_v {
						board[c] = int(puyo2.Yellow)
					}
					field <- board
				}
			}
		}
	}
}
