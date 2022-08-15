package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

func fields(chains int) [][]int {
	results := [][]int{}
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			for z := 0; z < 10; z++ {
				if x+y+z == chains*4 && y != 0 {
					results = append(results, []int{x, y, z, 0, 0, 0})
				}
			}
		}
	}
	return results
}

func main() {
	var wg sync.WaitGroup
	c := make(chan []int)
	grc := 8
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage %s chains\n", os.Args[0])
		return
	}
	chainc, _ := strconv.Atoi(os.Args[1])
	wg.Add(grc)
	for i := 0; i < grc; i++ {
		go FillSearch(c, chainc, &wg)
	}
	for _, field := range fields(chainc) {
		c <- field
	}
	for i := 0; i < grc; i++ {
		c <- []int{}
	}
	wg.Wait()
}
