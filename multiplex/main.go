package main

import (
	"fmt"
	"sync"

	"github.com/wata-gh/puyo2"
)

const CHAINC = 4
const FIELDC = 27

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

var GTR_R_18 = [...][2]int{
	{1, 5},
	{2, 5},
	{3, 5},
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

// multiplex27
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

// GTR 15
// func validEmpty(list []int) bool {
// 	for _, v := range list {
// 		if v < 4 {
// 			continue
// 		}
// 		if v < 8 {
// 			if contains(list, v-4) {
// 				continue
// 			}
// 			return false
// 		}
// 		if v < 11 {
// 			if contains(list, v-3) && contains(list, v-7) {
// 				continue
// 			}
// 			return false
// 		}
// 		return false
// 	}
// 	return true
// }

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

func index2field(idx int) [2]int {
	return MULTIPLEX_27[idx]
}

// GTR 15 check
// func check(field <-chan []int, wg *sync.WaitGroup) {
// 	for {
// 		puyos := <-field
// 		if len(puyos) == 0 {
// 			break
// 		}
// 		bf := puyo2.NewBitField()
// 		sort.Ints(puyos)
// 		for n, puyo := range puyos {
// 			for i := 0; puyo > 0; i++ {
// 				if puyo&1 == 1 {
// 					pos := index2field(i)
// 					color := []puyo2.Color{puyo2.Red, puyo2.Green, puyo2.Blue, puyo2.Yellow}[n%4]
// 					bf.SetColor(color, pos[0], pos[1])
// 				}
// 				puyo >>= 1
// 			}
// 		}
// 		result := bf.SimulateWithNewBitField()
// 		if result.BitField.Equals(bf) {
// 			nbf := bf.Clone()
// 			fb := puyo2.NewFieldBits()
// 			fb.SetOnebit(0, 2)
// 			nbf.Drop(fb)
// 			result = nbf.SimulateWithNewBitField()
// 			if result.Chains == CHAINC {
// 				fmt.Println(bf.MattulwanEditorParam())
// 				bf.ExportImage("results/" + bf.MattulwanEditorParam() + ".png")
// 			}
// 		}
// 	}
// 	wg.Done()
// }

// simple check
func check(field <-chan []int, wg *sync.WaitGroup) {
	for {
		puyos := <-field
		if len(puyos) == 0 {
			break
		}
		bf := puyo2.NewBitField()
		for n, puyo := range puyos {
			for i := 0; puyo > 0; i++ {
				if puyo&1 == 1 {
					pos := index2field(i)
					color := []puyo2.Color{puyo2.Red, puyo2.Green, puyo2.Blue, puyo2.Yellow}[n%4]
					bf.SetColor(color, pos[0], pos[1])
				}
				puyo >>= 1
			}
		}
		result := bf.SimulateWithNewBitField()
		if result.Chains == CHAINC {
			fmt.Println(bf.MattulwanEditorParam())
			bf.ExportImage("results/" + bf.MattulwanEditorParam() + ".png")
		}
	}
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	field := make(chan []int)
	wg.Add(1)
	go check(field, &wg)
	Gen(field, 1)
	wg.Wait()
}
