package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/wata-gh/puyo2"
)

func yoko3() {
	scanner := bufio.NewScanner(os.Stdin)
	target := puyo2.NewFieldBits()
	os.Mkdir("yoko3/", 0755)
	target.SetOnebit(0, 1)
	target.SetOnebit(1, 1)
	target.SetOnebit(2, 1)
	for scanner.Scan() {
		param := strings.TrimRight(scanner.Text(), "\n")
		sbf := puyo2.NewShapeBitFieldWithFieldString(param)
		fb := sbf.Shapes[len(sbf.Shapes)-1]
		if fb.And(target).Equals(target) {
			sbf.ExportImage("yoko3/" + sbf.FieldString() + ".png")
		}
	}
}

func groupByBottomShape() {
	scanner := bufio.NewScanner(os.Stdin)
	os.Mkdir("results/", 0755)
	target := puyo2.NewFieldBits()
	target.SetOnebit(0, 1)
	target.SetOnebit(1, 1)
	target.SetOnebit(2, 1)

	for scanner.Scan() {
		param := strings.TrimRight(scanner.Text(), "\n")
		sbf := puyo2.NewShapeBitFieldWithFieldString(param)
		fb := sbf.Shapes[len(sbf.Shapes)-1]
		can := fb.And(target)
		ex := can.Expand(fb)
		var s strings.Builder
		for _, i := range ex.ToIntArray() {
			fmt.Fprintf(&s, "_%d", i)
		}
		dir := fmt.Sprintf("results/%s", s.String())
		os.Mkdir(dir, 0755)
		fmt.Printf("%s/%s.png\n", dir, sbf.FieldString())
		sbf.ExportImage(fmt.Sprintf("%s/%s.png", dir, sbf.FieldString()))
	}
}

func main() {
	groupByBottomShape()
}
