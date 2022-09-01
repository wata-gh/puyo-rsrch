package main

import (
	"fmt"
	"math/bits"
	"os"
	"sync"
	"testing"

	"github.com/wata-gh/puyo2"
)

func TestFronMain(t *testing.T) {
	param := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa4aaaa434aaa421aaa311aaa313aaa222aaa"
	sbf := puyo2.NewShapeBitFieldWithFieldString(param)
	sbfb := sbf.Clone()

	r := adjacentColorCount(sbfb)
	sbfb = sbf.Clone()
	if fronConnectable(sbfb) {
		sbfb = sbf.Clone()
		fronableBase(sbfb)
		overall := sbfb.OverallShape()
		overall.SetOnebit(0, 0)
		overall.SetOnebit(1, 0)
		y1 := bits.Len64(overall.ColBits(0))
		y2 := bits.Len64(overall.ColBits(1) >> 16)

		keyPatterns := [][2][2]int{}
		if y1 <= 12 {
			keyPatterns = append(keyPatterns, [2][2]int{{0, y1}, {0, y1 + 1}})
		}
		keyPatterns = append(keyPatterns, [2][2]int{{0, y1}, {1, y2}})
		keyPatterns = append(keyPatterns, [2][2]int{{1, y2}, {0, y1}})
		if y2 <= 12 {
			keyPatterns = append(keyPatterns, [2][2]int{{1, y2}, {1, y2 + 1}})
		}
		found := false
		for _, keyPattern := range keyPatterns {
			sbfc := sbfb.Clone()
			s1 := sbfc.Shapes[len(sbfc.Shapes)-3]
			s3 := sbfc.Shapes[len(sbfc.Shapes)-1]
			s1.SetOnebit(keyPattern[0][0], keyPattern[0][1])
			s3.SetOnebit(keyPattern[1][0], keyPattern[1][1])
			s := fronFireable(sbfc)
			if s != UnFireable && colorable(sbfc) {
				sbfc.ShowDebug()
				fmt.Println(sbfc.FillChainableColor().MattulwanEditorUrl())
				fmt.Println(param, s, r)
				found = true
				break
			}
		}
		if found == false {
			fmt.Fprintf(os.Stderr, "%s\n", param)
		}
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", param)
	}
}
func TestMain(t *testing.T) {
	param := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa4aaaa434aaa421aaa311aaa313aaa222aaa"
	sbf := puyo2.NewShapeBitFieldWithFieldString(param)
	sbfc := sbf.Clone()

	r := adjacentColorCount(sbfc)
	s := fireable(sbfc)
	sbfc = sbf.Clone()
	if connectable(sbfc) && s != UnFireable {
		sbfc = sbf.Clone()
		if gtrable(sbfc) {
			if colorable(sbfc) {
				sbfc.ShowDebug()
				fmt.Println(param, s, r)
			} else {
				fmt.Fprintf(os.Stderr, "%s\n", param)
			}
		}
	}

}

func TestColor(t *testing.T) {
	sbf := puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa2aaaa12aaaa11aaaa212aaa")
	fronableBase(sbf)
	s1 := sbf.Shapes[len(sbf.Shapes)-3]
	s3 := sbf.Shapes[len(sbf.Shapes)-1]
	s1.SetOnebit(0, 7)
	s3.SetOnebit(0, 8)

	bf := sbf.FillChainableColor()
	bf.ShowDebug()
	// colorable(sbf)
}

func TestGtrColor(t *testing.T) {
	sbf := puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa4aaa443aaa112aaa411aaa322aaa332aaa")
	gtrable(sbf)
	sbf.ShowDebug()

	bf := sbf.FillChainableColor()
	bf.ShowDebug()
	// colorable(sbf)
}

func TestShift(t *testing.T) {
	for i := 0; i < 6; i++ {
		array := []int{1, 2, 3, 4}
		remainder := i % 4
		part := array[0:remainder]
		array = array[remainder:]
		array = append(array, part...)
		fmt.Println(array)
	}
}

func TestFronable(t *testing.T) {
	params := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go check(params, &wg)
	params <- "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa2aaaa12aaaa11aaaa212aaa"
	params <- ""
	wg.Wait()
}

func TestGtrable(t *testing.T) {
	sbf := puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaa3aaaaa2aaaaa2aaaaa2aaaaa1aaaaa1aaaaa11aaaa23aaaa33aaaa")
	// sbf.ShowDebug()
	gtrable(sbf)
	len := len(sbf.Shapes)
	tmp := sbf.Shapes[len-3]
	sbf.Shapes[len-3] = sbf.Shapes[len-1]
	sbf.Shapes[len-1] = tmp
	result := sbf.Clone().Simulate()
	fmt.Println(result)
	r := colorable(sbf)
	fmt.Println(r)
	sbf.ShowDebug()
}

func TestFireable(t *testing.T) {
	sbf := puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa2aaaa22aaa333aaa211aaa311aaa")
	sbf.ShowDebug()
	sbf.Simulate()
	s := fireable(sbf)
	if s != UnFireable {
		t.Fatalf("must be unfireable. %d", s)
	}

	sbf = puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa3aaaaa2aaaaa2aaaaa2aaaaa1aaaaa1aaaaa12aaaa13aaaa33aaaa")
	sbf.ShowDebug()
	sbf.Simulate()
	s = fireable(sbf)
	if s != UnFireable {
		t.Fatalf("must be unfireable. %d", s)
	}

	sbf = puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa3aaaaa2aaaaa2aaaaa2aaaaa1aaaaa1aaaaa11aaaa23aaaa33aaaa")
	sbf.ShowDebug()
	sbf.Simulate()
	s = fireable(sbf)
	if s != Normal {
		t.Fatalf("must be fireable. %d", s)
	}

	sbf = puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa3aaaaa3aaaaa2aaaaa1aaaaa12aaaa11aaaa32aaaa32aaaa")
	sbf.ShowDebug()
	sbf.Simulate()
	s = fireable(sbf)
	if s != Eighth1 {
		t.Fatalf("must be fireable from eighth. %d", s)
	}
}

// func TestAdjacentColorCount(t *testing.T) {
// 	sbf := puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa3aaaaa2aaaaa2aaaaa2aaaaa1aaaaa1aaaaa11aaaa23aaaa33aaaa")
// 	sbf.ShowDebug()
// 	r := adjacentColorCount(sbf)
// 	if r[0] != 1 || r[1] != 1 || r[2] != 0 {
// 		t.Fatalf("result must be [1 1 0] but %v", r)
// 	}
// 	sbf = puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa2aaaaa11aaaa133aaa123aaa322aaa")
// 	sbf.ShowDebug()
// 	r = adjacentColorCount(sbf)
// 	if r[0] != 1 || r[1] != 2 || r[2] != 2 {
// 		t.Fatalf("result must be [1 2 2] but %v", r)
// 	}
// 	sbf = puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa3aaaa22aaaa123aaa112aaa313aaa")
// 	sbf.ShowDebug()
// 	r = adjacentColorCount(sbf)
// 	if r[0] != 1 || r[1] != 3 || r[2] != 1 {
// 		t.Fatalf("result must be [1 3 1] but %v", r)
// 	}
// }
