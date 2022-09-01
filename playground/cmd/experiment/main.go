package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/dustin/go-humanize"
	"github.com/wata-gh/puyo2"
	"gonum.org/v1/gonum/stat/combin"
)

func checkFlippableCount() {
	fieldc := 27
	chainc := 6
	fmt.Println(len(genValidEmpties(fieldc, fieldc-(chainc*4))))
}

func checkValidEmptyCount() {
	validCount := 0
	fieldc := 27
	chainc := 4
	base := make([]int, fieldc)
	for i := 0; i < fieldc; i++ {
		base[i] = i
	}

	dir := fmt.Sprintf("chain%d", chainc)
	os.Mkdir(dir, 0755)
	combins := combin.Combinations(fieldc, fieldc-(chainc*4))
	for i, combin := range combins {
		if validEmpty(combin) {
			validCount++
			sbf := puyo2.NewShapeBitField()
			shape := puyo2.NewFieldBits()
			for _, c := range combin {
				pos := index2Field(c)
				shape.SetOnebit(pos[0], pos[1])
			}
			sbf.AddShape(shape)
			sbf.ExportImage(fmt.Sprintf("%s/%d.png", dir, i))
		}
		// _, cv := exclude(cloneArray(base), combin)
		// if validPlace(cv) {
		// 	validCount++
		// }
	}
	fmt.Println(validCount)
}

func cacheKeyExp() {
	fieldc := 27
	chainc := 3
	base := make([]int, fieldc)
	for i := 0; i < fieldc; i++ {
		base[i] = i
	}
	// base = []int{19, 20, 21, 22, 23, 24, 25, 26}
	base = []int{15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26}
	// base = []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26}
	// base = []int{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26}
	colorcs := []int{}
	for i := 0; i < chainc; i++ {
		colorcs = append(colorcs, 4)
	}

	allCache := make(map[string]int)
	cache := make([]map[int]struct{}, chainc)
	for i := 0; i < len(cache); i++ {
		cache[i] = map[int]struct{}{}
	}

	genColorCombinations(len(base), base, 0, allCache, cache, []int{}, colorcs)
	fmt.Println(len(allCache))
	for key, value := range allCache {
		fmt.Println(key, value)
	}
	fmt.Println(len(cache))
	printMemInfo()
}

func printMemInfo() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// ヒープ上に割り当てられたオブジェクト累積メモリ量
	fmt.Fprintf(os.Stderr, "MAlloc : %v\n", humanize.Bytes(m.Mallocs))
	// ヒープ上から開放されたオブジェクト数
	fmt.Fprintf(os.Stderr, "Frees : %v\n", m.Frees)

	// ヒープ上に割り当てられたオブジェクトメモリ量
	fmt.Fprintf(os.Stderr, "Alloc : %v\n", humanize.Bytes(m.Alloc))
	fmt.Fprintf(os.Stderr, "HeapAlloc : %v\n", humanize.Bytes(m.HeapAlloc))
	// ヒープ上に割り当てられたオブジェクトメモリ量。ただし開放されたオブジェクト分も含む
	fmt.Fprintf(os.Stderr, "TotalAlloc : %v\n", humanize.Bytes(m.TotalAlloc))

	// OSから割り当てられたプロセスの総メモリ量
	// ヒープ + スタック + その他
	fmt.Fprintf(os.Stderr, "Sys : %v\n", humanize.Bytes(m.Sys))

	// ポインタのルックアップ数
	fmt.Fprintf(os.Stderr, "Lookups : %v\n", m.Lookups)

	// 到達可能、あるいはGCによって解放されていないヒープオブジェクトメモリ量
	fmt.Fprintf(os.Stderr, "HeapAlloc : %v\n", humanize.Bytes(m.HeapAlloc))
	// 未使用ヒープ領域メモリ量
	fmt.Fprintf(os.Stderr, "HeapIdle : %v\n", humanize.Bytes(m.HeapIdle))
	// 使用中ヒープ領域メモリ量
	fmt.Fprintf(os.Stderr, "HeapInuse : %v\n", humanize.Bytes(m.HeapInuse))
	// OSに返却される物理メモリ量
	fmt.Fprintf(os.Stderr, "HeapReleased : %v\n", humanize.Bytes(m.HeapReleased))
	// ヒープに割り当てられたオブジェクト量
	fmt.Fprintf(os.Stderr, "HeapObjects : %v\n", m.HeapObjects)

	// 使用中スタック領域メモリ量
	fmt.Fprintf(os.Stderr, "StackInuse : %v\n", humanize.Bytes(m.StackInuse))
	// OSから割り当てられたスタック領域メモリ量
	fmt.Fprintf(os.Stderr, "StackSys : %v\n", humanize.Bytes(m.StackSys))

	// 割り当てられたmspan構造体バイト数
	fmt.Fprintf(os.Stderr, "MSpanInuse : %v\n", humanize.Bytes(m.MSpanInuse))
	// OSから取得したmspan構造体バイト数
	fmt.Fprintf(os.Stderr, "MSpanSys : %v\n", humanize.Bytes(m.MSpanSys))

	// 割り当てられたmcache構造体バイト数
	fmt.Fprintf(os.Stderr, "MCacheInuse : %v\n", humanize.Bytes(m.MCacheInuse))
	// OSから取得したmcache構造体バイト数
	fmt.Fprintf(os.Stderr, "MCacheSys : %v\n", humanize.Bytes(m.MCacheSys))
	grcn := runtime.NumGoroutine()
	fmt.Fprintf(os.Stderr, "Goroutine : %v\n", grcn)
}

func fields() [][]int {
	results := [][]int{}
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			for z := 0; z < 10; z++ {
				if x+y+z == 20 {
					results = append(results, []int{x, y, z, 0, 0, 0})
				}
			}
		}
	}
	return results
}

func fields2() [][]int {
	return [][]int{{3, 3, 2, 0, 0, 0}}
	results := [][]int{}
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			for z := 0; z < 10; z++ {
				if x+y+z == 8 {
					results = append(results, []int{x, y, z, 0, 0, 0})
				}
			}
		}
	}
	return results
}

func fields3() [][]int {
	results := [][]int{}
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			for z := 0; z < 10; z++ {
				if x+y+z == 12 {
					results = append(results, []int{x, y, z, 0, 0, 0})
				}
			}
		}
	}
	return results
}

func fields4() [][]int {
	results := [][]int{}
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			for z := 0; z < 10; z++ {
				if x+y+z == 16 {
					results = append(results, []int{x, y, z, 0, 0, 0})
				}
			}
		}
	}
	return results
}

func fillSearch() {
	var wg sync.WaitGroup
	c := make(chan []int)
	wg.Add(12)
	for i := 0; i < 12; i++ {
		go FillSearch(c, &wg)
	}
	for _, field := range fields2() { // CHANGE HERE
		c <- field
	}
	for i := 0; i < 12; i++ {
		c <- []int{}
	}
	wg.Wait()
}

func main() {
	// checkValidEmptyCount()
	// checkFlippableCount()
	// cacheKeyExp()
	// clusters := [][]int{
	// 	{1, 0, 2, 1, 1},
	// 	{1, 0, 2, 1, 1},
	// 	{1, 0, 2, 1, 1},
	// 	{2, 0, 1, 2, 1},
	// 	{3, 1, 2, 2},
	// }

	// Fill(clusters)
	// field := []int{7, 7, 6, 0, 0, 0}
	fillSearch()
}
