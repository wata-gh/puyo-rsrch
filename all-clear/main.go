package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/wata-gh/puyo2"
)

func cntColor(haipuyo []rune) map[rune]int {
	counter := map[rune]int{'r': 0, 'g': 0, 'b': 0, 'y': 0, 'p': 0}
	for _, c := range haipuyo {
		counter[c]++
	}
	return counter
}

func checkAllClear8Impossible(n int, haipuyo []rune, table map[rune]int) bool {
	for _, v := range cntColor(haipuyo) {
		if v != 0 && v < 4 {
			return false
		}
	}

	puyoSets := []puyo2.PuyoSet{}
	for i := 0; i < len(haipuyo); i += 2 {
		puyoSet := puyo2.PuyoSet{
			Axis: puyo2.Color(table[haipuyo[i]]),
			Child: puyo2.Color(table[haipuyo[i+1]]),
		}
		puyoSets = append(puyoSets, puyoSet)
	}

	allClear := false
	bf := puyo2.NewBitField()
	hands := []puyo2.Hand{}
	bf.SearchWithPuyoSets(puyoSets, hands, func(sr *puyo2.SearchResult) bool {
		if sr.RensaResult.BitField.IsEmpty() {
			allClear = true
			return false
		}
		return true
	}, 1)

	if allClear == false {
		fmt.Printf("%d/%v\n", n, puyoSets)
	}
	return allClear
}

func checkAllClearImpossible(n int, haipuyo []rune, table map[rune]int) bool {
	for _, v := range cntColor(haipuyo) {
		if v != 0 && v < 4 {
			return false
		}
	}

	puyoSets := []puyo2.PuyoSet{}
	for i := 0; i < len(haipuyo); i += 2 {
		puyoSet := puyo2.PuyoSet{
			Axis: puyo2.Color(table[haipuyo[i]]),
			Child: puyo2.Color(table[haipuyo[i+1]]),
		}
		puyoSets = append(puyoSets, puyoSet)
	}

	allClear := false
	bf := puyo2.NewBitField()
	hands := []puyo2.Hand{}
	bf.SearchWithPuyoSets(puyoSets, hands, func(sr *puyo2.SearchResult) bool {
		if sr.Depth == len(puyoSets) && sr.RensaResult.Chains == len(haipuyo)/4 && sr.RensaResult.BitField.IsEmpty() {
			fmt.Printf("%d %v\n", n, sr.Hands)
			allClear = true
			return false
		}
		return true
	}, 1)

	if allClear == false {
		fmt.Printf("%d/%v\n", n, puyoSets)
	}
	return allClear
}

func color2Str(puyo puyo2.Color) string {
	switch puyo {
	case puyo2.Red:
		return "r"
	case puyo2.Blue:
		return "b"
	case puyo2.Green:
		return "g"
	case puyo2.Yellow:
		return "y"
	}
	panic("color can be r,b,g,y")
}


func hands2SimpleHands(hands []puyo2.Hand) string {
	var s string
	for _, hand := range hands {
		s += fmt.Sprintf("%s%s%d%d", color2Str(hand.PuyoSet.Axis), color2Str(hand.PuyoSet.Child), hand.Position[0], hand.Position[1])
	}
	return s
}

func checkAllClear(n int, haipuyo []rune, table map[rune]int) bool {
	for _, v := range cntColor(haipuyo) {
		if v != 0 && v < 4 {
			return false
		}
	}

	puyoSets := []puyo2.PuyoSet{}
	for i := 0; i < len(haipuyo); i += 2 {
		puyoSet := puyo2.PuyoSet{
			Axis: puyo2.Color(table[haipuyo[i]]),
			Child: puyo2.Color(table[haipuyo[i+1]]),
		}
		puyoSets = append(puyoSets, puyoSet)
	}

	allClear := false
	bf := puyo2.NewBitField()
	hands := []puyo2.Hand{}
	bf.SearchWithPuyoSets(puyoSets, hands, func(sr *puyo2.SearchResult) bool {
		// 最後までいかないでぷよが消えた場合
		// if rr.Depth != len(puyoSets) && rr.RensaResult.Chains != 0 {
		// 	// // 全消しだった場合
		// 	// if rr.RensaResult.BitField.IsEmpty() {
		// 	// 	// fmt.Printf("%d %d %v\n", n, rr.Depth, puyoSets[:(rr.Depth+1)*2])
		// 	// 	return true
		// 	// }
		// 	return true
		// }
		// if rr.Depth == len(puyoSets)-1 && rr.RensaResult.BitField.IsEmpty() {
		// 	fmt.Printf("%d/%d/%v\n", n, rr.Depth, puyoSets)
		// 	// rr.BeforeSimulate.ShowDebug()
		// 	return false
		// }
		if sr.RensaResult.Chains >= 2 && sr.RensaResult.BitField.IsEmpty() {
			fmt.Printf("%d/%d/%v/%s/%d/%d\n", n, sr.Depth, puyoSets, hands2SimpleHands(sr.Hands), sr.RensaResult.Chains, sr.RensaResult.Score)
			// rr.BeforeSimulate.ShowDebug()
			allClear = true
			return true
		}
		return true
	}, 1)
	return allClear
}

