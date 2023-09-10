package main

import "fmt"

func main1() {
	var arr = []int{1, 2, 3, 4, 5}
	slice := [5]int{}

	for i, v := range arr {
		slice[i] = v
	}

	fmt.Println(slice)
}
