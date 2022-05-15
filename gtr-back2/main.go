package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	puyorsrch "github.com/wata-gh/puyo-rsrch/common"
	"github.com/wata-gh/puyo2"
)

func grouping(bf *puyo2.BitField) [2]uint64 {
	gbf := bf.Clone()
	gbf.Simulate1()
	rvb := gbf.Bits(puyo2.Red).FindVanishingBits()
	gvb := gbf.Bits(puyo2.Green).FindVanishingBits()
	yvb := gbf.Bits(puyo2.Yellow).FindVanishingBits()
	bvb := gbf.Bits(puyo2.Blue).FindVanishingBits()
	var m [2]uint64
	if !rvb.IsEmpty() {
		m = rvb.ToIntArray()
	} else if !gvb.IsEmpty() {
		m = gvb.ToIntArray()
	} else if !yvb.IsEmpty() {
		m = yvb.ToIntArray()
	} else if !bvb.IsEmpty() {
		m = bvb.ToIntArray()
	} else {
		fmt.Println("=== Argument BitField ===")
		bf.ShowDebug()
		fmt.Println("=== Simmulated BitField ===")
		gbf.ShowDebug()
		panic("no vanishing bits.")
	}
	return m
}

func handleResult(config *puyorsrch.Config, satisfy <-chan puyo2.BitField, wg *sync.WaitGroup) {
	for {
		bf := <-satisfy
		if bf.IsEmpty() {
			break
		}
		m := grouping(&bf)
		bf.SetColor(puyo2.Red, 0, 3)
		bf.SetColor(puyo2.Red, 0, 2)
		bf.SetColor(puyo2.Red, 1, 2)
		param := bf.MattulwanEditorParam()
		fmt.Println(param)
		if config.PuyoConfig.ExportImagePath != "" {
			path := fmt.Sprintf(config.PuyoConfig.ExportImagePath+"/%d_%d", m[0], m[1])
			os.Mkdir(path, 0755)
			bf.ExportImage(path + "/" + param + ".png")
		}
	}
	fmt.Fprintln(os.Stderr, "[handleResult]wg.Done")
	wg.Done()
}

func main() {
	config := puyorsrch.ReadConfig()
	if config.PuyoConfig.ExportImagePath != "" {
		fmt.Fprintf(os.Stderr, "exporting image to `%s`\n", config.PuyoConfig.ExportImagePath)
		os.MkdirAll(config.PuyoConfig.ExportImagePath, 0755)
	}
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	var w sync.WaitGroup
	field := make(chan []int)
	satisfy := make(chan puyo2.BitField)
	defer close(field)
	w.Add(2)
	go puyorsrch.Check(config, field, satisfy, 0, &w)
	go handleResult(config, satisfy, &w)
	Gen(field)
	w.Wait()
}
