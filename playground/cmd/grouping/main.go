package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/wata-gh/puyo2"
)

func genTargets() []*puyo2.FieldBits {
	targets := []*puyo2.FieldBits{}
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				for l := 0; l < 2; l++ {
					for m := 0; m < 2; m++ {
						for n := 0; n < 2; n++ {
							target := puyo2.NewFieldBits()
							target.SetOnebit(0, 1)
							target.SetOnebit(1, 1)
							target.SetOnebit(2, 1)
							if i == 1 {
								target.SetOnebit(0, 2)
							}
							if j == 1 {
								target.SetOnebit(1, 2)
							}
							if k == 1 {
								target.SetOnebit(2, 2)
							}
							if l == 1 {
								target.SetOnebit(0, 3)
							}
							if m == 1 {
								target.SetOnebit(1, 3)
							}
							if n == 1 {
								target.SetOnebit(2, 3)
							}
							if l == 1 && i != 1 {
								continue
							}
							if m == 1 && j != 1 {
								continue
							}
							if n == 1 && k != 1 {
								continue
							}
							targets = append(targets, target)
						}
					}
				}
			}
		}
	}
	return targets
}

func groupByPlace() {
	targets := genTargets()
	files := map[string]*os.File{}

	for _, target := range targets {
		m0 := target.ToIntArray()[0]
		name := fmt.Sprintf("%d_0.txt", m0)
		f, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		files[name] = f
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		param := strings.TrimRight(scanner.Text(), "\n")
		sbf := puyo2.NewShapeBitFieldWithFieldString(param)

		for _, target := range targets {
			m0 := target.ToIntArray()[0]
			name := fmt.Sprintf("%d_0.txt", m0)
			shapes := []uint64{}
			for _, shape := range sbf.Shapes {
				check := shape.And(target)
				if check.IsEmpty() == false {
					shapes = append(shapes, check.ToIntArray()[0])
				}
			}
			sort.Slice(shapes, func(i, j int) bool {
				return shapes[i] < shapes[j]
			})
			var b strings.Builder
			for i, shape := range shapes {
				if i != 0 {
					fmt.Fprint(&b, ":")
				}
				fmt.Fprintf(&b, "%d_0", shape)
			}
			fmt.Fprintln(files[name], b.String(), param)
		}
	}
}

func groupBySimulation() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		param := strings.TrimRight(scanner.Text(), "\n")
		sbf := puyo2.NewShapeBitFieldWithFieldString(param)

		var b strings.Builder
		for sbf.IsEmpty() == false {
			vfbn := sbf.FindVanishingFieldBitsNum()[0]
			v := sbf.Shapes[vfbn]
			ary := v.ToIntArray()
			if b.Len() != 0 {
				fmt.Fprint(&b, ":")
			}
			key := fmt.Sprintf("%d_%d", ary[0], ary[1])
			fmt.Fprint(&b, key)
			sbf.Simulate1()
		}
		fmt.Println(b.String(), param)
	}
}

func groupByChainOrder() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		param := strings.TrimRight(scanner.Text(), "\n")
		sbf := puyo2.NewShapeBitFieldWithFieldString(param)
		sbf.Simulate()
		var b strings.Builder
		for i := len(sbf.ChainOrderedShapes) - 1; i >= 0; i-- {
			ary := sbf.ChainOrderedShapes[i][0].ToIntArray()
			key := fmt.Sprintf("%d_%d", ary[0], ary[1])
			fmt.Fprint(&b, key)
			if i != 0 {
				fmt.Fprintf(&b, ":")
			}
		}
		fmt.Println(b.String(), param)
	}
}

func groupByOverall() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		param := strings.TrimRight(scanner.Text(), "\n")
		sbf := puyo2.NewShapeBitFieldWithFieldString(param)
		ary := sbf.OriginalOverallShape().ToIntArray()
		key := fmt.Sprintf("%d_%d", ary[0], ary[1])
		fmt.Println(key, param)
	}
}

func main() {
	// groupBySimulation()
	groupByPlace()
}
