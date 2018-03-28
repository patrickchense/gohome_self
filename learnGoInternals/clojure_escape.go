package main

func main()  {

}

type Cursor struct {
	X int
}

func f() *Cursor {
	var c Cursor
	c.X = 500
	noinline()
	return &c
}

func noinline() {
	println("noinline")
}