func checkAllClearAll(n int, min int, max int, haipuyo []rune, table map[rune]int) {
	for c := min; c <= max; c += 2 {
		if checkAllClear(n, haipuyo[:c], table) {
			return
		}
	}
}


type HaipuyoInfo struct {
	Haipuyo []rune
	No      int
	Table   map[rune]int
}

func checkAllClearImpossibleG(chaipuyo <-chan HaipuyoInfo, wg *sync.WaitGroup) {
	loop: for {
		haipuyoInfo := <-chaipuyo
		if len(haipuyoInfo.Haipuyo) == 0 {
			break
		}
		for _, v := range cntColor(haipuyoInfo.Haipuyo) {
			if v != 0 && v < 4 {
				continue loop
			}
		}

		puyoSets := []puyo2.PuyoSet{}
		for i := 0; i < len(haipuyoInfo.Haipuyo); i += 2 {
			puyoSet := puyo2.PuyoSet{
				Axis: puyo2.Color(haipuyoInfo.Table[haipuyoInfo.Haipuyo[i]]),
				Child: puyo2.Color(haipuyoInfo.Table[haipuyoInfo.Haipuyo[i+1]]),
			}
			puyoSets = append(puyoSets, puyoSet)
		}
	
		allClear := false
		bf := puyo2.NewBitField()
		hands := []puyo2.Hand{}
		bf.SearchWithPuyoSets(puyoSets, hands, func(sr *puyo2.SearchResult) bool {
			if sr.Depth == len(puyoSets) && sr.RensaResult.Erased == len(haipuyoInfo.Haipuyo) && sr.RensaResult.BitField.IsEmpty() {
				allClear = true
				return false
			}
			return true
		}, 1)

		if allClear == false {
			fmt.Printf("%d %v\n", haipuyoInfo.No, puyoSets)
		}
	}
	wg.Done()
}

func main() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	var wg sync.WaitGroup

	var fp *os.File
	var err error
	chaipuyo := make(chan HaipuyoInfo)

	if len(os.Args) < 2 {
		fp = os.Stdin
	} else {
		fmt.Printf(">> read file: %s\n", os.Args[1])
		fp, err = os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer fp.Close()
	}

	fmt.Fprintf(os.Stderr, ">> goroutines %d\n", cpus)

	wg.Add(cpus)
	for i := 0; i < cpus; i++ {
		go checkAllClearImpossibleG(chaipuyo, &wg)
	}
	reader := bufio.NewReaderSize(fp, 257)
	n := 0
	for line := ""; err == nil; line, err = reader.ReadString('\n') {
		if line == "" {
			continue
		}
		n++
		if n%1000 == 0 {
			fmt.Fprintf(os.Stderr, "%d\n", n)
		}

		line = strings.TrimRight(line, "\n")
		haipuyo := []rune(line)
		cnt := cntColor(haipuyo)
		table := map[rune]int{'r': int(puyo2.Red), 'g': int(puyo2.Green), 'b': int(puyo2.Blue), 'y': int(puyo2.Yellow)}
		for c, v := range cnt {
			if c != 'p' && v == 0 {
				table['p'] = table[c]
				break
			}
		}

		// 4-16 ぷよまでの全消しチェック
		checkAllClearAll(n, 4, 16, haipuyo, table)

		// 8 ぷよのでぷよの個数は足りているが全消しが取れないパターン
		// checkAllClearImpossible(n, haipuyo[:8], table)
		// checkAllClearImpossible(n, haipuyo[:10], table)
		// checkAllClearImpossible(n, haipuyo[:12], table)
		// checkAllClearImpossible(n, haipuyo[:14], table)
		// checkAllClearImpossible(n, haipuyo, table)

		// checkAllClearImpossibleG 用の channel データ送信
		// var haipuyoInfo HaipuyoInfo
		// haipuyoInfo.Haipuyo = haipuyo[:16]
		// haipuyoInfo.No = n
		// haipuyoInfo.Table = table
		// chaipuyo <- haipuyoInfo
	}

	for i := 0; i < cpus; i++ {
		var haipuyoInfo HaipuyoInfo
		haipuyoInfo.Haipuyo = []rune{}
		chaipuyo <- haipuyoInfo
	}

	if err != io.EOF {
		panic(err)
	}
	wg.Wait()
}
