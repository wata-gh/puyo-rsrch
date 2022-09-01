package main

import (
	"gonum.org/v1/gonum/stat/combin"
)

func genColorCombinations(fieldc int, base []int, n int, allCache map[string]int, cache []map[int]struct{}, bits []int, colorcs []int) {
	if len(colorcs) == 0 {
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
		if _, ok := cache[0][nbit]; ok {
			continue
		}

		// check valid place to put puyos
		if validPlace(ncv) == false {
			cache[0][nbit] = struct{}{}
			continue
		}

		newBits := append(cloneArray(bits), bit)

		key := cacheKey(newBits)
		_, ok := allCache[key]
		if ok {
			allCache[key] += 1
			return
		}
		genColorCombinations(fieldc, newBase, n+1, allCache, cache, newBits, colorcs[1:])
		if !ok && n != 0 && n < 3 {
			allCache[key] = 0
		}
	}
}
