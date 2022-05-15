package puyorsrch

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/wata-gh/puyo2"
)

func Check(config *Config, field <-chan []int, satisfy chan<- puyo2.BitField, num int, wg *sync.WaitGroup) {
	cnt := 0
	resultFields := make(map[string]struct{})
	bbf := puyo2.NewBitFieldWithMattulwan(config.PuyoConfig.Field)
	for {
		puyos := <-field
		if len(puyos) == 0 {
			break
		}
		if cnt%10000000 == 0 {
			fmt.Fprintf(os.Stderr, "[%d]%s %d\n", num, time.Now().String(), cnt)
			fmt.Fprintf(os.Stderr, "[%d]%v\n", num, puyos)
		}
		bf := bbf.Clone()
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
			result := bf.SimulateWithNewBitField()
			if result.Chains == config.PuyoConfig.ExpectedChains && result.BitField.IsEmpty() == config.PuyoConfig.AllClear {
				nbf := bf.Normalize()
				param := nbf.MattulwanEditorParam()
				_, exist := resultFields[param]
				if !exist {
					resultFields[param] = struct{}{}
					satisfy <- *nbf
					cnt++
				}
			}
		}
		cnt++
	}
	// 終了のために空フィールドを送る
	satisfy <- *puyo2.NewBitField()
	fmt.Fprintf(os.Stderr, "[%d]wg.Done\n", num)
	wg.Done()
}
