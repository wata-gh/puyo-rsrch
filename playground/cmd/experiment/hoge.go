package main

var clusters [][]int = [][]int{
	{0, 0, 1, 1, 1, 1},
	{0, 1, 1, 1, 1, 1},
	{0, 2, 1, 1, 1, 1},

	{1, 0, 2, 1, 1},
	{1, 1, 2, 1, 1},
	{1, 2, 2, 1, 1},
	{1, 3, 2, 1, 1},

	{2, 0, 1, 1, 2},
	{2, 1, 1, 1, 2},
	{2, 2, 1, 1, 2},
	{2, 3, 1, 1, 2},

	{3, 0, 1, 2, 1},
	{3, 1, 1, 2, 1},
	{3, 2, 1, 2, 1},
	{3, 3, 1, 2, 1},

	{4, 0, 3, 1},
	{4, 1, 3, 1},
	{4, 2, 3, 1},
	{4, 3, 3, 1},
	{4, 4, 3, 1},

	{5, 0, 1, 3},
	{5, 1, 1, 3},
	{5, 2, 1, 3},
	{5, 3, 1, 3},
	{5, 4, 1, 3},

	{6, 0, 2, 2},
	{6, 1, 2, 2},
	{6, 2, 2, 2},
	{6, 3, 2, 2},
	{6, 4, 2, 2},

	{7, 0, 4},
	{7, 1, 4},
	{7, 2, 4},
	{7, 3, 4},
	{7, 4, 4},
	{7, 5, 4},
}

func add(mem []int, cluster []int) {
	offset := cluster[1]
	for i, v := range cluster[2:] {
		mem[i+offset] += v
	}
}
func checkOverflow(mem []int, field []int) bool {
	for i, v := range mem {
		if field[i] < v {
			return false
		}
	}
	return true
}

func fill5(field []int) [][][]int {
	results := [][][]int{}
	cnt := 0
	for _, c1 := range clusters {
		m1 := []int{0, 0, 0, 0, 0, 0}
		add(m1, c1)
		if checkOverflow(m1, field) == false {
			continue
		}
		for _, c2 := range clusters {
			m2 := make([]int, len(m1))
			copy(m2, m1)
			add(m2, c2)
			if checkOverflow(m2, field) == false {
				continue
			}
			for _, c3 := range clusters {
				m3 := make([]int, len(m1))
				copy(m3, m2)
				add(m3, c3)
				if checkOverflow(m3, field) == false {
					continue
				}
				for _, c4 := range clusters {
					m4 := make([]int, len(m1))
					copy(m4, m3)
					add(m4, c4)
					if checkOverflow(m4, field) == false {
						continue
					}
					for _, c5 := range clusters {
						m5 := make([]int, len(m1))
						copy(m5, m4)
						add(m5, c5)
						if checkOverflow(m5, field) == false {
							continue
						}
						for i, v := range m5 {
							if field[i] != v {
								break
							}
						}
						cnt++
						results = append(results, [][]int{c1, c2, c3, c4, c5})
						// fmt.Println(c1, c2, c3, c4, c5, m5)
					}
				}
			}
		}
	}
	// fmt.Println(cnt)
	// fmt.Println(ng)
	return results
}

func fill4(field []int) [][][]int {
	results := [][][]int{}
	cnt := 0
	for _, c1 := range clusters {
		m1 := []int{0, 0, 0, 0, 0, 0}
		add(m1, c1)
		if checkOverflow(m1, field) == false {
			continue
		}
		for _, c2 := range clusters {
			m2 := make([]int, len(m1))
			copy(m2, m1)
			add(m2, c2)
			if checkOverflow(m2, field) == false {
				continue
			}
			for _, c3 := range clusters {
				m3 := make([]int, len(m1))
				copy(m3, m2)
				add(m3, c3)
				if checkOverflow(m3, field) == false {
					continue
				}
				for _, c4 := range clusters {
					m4 := make([]int, len(m1))
					copy(m4, m3)
					add(m4, c4)
					if checkOverflow(m4, field) == false {
						continue
					}
					for i, v := range m4 {
						if field[i] != v {
							break
						}
					}
					cnt++
					results = append(results, [][]int{c1, c2, c3, c4})
					// fmt.Println(c1, c2, c3, c4, m4)
				}
			}
		}
	}
	// fmt.Println(cnt)
	return results
}

func fill3(field []int) [][][]int {
	results := [][][]int{}
	cnt := 0
	for _, c1 := range clusters {
		m1 := []int{0, 0, 0, 0, 0, 0}
		add(m1, c1)
		if checkOverflow(m1, field) == false {
			continue
		}
		for _, c2 := range clusters {
			m2 := make([]int, len(m1))
			copy(m2, m1)
			add(m2, c2)
			if checkOverflow(m2, field) == false {
				continue
			}
			for _, c3 := range clusters {
				m3 := make([]int, len(m1))
				copy(m3, m2)
				add(m3, c3)
				if checkOverflow(m3, field) == false {
					continue
				}
				for i, v := range m3 {
					if field[i] != v {
						break
					}
				}
				results = append(results, [][]int{c1, c2, c3})
				cnt++
				// fmt.Println(c1, c2, c3, c4, m4)
			}
		}
	}
	// fmt.Println(cnt)
	return results
}

func fill2(field []int) [][][]int {
	results := [][][]int{}
	cnt := 0
	for _, c1 := range clusters {
		m1 := []int{0, 0, 0, 0, 0, 0}
		add(m1, c1)
		if checkOverflow(m1, field) == false {
			continue
		}
		for _, c2 := range clusters {
			m2 := make([]int, len(m1))
			copy(m2, m1)
			add(m2, c2)
			if checkOverflow(m2, field) == false {
				continue
			}
			for i, v := range m2 {
				if field[i] != v {
					break
				}
			}
			// fmt.Println(c1, c2, m4)
			results = append(results, [][]int{c1, c2})
			cnt++
		}
	}
	// fmt.Println(cnt)
	return results
}
