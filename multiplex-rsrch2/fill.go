package main

import (
	"fmt"
	"math/bits"
	"os"
	"sync"

	"github.com/wata-gh/puyo2"
)

var clusters [][]int = [][]int{
	{0, 0, 1, 1, 1, 1},
	{0, 1, 1, 1, 1, 1},
	{0, 2, 1, 1, 1, 1},

	{1, 0, 2, 1, 1},
	{1, 1, 2, 1, 1},
	{1, 2, 2, 1, 1},
	{1, 3, 2, 1, 1},

	{2, 0, 1, 1, 2},
	{2, 1, 1, 1, 2},
	{2, 2, 1, 1, 2},
	{2, 3, 1, 1, 2},

	{3, 0, 1, 2, 1},
	{3, 1, 1, 2, 1},
	{3, 2, 1, 2, 1},
	{3, 3, 1, 2, 1},

	{4, 0, 3, 1},
	{4, 1, 3, 1},
	{4, 2, 3, 1},
	{4, 3, 3, 1},
	{4, 4, 3, 1},

	{5, 0, 1, 3},
	{5, 1, 1, 3},
	{5, 2, 1, 3},
	{5, 3, 1, 3},
	{5, 4, 1, 3},

	{6, 0, 2, 2},
	{6, 1, 2, 2},
	{6, 2, 2, 2},
	{6, 3, 2, 2},
	{6, 4, 2, 2},

	{7, 0, 4},
	{7, 1, 4},
	{7, 2, 4},
	{7, 3, 4},
	{7, 4, 4},
	{7, 5, 4},
}

func add(mem []int, cluster []int) {
	offset := cluster[1]
	for i, v := range cluster[2:] {
		mem[i+offset] += v
	}
}
func checkOverflow(mem []int, field []int) bool {
	for i, v := range mem {
		if field[i] < v {
			return false
		}
	}
	return true
}

func fill(field []int, chainc int, idx int, mem []int, result [][]int, results Results) Results {
ClusterLoop:
	for _, c := range clusters {
		m := make([]int, len(mem))
		copy(m, mem)
		add(m, c)
		if checkOverflow(m, field) == false {
			continue
		}
		r := make([][]int, len(result))
		copy(r, result)
		r = append(r, c)
		var res Result
		res.results = r

		if idx == chainc {
			for i, v := range m {
				if field[i] != v {
					continue ClusterLoop
				}
			}
			results = append(results, res)
		} else {
			results = fill(field, chainc, idx+1, m, r, results)
		}
	}
	return results
}

func fillUp(field []int, chainc int) Results {
	var results Results
	results = fill(field, chainc, 1, []int{0, 0, 0, 0, 0, 0}, [][]int{}, results)
	return results
}

func FillSearch(fields chan []int, chainc int, wg *sync.WaitGroup) {
	for {
		field := <-fields
		if len(field) == 0 {
			break
		}
		fmt.Fprintln(os.Stderr, field)
		patterns := fillUp(field, chainc)
		fmt.Fprintf(os.Stderr, "%v %d\n", field, len(patterns))
		for i, pattern := range patterns {
			fmt.Fprintf(os.Stderr, "%v %d/%d\n", field, i, len(patterns))
			Fill(pattern, chainc)
		}
	}
	wg.Done()
}

func shapes(c int, x int) []*puyo2.FieldBits {
	var shapes []*puyo2.FieldBits
	switch c {
	case 0:
		s := puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+2, 1)
		s.SetOnebit(x+3, 1)
	case 1:
		s := puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x, 2)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+2, 1)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x, 2)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+2, 2)
		shapes = append(shapes, s)
	case 2:
		s := puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+2, 1)
		s.SetOnebit(x+2, 2)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 2)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+2, 1)
		s.SetOnebit(x+2, 2)
		shapes = append(shapes, s)
	case 3:
		s := puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+2, 1)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 2)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+2, 2)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 2)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+2, 1)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+2, 2)
		shapes = append(shapes, s)
	case 4:
		s := puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x, 2)
		s.SetOnebit(x, 3)
		s.SetOnebit(x+1, 1)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x, 2)
		s.SetOnebit(x, 3)
		s.SetOnebit(x+1, 2)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x, 2)
		s.SetOnebit(x, 3)
		s.SetOnebit(x+1, 3)
		shapes = append(shapes, s)
	case 5:
		s := puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+1, 3)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 2)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+1, 3)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 3)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+1, 3)
		shapes = append(shapes, s)
	case 6:
		s := puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x, 2)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+1, 2)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x, 2)
		s.SetOnebit(x+1, 2)
		s.SetOnebit(x+1, 3)
		shapes = append(shapes, s)
		s = puyo2.NewFieldBits()
		s.SetOnebit(x, 2)
		s.SetOnebit(x, 3)
		s.SetOnebit(x+1, 1)
		s.SetOnebit(x+1, 2)
		shapes = append(shapes, s)
	case 7:
		s := puyo2.NewFieldBits()
		s.SetOnebit(x, 1)
		s.SetOnebit(x, 2)
		s.SetOnebit(x, 3)
		s.SetOnebit(x, 4)
		shapes = append(shapes, s)
	}
	return shapes
}

func willNotDrop(y int, shape *puyo2.FieldBits, overall *puyo2.FieldBits) bool {
	for x := 0; x < 6; x++ {
		if shape.Onebit(x, y+1) > 0 { // shape あり
			if overall.Onebit(x, y) == 0 { // 下なし
				return false
			}
		} else {
			col := shape.ColBits(x)
			if col > 0 {
				for z := y; z < bits.TrailingZeros64(col>>(x*16)); z++ {
					if overall.Onebit(x, z) == 0 {
						return false
					}
				}
			}
		}
	}
	return true
}

func place(osbf *puyo2.ShapeBitField, clusters [][]int, chainc int, last *puyo2.FieldBits) {
	if len(clusters) == 0 {
		nsbf := osbf.Clone()
		result := nsbf.Simulate()
		if result.Chains == chainc {
			fmt.Println(nsbf.ChainOrderedFieldString())
		}
		return
	}
	cluster := clusters[0]
	overall := osbf.OverallShape()
	for x := 0; x < 6; x++ {
		overall.SetOnebit(x, 0)
	}
	for _, shape := range shapes(cluster[0], cluster[1]) {
		for yOffset := 0; yOffset < 13; yOffset++ {
			s := shape.FastLift(yOffset)
			if s.Equals(s.MaskField13()) == false {
				continue
			}
			if willNotDrop(yOffset, s, overall) {
				sfb := osbf.Clone()
				sfb.InsertShape(s)
				place(sfb, clusters[1:], chainc, shape)
			}
		}
	}
}

func Fill(result Result, chainc int) {
	sbf := puyo2.NewShapeBitField()
	place(sbf, result.results, chainc, puyo2.NewFieldBits())
}

// func removeDuplication(results Results) Results {
// 	m := map[string]bool{}
// 	var uniq Results
// 	for _, result := range results {
// 		sort.Sort(result)
// 		key := result.ToString()
// 		if m[key] == false {
// 			m[key] = true
// 			uniq = append(uniq, result)
// 		}
// 	}
// 	return uniq
// }
