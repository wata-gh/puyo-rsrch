package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/wata-gh/puyo2"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\n")
		m := [3][2]uint64{}
		idx := 0
		for _, part := range strings.Split(line, "_") {
			if part == "" {
				continue
			}
			ms := strings.Split(part, "-")
			f0, err := strconv.ParseUint(ms[0], 10, 64)
			if err != nil {
				panic(err)
			}
			f1, err := strconv.ParseUint(ms[1], 10, 64)
			if err != nil {
				panic(err)
			}
			m[idx] = [2]uint64{
				f0,
				f1,
			}
			idx++
		}
		fmt.Println(m)
		bf := puyo2.NewBitFieldWithM(m)
		bf.ShowDebug()
	}

}
