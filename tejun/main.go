package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/wata-gh/puyo2"
)

type Options struct {
	Haipuyo   string
	Threshold int
	BeamWidth int
	Param     string
	Shapes    []*puyo2.ShapeBitField
	Result    chan string
}

type SearchCondition struct {
	opt      Options
	bf       *puyo2.BitField
	fccs     []*FieldColorCandidate
	puyoSets []*puyo2.PuyoSet
	hands    []puyo2.Hand
	depth    int
}

func placeAndSetColor(bf *puyo2.BitField, place *puyo2.PuyoSetPlacement, fccs []*FieldColorCandidate, threshold int) []*FieldColorCandidate {
	bf.SetColor(place.PuyoSet.Axis, place.AxisX, place.AxisY)
	bf.SetColor(place.PuyoSet.Child, place.ChildX, place.ChildY)
	vfccs := []*FieldColorCandidate{}
	for _, fcc := range fccs {
		nfcc := fcc.Clone()
		// must place puyo from bottom.
		if place.AxisY <= place.ChildY {
			acc := fcc.GetColorCandidate(place.AxisX, place.AxisY)
			if acc.Contains(place.PuyoSet.Axis) == false {
				continue
			}
			nfcc.SetColorCandidate(place.AxisX, place.AxisY, []puyo2.Color{place.PuyoSet.Axis})

			ccc := nfcc.GetColorCandidate(place.ChildX, place.ChildY)
			if ccc.Contains(place.PuyoSet.Child) == false {
				continue
			}
			nfcc.SetColorCandidate(place.ChildX, place.ChildY, []puyo2.Color{place.PuyoSet.Child})
		} else { // child is lower than axis.
			ccc := fcc.GetColorCandidate(place.ChildX, place.ChildY)
			if ccc.Contains(place.PuyoSet.Child) == false {
				continue
			}
			nfcc.SetColorCandidate(place.ChildX, place.ChildY, []puyo2.Color{place.PuyoSet.Child})

			acc := nfcc.GetColorCandidate(place.AxisX, place.AxisY)
			if acc.Contains(place.PuyoSet.Axis) == false {
				continue
			}
			nfcc.SetColorCandidate(place.AxisX, place.AxisY, []puyo2.Color{place.PuyoSet.Axis})
		}
		if place.Chigiri {
			nfcc.ChigiriCount += 1
		}
		// remove candidate if it's over threshold
		if threshold != -1 && countOuterPlaced(bf, nfcc) >= threshold {
			continue
		}

		vfccs = append(vfccs, nfcc)
	}
	return vfccs
}

func setFirstTwoPuyoSets(bf *puyo2.BitField, puyoSets [2]*puyo2.PuyoSet, fccs []*FieldColorCandidate) [2][2]int {
	first := puyoSets[0]
	second := puyoSets[1]
	a := first.Axis
	firstPos := [2]int{}
	secondPos := [2]int{}
	if first.Axis == first.Child { // AA
		if second.Axis == a && second.Child == a { // AAAA
			firstPos = [2]int{2, 0}
			secondPos = [2]int{2, 0}
		} else if second.Axis == a { // AAAB
			firstPos = [2]int{0, 1}
			secondPos = [2]int{2, 2}
		} else if second.Child == a { // AABA
			firstPos = [2]int{0, 1}
			secondPos = [2]int{2, 0}
		} else if second.Axis == second.Child { // AABB
			firstPos = [2]int{0, 1}
			secondPos = [2]int{0, 1}
		} else { // AABC
			firstPos = [2]int{0, 1}
			secondPos = [2]int{2, 1}
		}
	} else { // AB
		b := first.Child
		if second.Axis == a && second.Child == a { // ABAA
			firstPos = [2]int{2, 2}
			secondPos = [2]int{0, 1}
		} else if second.Axis == a && second.Child == b { // ABAB
			firstPos = [2]int{0, 2}
			secondPos = [2]int{1, 2}
		} else if second.Axis == b && second.Child == a { // ABBA
			firstPos = [2]int{0, 2}
			secondPos = [2]int{1, 0}
		} else if second.Axis == b && second.Child == b { // ABBB
			firstPos = [2]int{2, 0}
			secondPos = [2]int{0, 1}
		} else if second.Axis == b && second.Child != a && second.Child != b { // ABBC
			firstPos = [2]int{0, 2}
			secondPos = [2]int{1, 1}
		} else if second.Axis == a && second.Child != a && second.Child != b { // ABAC
			firstPos = [2]int{1, 1}
			secondPos = [2]int{0, 0}
		} else if second.Axis != a && second.Axis != b && second.Child == a { // ABCA
			firstPos = [2]int{1, 1}
			secondPos = [2]int{0, 2}
		} else if second.Axis != a && second.Axis != b && second.Child != a && second.Child != b { // ABCC
			firstPos = [2]int{2, 1}
			secondPos = [2]int{0, 1}
		} else if second.Axis != a && second.Axis != b && second.Child == b { // ABCB
			firstPos = [2]int{2, 3}
			secondPos = [2]int{0, 2}
		} else {
			panic(fmt.Sprintf("first: %+v second: %+v\n", first, second))
		}
	}
	return [2][2]int{firstPos, secondPos}
}

