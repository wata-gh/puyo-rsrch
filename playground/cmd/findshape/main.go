package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/wata-gh/puyo2"
)

func main() {
	target := puyo2.NewFieldBits()
	target.SetOnebit(1, 1)
	target.SetOnebit(1, 2)
	target.SetOnebit(1, 3)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		param := strings.TrimRight(scanner.Text(), "\n")
		sbf := puyo2.NewShapeBitFieldWithFieldString(param)
		for _, shape := range sbf.Shapes {
			and := shape.And(target)
			if and.Equals(target) {
				fmt.Println(param)
			}
		}
	}
}
