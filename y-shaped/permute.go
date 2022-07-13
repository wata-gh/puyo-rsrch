package main

import (
	"sync"
)

func ParallelPermute(nums []int, suffix int, field chan<- []int, wg *sync.WaitGroup) {
	n := len(nums)
	field <- append(append([]int{}, nums...), suffix)
	p := make([]int, n+1)
	for i := 0; i <= n; i++ {
		p[i] = i
	}
	for i := 1; i < n; {
		p[i]--
		j := i % 2 * p[i]

		if nums[i] != nums[j] {
			// fmt.Printf("%v\n", nums)
			nums[i], nums[j] = nums[j], nums[i]
			field <- append(append([]int{}, nums...), suffix)
		}
		for i = 1; p[i] == 0; i++ {
			p[i] = i
		}
	}
	field <- []int{}
	wg.Done()
}

func makeCopy(nums []int) []int {
	return append([]int{}, nums...)
}
