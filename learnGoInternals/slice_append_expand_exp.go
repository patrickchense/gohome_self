package main

import "fmt"

func main() {
	b := make([]int, 1024)
	b = append(b, 99)
	fmt.Println("len:", len(b), "cap:", cap(b))
}
/*
len: 1025 cap: 1280 接近25%增加
 */