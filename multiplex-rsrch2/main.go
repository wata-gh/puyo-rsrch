package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/wata-gh/puyo2"
)

func embeddedFields(chains int) [][]int {
	results := [][]int{}
	for x := 0; x < 10; x++ {
		for y := 0; y <= 10; y++ {
			for z := 0; z <= 10; z++ {
				if x+y+z == chains*4 {
					results = append(results, []int{x + 1, y, z, 0, 0, 0})
				}
			}
		}
	}
	return results
}

func fields(chains int) [][]int {
	results := [][]int{}
	for x := 0; x < 10; x++ {
		for y := 0; y <= 10; y++ {
			for z := 0; z <= 10; z++ {
				if x+y+z == chains*4 {
					results = append(results, []int{x, y, z, 0, 0, 0})
				}
			}
		}
	}
	return results
}

func rensabiFields(chains int) [][]int {
	results := [][]int{}
	for w := 2; w <= 3; w++ {
		for x := 0; x <= 5; x++ {
			for y := 0; y <= 5; y++ {
				for z := 0; z <= 5; z++ {
					if w+x+y+z == chains*4 {
						results = append(results, []int{0, 0, w, x, y, z})
					}
				}
			}
		}
	}
	return results
}

func gtrKeyPuyo() *puyo2.FieldBits {
	keyPuyo := puyo2.NewFieldBits()
	keyPuyo.SetOnebit(0, 1)
	return keyPuyo
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
	initSbf := puyo2.NewShapeBitField()
	// initSbf.AddShape(gtrKeyPuyo())
	for i := 0; i < grc; i++ {
		go FillSearch(c, chainc, []int{0, 0, 0, 0, 0, 0}, initSbf, &wg)
	}
	for _, field := range rensabiFields(chainc) {
		c <- field
	}
	for i := 0; i < grc; i++ {
		c <- []int{}
	}
	wg.Wait()
}
