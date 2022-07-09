package main

import (
	"os"
	"sync"
)

// GTR 18
// func validEmpty(list []int) bool {
// 	for _, v := range list {
// 		if v <= 3 {
// 			continue
// 		}
// 		if v <= 6 {
// 			if contains(list, v-4) {
// 				continue
// 			}
// 			return false
// 		}
// 		if v <= 10 {
// 			if contains(list, v-4) && (v == 7 || contains(list, v-8)) {
// 				continue
// 			}
// 			return false
// 		}
// 		if v <= 13 {
// 			if contains(list, v-3) && contains(list, v-7) && contains(list, v-11) {
// 				continue
// 			}
// 			return false
// 		}
// 		if v == 14 {
// 			return false
// 		}
// 		if v <= 17 {
// 			if contains(list, v-4) && contains(list, v-7) && contains(list, v-11) && contains(list, v-15) {
// 				continue
// 			}
// 			return false
// 		}
// 		panic(list)
// 	}
// 	return true
// }

// func validEmpty(list []int) bool {
// 	fb := puyo2.NewFieldBitsWithM([2]uint64{18446744073709551615, 18446744073709551615})
// 	for _, v := range list {
// 		pos := GTR_R_18[v]
// 		fb.SetOnebit(pos[0], pos[1])
// 	}
// 	for _, v := range list {
// 		pos := GTR_R_18[v]
// 		fb.Onebit(pos[0], pos[1] + 1)
// 	}
// 	for _, v := range list {
// 		for up := v - 3; up > 0; up -= 3 {
// 			if contains(list, up) {
// 				continue
// 			}
// 			return false
// 		}
// 	}
// 	return true
// }

// func index2field(idx int) [2]int {
// 	return MULTIPLEX_27[idx]
// }

func main() {
	os.Mkdir("results", 0755)
	var wg sync.WaitGroup
	field := make(chan []int)
	wg.Add(1)
	// var pattern Pattern = &Gtr15Pattern{}
	// var pattern Pattern = &Multi4Pattern{}
	// var pattern Pattern = &Multi9Pattern{}
	var pattern Pattern = &Multi27Pattern{
		ChainCount: 4,
	}
	go pattern.Check(field, &wg)
	Gen(&pattern, field, 1)
	wg.Wait()
	pattern.ShowResult()
}
