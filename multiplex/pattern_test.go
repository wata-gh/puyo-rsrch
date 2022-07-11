package main

import (
	"testing"
)

func TestCheckShape(t *testing.T) {
	var opt options
	opt.Chains = 1
	opt.Dir = "test"
	puyos := []int{119537664}
	var pattern Pattern
	pattern = &Multi27Pattern{
		ChainCount: opt.Chains,
	}
	checkShape(&pattern, puyos, opt)
}
