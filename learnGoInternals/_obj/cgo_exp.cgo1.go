// Created by cgo - DO NOT EDIT

//line /home/go/src/self_go_home/learnGoInternals/cgo_exp.go:1
package main

//line /home/go/src/self_go_home/learnGoInternals/cgo_exp.go:7
import "fmt"

func Random() int {
	return int(_Cfunc_random())
}

func Seed(i int) {
	_Cfunc_srandom(_Ctype_uint(i))
}

func main() {
	fmt.Println(Random())
}
