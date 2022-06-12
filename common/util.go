package puyorsrch

import (
	"fmt"

	"github.com/wata-gh/puyo2"
)

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
	case puyo2.Purple:
		return "p"
	}
	panic("color can be r,b,g,y,p")
}

func Hands2SimpleHands(hands []puyo2.Hand) string {
	var s string
	for _, hand := range hands {
		s += fmt.Sprintf("%s%s%d%d", color2Str(hand.PuyoSet.Axis), color2Str(hand.PuyoSet.Child), hand.Position[0], hand.Position[1])
	}
	return s
}
