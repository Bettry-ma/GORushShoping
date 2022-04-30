package main

import (
	"fmt"
)

func main() {
	sl := make([]int, 8)
	sl = []int{1, 2, 3}
	fmt.Println(sl)
	sl = append(sl, 9, 8, 7)
	fmt.Println(sl)
	sts := make([]byte, 99)
	sts = []byte{'h', 'e', 'l', 'l', 'o'}
	for _, c := range sts {
		fmt.Printf("%c", c)
	}
	fmt.Println()
	sts = append(sts, "world"...)
	for _, c := range sts {
		fmt.Printf("%c", c)
	}
	fmt.Println()
	newsl := make([]byte, 99)
	i := copy(newsl, sts)
	fmt.Println(i)
	for _, c := range newsl {
		fmt.Printf("%c", c)
	}
}
