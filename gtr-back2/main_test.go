package main

import (
	"testing"

	"github.com/wata-gh/puyo2"
)

func TestMain(t *testing.T) {
	puyo2.NewFieldBitsWithM([2]uint64{5066575350595584, 1048594}).ShowDebug()
}
