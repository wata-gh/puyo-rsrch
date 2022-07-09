package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"gonum.org/v1/gonum/stat/combin"
)

type conditions struct {
	pattern *Pattern
	base    []int
	colorc  []int
	field   chan<- []int
	end     bool
}

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

func intArray2Bit(cv []int) int {
	bit := 0
	for _, c := range cv {
		bit |= 1 << c
	}
	return bit
}

func cloneArray(ints []int) []int {
	newInts := make([]int, len(ints))
	copy(newInts, ints)
	return newInts
}

func genColorCombinations(pattern *Pattern, fieldc int, base []int, n int, allCache map[string]struct{}, cache []map[int]struct{}, bits []int, colorcs []int, field chan<- []int) {
	if len(colorcs) == 0 {
		(*pattern).AddExecCombi()
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
		newBase, cv := exclude(cloneArray(base), combin)
		bit := intArray2Bit(cv)

		// invalid place cache check
		if _, ok := cache[(*pattern).ChainC()-1][bit]; ok {
			(*pattern).AddExecCombi()
			(*pattern).AddCacheSkip()
			continue
		}

		if n != 0 {
			if _, ok := cache[n-1][bit]; ok {
				(*pattern).AddExecCombi()
				(*pattern).AddCacheSkip()
				continue
			}
		}

		// check valid place to put puyos
		if (*pattern).ValidPlace(cv) == false {
			cache[(*pattern).ChainC()-1][bit] = struct{}{}
			(*pattern).AddExecCombi()
			(*pattern).AddInvalidPlace()
			continue
		}
		newBits := append(cloneArray(bits), bit)
		if n < (*pattern).ChainC()-1 {
			cache[n][bit] = struct{}{}
		}

		genColorCombinations(pattern, fieldc, newBase, n+1, allCache, cache, newBits, colorcs[1:], field)
	}
}

func handleCondition(condition <-chan conditions, wg *sync.WaitGroup) {
	for {
		cond := <-condition
		if cond.end {
			break
		}
		cache := make([]map[int]struct{}, (*cond.pattern).ChainC())
		for i := 0; i < len(cache); i++ {
			cache[i] = map[int]struct{}{}
		}
		allCache := make(map[string]struct{}, (*cond.pattern).ChainC())
		genColorCombinations(cond.pattern, (*cond.pattern).FieldC(), cond.base, 0, allCache, cache, []int{}, cond.colorc, cond.field)
	}
	wg.Done()
}

func genCombinations(pattern *Pattern, base []int, colorc []int, condition chan<- conditions, field chan<- []int, grc int) {
	ctotal := 0
	for _, c := range colorc {
		ctotal += c
	}
	fieldc := (*pattern).FieldC()
	(*pattern).AddCombi(calcCombiCount(fieldc, ctotal, colorc))

	if fieldc > ctotal {
		emptycs := combin.Combinations(len(base), fieldc-ctotal)
		for _, emptyc := range emptycs {
			if (*pattern).ValidEmpty(emptyc) == false {
				(*pattern).AddInvalidEmpty()
				continue
			}
			base, _ := exclude(base, emptyc)
			cond := conditions{
				pattern: pattern,
				base:    base,
				colorc:  colorc,
				field:   field,
				end:     false,
			}
			condition <- cond
		}
	} else {
		cond := conditions{
			pattern: pattern,
			base:    base,
			colorc:  colorc,
			field:   field,
			end:     false,
		}
		condition <- cond
	}
	for i := 0; i < grc; i++ {
		cond := conditions{
			end: true,
		}
		condition <- cond
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
	var wg sync.WaitGroup
	fieldc := (*pattern).FieldC()
	colorc := []int{}
	base := []int{}
	for i := 0; i < fieldc; i++ {
		base = append(base, i)
	}
	for i := 0; i < (*pattern).ChainC(); i++ {
		colorc = append(colorc, 4)
	}

	condition := make(chan conditions)
	conditionGrc := 2
	wg.Add(conditionGrc)
	for i := 0; i < conditionGrc; i++ {
		go handleCondition(condition, &wg)
	}

	genCombinations(pattern, base, colorc, condition, field, conditionGrc)

	for i := 0; i < grc; i++ {
		field <- []int{}
	}
}
