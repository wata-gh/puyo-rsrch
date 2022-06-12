package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/wata-gh/puyo2"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\n")
		bf := puyo2.NewBitFieldWithMattulwan(line)
		c := bf.Color(0, 3)
		bf.SetColor(c, 0, 4)
		var s strings.Builder
		for _, i := range bf.ToChainShapesUInt64Array() {
			fmt.Fprintf(&s, "_%d-%d", i[0], i[1])
		}
		fmt.Println(line, s.String())
	}
}
