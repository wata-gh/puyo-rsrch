package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

type options struct {
	Dir       string
	Field     string
	Chains    int
	ShapeOnly bool
}

// GTR 18
// func validEmpty(list []int) bool {
// 	for _, v := range list {
// 		if v <= 3 {
// 			continue
// 		}
// 		if v <= 6 {
// 			if contains(list, v-4) {
// 				continue
// 			}
// 			return false
// 		}
// 		if v <= 10 {
// 			if contains(list, v-4) && (v == 7 || contains(list, v-8)) {
// 				continue
// 			}
// 			return false
// 		}
// 		if v <= 13 {
// 			if contains(list, v-3) && contains(list, v-7) && contains(list, v-11) {
// 				continue
// 			}
// 			return false
// 		}
// 		if v == 14 {
// 			return false
// 		}
// 		if v <= 17 {
// 			if contains(list, v-4) && contains(list, v-7) && contains(list, v-11) && contains(list, v-15) {
// 				continue
// 			}
// 			return false
// 		}
// 		panic(list)
// 	}
// 	return true
// }

var opt options

func init() {
	flag.StringVar(&opt.Dir, "dir", "results", "output directory path")
	flag.StringVar(&opt.Field, "field", "multi27", "field pattern")
	flag.IntVar(&opt.Chains, "chains", 0, "chain count")
	flag.BoolVar(&opt.ShapeOnly, "shape-only", false, "use shape only")
}

func main() {
	now := time.Now()
	flag.Parse()

	os.Mkdir(opt.Dir, 0755)
	var wg sync.WaitGroup
	field := make(chan []int, 1000)
	patterns := map[string]Pattern{
		"gtr15": &Gtr15Pattern{},
		"multi27": &Multi27Pattern{
			ChainCount: opt.Chains,
		},
	}
	pattern, ok := patterns[opt.Field]
	if !ok {
		fmt.Fprintln(os.Stderr, "no such field. "+opt.Field)
		return
	}
	pattern.Init()
	grc := 8
	wg.Add(grc)
	for i := 0; i < grc; i++ {
		go pattern.Check(field, opt, &wg)
	}
	Gen(&pattern, field, grc)
	wg.Wait()
	fmt.Fprintf(os.Stderr, "%+v\n", opt)
	pattern.ShowResult()
	fmt.Fprintf(os.Stderr, "経過: %vms\n", time.Since(now).Milliseconds())
}
