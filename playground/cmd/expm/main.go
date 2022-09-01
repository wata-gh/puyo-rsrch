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
		param := strings.TrimRight(scanner.Text(), "\n")
		bf := puyo2.NewBitFieldWithMattulwan(param)
		fmt.Println(bf.MattulwanEditorParam())
	}
}
