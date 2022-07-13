package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/wata-gh/puyo2"
)

type Config struct {
	PuyoConfig PuyoConfig
}

type PuyoConfig struct {
	SearchLocations [][2]int
	Field           string
	Mappings        map[string]int
}

func check2(field <-chan [15]int) {
	for {
		puyos := <-field
		fmt.Printf("%v\n", puyos)
	}
}

func main() {
	result := [15]int{}
	field := make(chan [15]int)
	go check2(field)
	Gen(&result, 0, field)
}

func main2() {
	param := flag.String("param", "a78", "puyofu")
	out := flag.String("out", "", "output file path")
	flag.Parse()
	bf := puyo2.NewBitFieldWithMattulwan(*param)
	bf.ShowDebug()
	if *out == "" {
		*out = *param + ".png"
	}
	bf.ExportImage(*out)
	fmt.Println(bf.MattulwanEditorUrl())
}

func readConfig() *Config {
	var conf Config
	toml.DecodeFile("puyo.toml", &conf)
	if conf.PuyoConfig.SearchLocations == nil {
		conf.PuyoConfig.SearchLocations = [][2]int{{2, 4}, {3, 4}, {4, 4}, {5, 4}, {2, 3}, {3, 3}, {4, 3}, {5, 3}, {3, 2}, {4, 2}, {5, 2}, {2, 1}, {3, 1}, {4, 1}, {5, 1}}
	}
	if conf.PuyoConfig.Mappings == nil {
		conf.PuyoConfig.Mappings = map[string]int{"a": 4, "b": 5, "c": 6, "d": 7, " ": 0}
	}
	return &conf
}

func check(num int, config *Config, field <-chan []int, wg *sync.WaitGroup) {
	cnt := 0
	resultFields := make(map[string]struct{})
	for {
		puyos := <-field
		if len(puyos) == 0 {
			break
		}
		if cnt%10000000 == 0 {
			fmt.Fprintf(os.Stderr, "[%d] %s %d\n", num, time.Now().String(), cnt)
		}
		bf := puyo2.NewBitFieldWithMattulwan(config.PuyoConfig.Field)
		if cnt%10000000 == 0 {
			fmt.Fprintf(os.Stderr, "%v\n", puyos)
		}
		vanish := puyo2.NewFieldBits()
		for i := 0; i < len(puyos); i++ {
			if puyos[i] != 0 {
				bf.SetColor(puyo2.Color(puyos[i]), config.PuyoConfig.SearchLocations[i][0], config.PuyoConfig.SearchLocations[i][1])
			} else {
				vanish.SetOnebit(config.PuyoConfig.SearchLocations[i][0], config.PuyoConfig.SearchLocations[i][1])
			}
		}
		beforeDrop := bf.Clone()
		bf.Drop(vanish)
		if bf.Equals(beforeDrop) {
			sbf := bf.Clone()
			result := sbf.Simulate()
			if result.Chains == 4 {
				nbf := bf.Normalize()
				param := nbf.MattulwanEditorParam()
				_, exist := resultFields[param]
				if !exist {
					resultFields[param] = struct{}{}
					fmt.Println(nbf.MattulwanEditorParam())
					nbf.ExportImage("verify/" + param + ".png")
					cnt++
				}

			}
		}
		cnt++
	}
	fmt.Fprintf(os.Stderr, "[%d], wg.Done\n", num)
	wg.Done()
}

func main_bak() {
	config := readConfig()
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	var w sync.WaitGroup
	field4 := make(chan []int)
	field5 := make(chan []int)
	defer close(field4)
	defer close(field5)

	w.Add(4)
	go check(4, config, field4, &w)
	go ParallelPermute([]int{4, 4, 4, 0, 5, 5, 5, 5, 0, 0, 0}, 4, field4, &w)
	go check(5, config, field5, &w)
	go ParallelPermute([]int{4, 4, 4, 4, 5, 5, 5, 0, 0, 0, 0}, 5, field5, &w)

	w.Wait()
}
