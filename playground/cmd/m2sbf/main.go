package main

import (
	"bufio"
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/wata-gh/puyo2"
)

func main() {
	export := flag.Bool("e", false, "export image")
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)
	sbf := puyo2.NewShapeBitField()
	for scanner.Scan() {
		param := strings.TrimRight(scanner.Text(), "\n")
		ms := strings.Split(param, "_")
		m := [2]uint64{}
		m[0], _ = strconv.ParseUint(ms[0], 10, 64)
		m[1], _ = strconv.ParseUint(ms[1], 10, 64)
		shape := puyo2.NewFieldBitsWithM(m)
		sbf.AddShape(shape)
	}
	sbf.ShowDebug()
	if *export {
		sbf.ExportImage(sbf.FieldString() + ".png")
	}
}
