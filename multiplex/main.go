package main

import (
	"fmt"
	"sync"

	"github.com/wata-gh/puyo2"
)

const CHAINC = 2

func index2field(idx int) return [2]int {
	return [2]int{
		i / 4,
		i % 4,
	}
}

func check(field <-chan []int, wg *sync.WaitGroup) {
	for {
		puyos := <-field
		if len(puyos) == 0 {
			break
		}
		bf := puyo2.NewBitField()
		for _, puyo := range puyos {
			for i := 0; puyo > 0; i++ {
				pos := index2field(i)

				puyo <<= 1
			}
		}
		fmt.Println(puyos)
	}
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	field := make(chan []int)
	wg.Add(1)
	go check(field, &wg)
	Gen(field, 1)
	wg.Wait()
}
