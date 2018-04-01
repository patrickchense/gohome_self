package main

import "fmt"

type Point struct { X, Y int}

const W = 640
const H = 480

func Center(point * Point)  {
	point.X = W /2
	point.Y = H/2
}
func NewPoint()  {
	p:=new(Point)
	Center(p)
	fmt.Println(p.X, p.Y)
}
/*
-bash-3.2$ go build -gcflags=-m escape_exp1.go
# command-line-arguments
./escape_exp1.go:10:6: can inline Center
./escape_exp1.go:16:8: inlining call to Center
./escape_exp1.go:10:21: Center point does not escape
./escape_exp1.go:17:15: p.X escapes to heap
./escape_exp1.go:17:20: p.Y escapes to heap
./escape_exp1.go:15:8: NewPoint new(Point) does not escape
./escape_exp1.go:17:13: NewPoint ... argument does not escape
# command-line-arguments
runtime.main_main·f: relocation target main.main not defined
runtime.main_main·f: undefined: "main.main"
 */
