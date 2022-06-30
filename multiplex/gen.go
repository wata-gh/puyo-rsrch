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

func genCombinations(fieldc int, base []int, cache map[[CHAINC]int]struct{}, colorc []int, field chan<- []int) {
	ctotal := 0
	for _, c := range colorc {
		ctotal += c
	}

	if fieldc > ctotal {
		emptycs := combin.Combinations(len(base), fieldc-ctotal)
		for _, emptyc := range emptycs {
			if validEmpty(emptyc) {
				base, _ := exclude(base, emptyc)
				genColorCombinations(fieldc, base, cache, []int{}, colorc, field)
			}
		}
	} else {
		genColorCombinations(fieldc, base, cache, []int{}, colorc, field)
	}
}

func genColorCombinations(fieldc int, base []int, cache map[[CHAINC]int]struct{}, bits []int, colorcs []int, field chan<- []int) {
	if len(colorcs) == 0 {
		sort.Ints(bits)
		key := *(*[CHAINC]int)(bits)
		if _, ok := cache[key]; !ok {
			cache[key] = struct{}{}
			field <- bits
		}
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
		newBits = append(newBits, bit)
		genColorCombinations(fieldc, newBase, cache, newBits, colorcs[1:], field)
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

	genCombinations(fieldc, base, map[[CHAINC]int]struct{}{}, colorc, field)

	for i := 0; i < grc; i++ {
		field <- []int{}
	}
}
