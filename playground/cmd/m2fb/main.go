package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"

	"github.com/wata-gh/puyo2"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		param := strings.TrimRight(scanner.Text(), "\n")
		ms := strings.Split(param, "_")
		m := [2]uint64{}
 		m[0], _ = strconv.ParseUint(ms[0], 10, 64)
 		m[1], _ = strconv.ParseUint(ms[1], 10, 64)
		sbf := puyo2.NewFieldBitsWithM(m)
		fmt.Println(param)
		sbf.ShowDebug()
	}
}
