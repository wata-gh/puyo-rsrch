package main

import (
	"sort"

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

func genColorCombinations(fieldc int, base []int, n int, cache []map[int]struct{}, bits []int, colorcs []int, field chan<- []int) {
	if len(colorcs) == 0 {
		sort.Ints(bits)
		field <- bits
		return
	}

	combins := combin.Combinations(len(base), colorcs[:1][0])
	for _, combin := range combins {
		newBits := make([]int, len(bits))
		copy(newBits, bits)
		newBase := make([]int, len(base))
		copy(newBase, base)
		newBase, cv := exclude(newBase, combin)
		bit := 0
		for _, c := range cv {
			bit |= 1 << c
		}
		if n > 0 {
			if _, ok := cache[n-1][bit]; ok {
				// fmt.Println("cache hit.", cv)
				continue
			}
		}
		cache[n][bit] = struct{}{}

		// if n == 0 {
		// 	cache[0][bit] = struct{}{}
		// 	// fmt.Println("cache add.", cv)
		// } else {
		// 	if _, ok := cache[0][bit]; ok {
		// 		// fmt.Println("cache hit.", cv)
		// 		continue
		// 	}
		// }
		newBits = append(newBits, bit)
		genColorCombinations(fieldc, newBase, n+1, cache, newBits, colorcs[1:], field)
	}
}

func genCombinations(fieldc int, base []int, colorc []int, field chan<- []int) {
	ctotal := 0
	for _, c := range colorc {
		ctotal += c
	}

	if fieldc > ctotal {
		emptycs := combin.Combinations(len(base), fieldc-ctotal)
		for _, emptyc := range emptycs {
			if validEmpty(emptyc) {
				cache := make([]map[int]struct{}, CHAINC)
				for i := 0; i < len(cache); i++ {
					cache[i] = map[int]struct{}{}
				}
				base, _ := exclude(base, emptyc)
				genColorCombinations(fieldc, base, 0, cache, []int{}, colorc, field)
			}
		}
	} else {
		cache := make([]map[int]struct{}, CHAINC)
		for i := 0; i < len(cache); i++ {
			cache[i] = map[int]struct{}{}
		}
		genColorCombinations(fieldc, base, 0, cache, []int{}, colorc, field)
	}
}

func Gen(field chan<- []int, grc int) {
	fieldc := FIELDC
	colorc := []int{}
	base := []int{}
	for i := 0; i < fieldc; i++ {
		base = append(base, i)
	}
	for i := 0; i < CHAINC; i++ {
		colorc = append(colorc, 4)
	}

	cache := make([]map[int]struct{}, CHAINC)
	for i := 0; i < len(cache); i++ {
		cache[i] = map[int]struct{}{}
	}
	genCombinations(fieldc, base, colorc, field)

	for i := 0; i < grc; i++ {
		field <- []int{}
	}
}
