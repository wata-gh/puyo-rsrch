package main

func Gen(result *[15]int, i int, field chan<- [15]int) {
	if i == 15 {
		// fmt.Printf("%v\n", result)
		field <- *result
		return
	}
	for j := 0; j < 4; j++ {
		result[i] = j
		Gen(result, i+1, field)
	}
}
