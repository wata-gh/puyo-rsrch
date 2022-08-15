package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/wata-gh/puyo2"
)

func main() {
	puyo2.NewShapeBitField()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		param := strings.TrimRight(scanner.Text(), "\n")
		sbf := puyo2.NewShapeBitFieldWithFieldString(param)

		fs := sbf.FieldString()
		sbf.ExportShapeImage(0, fmt.Sprintf("1/%s_1.png", fs))
		sbf.Simulate1()
		afs := sbf.FieldString()
		fmt.Println(afs, fs)
		sbf.ExportImage(fmt.Sprintf("after/%s.png", afs))
	}
}
