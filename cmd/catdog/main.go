package main

import "fmt"

func main() {
	var s = make(map[int][]int)
	s[0] = append(s[0], 1, 2, 3, 4)

	s0 := s[0]
	s0 = append(s0, 1234)
	fmt.Println(s[0])
}
