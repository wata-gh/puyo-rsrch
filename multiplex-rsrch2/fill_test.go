package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/wata-gh/puyo2"
)

// func TestFillUp(t *testing.T) {
// 	results := fillUp([]int{3, 3, 2, 0, 0, 0}, 2)
// 	fmt.Println(results)
// 	removeDuplication(results)
// }

// func TestFill(t *testing.T) {
// 	var pattern Result
// 	pattern.results = [][]int{{1, 0, 2, 1, 1}, {3, 0, 1, 2, 1}}
// 	Fill(pattern, 2)
// }
func TestPlace(t *testing.T) {
	backup := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	sbf := puyo2.NewShapeBitField()
	clusters := [][]int{
		{1, 3, 2, 1, 1},
		{1, 2, 2, 1, 1},
		{3, 3, 1, 2, 1},
		{0, 2, 1, 1, 1, 1},
	}
	// clusters := [][]int{
	// 	{2, 3, 1, 1, 2},
	// 	{2, 3, 1, 1, 2},
	// 	{3, 3, 1, 2, 1},
	// 	{1, 2, 2, 1, 1},
	// }
	chainc := 4
	place(sbf, clusters, chainc, nil)
	w.Close()
	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(r)
	if err != nil {
		t.Fatalf("fail read buf: %v", err)
	}
	os.Stdout = backup
	result := buffer.String()
	expect := `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa11aaa212a654231665444552333
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa11aaa212a654333665444552231
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa11aaa212a654331665444552233
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa11aaa212a654233665444552331
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa12aaa221a654131665444552333
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa12aaa423a652444665333552111
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa12aaa221a654333665444552131
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa12aaa221a654331665444552133
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa12aaa433a652444665233552111
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa12aaa221a654133665444552331
`
	if expect != result {
		t.Fatalf("unexpected results.\nexpected => %s\nresult =>%s", expect, result)
	}
}
