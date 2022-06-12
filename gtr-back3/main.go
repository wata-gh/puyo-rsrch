package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

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

func createGtr(bf *puyo2.BitField) {
	gtrColor := bf.Color(2, 2)
	for _, c := range []puyo2.Color{puyo2.Red, puyo2.Green, puyo2.Yellow, puyo2.Blue} {
		if gtrColor != c {
			bf.SetColor(c, 0, 3)
			bf.SetColor(c, 0, 2)
			bf.SetColor(c, 1, 2)
			return
		}
	}
	bf.ShowDebug()
	panic("can not create GTR.")
}

func handleResult(config *puyorsrch.Config, satisfy <-chan puyo2.BitField, wg *sync.WaitGroup) {
	for {
		bf := <-satisfy
		if bf.IsEmpty() {
			break
		}
		m := grouping(&bf)
		createGtr(&bf)
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

func checkAvailColor(bf *puyo2.BitField) puyo2.Color {
	avail := map[puyo2.Color]bool{
		puyo2.Red:    true,
		puyo2.Green:  true,
		puyo2.Yellow: true,
		puyo2.Blue:   true,
	}
	for _, pos := range [3][2]int{{2, 3}, {2, 1}, {3, 2}} {
		for c := range avail {
			if bf.Color(pos[0], pos[1]) == c {
				avail[c] = false
				break
			}
		}
	}
	for c := range avail {
		if avail[c] {
			return c
		}
	}
	bf.ShowDebug()
	panic("no avaliable color.")
}

func checkBlueOnTop(bf *puyo2.BitField) bool {
	fb := bf.Bits(puyo2.Blue)
	m := fb.ToIntArray()
	m[0] <<= 1
	m[1] <<= 1
	puyos := bf.Bits(puyo2.Red).Or(bf.Bits(puyo2.Green)).Or(bf.Bits(puyo2.Yellow)).ToIntArray()
	return ((m[0]&puyos[0]) == m[0] && (m[1]&puyos[1]) == m[1]) == false
}

func check(config *puyorsrch.Config, field <-chan []int, satisfy chan<- puyo2.BitField, num int, wg *sync.WaitGroup) {
	cnt := 0
	resultFields := make(map[string]struct{})
	bbf := puyo2.NewBitField()
	for {
		puyos := <-field
		if len(puyos) == 0 {
			break
		}
		// fmt.Fprintf(os.Stderr, "[%d]%v\n", num, puyos)
		if cnt%10000000 == 0 {
			fmt.Fprintf(os.Stderr, "[%d]%s %d\n", num, time.Now().String(), cnt)
			fmt.Fprintf(os.Stderr, "[%d]%v\n", num, puyos)
		}
		bf := bbf.Clone()
		vanish := puyo2.NewFieldBits()
		for i := 0; i < len(puyos); i++ {
			if puyos[i] != 0 {
				bf.SetColor(puyo2.Color(puyos[i]), config.PuyoConfig.SearchLocations[i][0], config.PuyoConfig.SearchLocations[i][1])
			} else {
				vanish.Onebit(config.PuyoConfig.SearchLocations[i][0], config.PuyoConfig.SearchLocations[i][1])
			}
		}

		availColor := checkAvailColor(bf)
		for _, pos := range [5][2]int{{0, 1}, {1, 3}, {1, 2}, {1, 1}, {2, 2}} {
			bf.SetColor(availColor, pos[0], pos[1])
		}

		beforeDrop := bf.Clone()
		bf.Drop(vanish)
		if bf.Equals(beforeDrop) && checkBlueOnTop(bf) == false {
			result := bf.SimulateWithNewBitField()
			if result.Chains == config.PuyoConfig.ExpectedChains {
				clear := true
				for _, c := range []puyo2.Color{puyo2.Red, puyo2.Green, puyo2.Yellow} {
					if result.BitField.Bits(c).IsEmpty() == false {
						clear = false
					}
				}
				if clear {
					nbf := bf.Normalize()
					param := nbf.MattulwanEditorParam()
					if param == "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbcbdeddbbbdcbdddeeed" {
						fmt.Println("==>")
						beforeDrop.ShowDebug()
					}
					_, exist := resultFields[param]
					if !exist {
						resultFields[param] = struct{}{}
						satisfy <- *nbf
					}
				}
			}
		}
		cnt++
	}
	// 終了のために空フィールドを送る
	satisfy <- *puyo2.NewBitField()
	fmt.Fprintf(os.Stderr, "[%d]wg.Done\n", num)
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
	w.Add(8)
	go check(config, field, satisfy, 0, &w)
	go handleResult(config, satisfy, &w)
	go check(config, field, satisfy, 1, &w)
	go handleResult(config, satisfy, &w)
	go check(config, field, satisfy, 2, &w)
	go handleResult(config, satisfy, &w)
	go check(config, field, satisfy, 3, &w)
	go handleResult(config, satisfy, &w)
	Gen(field, 4)
	w.Wait()
}
