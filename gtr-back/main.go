package main

import (
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

func check(config *Config, field <-chan [15]int, wg *sync.WaitGroup) {
	cnt := 0
	resultFields := make(map[string]struct{})
	for {
		puyos := <-field
		if len(puyos) == 0 {
			break
		}
		if cnt%10000000 == 0 {
			fmt.Fprintf(os.Stderr, "%s %d\n", time.Now().String(), cnt)
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
				vanish.Onebit(config.PuyoConfig.SearchLocations[i][0], config.PuyoConfig.SearchLocations[i][1])
			}
		}
		beforeDrop := bf.Clone()
		bf.Drop(vanish)
		if bf.Equals(beforeDrop) {
			sbf := bf.Clone()
			result := sbf.Simulate()
			if result.Chains == 4 && result.BitField.IsEmpty() {
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
	fmt.Fprintf(os.Stderr, "wg.Done\n")
	wg.Done()
}

func main() {
	config := readConfig()
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	var w sync.WaitGroup
	field := make(chan [15]int)
	defer close(field)
	go check(config, field, &w)
	Gen(field)
}