func searchPlacement(fccs []*FieldColorCandidate, bf *puyo2.BitField, puyoSet *puyo2.PuyoSet) []*puyo2.PuyoSetPlacement {
	placements := []*puyo2.PuyoSetPlacement{}
	for _, pos := range puyo2.SetupPositions {
		placement := bf.SearchPlacementForPos(puyoSet, pos)
		if placement == nil {
			continue
		}
		for _, fcc := range fccs {
			sbf := fcc.ShapeBitField
			an := sbf.ShapeNum(placement.AxisX, placement.AxisY)
			acc := fcc.GetColorCandidate(placement.AxisX, placement.AxisY)
			if aok := acc.Contains(puyoSet.Axis); aok {
				acc = NewColorCandidate([]puyo2.Color{puyoSet.Axis})
			} else {
				continue
			}

			cn := sbf.ShapeNum(placement.ChildX, placement.ChildY)
			ccc := fcc.GetColorCandidate(placement.ChildX, placement.ChildY)
			if an != -1 && cn != -1 && cn == an { // both In-Shape and same shape
				ccc = acc
			}
			if cok := ccc.Contains(puyoSet.Child); cok {
				placements = append(placements, placement)
				break
			}
		}
	}
	return placements
}

func countPlaced(bf *puyo2.BitField, fcc *FieldColorCandidate) int {
	fb := bf.Bits(puyo2.Empty)
	fb.M[0] = ^fb.M[0]
	fb.M[1] = ^fb.M[1]
	return fcc.ShapeBitField.OverallShape().And(fb).PopCount()
}

func countOuterPlaced(bf *puyo2.BitField, fcc *FieldColorCandidate) int {
	empty := bf.Bits(puyo2.Empty)
	empty.M[0] = ^empty.M[0]
	empty.M[1] = ^empty.M[1]
	empty = empty.MaskField13()

	overall := fcc.ShapeBitField.OverallShape()
	overall.M[0] = ^overall.M[0]
	overall.M[1] = ^overall.M[1]
	overall = overall.MaskField13()

	return empty.And(overall).PopCount()
}

func beamSearch(opt Options, conds []*SearchCondition) {
	if len(conds) == 0 {
		return
	}

	sort.Slice(conds, func(i, j int) bool {
		return len(conds[i].fccs) > len(conds[j].fccs)
	})
	width := opt.BeamWidth
	if width == -1 || len(conds) < width {
		width = len(conds)
	} else if len(conds) > width {
		cutoff := len(conds[width-1].fccs)
		for i := width; i < len(conds); i++ {
			len := len(conds[i].fccs)
			if cutoff > len {
				width = i
				break
			}
		}
	}
	fmt.Fprintf(os.Stderr, "width: %d\n", width)
	nextCondsCandidates := make([][]*SearchCondition, width)
	wg := sync.WaitGroup{}
	for i, cond := range conds[:width] {
		nextCondsCandidates[i] = []*SearchCondition{}
		wg.Add(1)
		go search(cond, &nextCondsCandidates[i], &wg)
	}
	wg.Wait()
	nextConds := []*SearchCondition{}
	for _, candidate := range nextCondsCandidates {
		nextConds = append(nextConds, candidate...)
	}
	beamSearch(opt, nextConds)
}

func search(cond *SearchCondition, nextConds *[]*SearchCondition, wg *sync.WaitGroup) {
	if len(cond.puyoSets) == 0 {
		var b strings.Builder
		sort.Slice(cond.fccs, func(i, j int) bool {
			return countPlaced(cond.bf, cond.fccs[i]) > countPlaced(cond.bf, cond.fccs[j])
		})
		for _, fcc := range cond.fccs {
			fmt.Fprintf(&b, " %s:%d(%d)", fcc.ShapeBitField.FieldString(), countPlaced(cond.bf, fcc), fcc.ChigiriCount)
		}
		cond.opt.Result <- fmt.Sprintf("%s %s%s", cond.bf.MattulwanEditorParam(), puyo2.ToSimpleHands(cond.hands), b.String())
		wg.Done()
		return
	}

	placements := searchPlacement(cond.fccs, cond.bf, cond.puyoSets[0])
	for _, place := range placements {
		bfc := cond.bf.Clone()
		nfccs := placeAndSetColor(bfc, place, cond.fccs, cond.opt.Threshold)
		if len(nfccs) == 0 {
			continue
		}

		nhands := make([]puyo2.Hand, len(cond.hands))
		copy(nhands, cond.hands)
		nhands = append(nhands, puyo2.Hand{
			PuyoSet:  *place.PuyoSet,
			Position: place.Pos,
		})

		newCond := SearchCondition{
			opt:      cond.opt,
			bf:       bfc,
			fccs:     nfccs,
			puyoSets: cond.puyoSets[1:],
			hands:    nhands,
			depth:    cond.depth + 1,
		}
		*nextConds = append(*nextConds, &newCond)
	}
	wg.Done()
}

