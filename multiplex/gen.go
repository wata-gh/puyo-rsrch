package main

import (
	"fmt"
	"sort"
	"strings"

	"gonum.org/v1/gonum/stat/combin"
)

func cacheKey(bits []int) string {
	sort.Ints(bits)
	var b strings.Builder
	for _, bit := range bits {
		fmt.Fprintf(&b, "_%d", bit)
	}
	return b.String()
}

func calcCombiCount(fieldc int, ctotal int, colorc []int) int {
	f := fieldc
	emptyc := f - ctotal
	combi := combination(f, emptyc)
	f -= emptyc
	for _, c := range colorc {
		combi *= combination(f, c)
		f -= c
	}
	return combi
}

func combination(n int, k int) int {
	return permutation(n, k) / factorial(k)
}

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

func factorial(n int) int {
	return permutation(n, n-1)
}

func genColorCombinations(pattern *Pattern, fieldc int, base []int, n int, allCache map[string]struct{}, cache []map[int]struct{}, bits []int, colorcs []int, field chan<- []int) {
	if len(colorcs) == 0 {
		key := cacheKey(bits)
		if _, ok := allCache[key]; ok {
			(*pattern).AddCacheSkip()
			return
		}
		allCache[key] = struct{}{}
		field <- bits
		return
	}

	combins := combin.Combinations(len(base), colorcs[:1][0])
	for _, combin := range combins {
		if n == (*pattern).ChainC()-1 {
			(*pattern).AddExecCombi()
		}
		newBits := make([]int, len(bits))
		copy(newBits, bits)
		newBase := make([]int, len(base))
		copy(newBase, base)
		newBase, cv := exclude(newBase, combin)

		bit := 0
		for _, c := range cv {
			bit |= 1 << c
		}
		// invalid place cache check
		if _, ok := cache[(*pattern).ChainC()-1][bit]; ok {
			(*pattern).AddCacheSkip()
			continue
		}

		if n > 0 {
			if _, ok := cache[n-1][bit]; ok {
				if n != (*pattern).ChainC()-1 {
					(*pattern).AddExecCombi()
				}
				(*pattern).AddCacheSkip()
				continue
			}
		}

		if (*pattern).ValidPlace(cv) == false {
			cache[(*pattern).ChainC()-1][bit] = struct{}{}
			if n != (*pattern).ChainC()-1 {
				(*pattern).AddExecCombi()
			}
			(*pattern).AddInvalidPlace()
			continue
		}
		newBits = append(newBits, bit)
		if n < (*pattern).ChainC()-1 {
			cache[n][bit] = struct{}{}
		}

		genColorCombinations(pattern, fieldc, newBase, n+1, allCache, cache, newBits, colorcs[1:], field)
	}
}

func genCombinations(pattern *Pattern, base []int, colorc []int, field chan<- []int) {
	ctotal := 0
	for _, c := range colorc {
		ctotal += c
	}
	fieldc := (*pattern).FieldC()
	(*pattern).AddCombi(calcCombiCount(fieldc, ctotal, colorc))

	if fieldc > ctotal {
		emptycs := combin.Combinations(len(base), fieldc-ctotal)
		for _, emptyc := range emptycs {
			if (*pattern).ValidEmpty(emptyc) {
				cache := make([]map[int]struct{}, (*pattern).ChainC())
				for i := 0; i < len(cache); i++ {
					cache[i] = map[int]struct{}{}
				}
				allCache := make(map[string]struct{}, (*pattern).ChainC())
				base, _ := exclude(base, emptyc)
				genColorCombinations(pattern, fieldc, base, 0, allCache, cache, []int{}, colorc, field)
			} else {
				(*pattern).AddInvalidEmpty()
			}
		}
	} else {
		cache := make([]map[int]struct{}, (*pattern).ChainC())
		for i := 0; i < len(cache); i++ {
			cache[i] = map[int]struct{}{}
		}
		allCache := make(map[string]struct{}, (*pattern).ChainC())
		genColorCombinations(pattern, fieldc, base, 0, allCache, cache, []int{}, colorc, field)
	}
}

func permutation(n int, k int) int {
	v := 1
	if 0 < k && k <= n {
		for i := 0; i < k; i++ {
			v *= (n - i)
		}
	} else if k > n {
		v = 0
	}
	return v
}

func Gen(pattern *Pattern, field chan<- []int, grc int) {
	fieldc := (*pattern).FieldC()
	colorc := []int{}
	base := []int{}
	for i := 0; i < fieldc; i++ {
		base = append(base, i)
	}
	for i := 0; i < (*pattern).ChainC(); i++ {
		colorc = append(colorc, 4)
	}

	genCombinations(pattern, base, colorc, field)

	for i := 0; i < grc; i++ {
		field <- []int{}
	}
}
