package main

import (
	"fmt"
	"testing"

	"github.com/wata-gh/puyo2"
)

func TestRun(t *testing.T) {
	// haipuyo := "rrprpypyrbbbpbbprrrybbyppbyyrppybybrbbbppppbyypybppyypbyrbyyyppppbpppyryyrpyybpbryrbbrpybrrbrrbbpypyrryrrybrbpbbybrrpppyrprrryrrbybrbbrbrybprpyppybyrpprpbbyybbyybrprbybryrrbrybyppbbbpyybprpyyrryppyrrbppybyyypprpryrpbpbpbyrpyprybrybrrbppyrbyypryrbbprrbprprb"
	// haipuyo := "rrprpypy" // rbbbpbbprrrybbyp
	haipuyo := "ppbbygyyppppybbygypyggpb"
	opt := Options{
		Haipuyo:   haipuyo,
		Threshold: 5,
	}
	opt.Shapes = append(opt.Shapes, puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1a1aaa212654223665444551333"))
	run(opt)
}

func TestCountOuterPlaced(t *testing.T) {
	sbf := puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa5123445112334223554")
	bf := puyo2.NewBitFieldWithMattulwan("a54ba5ba5ba5ba5")
	fcc := NewFieldColorCandidate([]puyo2.Color{}, sbf)
	cnt := countOuterPlaced(bf, fcc)
	if cnt != 1 {
		sbf.ShowDebug()
		bf.ShowDebug()
		t.Fatalf("countPlaced must be 1. but %d.\n", cnt)
	}

	bf = puyo2.NewBitFieldWithMattulwan("a60ba5ba5ba5")
	fcc = NewFieldColorCandidate([]puyo2.Color{}, sbf)
	cnt = countOuterPlaced(bf, fcc)
	if cnt != 0 {
		sbf.ShowDebug()
		bf.ShowDebug()
		t.Fatalf("countPlaced must be 0. but %d.\n", cnt)
	}

	bf = puyo2.NewBitFieldWithMattulwan("a53b2a4b2a5ba5ba5")
	fcc = NewFieldColorCandidate([]puyo2.Color{}, sbf)
	cnt = countOuterPlaced(bf, fcc)
	if cnt != 2 {
		sbf.ShowDebug()
		bf.ShowDebug()
		t.Fatalf("countPlaced must be 2. but %d.\n", cnt)
	}
}

func TestRemoveColorCandidate(t *testing.T) {
	sbf := puyo2.NewShapeBitField()
	fcc := NewFieldColorCandidate([]puyo2.Color{puyo2.Red, puyo2.Blue, puyo2.Yellow, puyo2.Green}, sbf)
	fcc.RemoveColorCandidate(0, 1, []puyo2.Color{puyo2.Red})
	cc := fcc.GetColorCandidate(0, 1)
	if cc.Contains(puyo2.Red) || len(cc.colors) != 3 {
		fmt.Printf("%t %t\n", cc.Contains(puyo2.Red), len(cc.colors) != 3)
		t.Fatal("Remove colors candidate.")
	}
}

func TestSetAdjacent(t *testing.T) {
	colors := []puyo2.Color{puyo2.Red, puyo2.Blue, puyo2.Yellow, puyo2.Green}
	sbf := puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1a1aaa212654223665444551333")
	fcc := NewFieldColorCandidate(colors, sbf)
	fcc.setShapeAdjacent()
	if len(fcc.ShapeAdjacent[0]) != 4 {
		t.Fatalf("0 must be len == 4 %v\n", fcc.ShapeAdjacent[0])
	}
	if len(fcc.ShapeAdjacent[1]) != 3 {
		t.Fatalf("1 must be len == 3 %v\n", fcc.ShapeAdjacent[1])
	}
	if len(fcc.ShapeAdjacent[2]) != 3 {
		t.Fatalf("2 must be len == 3 %v\n", fcc.ShapeAdjacent[2])
	}
	if len(fcc.ShapeAdjacent[3]) != 4 {
		t.Fatalf("3 must be len == 4 %v\n", fcc.ShapeAdjacent[3])
	}
	if len(fcc.ShapeAdjacent[4]) != 3 {
		t.Fatalf("4 must be len == 3 %v\n", fcc.ShapeAdjacent[4])
	}
	if len(fcc.ShapeAdjacent[5]) != 1 {
		t.Fatalf("5 must be len == 1 %v\n", fcc.ShapeAdjacent[5])
	}
}

func TestSetFirstTwoPuyoSets(t *testing.T) {
	sbf := puyo2.NewShapeBitFieldWithFieldString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa6a6aaa565123554112333226444")
	bf := puyo2.NewBitField()
	fcc := NewFieldColorCandidate([]puyo2.Color{puyo2.Red, puyo2.Blue, puyo2.Yellow, puyo2.Green}, sbf)
	fcc.BitField = bf
	fccs := []*FieldColorCandidate{
		fcc,
	}
	poss := setFirstTwoPuyoSets(bf, [2]*puyo2.PuyoSet{{Axis: puyo2.Red, Child: puyo2.Red}, {Axis: puyo2.Red, Child: puyo2.Red}}, fccs)
	if poss[0] != [2]int{2, 0} || poss[1] != [2]int{2, 0} {
		bf.ShowDebug()
		t.Fatal("aaaa failed")
	}
	bf = puyo2.NewBitField()
	fccs[0] = NewFieldColorCandidate([]puyo2.Color{puyo2.Red, puyo2.Blue, puyo2.Yellow, puyo2.Green}, sbf)
	poss = setFirstTwoPuyoSets(bf, [2]*puyo2.PuyoSet{{Axis: puyo2.Red, Child: puyo2.Red}, {Axis: puyo2.Red, Child: puyo2.Blue}}, fccs)
	if poss[0] != [2]int{0, 1} || poss[1] != [2]int{2, 2} {
		bf.ShowDebug()
		t.Fatal("aaab failed")
	}
	bf = puyo2.NewBitField()
	poss = setFirstTwoPuyoSets(bf, [2]*puyo2.PuyoSet{{Axis: puyo2.Red, Child: puyo2.Red}, {Axis: puyo2.Blue, Child: puyo2.Red}}, fccs)
	if poss[0] != [2]int{0, 1} || poss[1] != [2]int{2, 0} {
		bf.ShowDebug()
		fmt.Println(bf.MattulwanEditorParam())
		t.Fatal("aaba failed")
	}
	bf = puyo2.NewBitField()
	poss = setFirstTwoPuyoSets(bf, [2]*puyo2.PuyoSet{{Axis: puyo2.Red, Child: puyo2.Red}, {Axis: puyo2.Blue, Child: puyo2.Blue}}, fccs)
	if poss[0] != [2]int{0, 1} || poss[1] != [2]int{0, 1} {
		bf.ShowDebug()
		fmt.Println(bf.MattulwanEditorParam())
		t.Fatal("aabb failed")
	}
	bf = puyo2.NewBitField()
	poss = setFirstTwoPuyoSets(bf, [2]*puyo2.PuyoSet{{Axis: puyo2.Red, Child: puyo2.Red}, {Axis: puyo2.Blue, Child: puyo2.Green}}, fccs)
	if poss[0] != [2]int{0, 1} || poss[1] != [2]int{2, 1} {
		bf.ShowDebug()
		fmt.Println(bf.MattulwanEditorParam())
		t.Fatal("aabc failed")
	}
	bf = puyo2.NewBitField()
	poss = setFirstTwoPuyoSets(bf, [2]*puyo2.PuyoSet{{Axis: puyo2.Red, Child: puyo2.Blue}, {Axis: puyo2.Red, Child: puyo2.Red}}, fccs)
	if poss[0] != [2]int{2, 2} || poss[1] != [2]int{0, 1} {
		bf.ShowDebug()
		fmt.Println(bf.MattulwanEditorParam())
		t.Fatal("abaa failed")
	}
	bf = puyo2.NewBitField()
	poss = setFirstTwoPuyoSets(bf, [2]*puyo2.PuyoSet{{Axis: puyo2.Red, Child: puyo2.Blue}, {Axis: puyo2.Red, Child: puyo2.Blue}}, fccs)
	if poss[0] != [2]int{0, 2} || poss[1] != [2]int{1, 2} {
		bf.ShowDebug()
		fmt.Println(bf.MattulwanEditorParam())
		t.Fatal("abab failed")
	}
	bf = puyo2.NewBitField()
	poss = setFirstTwoPuyoSets(bf, [2]*puyo2.PuyoSet{{Axis: puyo2.Red, Child: puyo2.Blue}, {Axis: puyo2.Blue, Child: puyo2.Blue}}, fccs)
	if poss[0] != [2]int{2, 0} || poss[1] != [2]int{0, 1} {
		bf.ShowDebug()
		fmt.Println(bf.MattulwanEditorParam())
		t.Fatal("abbb failed")
	}
	bf = puyo2.NewBitField()
	poss = setFirstTwoPuyoSets(bf, [2]*puyo2.PuyoSet{{Axis: puyo2.Red, Child: puyo2.Blue}, {Axis: puyo2.Blue, Child: puyo2.Green}}, fccs)
	if poss[0] != [2]int{0, 2} || poss[1] != [2]int{1, 1} {
		bf.ShowDebug()
		fmt.Println(bf.MattulwanEditorParam())
		t.Fatal("abbc failed")
	}
	bf = puyo2.NewBitField()
	poss = setFirstTwoPuyoSets(bf, [2]*puyo2.PuyoSet{{Axis: puyo2.Red, Child: puyo2.Blue}, {Axis: puyo2.Blue, Child: puyo2.Red}}, fccs)
	if poss[0] != [2]int{0, 2} || poss[1] != [2]int{1, 0} {
		bf.ShowDebug()
		fmt.Println(bf.MattulwanEditorParam())
		t.Fatal("abba failed")
	}
	bf = puyo2.NewBitField()
	poss = setFirstTwoPuyoSets(bf, [2]*puyo2.PuyoSet{{Axis: puyo2.Red, Child: puyo2.Blue}, {Axis: puyo2.Red, Child: puyo2.Green}}, fccs)
	if poss[0] != [2]int{1, 1} || poss[1] != [2]int{0, 0} {
		bf.ShowDebug()
		fmt.Println(bf.MattulwanEditorParam())
		t.Fatal("abac failed")
	}
	bf = puyo2.NewBitField()
	poss = setFirstTwoPuyoSets(bf, [2]*puyo2.PuyoSet{{Axis: puyo2.Red, Child: puyo2.Blue}, {Axis: puyo2.Green, Child: puyo2.Red}}, fccs)
	if poss[0] != [2]int{1, 1} || poss[1] != [2]int{0, 2} {
		bf.ShowDebug()
		fmt.Println(bf.MattulwanEditorParam())
		t.Fatal("abac failed")
	}
	bf = puyo2.NewBitField()
	poss = setFirstTwoPuyoSets(bf, [2]*puyo2.PuyoSet{{Axis: puyo2.Red, Child: puyo2.Blue}, {Axis: puyo2.Green, Child: puyo2.Green}}, fccs)
	if poss[0] != [2]int{2, 1} || poss[1] != [2]int{0, 1} {
		bf.ShowDebug()
		fmt.Println(bf.MattulwanEditorParam())
		t.Fatal("abcc failed")
	}
	bf = puyo2.NewBitField()
	poss = setFirstTwoPuyoSets(bf, [2]*puyo2.PuyoSet{{Axis: puyo2.Red, Child: puyo2.Blue}, {Axis: puyo2.Green, Child: puyo2.Blue}}, fccs)
	if poss[0] != [2]int{2, 3} || poss[1] != [2]int{0, 2} {
		bf.ShowDebug()
		fmt.Println(bf.MattulwanEditorParam())
		t.Fatal("abcb failed")
	}
	bf = puyo2.NewBitField()
	poss = setFirstTwoPuyoSets(bf, [2]*puyo2.PuyoSet{{Axis: puyo2.Red, Child: puyo2.Blue}, {Axis: puyo2.Green, Child: puyo2.Red}}, fccs)
	if poss[0] != [2]int{1, 1} || poss[1] != [2]int{0, 2} {
		bf.ShowDebug()
		fmt.Println(bf.MattulwanEditorParam())
		t.Fatal("abca failed")
	}
}
