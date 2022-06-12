package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	common "github.com/wata-gh/puyo-rsrch/common"
	"github.com/wata-gh/puyo2"
)

func allClearShape(sr *puyo2.SearchResult, cnt *int, cache map[string]string, line *string) bool {
	if sr.RensaResult.Chains == 2 {
		*cnt++
		sr.BeforeSimulate.TrimLeft()
		shapes := sr.BeforeSimulate.ToChainShapesUInt64Array()
		s := ""
		for _, shape := range shapes {
			s += fmt.Sprintf("_%d-%d_", shape[0], shape[1])
		}
		_, ok := cache[s]
		if ok == false {
			cache[s] = s
			clone := sr.BeforeSimulate.Clone()
			clone.FlipHorizontal()
			cshapes := clone.ToChainShapesUInt64Array()
			cs := ""
			for _, shape := range cshapes {
				cs += fmt.Sprintf("_%d-%d_", shape[0], shape[1])
			}
			_, ok = cache[cs]
			if ok == false {
				cache[cs] = s
				sr.BeforeSimulate.ExportImage(fmt.Sprintf("%s/%s.png", *line, s))
			}
		}
		// sr.BeforeSimulate.ShowDebug()
	}
	return true
}

func allClearShapeCount(chains int, depth int, sr *puyo2.SearchResult, cnt *int, cache map[string]string, line *string) bool {
	if sr.RensaResult.BitField.IsEmpty() && sr.RensaResult.Chains == chains && sr.Depth == depth {
		*cnt++
		tlbf := sr.BeforeSimulate.Clone().TrimLeft()
		if tlbf.Equals(sr.BeforeSimulate) == false { // 単に横移動したパターンは除外
			return true
		}
		s := ""
		for _, shape := range sr.BeforeSimulate.ToChainShapesUInt64Array() {
			s += fmt.Sprintf("_%d-%d_", shape[0], shape[1])
		}
		_, ok := cache[s]
		if ok == false {
			fs := ""
			for _, shape := range sr.BeforeSimulate.Clone().FlipHorizontal().ToChainShapesUInt64Array() {
				fs += fmt.Sprintf("_%d-%d_", shape[0], shape[1])
			}
			_, ok := cache[fs]
			if ok { // 左右反転パターンも除外
				return true
			}
		}
		fmt.Printf("%s,%s,%v,%s\n", *line, s, common.Hands2SimpleHands(sr.Hands), sr.BeforeSimulate.MattulwanEditorParam())
	} else {
		// fmt.Println("allClearShapeCount: ", sr.RensaResult.Chains, sr.Depth)
		// sr.BeforeSimulate.ShowDebug()
	}
	return true
}

func main() {
	table := map[string]int{"r": int(puyo2.Red), "g": int(puyo2.Green), "b": int(puyo2.Blue), "y": int(puyo2.Yellow), "p": int(puyo2.Purple)}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cache := map[string]string{}
		line := strings.TrimRight(scanner.Text(), "\n")
		bfTable := map[puyo2.Color]puyo2.Color{
			puyo2.Red:    puyo2.Purple,
			puyo2.Green:  puyo2.Purple,
			puyo2.Blue:   puyo2.Purple,
			puyo2.Yellow: puyo2.Purple,
			puyo2.Empty:  puyo2.Empty,
		}
		for _, puyo := range strings.Split(line, "") {
			fmt.Println(puyo)
			switch puyo {
			case "r":
				bfTable[puyo2.Color(table[puyo])] = puyo2.Color(table[puyo])
			case "g":
				bfTable[puyo2.Color(table[puyo])] = puyo2.Color(table[puyo])
			case "b":
				bfTable[puyo2.Color(table[puyo])] = puyo2.Color(table[puyo])
			case "y":
				bfTable[puyo2.Color(table[puyo])] = puyo2.Color(table[puyo])
			}
			fmt.Println(bfTable)
		}
		for _, puyo := range []puyo2.Color{puyo2.Red, puyo2.Green, puyo2.Blue, puyo2.Yellow} {
			if bfTable[puyo] == puyo2.Purple {
				bfTable[puyo2.Purple] = puyo
			}
		}
		fmt.Println(bfTable)
		bf := puyo2.NewBitFieldWithTable(bfTable)
		fmt.Fprintf(os.Stderr, "%s\n", line)
		haipuyo := strings.Split(line, "")
		puyoSets := []puyo2.PuyoSet{}
		for i := 0; i < len(haipuyo); i += 2 {
			puyoSet := puyo2.PuyoSet{
				Axis:  puyo2.Color(table[haipuyo[i]]),
				Child: puyo2.Color(table[haipuyo[i+1]]),
			}
			puyoSets = append(puyoSets, puyoSet)
		}
		fmt.Println(puyoSets)
		hands := []puyo2.Hand{}
		cnt := 0
		os.Mkdir(line, 0755)
		bf.SearchWithPuyoSets(puyoSets, hands, func(sr *puyo2.SearchResult) bool {
			return allClearShapeCount(len(puyoSets)/2, len(puyoSets), sr, &cnt, cache, &line)
			// return allClearShape(sr, &cnt, cache, &line)
		}, 1)
		fmt.Fprintf(os.Stderr, "%s %d\n", line, cnt)
	}

}
