package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/dustin/go-humanize"
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

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP)
	go func() {
		for range c {
			t := time.Now()
			fmt.Fprintf(os.Stderr, "[%s] %+v\n", t.Format("2006-01-02 15:04:05"), pattern)
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			// ヒープ上に割り当てられたオブジェクト累積メモリ量
			fmt.Fprintf(os.Stderr, "MAlloc : %v\n", humanize.Bytes(m.Mallocs))
			// ヒープ上から開放されたオブジェクト数
			fmt.Fprintf(os.Stderr, "Frees : %v\n", m.Frees)

			// ヒープ上に割り当てられたオブジェクトメモリ量
			fmt.Fprintf(os.Stderr, "Alloc : %v\n", humanize.Bytes(m.Alloc))
			fmt.Fprintf(os.Stderr, "HeapAlloc : %v\n", humanize.Bytes(m.HeapAlloc))
			// ヒープ上に割り当てられたオブジェクトメモリ量。ただし開放されたオブジェクト分も含む
			fmt.Fprintf(os.Stderr, "TotalAlloc : %v\n", humanize.Bytes(m.TotalAlloc))

			// OSから割り当てられたプロセスの総メモリ量
			// ヒープ + スタック + その他
			fmt.Fprintf(os.Stderr, "Sys : %v\n", humanize.Bytes(m.Sys))

			// ポインタのルックアップ数
			fmt.Fprintf(os.Stderr, "Lookups : %v\n", m.Lookups)

			// 到達可能、あるいはGCによって解放されていないヒープオブジェクトメモリ量
			fmt.Fprintf(os.Stderr, "HeapAlloc : %v\n", humanize.Bytes(m.HeapAlloc))
			// 未使用ヒープ領域メモリ量
			fmt.Fprintf(os.Stderr, "HeapIdle : %v\n", humanize.Bytes(m.HeapIdle))
			// 使用中ヒープ領域メモリ量
			fmt.Fprintf(os.Stderr, "HeapInuse : %v\n", humanize.Bytes(m.HeapInuse))
			// OSに返却される物理メモリ量
			fmt.Fprintf(os.Stderr, "HeapReleased : %v\n", humanize.Bytes(m.HeapReleased))
			// ヒープに割り当てられたオブジェクト量
			fmt.Fprintf(os.Stderr, "HeapObjects : %v\n", m.HeapObjects)

			// 使用中スタック領域メモリ量
			fmt.Fprintf(os.Stderr, "StackInuse : %v\n", humanize.Bytes(m.StackInuse))
			// OSから割り当てられたスタック領域メモリ量
			fmt.Fprintf(os.Stderr, "StackSys : %v\n", humanize.Bytes(m.StackSys))

			// 割り当てられたmspan構造体バイト数
			fmt.Fprintf(os.Stderr, "MSpanInuse : %v\n", humanize.Bytes(m.MSpanInuse))
			// OSから取得したmspan構造体バイト数
			fmt.Fprintf(os.Stderr, "MSpanSys : %v\n", humanize.Bytes(m.MSpanSys))

			// 割り当てられたmcache構造体バイト数
			fmt.Fprintf(os.Stderr, "MCacheInuse : %v\n", humanize.Bytes(m.MCacheInuse))
			// OSから取得したmcache構造体バイト数
			fmt.Fprintf(os.Stderr, "MCacheSys : %v\n", humanize.Bytes(m.MCacheSys))
			grcn := runtime.NumGoroutine()
			fmt.Fprintf(os.Stderr, "Goroutine : %v\n", grcn)

			runtime.ReadMemStats(&m)

			sr := make([]runtime.StackRecord, grcn)
			n, ok := runtime.GoroutineProfile(sr)
			fmt.Printf("n-GoroutineProfile : %d\n", n)
			if ok {
				for i, p := range sr {
					fmt.Fprintf(os.Stderr, "GoroutineProfile-%d: %v\n", i, p.Stack0)
				}
			}
		}
	}()

	grc := 8
	wg.Add(grc)
	for i := 0; i < grc; i++ {
		go pattern.Check(field, opt, &wg)
	}
	Gen(&pattern, field, grc)
	wg.Wait()
	fmt.Fprintf(os.Stderr, "%+v\n", opt)
	pattern.ShowResult()
	fmt.Fprintf(os.Stderr, "elapsed: %vms\n", time.Since(now).Milliseconds())
}
