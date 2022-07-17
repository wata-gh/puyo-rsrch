package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

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
		fmt.Fprintf(&b, "_%x", bit)
	}
	return b.String()
}

func calcCombiCount(fieldc int, ctotal int, colorc []int) int {
	f := fieldc
	emptyc := f - ctotal
	combi := combin.Binomial(f, emptyc)
	f -= emptyc
	for _, c := range colorc {
		combi *= combin.Binomial(f, c)
		f -= c
	}
	return combi
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

func genColorCombinations(pattern *Pattern, fieldc int, base []int, n int, allCache map[string]int, cache []map[int]struct{}, bits []int, colorcs []int, field chan<- []int) {
	if len(colorcs) == 0 {
		(*pattern).Incr("ExecCombiCount")
		field <- bits
		return
	}

	combins := combin.Combinations(len(base), colorcs[:1][0])
	for _, combin := range combins {
		newBase, cv := exclude(cloneArray(base), combin)
		bit := intArray2Bit(cv)
		ncv := cloneArray(cv)
		k := ncv[0]
		ncv[0] %= 3
		for i := 1; i < len(ncv); i++ {
			ncv[i] -= k - ncv[0]
		}
		nbit := intArray2Bit(ncv)

		// invalid place cache check
		if _, ok := cache[(*pattern).ChainC()-1][nbit]; ok {
			(*pattern).Incr("ExecCombiCount")
			(*pattern).Incr("CacheSkipCount")
			continue
		}

		// check valid place to put puyos
		if (*pattern).ValidPlace(ncv) == false {
			cache[(*pattern).ChainC()-1][nbit] = struct{}{}
			(*pattern).Incr("ExecCombiCount")
			(*pattern).Incr("InvalidPlaceCount")
			continue
		}

		newBits := append(cloneArray(bits), bit)

		key := cacheKey(newBits)
		c, ok := allCache[key]
		if ok {
			(*pattern).Incr("AllCacheSkipCount")
			allCache[key] += 1
			if c == n {
				delete(allCache, key)
			}
			return
		}

		genColorCombinations(pattern, fieldc, newBase, n+1, allCache, cache, newBits, colorcs[1:], field)
		if !ok && n < (*pattern).ChainC()-2 {
			allCache[key] = 0
		}
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
		allCache := make(map[string]int)
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
	(*pattern).Add("CombiCount", calcCombiCount(fieldc, ctotal, colorc))

	if fieldc > ctotal {
		emptycs := (*pattern).GenValidEmpties(pattern, base, fieldc, ctotal)
		emptycsLen := len(emptycs)
		fmt.Fprintf(os.Stderr, "valid empties: %d\n", emptycsLen)
		for i, emptyc := range emptycs {
			t := time.Now()
			fmt.Fprintf(os.Stderr, "[%s] emptyc: %d / %d (%0.2f%%)\n", t.Format("2006-01-02 15:04:05"), i+1, emptycsLen, float64((i+1)*100)/float64(emptycsLen))
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

	wg.Wait()
	for i := 0; i < grc; i++ {
		field <- []int{}
	}
}