func createTableAndColors(puyoSets []*puyo2.PuyoSet) (map[puyo2.Color]puyo2.Color, []puyo2.Color) {
	colors := []puyo2.Color{}
	table := map[puyo2.Color]puyo2.Color{
		puyo2.Red:    puyo2.Empty,
		puyo2.Blue:   puyo2.Empty,
		puyo2.Green:  puyo2.Empty,
		puyo2.Yellow: puyo2.Empty,
		puyo2.Purple: puyo2.Empty,
	}

	for _, puyoSet := range puyoSets {
		table[puyoSet.Axis] = puyoSet.Axis
		table[puyoSet.Child] = puyoSet.Child
	}
	// contains purple
	if table[puyo2.Purple] == puyo2.Purple {
		for _, c := range []puyo2.Color{puyo2.Red, puyo2.Blue, puyo2.Yellow, puyo2.Green} {
			if table[c] == puyo2.Empty {
				table[c] = puyo2.Purple
				table[puyo2.Purple] = c
			} else {
				colors = append(colors, c)
			}
		}
		colors = append(colors, puyo2.Purple)
	} else {
		for k, v := range table {
			if v != puyo2.Empty {
				colors = append(colors, k)
			}
		}
	}
	return table, colors
}

func run(opt Options) {
	// f, err := os.Create("cpu.pprof")
	// if err != nil {
	// 	return
	// }
	// if err := pprof.StartCPUProfile(f); err != nil {
	// 	return
	// }
	// defer pprof.StopCPUProfile()

	puyoSets := puyo2.Haipuyo2PuyoSets(opt.Haipuyo)
	table, colors := createTableAndColors(puyoSets)
	bf := puyo2.NewBitFieldWithTableAndColors(table, colors)
	fccs := make([]*FieldColorCandidate, len(opt.Shapes))
	for i, sbf := range opt.Shapes {
		fccs[i] = NewFieldColorCandidate(colors, sbf)
	}

	poss := setFirstTwoPuyoSets(bf, [2]*puyo2.PuyoSet{puyoSets[0], puyoSets[1]}, fccs)
	hands := []puyo2.Hand{}
	for i, pos := range poss {
		placement := bf.SearchPlacementForPos(puyoSets[i], pos)
		fccs = placeAndSetColor(bf, placement, fccs, opt.Threshold)
		if len(fccs) == 0 {
			return
		}
		hands = append(hands, puyo2.Hand{
			PuyoSet:  *placement.PuyoSet,
			Position: placement.Pos,
		})
	}

	if bf.Clone().Simulate().BitField.IsEmpty() { // skip All-Clear hands
		fmt.Fprintln(os.Stderr, "[end] all clear.")
		return
	}

	newCond := SearchCondition{
		opt:      opt,
		bf:       bf,
		fccs:     fccs,
		puyoSets: puyoSets[2:],
		hands:    hands,
		depth:    0,
	}
	beamSearch(opt, []*SearchCondition{&newCond})
}

func handleResult(result chan string, wg *sync.WaitGroup) {
	for {
		r := <-result
		if r == "" {
			break
		}
		fmt.Println(r)
	}
	wg.Done()
}

func main() {
	wg := sync.WaitGroup{}
	now := time.Now()
	opt := Options{}
	flag.StringVar(&opt.Param, "param", "a78", "field parameter")
	flag.StringVar(&opt.Haipuyo, "haipuyo", "", "haipuyo")
	flag.IntVar(&opt.Threshold, "threshold", -1, "threshold of out of placements")
	flag.IntVar(&opt.BeamWidth, "width", -1, "beam search width")
	flag.Parse()

	fmt.Fprintf(os.Stderr, "%+v\n", opt)
	opt.Shapes = []*puyo2.ShapeBitField{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		shapeStr := strings.TrimRight(scanner.Text(), "\n")
		shape := puyo2.NewShapeBitFieldWithFieldString(shapeStr)
		opt.Shapes = append(opt.Shapes, shape)
	}
	fmt.Fprintf(os.Stderr, "Shapes.Len(): %d\n", len(opt.Shapes))

	opt.Result = make(chan string)
	wg.Add(1)
	go handleResult(opt.Result, &wg)
	run(opt)
	opt.Result <- ""
	wg.Wait()
	fmt.Fprintf(os.Stderr, "elapsed: %v ms\n", time.Since(now).Milliseconds())
}
