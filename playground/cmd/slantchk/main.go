package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"

	"github.com/wata-gh/puyo2"
)

func check(bf *puyo2.BitField) bool {
	fb := bf.Bits(puyo2.Red).Or(bf.Bits(puyo2.Blue).Or(bf.Bits(puyo2.Green).Or(bf.Bits(puyo2.Yellow))))
	last := 6
	for i := 5; i >= 0; i-- {
		b := fb.ColBits(i)
		c := bits.OnesCount64(b)
		if last >= c {
			last = c
		} else {
			return false
		}
	}
	return true
}
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		param := strings.TrimRight(scanner.Text(), "\n")
		bf := puyo2.NewBitFieldWithMattulwan(param)
		if check(bf) {
			fmt.Println(param)
		}
	}
}
